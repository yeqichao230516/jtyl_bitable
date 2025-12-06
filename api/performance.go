package api

import (
	"jtyl_bitable/model"
	"jtyl_bitable/service"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

type PerformanceRequest struct {
	AppToken        string `json:"appToken"`
	ParameterTable  string `json:"参数表"`
	StatisticsTable string `json:"统计表"`
	RecordID        string `json:"记录ID"`
	Specification   string `json:"型号"`
	Process         string `json:"工序"`
	Alarm           string `json:"告警"`
}
type PerformanceResponse struct {
	PerformanceLevel  string  `json:"绩效等级"`
	PerformanceAmount float64 `json:"绩效金额"`
}

func Performance(c *gin.Context) {
	var req PerformanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp{
			Code: http.StatusBadRequest,
			Msg:  "bad request",
		})
		return
	}

	sort := []*larkbitable.Sort{
		larkbitable.NewSortBuilder().
			FieldName(`日期`).
			Desc(true).
			Build(),
	}
	filter := larkbitable.NewFilterInfoBuilder().
		Conjunction(`and`).
		Conditions([]*larkbitable.Condition{
			larkbitable.NewConditionBuilder().
				FieldName(`工序`).
				Operator(`is`).
				Value([]string{req.Process}).
				Build(),
			larkbitable.NewConditionBuilder().
				FieldName(`型号`).
				Operator(`is`).
				Value([]string{req.Specification}).
				Build(),
		}).
		Build()

	body_statistics := larkbitable.NewSearchAppTableRecordReqBodyBuilder().
		FieldNames([]string{`记录ID`, `告警`, `绩效等级`}).
		Sort(sort).
		Filter(filter).
		Build()
	body_parameter := larkbitable.NewSearchAppTableRecordReqBodyBuilder().
		FieldNames([]string{`1级绩效`, `2级绩效`, `3级绩效`, `4级绩效`, `5级绩效`}).
		Filter(filter).
		Build()

	records_statistics := service.SearchRecords(req.AppToken, req.StatisticsTable, 8, body_statistics)

	records_parameter := service.SearchRecords(req.AppToken, req.ParameterTable, 1, body_parameter)

	if len(records_statistics) < 3 {
		service.UpdateRecord(req.AppToken, req.StatisticsTable, req.RecordID, map[string]any{
			"绩效等级": "无",
		})
		c.JSON(http.StatusOK, model.SuccessResp{
			Code: http.StatusOK,
			Msg:  "success",
			Data: PerformanceResponse{
				PerformanceLevel: "无",
			},
		})
		return
	}
	if req.Alarm == "正常" {
		service.UpdateRecord(req.AppToken, req.StatisticsTable, req.RecordID, map[string]any{
			"绩效等级": "无",
		})
		c.JSON(http.StatusOK, model.SuccessResp{
			Code: http.StatusOK,
			Msg:  "success",
			Data: PerformanceResponse{
				PerformanceLevel: "无",
			},
		})
		return
	}
	hasPerformance := false
	for i := 1; i < len(records_statistics); i++ {
		if records_statistics[i]["绩效等级"].(string) != "无" {
			hasPerformance = true
			break
		}
	}
	if hasPerformance {
		service.UpdateRecord(req.AppToken, req.StatisticsTable, req.RecordID, map[string]any{
			"绩效等级": "无",
		})
		c.JSON(http.StatusOK, model.SuccessResp{
			Code: http.StatusOK,
			Msg:  "success",
			Data: PerformanceResponse{
				PerformanceLevel: "无",
			},
		})
		return
	}
	record1_alarm := records_statistics[0]["告警"].(map[string]any)["value"].([]any)[0].(string)
	record2_alarm := records_statistics[1]["告警"].(map[string]any)["value"].([]any)[0].(string)
	record3_alarm := records_statistics[2]["告警"].(map[string]any)["value"].([]any)[0].(string)

	alertLevels := map[string]int{
		"无":    0,
		"1级警报": 1,
		"2级警报": 2,
		"3级警报": 3,
		"4级警报": 4,
		"5级警报": 5,
	}
	record1_level := alertLevels[record1_alarm]
	record2_level := alertLevels[record2_alarm]
	record3_level := alertLevels[record3_alarm]

	max_level := math.Max(math.Max(float64(record1_level), float64(record2_level)), float64(record3_level))
	performance := ""
	switch int(max_level) {
	case 1:
		performance = "1级绩效"
	case 2:
		performance = "2级绩效"
	case 3:
		performance = "3级绩效"
	case 4:
		performance = "4级绩效"
	case 5:
		performance = "5级绩效"
	}
	service.UpdateRecord(req.AppToken, req.StatisticsTable, req.RecordID, map[string]any{
		"绩效等级": performance,
		"绩效金额": records_parameter[0][performance].(float64),
	})

	c.JSON(http.StatusOK, PerformanceResponse{
		PerformanceLevel:  performance,
		PerformanceAmount: records_parameter[0][performance].(float64),
	})
}
