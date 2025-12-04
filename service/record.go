package service

import (
	"context"
	"fmt"
	"jtyl_bitable/global"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func SearchRecords(appToken, tableId string, pageSize int, body *larkbitable.SearchAppTableRecordReqBody) []map[string]any {
	if pageSize <= 0 || pageSize > 500 {
		pageSize = 500
	}
	records := []map[string]any{}
	pageToken := ""
	for {
		req := larkbitable.NewSearchAppTableRecordReqBuilder().
			AppToken(appToken).
			TableId(tableId).
			PageToken(pageToken).
			PageSize(pageSize).
			Body(body).
			Build()
		resp, err := global.FEISHU.Bitable.V1.AppTableRecord.Search(context.Background(), req)
		if err != nil {
			global.LOGGER.Errorf("SearchAppTableRecord error: %v", err)
			return nil
		}
		if resp == nil || resp.Data == nil {
			global.LOGGER.Error("SearchAppTableRecord resp or resp.Data is nil")
			return nil
		}

		for _, item := range resp.Data.Items {
			if item == nil || item.Fields == nil {
				continue
			}
			records = append(records, item.Fields)
		}

		hasMore := false
		if resp.Data.HasMore != nil {
			hasMore = *resp.Data.HasMore
		}
		nextPageToken := ""
		if resp.Data.PageToken != nil {
			nextPageToken = *resp.Data.PageToken
		}

		if !hasMore || nextPageToken == "" {
			break
		}
		pageToken = nextPageToken
	}
	return records
}

func UpdateRecord(appToken, tableId, recordId string, fields map[string]any) {
	req := larkbitable.NewUpdateAppTableRecordReqBuilder().
		AppToken(appToken).
		TableId(tableId).
		RecordId(recordId).
		AppTableRecord(larkbitable.NewAppTableRecordBuilder().
			Fields(fields).
			Build()).
		Build()

	_, err := global.FEISHU.Bitable.V1.AppTableRecord.Update(context.Background(), req)

	if err != nil {
		global.LOGGER.Errorf("UpdateAppTableRecord error: %v", err)
		return
	}
}

func BatchCreateRecords(appToken, tableId string, records []*larkbitable.AppTableRecord) []string {
	// 每1000条分批处理
	batchSize := 1000
	recordIds := []string{}
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		req := larkbitable.NewBatchCreateAppTableRecordReqBuilder().
			AppToken(appToken).
			TableId(tableId).
			Body(larkbitable.NewBatchCreateAppTableRecordReqBodyBuilder().
				Records(records[i:end]).
				Build()).
			Build()

		resp, err := global.FEISHU.Bitable.V1.AppTableRecord.BatchCreate(context.Background(), req)
		if err != nil {
			global.LOGGER.Errorf("BatchCreateAppTableRecord error: %v", err)
			return []string{}
		}
		if !resp.Success() {
			fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
			return []string{}
		}

		if resp == nil || resp.Data == nil || len(resp.Data.Records) == 0 {
			continue
		}

		for _, record := range resp.Data.Records {
			if record != nil && record.RecordId != nil && *record.RecordId != "" {
				recordIds = append(recordIds, *record.RecordId)
			}
		}
	}
	return recordIds
}

func BatchDeleteRecords(appToken, tableId string, recordIds []string) error {
	req := larkbitable.NewBatchDeleteAppTableRecordReqBuilder().
		AppToken(appToken).
		TableId(tableId).
		Body(larkbitable.NewBatchDeleteAppTableRecordReqBodyBuilder().
			Records(recordIds).
			Build()).
		Build()

	_, err := global.FEISHU.Bitable.V1.AppTableRecord.BatchDelete(context.Background(), req)
	if err != nil {
		global.LOGGER.Errorf("BatchDeleteAppTableRecord error: %v", err)
		return err
	}
	return nil
}
