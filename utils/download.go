package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DownloadFileFromURL(url string) (*os.File, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	fileName := getFileNameFromResponse(resp)

	if err := os.MkdirAll("temp", 0755); err != nil {
		return nil, "", fmt.Errorf("创建 temp 目录失败: %v", err)
	}

	filePath := filepath.Join("temp", fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("创建文件失败: %v", err)
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		file.Close()
		os.Remove(filePath)
		return nil, "", fmt.Errorf("写入文件失败: %v", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, "", fmt.Errorf("重置文件指针失败: %v", err)
	}

	return file, filePath, nil
}

func CleanupTmpFile(f *os.File) {
	if f == nil {
		return
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
}

func getFileNameFromResponse(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		return ""
	}

	// 尝试匹配 filename="..." 或 filename*=...
	patterns := []string{
		`filename\*?=(?:\"?UTF-8''(?P<filename>[^;\"]+)\"?|(?:\"?(?P<filename2>[^;\"]+)\"?))`,
		`filename=(?:\"([^\"]+)\"|([^;\"]+))`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(contentDisposition)

		if len(matches) > 0 {
			// 获取非空的捕获组
			for i, match := range matches {
				if i > 0 && match != "" && match != "filename" && match != "filename2" {
					// 如果是URL编码的文件名，需要解码
					if strings.HasPrefix(contentDisposition, "filename*=UTF-8''") {
						decoded, err := decodeURLEncodedFilename(match)
						if err == nil {
							return decoded
						}
					}
					return match
				}
			}
		}
	}

	return ""
}

func decodeURLEncodedFilename(filename string) (string, error) {
	decoded, err := url.QueryUnescape(filename)
	if err != nil {
		return filename, err
	}
	return decoded, nil
}
