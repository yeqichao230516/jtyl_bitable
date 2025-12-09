package service

import (
	"context"
	"encoding/json"
	"fmt"
	"jtyl_bitable/global"

	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
)

func CreateApproval(recordId, approvalDate, process, specification, performanceLevel string, performanceAmount float64, userIds []string, approvalCode string) error {
	form := []map[string]any{
		{ // 记录ID
			"id":    "widget17652472161290001",
			"type":  "input",
			"value": recordId,
		},
		{ // 绩效日期
			"id":    "widget17645649856020001",
			"type":  "date",
			"value": approvalDate,
		},

		{ // 工序
			"id":    "widget17651588323810001",
			"type":  "input",
			"value": process,
		},
		{ // 型号
			"id":    "widget17645649600730001",
			"type":  "input",
			"value": specification,
		},
		{ // 绩效等级
			"id":    "widget17645649988530001",
			"type":  "input",
			"value": performanceLevel,
		},
		{ // 绩效金额（元）
			"id":    "widget17651588509370001",
			"type":  "number",
			"value": performanceAmount,
		},
		{ // 责任人
			"id":    "widget17645650425740001",
			"type":  "contact",
			"value": userIds,
		},
	}
	formByte, _ := json.Marshal(form)
	req := larkapproval.NewCreateInstanceReqBuilder().
		InstanceCreate(larkapproval.NewInstanceCreateBuilder().
			ApprovalCode(`20EBFA89-CD89-47A9-AA64-FBFB9EC92DF9`).
			UserId(`6424c1g6`).
			Form(string(formByte)).
			Build()).
		Build()
	resp, err := global.FEISHU.Approval.V4.Instance.Create(context.Background(), req)
	if err != nil {
		global.LOGGER.Error("创建审批实例失败", err)
		return err
	}
	if !resp.Success() {
		global.LOGGER.Error("创建审批实例失败", resp.CodeError)
		return fmt.Errorf("创建审批实例失败, %s", resp.CodeError)
	}
	return nil
}

func GetInstanceForm(instanceCode string) (string, error) {
	req := larkapproval.NewGetInstanceReqBuilder().
		InstanceId(instanceCode).
		Locale(`zh-CN`).
		Build()
	resp, err := global.FEISHU.Approval.V4.Instance.Get(context.Background(), req)
	if err != nil {
		global.LOGGER.Error("获取审批实例失败", err)
		return "", err
	}
	if !resp.Success() {
		return "", fmt.Errorf("获取审批实例失败, %s", resp.CodeError)
	}
	return *resp.Data.Form, nil
}
