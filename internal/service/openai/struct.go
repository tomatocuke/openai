package openai

import (
	"openai/internal/config"
)

type request struct {
	Model       string               `json:"model"`
	Messages    []requestMessageItem `json:"messages"`
	Temperature float32              `json:"temperature"`
	MaxTokens   uint16               `json:"max_tokens"`
	Stream      bool                 `json:"stream"`
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
	r.Stream = true
	r.Messages = []requestMessageItem{
		{Role: "user", Content: msg},
	}
	if conf.Prompt != "" {
		r.Messages = append(r.Messages, requestMessageItem{Role: "system", Content: conf.Prompt})
	}
	return
}

type response struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Delta delta `json:"delta"`
}

type delta struct {
	Content string `json:"content"`
}

// type response struct {
// 	ID    string `json:"id"`
// 	Usage struct {
// 		PromptTokens     int `json:"prompt_tokens"`
// 		CompletionTokens int `json:"completion_tokens"`
// 		TotalTokens      int `json:"total_tokens"`
// 	} `json:"usage"`
// 	// Object  string                 `json:"object"`
// 	// Created int                    `json:"created"`
// 	// Model   string                 `json:"model"`
// 	Choices []choiceItem `json:"choices"`
// 	// Usage   map[string]interface{} `json:"usage"`
// 	Error struct {
// 		Message string `json:"message"`
// 	} `json:"error"`
// }

// type choiceItem struct {
// 	Message struct {
// 		Content string `json:"content"`
// 	} `json:"message"`
// 	Delta struct {
// 		Content string `json:"content"`
// 	} `json:"delta"`
// }
