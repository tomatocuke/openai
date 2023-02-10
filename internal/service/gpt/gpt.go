package gpt

import (
	"bytes"
	"chatgpt/config"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	api = "https://api.openai.com/v1/completions"
)

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
func Query(isFast bool, msg string, timeout time.Duration) string {
	start := time.Now()
	ch := make(chan string, 1)
	ctx, candel := context.WithTimeout(context.Background(), timeout)
	defer candel()

	go func() {
		defer close(ch)
		result, err := completions(isFast, msg, time.Second*100)
		if err != nil {
			result = "发生错误「" + err.Error() + "」，您重试一下"
		}
		ch <- result
		// 超时，内容未通过接口及时回复，打印内容及总用时
		since := time.Since(start)
		if since > timeout {
			log.Printf("超时%ds，「%s」-「%s」\n", int(since.Seconds()), msg, result)
		}
	}()

	var result string
	select {
	case result = <-ch:
	case <-ctx.Done():
		result = "超时啦"
	}

	log.Printf("用时%ds，「%s」-「%s」\n", int(time.Since(start).Seconds()), msg, result)

	return result
}

// https://beta.openai.com/docs/api-reference/making-requests
func completions(isFast bool, msg string, timeout time.Duration) (string, error) {
	wordSize := 35 // 中文字符数量
	temperature := 0.3

	if !isFast {
		wordSize = 800
		temperature = 0.8
	}

	// start := time.Now()
	params := map[string]interface{}{
		"model":  "text-davinci-003",
		"prompt": msg,
		// 影响回复速度和内容长度。小则快，但内容短，可能是截断的。
		"max_tokens": wordSize * 3,
		// 0-1，默认1，越高越有创意
		"temperature": temperature,
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
		result := strings.TrimPrefix(data.Choices[0].Text, "？")
		return strings.TrimSpace(result), nil
	}

	return data.Error.Message, nil
}
