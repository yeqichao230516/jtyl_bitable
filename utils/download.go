package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadURLToTempFile(url string) (*os.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "feishu.*.tmp")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, err
	}

	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, err
	}

	return tmp, nil
}

func CleanupTmpFile(f *os.File) {
	if f == nil {
		return
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
}
