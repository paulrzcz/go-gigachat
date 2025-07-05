package gigachat

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func Test_UploadFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("file not found", func(t *testing.T) {
		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		if !assert.NoError(t, err) {
			return
		}

		_, err = c.UploadFile(context.Background(), "test")
		assert.EqualError(t, err, "file \"test\" niot found")
	})
	t.Run("bad type", func(t *testing.T) {
		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		if !assert.NoError(t, err) {
			return
		}

		f, _ := os.CreateTemp("", "*.doc")
		_ = f.Close()
		defer os.Remove(f.Name())

		_, err = c.UploadFile(context.Background(), f.Name())
		assert.EqualError(t, err, "only jpeg, png, tiff, bmp file types are supported")
	})
	t.Run("bad size", func(t *testing.T) {
		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		if !assert.NoError(t, err) {
			return
		}

		f, _ := os.CreateTemp("", "*.doc")
		_, _ = f.WriteString("dcsnkcskjdvndfnvdknvkjgfnvkjngfkjvn")
		_ = f.Close()
		defer os.Remove(f.Name())

		maxImgSize = 10
		_, err = c.UploadFile(context.Background(), f.Name())
		assert.EqualError(t, err, "the maximum allowed file size is 10 bytes")
	})
	t.Run("error req", func(t *testing.T) {

		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		if !assert.NoError(t, err) {
			return
		}

		f, _ := os.CreateTemp("", "*.png")
		_ = f.Close()
		defer os.Remove(f.Name())

		httpmock.RegisterResponder(http.MethodPost, c.config.BaseUrl+Files, httpmock.NewStringResponder(http.StatusBadGateway, ""))

		c.client = http.DefaultClient
		_, err = c.UploadFile(context.Background(), f.Name())
		assert.EqualError(t, err, "EOF")
	})
	t.Run("pass", func(t *testing.T) {

		c, err := NewInsecureClient("api-client-id", "api-client-secret")
		if !assert.NoError(t, err) {
			return
		}

		f, _ := os.CreateTemp("", "*.png")
		_ = f.Close()
		defer os.Remove(f.Name())

		httpmock.RegisterResponder(http.MethodPost, c.config.BaseUrl+Files, httpmock.NewStringResponder(http.StatusOK, `{
  "bytes": 120000,
  "created_at": 1677610602,
  "filename": "file123",
  "id": "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e",
  "object": "file",
  "purpose": "general",
  "access_policy": "private"
}`))

		c.client = http.DefaultClient
		id, err := c.UploadFile(context.Background(), f.Name())
		if assert.NoError(t, err) {
			assert.Equal(t, "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e", id)
		}
	})
}

func Test_GetFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, err := NewClient("api-client-id", "api-client-secret")
	if !assert.NoError(t, err) {
		return
	}

	t.Run("error", func(t *testing.T) {
		httpmock.RegisterResponder(http.MethodGet, c.config.BaseUrl+Files, httpmock.NewStringResponder(http.StatusBadGateway, ""))

		c.client = http.DefaultClient
		_, err = c.GetFiles(context.Background())
		assert.EqualError(t, err, "EOF")
	})
	t.Run("pass", func(t *testing.T) {
		httpmock.RegisterResponder(http.MethodGet, c.config.BaseUrl+Files, httpmock.NewStringResponder(http.StatusOK, `{
  "data": [
    {
      "bytes": 120000,
      "created_at": 1677610602,
      "filename": "file123",
      "id": "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e",
      "object": "file",
      "purpose": "general",
      "access_policy": "private"
    }
  ]
}`))

		c.client = http.DefaultClient
		data, err := c.GetFiles(context.Background())
		if assert.NoError(t, err) {
			assert.Len(t, data.Data, 1)
		}
	})
}

func Test_DeleteFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, err := NewClient("api-client-id", "api-client-secret")
	if !assert.NoError(t, err) {
		return
	}

	t.Run("error", func(t *testing.T) {
		httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/%s/delete", c.config.BaseUrl+Files, "12334"), httpmock.NewStringResponder(http.StatusBadGateway, ""))

		c.client = http.DefaultClient
		err = c.DeleteFile(context.Background(), "12334")
		assert.EqualError(t, err, "EOF")
	})
	t.Run("pass", func(t *testing.T) {
		httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/%s/delete", c.config.BaseUrl+Files, "12334"), httpmock.NewStringResponder(http.StatusOK, `{
  "id": "d3277ca1-a140-484a-a3b4-9a121bea4bdc",
  "deleted": true,
  "access_policy": "private"
}`))

		c.client = http.DefaultClient
		err = c.DeleteFile(context.Background(), "12334")
		assert.NoError(t, err)
	})
}
