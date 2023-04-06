package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"openai/internal/config"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	serverErrMsg = []byte("服务出错啦，稍后会修复")
	scrollMsg    = []byte("【回复“继续”以滚动查看】")
	resultCache  sync.Map
)

func init() {

}

// OpenAI可能无法在希望的时间内做出回复
// 使用goroutine + channel 的形式，不管是否能及时回复用户，后台都打印结果
func Query(uid string, msg string, timeout time.Duration) string {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(3)
			log.Println("ERROR:", err, file, line)
		}
	}()

	msg = strings.TrimSpace(msg)
	length := len([]rune(msg))

	if length <= 1 {
		return "请说详细些..."
	}
	if length > config.OpenAI.MaxQuestionLength {
		return "问题字数超出设定限制，请精简问题"
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var reply bytes.Buffer
	reply.Grow(4096)

	ch := make(chan []byte)
	v, ok := resultCache.Load(uid)

	if ok && msg == "继续" {
		ch = v.(chan []byte)
		b, ok := <-ch
		if !ok {
			return "没有啦，问下一个问题吧"
		}
		reply.Write(b)
	} else {
		resultCache.Store(uid, ch)
		go completions(ch, msg)
	}

loop:
	for {
		select {
		// 超时控制
		case <-ctx.Done():
			if reply.Len() <= 1 {
				log.Println("openai请求超时")
				reply.Write(serverErrMsg)
			} else {
				if reply.Bytes()[reply.Len()-1] != '\n' {
					reply.Write([]byte("\r\n\n"))
				}
				reply.Write(scrollMsg)
			}
			break loop
		// 不断读取
		case b, ok := <-ch:
			if !ok {
				break loop
			}
			reply.Write(b)
		}
	}

	s := reply.String()
	log.Printf(
		"用时:%ds \nQ: %s \nA: %s\n\n",
		int(time.Since(start).Seconds()),
		msg,
		s,
	)
	return s
}

// https://beta.openai.com/docs/api-reference/making-requests
func completions(ch chan []byte, msg string) {
	defer close(ch)
	timeout := 100 * time.Second
	// 防止ch泄漏
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	r := newRequest(msg)
	bs, _ := json.Marshal(&r)

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("POST", config.OpenAI.Params.Api, bytes.NewReader(bs))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.OpenAI.Key)

	// 设置代理
	if config.Http.Proxy != "" {
		proxyURL, _ := url.Parse(config.Http.Proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		ch <- serverErrMsg
		log.Println("ERROR", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)

	var buff bytes.Buffer
	buff.Grow(512)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			bs := scanner.Bytes()
			// data: {"id":"chatcmpl-706X266U3DAVJiwmo3gjPBfF0ZcZH","object":"chat.completion.chunk","created":1680259464,"model":"gpt-3.5-turbo-0301","choices":[{"delta":{"content":"充"},"index":0,"finish_reason":null}]}
			if len(bs) > 100 {
				bs = bs[6:]
				var r response
				json.Unmarshal(bs, &r)
				if len(r.Choices) == 0 {
					continue
				}

				tokenContent := r.Choices[0].Delta.Content
				// fmt.Print(tokenContent)
				buff.WriteString(tokenContent)
				if strings.Contains(tokenContent, "。") || strings.Contains(tokenContent, "\n") {
					ch <- buff.Bytes()
					buff.Reset()
				}
			}
		}
	}

	if buff.Len() > 0 {
		ch <- buff.Bytes()
	}
}
