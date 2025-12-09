package api

import (
	"jtyl_bitable/model"
	"jtyl_bitable/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type createApprovalRequest struct {
	RecordId          string  `json:"recordId"`
	ApprovalDate      string  `json:"绩效日期"`
	Process           string  `json:"工序"`
	Specification     string  `json:"型号"`
	PerformanceLevel  string  `json:"绩效等级"`
	PerformanceAmount float64 `json:"绩效金额"`
	Assignees         string  `json:"责任人"`
}

func CreateApproval(c *gin.Context) {
	var req createApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp{
			Msg:   "bad request",
			Error: err.Error(),
		})
		return
	}
	ms, _ := strconv.ParseInt(req.ApprovalDate, 10, 64)
	t := time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond))
	approvalDate := t.Format(time.RFC3339)

	userIds, err := service.GetUserIdsFromUnionIds(strings.Split(req.Assignees, ","))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp{
			Msg:   "user id not found",
			Error: err.Error(),
		})
		return
	}

	err = service.CreateApproval(req.RecordId, approvalDate, req.Process, req.Specification, req.PerformanceLevel, req.PerformanceAmount, userIds, req.ApprovalDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp{
			Msg:   "internal server error",
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, model.SuccessResp{
		Msg: "success",
	})
}
