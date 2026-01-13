package services_feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TaskService interface {
	CreateTask(app_id, app_secret string, data map[string]any) (guid string, url string, err error)
}

func CreateTask(app_id, app_secret string, data map[string]any) (guid string, url string, err error) {
	accessToken := GetToken()
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://open.feishu.cn/open-apis/task/v2/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}
	var respBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", "", err
	}
	if respBody["code"].(float64) != 0 {
		return "", "", fmt.Errorf("请求失败，状态码: %d", int(respBody["code"].(float64)))
	}

	return respBody["data"].(map[string]any)["task"].(map[string]any)["guid"].(string), respBody["data"].(map[string]any)["task"].(map[string]any)["url"].(string), nil
}

// DeleteTask 删除任务
func DeleteTask(app_id, app_secret string, guid string) (err error) {
	accessToken := GetToken()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://open.feishu.cn/open-apis/task/v2/tasks/%s", guid), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}
	var respBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	if respBody["code"].(float64) != 0 {
		return fmt.Errorf("请求失败，状态码: %d", int(respBody["code"].(float64)))
	}
	return nil
}

// PostCreateSubtasks 创建子任务
func PostCreateSubtasks(app_id, app_secret string, guid string, data map[string]any) (err error) {
	accessToken := GetToken()
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://open.feishu.cn/open-apis/task/v2/tasks/%s/subtasks?user_id_type=open_id", guid), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}
	var respBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	if respBody["code"].(float64) != 0 {
		return fmt.Errorf("请求失败，状态码: %d", int(respBody["code"].(float64)))
	}
	return nil
}
