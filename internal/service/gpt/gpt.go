package gpt

import (
	"bytes"
	"chatgpt/config"
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

// https://beta.openai.com/docs/api-reference/making-requests
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

	bs, _ := json.Marshal(params)
	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("POST", api, bytes.NewReader(bs))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if !strings.Contains(err.Error(), "Timeout") {
			log.Println("ERROR:", err)
		} else {
			log.Println("超时:", msg)
		}
		return ""
	}

	var data response
	json.Unmarshal(body, &data)
	if data.Error.Message != "" {
		log.Println("ERROR:", data.Error.Message)
		return ""
	}

	if len(data.Choices) > 0 {
		return strings.TrimSpace(data.Choices[0].Text)
	}

	return ""
}
