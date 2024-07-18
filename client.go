package gigachat

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	client *http.Client
	config *Config
	token  *Token
}

type Config struct {
	AuthUrl      string
	BaseUrl      string
	ClientId     string
	ClientSecret string
	Scope        string
	Insecure     bool
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func NewClient(clientId string, clientSecret string) (*Client, error) {
	var conf = &Config{
		AuthUrl:      AuthUrl,
		BaseUrl:      BaseUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        ScopeApiIndividual,
		Insecure:     false,
	}
	return NewClientWithConfig(conf)
}

// NewInsecureClient creates a new GigaChat client with InsecureSkipVerify because GigaChat uses a weird certificate authority.
func NewInsecureClient(clientId string, clientSecret string) (*Client, error) {
	var conf = &Config{
		AuthUrl:      AuthUrl,
		BaseUrl:      BaseUrl,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        ScopeApiIndividual,
		Insecure:     true,
	}
	return NewClientWithConfig(conf)
}

// NewClientWithConfig creates a new GigaChat client with the specified configuration.
func NewClientWithConfig(config *Config) (*Client, error) {
	var customTransport *http.Transport

	if config.Insecure {
		if dt, ok := http.DefaultTransport.(*http.Transport); ok {
			customTransport = dt.Clone()
			customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}

	return &Client{
		client: &http.Client{Transport: customTransport},
		config: config,
		token:  new(Token),
	}, nil
}

func (c *Client) Auth() error {
	return c.AuthWithContext(context.Background())
}

func (c *Client) AuthWithContext(ctx context.Context) error {
	if c.token.Active() {
		return nil
	}

	payload := strings.NewReader("scope=" + c.config.Scope)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.AuthUrl+OAuthPath, payload)
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

	c.token.Set(oauth.AccessToken, time.UnixMilli(oauth.ExpiresAt))
	return nil
}

func (c *Client) sendRequest(_ context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.Get()))
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
