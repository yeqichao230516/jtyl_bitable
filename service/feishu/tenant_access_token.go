package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func PostTenantAccessToken(app_id, app_secret string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"app_id":     app_id,
		"app_secret": app_secret,
	})
	req, _ := http.NewRequest("POST", "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	respBody := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}
	if respBody["code"].(float64) != 0 {
		return "", fmt.Errorf("请求失败，状态码: %d", int(respBody["code"].(float64)))
	}
	return respBody["tenant_access_token"].(string), nil
}
