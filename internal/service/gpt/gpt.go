package gpt

import (
	"chatgpt/config"
	"chatgpt/util"
	"strings"
	"time"
)

const (
	// https://beta.openai.com/docs/api-reference/making-requests
	api = "https://api.openai.com/v1/completions"
)

type response struct {
	ID string `json:"id"`
	// Object  string                 `json:"object"`
	// Created int                    `json:"created"`
	// Model   string                 `json:"model"`
	Choices []choiceItem `json:"choices"`
	// Usage   map[string]interface{} `json:"usage"`
}

type choiceItem struct {
	Text string `json:"text"`
	// Index        int    `json:"index"`
	// Logprobs     int    `json:"logprobs"`
	// FinishReason string `json:"finish_reason"`
}

func Completions(msg string, timeout time.Duration) string {
	// start := time.Now()
	params := map[string]interface{}{
		"model": "text-davinci-003",
		// "model":       "text-curie-001",
		"prompt":      msg,
		"max_tokens":  1024,
		"temperature": 0.2,
		// "top_p":             1,
		// "frequency_penalty": 0,
		// "presence_penalty":  0,
		// "stop":              "\n",
	}
	var resp response
	err := util.HttpPostJson(api, params).
		AddHeader("Authorization", "Bearer "+config.ApiKey).
		SetTimeout(timeout).DoTo(&resp)

	// fmt.Println("耗时:", int(time.Since(start).Seconds()))
	if err != nil {
		return ""
	}

	if len(resp.Choices) > 0 {
		return strings.TrimSpace(resp.Choices[0].Text)
	}

	return ""
}
