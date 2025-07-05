package gigachat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

var (
	maxImgSize int64 = 15 * 1024 * 1024
)

// https://developers.sber.ru/docs/ru/gigachat/guides/working-with-files?ext=image

func (c *Client) checkImgType(filePath string) bool {
	if lo.Contains([]string{".jpg", ".jpeg", ".png", ".tiff", ".bmp"}, filepath.Ext(filePath)) {
		return true
	}
	return false
}

func (c *Client) UploadFile(ctx context.Context, filePath string) (string, error) {
	if info, err := os.Stat(filePath); err != nil {
		return "", errors.New(fmt.Sprintf("file %q niot found", filePath))
	} else if info.Size() > maxImgSize {
		return "", errors.New(fmt.Sprintf("the maximum allowed file size is %d bytes", maxImgSize))
	}

	if !c.checkImgType(filePath) {
		return "", errors.New("only jpeg, png, tiff, bmp file types are supported")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	resp, err := c.uploadFile(ctx, file)
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}

func (c *Client) DeleteFile(ctx context.Context, fileID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s/delete", c.config.BaseUrl+Files, fileID), nil)
	if err != nil {
		return err
	}

	res, err := c.sendRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (c *Client) GetFiles(ctx context.Context) (*FilesInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.config.BaseUrl+Files, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp FilesInfo
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) uploadFile(ctx context.Context, file *os.File) (*FileResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Поле 'file'
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(file.Name())))
	h.Set("Content-Type", mimeType(file))

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)

	// Поле 'purpose'
	_ = writer.WriteField("purpose", "general")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseUrl+Files, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.Get()))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp FileResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func mimeType(file *os.File) string {
	defer file.Seek(0, 0)

	// Определяем MIME-тип по первым 512 байтам
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	return http.DetectContentType(buf[:n])
}
