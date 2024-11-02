package gigachat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ChatRequest struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Temperature       *float64  `json:"temperature"`
	TopP              *float64  `json:"top_p"`
	N                 *int64    `json:"n"`
	Stream            *bool     `json:"stream"`
	MaxTokens         *int64    `json:"max_tokens"`
	RepetitionPenalty *float64  `json:"repetition_penalty"`
	UpdateInterval    *int64    `json:"update_interval"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Model   string   `json:"model"`
	Created int64    `json:"created"`
	Method  string   `json:"object"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int64  `json:"index"`
	FinishReason string `json:"finish_reason"`
	Message      Message
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

func (c *Client) Chat(in *ChatRequest) (*ChatResponse, error) {
	return c.ChatWithContext(context.Background(), in)
}

func (c *Client) ChatWithContext(ctx context.Context, in *ChatRequest) (*ChatResponse, error) {
	reqBytes, _ := json.Marshal(in)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseUrl+ChatPath, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var chatResponse ChatResponse
	if err := json.NewDecoder(res.Body).Decode(&chatResponse); err != nil {
		return nil, err
	}

	return &chatResponse, nil
}
