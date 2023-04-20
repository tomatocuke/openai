package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"openai/internal/config"
	"runtime"
	"strings"
	"sync"
	"time"
)

type user struct {
	// id string

	question struct {
		counter int64
		value   string
		doing   bool
	}

	answer struct {
		counter int64
		mu      sync.Mutex
		buffer  bytes.Buffer
	}
}

var (
	userCache sync.Map
)

func Query(uid string, msg string, timeout time.Duration) (reply string) {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(3)
			log.Println("ERROR:", err, file, line)
		}
	}()
	defer printLog(msg, &reply, start)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	msg = strings.TrimSpace(msg)
	length := len([]rune(msg))

	if length <= 1 {
		reply = "请说详细些..."
		return
	}
	if length > config.OpenAI.MaxQuestionLength {
		reply = "问题字数超出设定限制，请精简问题"
		return
	}

	var u *user
	v, ok := userCache.Load(uid)
	if !ok {
		u = &user{}
		userCache.Store(uid, u)
	} else {
		u = v.(*user)
	}

	if u.question.doing {
		reply = "上个问题正在处理中，请稍等..."
		return
	}

	u.question.doing = true
	defer func() {
		u.question.doing = false
	}()

	if msg != "继续" {
		// ++后，answer自动停止
		u.question.counter++
		u.question.value = msg

		// 发起请求
		err := completions(u)
		if err != nil {
			return err.Error()
		}
	} else {
		if u.question.counter == u.answer.counter && u.answer.buffer.Len() == 0 {
			reply = "没有更多内容啦，下一个问题。"
			return
		}
	}

	// ticker := time.NewTicker(1 * time.Second)
	var done bool
	for !done {
		select {
		// 超时结束
		case <-ctx.Done():
			done = true
		// 每秒检测结果是否完全返回，用于提前结束
		default:
			if u.question.counter == u.answer.counter {
				done = true
			}
		}
	}

	u.answer.mu.Lock()
	defer u.answer.mu.Unlock()

	str := u.answer.buffer.String()
	u.answer.buffer.Reset()

	// 出错
	if len(str) == 0 {
		reply = "openai请求超时"
		return
	}

	// 完全完成
	if u.question.counter == u.answer.counter {
		reply = str
		return
	}

	// 优化为以。？！\n结尾
	runes := []rune(str)
	for i := len(runes) - 1; i > 0; i-- {
		r := runes[i]
		if r == '。' || r == '\n' || r == '！' {
			u.answer.buffer.WriteString(string(runes[i+1:]))
			runes = runes[:i+1]
			break
		}
	}
	if runes[len(runes)-1] != '\n' {
		runes = append(runes, '\n')
	}
	runes = append(runes, []rune("【回复“继续”以滚动查看】")...)
	reply = string(runes)

	return
}

// https://beta.openai.com/docs/api-reference/making-requests
// 同步返回ch，异步读取数据流
func completions(u *user) error {

	respBody, err := postApi(u.question.value)
	if err != nil {
		return err
	}

	// 读取数据
	go func(u *user, respBody io.ReadCloser) {
		defer func() {

			u.answer.counter++
			respBody.Close()

			if err := recover(); err != nil {
				_, file, line, _ := runtime.Caller(3)
				log.Println("ERROR:", err, file, line)
			}
		}()

		scanner := bufio.NewScanner(respBody)
		scanner.Split(bufio.ScanLines)

		x := 0

		u.answer.buffer.Reset()

		for scanner.Scan() {
			bs := scanner.Bytes()
			// fmt.Println(string(bs))
			if len(bs) > 100 {
				x++
				bs = bs[6:]
				var r response
				json.Unmarshal(bs, &r)
				if len(r.Choices) == 0 {
					continue
				}
				tokenContent := r.Choices[0].Delta.Content
				// if config.Debug {
				// 	fmt.Print(tokenContent)
				// }

				u.answer.mu.Lock()
				if u.question.counter-u.answer.counter <= 1 {
					u.answer.buffer.WriteString(tokenContent)
				}
				u.answer.mu.Unlock()
				if u.question.counter-u.answer.counter > 1 {
					return
				}
			}
		}

		// if config.Debug {
		// 	fmt.Println("\n问题结束:", u.question.value, "回答完成， 行数：", x)
		// }

	}(u, respBody)

	return nil
}

func postApi(msg string) (io.ReadCloser, error) {
	r := newRequest(msg)
	bs, _ := json.Marshal(&r)

	client := &http.Client{Timeout: time.Second * 200}
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
	if err == nil && resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
	}

	return resp.Body, err
}

func printLog(question string, answer *string, start time.Time) {
	log.Printf(
		"用时:%ds \nQ: %s \nA: %s\n\n",
		int(time.Since(start).Seconds()),
		question,
		*answer,
	)
}
