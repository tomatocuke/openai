package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"openai/internal/config"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	api          = "https://api.openai.com/v1/chat/completions"
	MsgWait      = "这个问题比较复杂，再稍等一下～"
	exchangeRate = 6.9
)

var totaltokens int64

var (
	// 结果缓存（主要用于超时，用户重新提问后能给出答案）
	resultCache sync.Map
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Hour * 24)
		for range ticker.C {
			resultCache = sync.Map{}
		}
	}()
}

type request struct {
	Model    string       `json:"model"`
	Messages []reqMessage `json:"messages"`
}
type reqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	ID    string `json:"id"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	// Object  string                 `json:"object"`
	// Created int                    `json:"created"`
	// Model   string                 `json:"model"`
	Choices []choiceItem `json:"choices"`
	// Usage   map[string]interface{} `json:"usage"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type choiceItem struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// OpenAI可能无法在希望的时间内做出回复
// 使用goroutine + channel 的形式，不管是否能及时回复用户，后台都打印结果
func Query(msg string, timeout time.Duration) string {
	ch := make(chan string, 1)
	ctx, candel := context.WithTimeout(context.Background(), timeout)
	defer candel()

	cacheVal, ok := resultCache.Load(msg)
	// 非第一次问这个问题，得到问题的缓存，直接返回。
	if ok {
		s := cacheVal.(string)
		if s != MsgWait {
			resultCache.Delete(msg)
		}
		return s
	}

	// 第一次问，缓存「等待...」，发起openai请求。
	// 如果在限定时间内得到答复，返回给用户。没及时答复，告诉用户超时。

	resultCache.Store(msg, MsgWait)

	go func(msg string, ctx context.Context, ch chan string) {
		result, err := completions(msg, time.Second*180)
		if err != nil {
			result = "发生错误「" + err.Error() + "」，您重试一下"
		}
		select {
		case <-ctx.Done():
			resultCache.Store(msg, result)
		default:
			ch <- result
			resultCache.Delete(msg)
		}
		close(ch)
	}(msg, ctx, ch)

	var result string
	select {
	case result = <-ch:
	case <-ctx.Done():
		result = "请稍等10s后复制问题再问我一遍"
	}

	return result
}

// https://beta.openai.com/docs/api-reference/making-requests
func completions(msg string, timeout time.Duration) (string, error) {
	start := time.Now()
	msg = strings.TrimSpace(msg)
	length := len([]rune(msg))
	if length <= 1 {
		return "请说详细些...", nil
	}
	var r request
	r.Model = "gpt-3.5-turbo"
	r.Messages = []reqMessage{{
		Role:    "user",
		Content: msg,
	}}

	bs, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("POST", api, bytes.NewReader(bs))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.C.OpenAI.Key)

	// 设置代理
	if config.C.Http.Proxy != "" {
		proxyURL, _ := url.Parse(config.C.Http.Proxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data response
	json.Unmarshal(body, &data)
	if len(data.Choices) > 0 {
		atomic.AddInt64(&totaltokens, int64(data.Usage.TotalTokens))

		reply := replyMsg(data.Choices[0].Message.Content)
		log.Printf("本次:用时:%ds,花费约:%f¥,token:%d,请求:%d,回复:%d。 服务启动至今累计花费约:%f¥ \nQ:%s \nA:%s \n",
			int(time.Since(start).Seconds()),
			float32(data.Usage.TotalTokens/1000)*0.002*exchangeRate,
			data.Usage.TotalTokens,
			data.Usage.PromptTokens,
			data.Usage.CompletionTokens,
			float32(totaltokens/1000)*0.002*exchangeRate,
			msg,
			reply,
		)

		return reply, nil
	}

	return data.Error.Message, nil
}

func queryMsg(prompt string) (string, int) {
	msg := strings.TrimSpace(prompt)
	wordSize := 0
	length := len([]rune(msg))
	if length <= 1 {
		wordSize = 0
	} else if length <= 3 {
		msg = "30字以内说说:" + msg
		wordSize = 30
	} else if length <= 5 {
		msg = "99字以内说说:" + msg
		wordSize = 100
	} else {
		// 默认400字以内
		wordSize = 400
	}

	// 检查规定 xx字
	if idx := strings.IndexRune(msg, '字'); idx > -1 {
		if idx > 3 && string(msg[idx-3:idx]) == "个" {
			idx -= 3
		}
		end := idx
		start := idx
		for i := idx - 1; i >= 0; i-- {
			if msg[i] <= '9' && msg[i] >= '0' {
				start = i
			} else {
				break
			}
		}

		if start != end {
			wordSize, _ = strconv.Atoi(msg[start:end])
			if wordSize == 0 {
				wordSize = 400
			} else if wordSize > 800 {
				wordSize = 800
			}
		}

	}

	return msg, wordSize
}

func replyMsg(reply string) string {
	idx := strings.Index(reply, "\n\n")
	if idx > -1 && reply[len(reply)-2] != '\n' {
		reply = reply[idx+2:]
	}
	start := 0
	for i, v := range reply {
		if !isSymbaol(v) {
			start = i
			break
		}
	}

	return reply[start:]
}

var symbols = []rune{'\n', ' ', '，', '。', '？', '?', ',', '.', '!', '！', ':', '：'}

func isSymbaol(w rune) bool {
	for _, v := range symbols {
		if v == w {
			return true
		}
	}
	return false
}
