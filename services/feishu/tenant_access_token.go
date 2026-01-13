package services_feishu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yeqichao230516/jtyl-api/pkg/redis"
)

func RefreshToken() error {
	body, _ := json.Marshal(map[string]any{
		"app_id":     "cli_a81807b812b7901c",
		"app_secret": "wGTNLAxJiZBCoBvht4b7UbeBmSkWprYw",
	})
	req, _ := http.NewRequest("POST", "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	respBody := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	if respBody["code"].(float64) != 0 {
		return fmt.Errorf("请求失败，状态码: %d", int(respBody["code"].(float64)))
	}
	redis.Client().Set(context.Background(), "feishu_tenant_access_token", respBody["tenant_access_token"].(string), 7200*time.Second)
	return nil
}

func GetToken() string {
	return redis.Client().Get(context.Background(), "feishu_tenant_access_token").String()
}
