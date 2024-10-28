package gigachat

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c, err := NewInsecureClient("api-client-id", "api-client-secret")
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
	}
	assert.Equal(t, c.config.ClientId, "api-client-id")
	assert.Equal(t, c.config.ClientSecret, "api-client-secret")
	assert.Equal(t, c.config.BaseUrl, BaseUrl)
	assert.NotNil(t, c.token)
}

func TestAuth(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("req error", func(t *testing.T) {
		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		assert.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, c.config.AuthUrl+OAuthPath, httpmock.NewStringResponder(http.StatusBadGateway, ""))

		c.client = http.DefaultClient
		err = c.Auth()
		assert.EqualError(t, err, "unexpected status code 502")
	})
	t.Run("json error", func(t *testing.T) {
		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		assert.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, c.config.AuthUrl+OAuthPath, httpmock.NewStringResponder(http.StatusOK, "{ddd"))

		c.client = http.DefaultClient
		err = c.Auth()
		assert.EqualError(t, err, "invalid character 'd' looking for beginning of object key string")
	})
	t.Run("pass", func(t *testing.T) {
		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		assert.NoError(t, err)

		n := time.Now().Add(time.Minute)

		httpmock.RegisterResponder(http.MethodPost, c.config.AuthUrl+OAuthPath, httpmock.NewStringResponder(http.StatusOK, fmt.Sprintf(`{"access_token": "test", "expires_at": %d}`, n.UnixMilli())))

		c.client = http.DefaultClient
		err = c.Auth()

		assert.NoError(t, err)
		assert.True(t, c.token.Active())
		assert.Equal(t, "test", c.token.Get())
	})
}
