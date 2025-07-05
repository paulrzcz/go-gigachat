package gigachat

import (
	"context"
	"encoding/json"
	"net/http"
)

type ModelListResponse struct {
	Models []Model `json:"data"`
	Type   string  `json:"object"`
}

type Model struct {
	Id      string `json:"id"`
	Type    string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

type FileResponse struct {
	Bytes        int    `json:"bytes,omitempty"`
	CreatedAt    int    `json:"created_at,omitempty"`
	Filename     string `json:"filename,omitempty"`
	Id           string `json:"id,omitempty"`
	Object       string `json:"object,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
	AccessPolicy string `json:"access_policy,omitempty"`
	Status       int    `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
}

type FilesInfo struct {
	Data []struct {
		Id           string `json:"id"`
		Object       string `json:"object"`
		Bytes        int    `json:"bytes"`
		AccessPolicy string `json:"access_policy"`
		CreatedAt    int    `json:"created_at"`
		Filename     string `json:"filename"`
		Purpose      string `json:"purpose"`
	} `json:"data"`
}

func (c *Client) Models() (*ModelListResponse, error) {
	return c.ModelsWithContext(context.Background())
}

func (c *Client) ModelsWithContext(ctx context.Context) (*ModelListResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.config.BaseUrl+ModelsPath, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var modelListResponse ModelListResponse
	if err := json.NewDecoder(res.Body).Decode(&modelListResponse); err != nil {
		return nil, err
	}

	return &modelListResponse, nil
}

func (c *Client) Model(model string) (*Model, error) {
	return c.ModelWithContext(context.Background(), model)
}

func (c *Client) ModelWithContext(ctx context.Context, model string) (*Model, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ModelsPath+"/"+model, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var modelResponse Model
	if err := json.NewDecoder(res.Body).Decode(&modelResponse); err != nil {
		return nil, err
	}

	return &modelResponse, nil
}
