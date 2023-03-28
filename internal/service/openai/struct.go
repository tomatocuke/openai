package openai

import "openai/internal/config"

type request struct {
	Model       string               `json:"model"`
	Messages    []requestMessageItem `json:"messages"`
	Temperature float32              `json:"temperature"`
	MaxTokens   uint16               `json:"max_tokens"`
}

type requestMessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func newRequest(msg string) (r request) {
	conf := config.OpenAI.Params
	r.Model = conf.Model
	r.Temperature = conf.Temperature
	r.MaxTokens = conf.MaxTokens
	r.Messages = []requestMessageItem{
		{Role: "user", Content: msg},
	}
	if conf.Prompt != "" {
		r.Messages = append(r.Messages, requestMessageItem{Role: "system", Content: conf.Prompt})
	}
	return
}
