package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"openai/internal/config"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	api     = "https://api.openai.com/v1/completions"
	MsgWait = "这个问题比较复杂，再稍等一下～"
)

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

type response struct {
	ID string `json:"id"`
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
	Text string `json:"text"`
	// Index        int    `json:"index"`
	// Logprobs     int    `json:"logprobs"`
	// FinishReason string `json:"finish_reason"`
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
		start := time.Now()
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
		log.Printf("用时%ds，「%s」 \n %s \n\n", int(time.Since(start).Seconds()), msg, result)
	}(msg, ctx, ch)

	var result string
	select {
	case result = <-ch:
	case <-ctx.Done():
		result = "超时啦，请稍等20-60s后再问我，就告诉你。"
	}

	return result
}

// https://beta.openai.com/docs/api-reference/making-requests
func completions(msg string, timeout time.Duration) (string, error) {
	msg, wordSize := queryMsg(msg)
	if wordSize == 0 {
		return "请说详细些...", nil
	}

	// fmt.Println("openai请求内容：", msg)
	params := map[string]interface{}{
		"model":  "text-davinci-003",
		"prompt": msg,
		// 影响回复速度和内容长度。  回复长度耗费token，影响花费的金额
		"max_tokens": wordSize * 3,
		// 0-1，默认1，越高越有创意
		"temperature": 0.8,
		// "top_p":             1,
		// "frequency_penalty": 0,
		// "presence_penalty":  0,
		// "stop": "。",
	}

	bs, _ := json.Marshal(params)

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("POST", api, bytes.NewReader(bs))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.ApiKey)

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
		return replyMsg(data.Choices[0].Text), nil
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
