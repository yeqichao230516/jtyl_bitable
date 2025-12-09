package api

import (
	"context"
	"encoding/json"
	"jtyl_bitable/global"
	"jtyl_bitable/service"
	"jtyl_bitable/utils"
	"strings"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
)

func Handler(ctx context.Context, event *larkevent.EventReq) error {
	body := make([]byte, len(event.Body))
	copy(body, event.Body)

	go func() {
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			global.LOGGER.Errorf("approvalEvent转译失败: %v", err)
		}
		decrypted := utils.Decrypt(req["encrypt"].(string))

		var data map[string]any
		if err := json.Unmarshal([]byte(decrypted), &data); err != nil {
			global.LOGGER.Errorf("解密 payload 解析失败: %v", err)
		}

		eventObj := data["event"].(map[string]any)
		status := eventObj["status"].(string)
		approvalCode := eventObj["approval_code"].(string)
		instanceCode := eventObj["instance_code"].(string)

		switch status {
		case "PENDING":
			global.LOGGER.Infof("正再处理审批%v", instanceCode)
		case "APPROVED":
			switch approvalCode {
			case global.CONFIG.Approval.BlgsCode:
				global.LOGGER.Infof("审批:%v同意", instanceCode)
				go handleAgreeApproved(instanceCode)
			}
		case "REJECTED":
			switch approvalCode {
			case global.CONFIG.Approval.BlgsCode:
				global.LOGGER.Infof("审批:%v拒绝", instanceCode)
				go handleRefuseApproved(instanceCode)
			}
		default:
			global.LOGGER.Warnf("Unhandled approval status: %s for instance: %s", status, instanceCode)
		}
	}()
	return nil
}

func handleRefuseApproved(instanceCode string) error {
	instance, err := service.GetInstanceForm(instanceCode)
	if err != nil {
		global.LOGGER.Error("获取审批实例失败", err)
		return err
	}
	var form []map[string]any
	err = json.Unmarshal([]byte(instance), &form)
	if err != nil {
		global.LOGGER.Error("解析审批实例失败", err)
		return err
	}
	var recordId string
	var tableId string

	for _, item := range form {
		if item["name"] == "记录ID" {
			tableId = strings.Split(item["value"].(string), "-")[0]
			recordId = strings.Split(item["value"].(string), "-")[1]
			break
		}
	}
	if err = service.UpdateRecord("KUD4bR614aak70s82iwc31LpnFf", tableId, recordId, map[string]any{
		"审批结果": "绩效处理",
	}); err != nil {
		global.LOGGER.Error("更新记录失败", err)
		return err
	}
	return nil
}

func handleAgreeApproved(instanceCode string) error {
	instance, err := service.GetInstanceForm(instanceCode)
	if err != nil {
		global.LOGGER.Error("获取审批实例失败", err)
		return err
	}
	var form []map[string]any
	err = json.Unmarshal([]byte(instance), &form)
	if err != nil {
		global.LOGGER.Error("解析审批实例失败", err)
		return err
	}
	var recordId string
	var tableId string

	for _, item := range form {
		if item["name"] == "记录ID" {
			tableId = strings.Split(item["value"].(string), "-")[0]
			recordId = strings.Split(item["value"].(string), "-")[1]
			break
		}
	}
	if err = service.UpdateRecord("KUD4bR614aak70s82iwc31LpnFf", tableId, recordId, map[string]any{
		"审批结果": "改善通过",
	}); err != nil {
		global.LOGGER.Error("更新记录失败", err)
		return err
	}
	return nil
}
