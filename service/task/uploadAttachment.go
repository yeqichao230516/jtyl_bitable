package task

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadAttachment(tenantAccessToken, resourceID, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("resource_type", "task")
	_ = w.WriteField("resource_id", resourceID)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, filename := filepath.Split(filePath)
	part, err := w.CreateFormFile("files", filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://open.feishu.cn/open-apis/task/v2/attachments/upload", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tenantAccessToken)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
