package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (a *App) AskAI(provider, model, sysPrompt, userPrompt string) (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}
	key := cfg.APIKeys[provider]
	if key == "" {
		return "", fmt.Errorf("API key for %s is not configured in Settings", provider)
	}
	if provider == "openai" {
		return askOpenAI(key, model, sysPrompt, userPrompt)
	} else if provider == "anthropic" {
		return askAnthropic(key, model, sysPrompt, userPrompt)
	}
	return "", fmt.Errorf("unsupported provider: %s", provider)
}

func askOpenAI(key, model, sysPrompt, userPrompt string) (string, error) {
	if model == "" {
		model = "gpt-4o-mini"
	}
	reqBody := AIRequest{
		Model: model,
		Messages: []AIMessage{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: userPrompt},
		},
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errData)
		return "", fmt.Errorf("OpenAI API error (HTTP %d): %v", resp.StatusCode, errData)
	}
	var respData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}
	if len(respData.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}
	return respData.Choices[0].Message.Content, nil
}

func askAnthropic(key, model, sysPrompt, userPrompt string) (string, error) {
	if model == "" {
		model = "claude-3-5-haiku-20241022"
	}
	type AnthropicRequest struct {
		Model     string      `json:"model"`
		MaxTokens int         `json:"max_tokens"`
		System    string      `json:"system,omitempty"`
		Messages  []AIMessage `json:"messages"`
	}
	reqBody := AnthropicRequest{
		Model:     model,
		MaxTokens: 2048,
		System:    sysPrompt,
		Messages: []AIMessage{
			{Role: "user", Content: userPrompt},
		},
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errData)
		return "", fmt.Errorf("Anthropic API error (HTTP %d): %v", resp.StatusCode, errData)
	}
	var respData struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}
	if len(respData.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic")
	}
	return respData.Content[0].Text, nil
}
