package service

import (
	"context"
	"fmt"
	"jtyl_bitable/global"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larktask "github.com/larksuite/oapi-sdk-go/v3/service/task/v2"
)

func CreateTask(summary, description string, startTimestamp, dueTimestamp string, members []*larktask.Member, tasklistGuid, sectionGuid string) (*larktask.CreateTaskRespData, error) {
	req := larktask.NewCreateTaskReqBuilder().
		UserIdType(`union_id`).
		InputTask(larktask.NewInputTaskBuilder().
			Summary(summary).
			Description(description).
			Start(larktask.NewStartBuilder().
				Timestamp(startTimestamp).
				IsAllDay(true).
				Build()).
			Due(larktask.NewDueBuilder().
				Timestamp(dueTimestamp).
				IsAllDay(true).
				Build()).
			Members(members).
			Tasklists([]*larktask.TaskInTasklistInfo{
				larktask.NewTaskInTasklistInfoBuilder().
					TasklistGuid(tasklistGuid).
					SectionGuid(sectionGuid).
					Build(),
			}).
			Build()).
		Build()

	resp, err := global.FEISHU.Task.V2.Task.Create(context.Background(), req)
	if err != nil {
		global.LOGGER.Error("创建任务失败", err)
		return nil, err
	}

	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return nil, fmt.Errorf("创建任务失败, logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}
	return resp.Data, nil
}

func GetTask(taskId string) (*larktask.GetTaskRespData, error) {
	req := larktask.NewGetTaskReqBuilder().
		TaskGuid(taskId).
		Build()

	resp, err := global.FEISHU.Task.V2.Task.Get(context.Background(), req)
	if err != nil {
		global.LOGGER.Error("获取任务失败", err)
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("获取任务失败, logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}
	return resp.Data, nil
}

func DeleteTask(taskId string) error {
	req := larktask.NewDeleteTaskReqBuilder().
		TaskGuid(taskId).
		Build()

	resp, err := global.FEISHU.Task.V2.Task.Delete(context.Background(), req)

	if err != nil {
		global.LOGGER.Error("删除任务失败", err)

		return err
	}
	if !resp.Success() {
		global.LOGGER.Error("删除任务失败", larkcore.Prettify(resp.CodeError))
		return fmt.Errorf("删除任务失败, logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}
	return nil
}

func GetDownloadUrl(fileToken []string) ([]*larkdrive.TmpDownloadUrl, error) {
	req := larkdrive.NewBatchGetTmpDownloadUrlMediaReqBuilder().
		FileTokens(fileToken).
		Build()
	resp, err := global.FEISHU.Drive.V1.Media.BatchGetTmpDownloadUrl(context.Background(), req)
	if err != nil {
		global.LOGGER.Error("获取下载链接失败", err)
		return nil, err
	}

	if !resp.Success() {
		global.LOGGER.Error("获取下载链接失败", larkcore.Prettify(resp.CodeError))
		return nil, fmt.Errorf("获取下载链接失败, logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	if len(resp.Data.TmpDownloadUrls) == 0 || resp.Data.TmpDownloadUrls[0].TmpDownloadUrl == nil {
		return nil, fmt.Errorf("未获取到下载链接")
	}
	return resp.Data.TmpDownloadUrls, nil
}
