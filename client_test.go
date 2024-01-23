package gigachat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("api-client-id", "api-client-secret")
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
	}
	assert.Equal(t, c.config.ClientId, "api-client-id")
	assert.Equal(t, c.config.ClientSecret, "api-client-secret")
	assert.Equal(t, c.config.BaseUrl, BaseUrl)
	assert.Nil(t, c.token)
	assert.Nil(t, c.exiresAt, nil)
}

func TestAuth(t *testing.T) {
	c, err := NewClient("api-client-id", "api-client-secret")
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
	}
	err = c.Auth()
	if err != nil {
		t.Errorf("Auth() error = %v", err)
	}
}
