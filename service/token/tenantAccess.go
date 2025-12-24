package token

import (
	"bytes"
	"encoding/json"
	"io"
	"jtyl_bitable/global"
	"net/http"
	"time"
)

const (
	url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
)

type reqBody struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type rspBody struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func GetTenantAccessToken(appID, appSecret string) (string, error) {
	reqData, err := json.Marshal(reqBody{
		AppID:     appID,
		AppSecret: appSecret,
	})
	if err != nil {
		global.LOGGER.Errorf("marshal req fail: %v", err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqData))
	if err != nil {
		global.LOGGER.Errorf("new request fail: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		global.LOGGER.Errorf("do request fail: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOGGER.Errorf("read body fail: %v", err)
		return "", err
	}

	var r rspBody
	if err := json.Unmarshal(body, &r); err != nil {
		global.LOGGER.Errorf("unmarshal resp fail: %v", err)
		return "", err
	}

	if r.Code != 0 {
		global.LOGGER.Errorf("feishu err: code=%d msg=%s", r.Code, r.Msg)
		return "", err
	}

	global.LOGGER.Infof("tenant_access_token = %s (expire %ds)", r.TenantAccessToken, r.Expire)
	return r.TenantAccessToken, nil
}
