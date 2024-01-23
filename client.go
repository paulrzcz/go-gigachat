package gigachat

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Client struct {
	client   *http.Client
	config   *Config
	token    *string
	exiresAt *int64
}

func NewClient(clientId string, clientSecret string) (*Client, error) {
	var conf = &Config{
		AuthUrl:      AuthUrl,
		BaseUrl:      BaseUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        ScopeApiIndividual,
	}
	return NewClientWithConfig(conf)
}

func NewClientWithConfig(config *Config) (*Client, error) {
	return &Client{
		client: &http.Client{},
		config: config,
	}, nil
}

func (c *Client) Auth() error {
	return c.AuthWithContext(context.Background())
}

func (c *Client) AuthWithContext(ctx context.Context) error {
	if c.token != nil {
		return nil
	}

	payload := strings.NewReader("scope=" + c.config.Scope)
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.AuthUrl+OAuthPath, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("RqUID", uuid.NewString())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.config.ClientId+":"+c.config.ClientSecret)))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var oauth OAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&oauth)
	if err != nil {
		return err
	}

	c.token = &oauth.AccessToken
	c.exiresAt = &oauth.ExpiresAt

	return nil
}

func (c *Client) sendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		var errMessage interface{}
		if err := json.NewDecoder(res.Body).Decode(&errMessage); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("GigaCode API request failed: status Code: %d %s %s Message: %+v", res.StatusCode, res.Status, res.Request.URL, errMessage)
	}

	return res, nil
}
