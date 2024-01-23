package gigachat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) Chat(in *ChatRequest) (*ChatResponse, error) {
	return c.ChatWithContext(context.Background(), in)
}

func (c *Client) ChatWithContext(ctx context.Context, in *ChatRequest) (*ChatResponse, error) {

	reqBytes, _ := json.Marshal(in)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ChatPath, bytes.NewReader(reqBytes))
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
