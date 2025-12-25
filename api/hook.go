package api

import (
	"context"
	"encoding/json"
	"jtyl_bitable/global"
	"jtyl_bitable/service"
	"jtyl_bitable/utils"
	"strings"
	"time"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

func ApprovalHandler(ctx context.Context, event *larkevent.EventReq) error {
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

func BitableRecordChangeHandler(ctx context.Context, event *larkdrive.P2FileBitableRecordChangedV1) error {
	body := make([]byte, len(event.Body))
	copy(body, event.Body)

	go func() {
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			global.LOGGER.Errorf("approvalEvent转译失败: %v", err)
			return
		}

		encrypted, ok := req["encrypt"].(string)
		if !ok {
			global.LOGGER.Error("encrypt字段不存在或类型错误")
			return
		}

		decrypted := utils.Decrypt(encrypted)
		var data map[string]any
		if err := json.Unmarshal([]byte(decrypted), &data); err != nil {
			global.LOGGER.Errorf("解密 payload 解析失败: %v", err)
			return
		}

		if !isRecordAddedEvent(data) {
			return
		}

		newTime := extractTimeFromEvent(data)

		if newTime == "" {
			return
		}

		global.LOGGER.Infof("开始处理，新时间值: %s", newTime)

		// 持续删除直到所有记录的时间条件都与新时间相同
		totalDeleted := 0
		maxIterations := 100 // 防止无限循环
		batchSize := 50

		for i := 0; i < maxIterations; i++ {
			// 搜索不等于新时间值的记录
			ids := service.SearchNotEqual("B3nEbeQTRa2Dg3sNxR1cPiDXnGC", "tbl32AgBs7GuN60Q", newTime)

			// 如果没有找到不匹配的记录，说明所有记录的时间条件都已更新
			if len(ids) == 0 {
				global.LOGGER.Infof("所有记录的时间条件都已更新为: %s，停止删除", newTime)
				break
			}
			// 批量删除
			for j := 0; j < len(ids); j += batchSize {
				end := j + batchSize
				if end > len(ids) {
					end = len(ids)
				}
				batch := ids[j:end]
				if err := service.BatchDelete("B3nEbeQTRa2Dg3sNxR1cPiDXnGC", "tbl32AgBs7GuN60Q", batch); err != nil {
					global.LOGGER.Errorf("批量删除失败: %v", err)
					continue
				}
			}

			totalDeleted += len(ids)
			// 添加短暂延迟，避免API限流
			time.Sleep(500 * time.Millisecond)
		}

		if totalDeleted > 0 {
			global.LOGGER.Infof("删除完成！总共删除 %d 条记录，表格中所有记录的时间条件均为: %s", totalDeleted, newTime)
		} else {
			global.LOGGER.Infof("无需删除，表格中所有记录的时间条件已为: %s", newTime)
		}
	}()
	return nil
}

// 检查是否为 record_added 事件
func isRecordAddedEvent(data map[string]any) bool {
	event, ok := data["event"].(map[string]any)
	if !ok {
		return false
	}

	actionList, ok := event["action_list"].([]any)
	if !ok {
		return false
	}

	for _, action := range actionList {
		actionMap, ok := action.(map[string]any)
		if !ok {
			continue
		}

		actionType, _ := actionMap["action"].(string)
		if actionType == "record_added" {
			return true
		}
	}
	return false
}

func extractTimeFromEvent(data map[string]any) string {
	event, ok := data["event"].(map[string]any)
	if !ok {
		global.LOGGER.Warn("event字段格式错误")
		return ""
	}

	actionList, ok := event["action_list"].([]any)
	if !ok {
		global.LOGGER.Warn("action_list字段格式错误")
		return ""
	}

	for _, action := range actionList {
		actionMap, ok := action.(map[string]any)
		if !ok {
			continue
		}

		actionType, _ := actionMap["action"].(string)
		if actionType != "record_added" {
			continue
		}

		afterValue, ok := actionMap["after_value"].([]any)
		if !ok {
			global.LOGGER.Warn("after_value字段格式错误")
			continue
		}

		for _, field := range afterValue {
			fieldMap, ok := field.(map[string]any)
			if !ok {
				continue
			}

			fieldID, _ := fieldMap["field_id"].(string)
			if fieldID != "fldDjFPbeC" {
				continue
			}

			// 处理 field_value 字段，它可能是一个 JSON 字符串
			fieldValueStr, ok := fieldMap["field_value"].(string)
			if !ok {
				global.LOGGER.Warn("field_value字段不是字符串类型")
				continue
			}

			if fieldValueStr == "" {
				global.LOGGER.Warn("field_value为空字符串")
				continue
			}

			// 解析 JSON 字符串
			var fieldValue []map[string]any
			if err := json.Unmarshal([]byte(fieldValueStr), &fieldValue); err != nil {
				global.LOGGER.Warnf("解析field_value JSON失败: %v, 原始值: %s", err, fieldValueStr)
				continue
			}

			if len(fieldValue) == 0 {
				global.LOGGER.Warn("解析后的field_value为空数组")
				continue
			}

			firstItem := fieldValue[0]
			time, ok := firstItem["text"].(string)
			if !ok || time == "" {
				global.LOGGER.Warn("text字段为空或格式错误")
				continue
			}

			global.LOGGER.Infof("成功提取到时间值: %s", time)
			return time
		}
	}

	global.LOGGER.Info("未找到时间值")
	return ""
}
