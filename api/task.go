package api

import (
	"jtyl_bitable/global"
	"jtyl_bitable/model"
	"jtyl_bitable/service"
	"jtyl_bitable/service/task"
	"jtyl_bitable/service/token"
	"jtyl_bitable/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	larktask "github.com/larksuite/oapi-sdk-go/v3/service/task/v2"
)

type createTaskRequest struct {
	Summary          string `json:"任务标题"`
	Description      string `json:"任务摘要"`
	AttachmentTokens string `json:"任务附件"`
	StartTimestamp   string `json:"开始日期"`
	DueTimestamp     string `json:"截止日期"`
	Assignees        string `json:"负责人"`
	Fllowers         string `json:"关注人"`
	TasklistGuid     string `json:"清单ID"`
	SectionGuid      string `json:"分组ID"`
}

type createTaskResponse struct {
	Guid   string `json:"guid"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

func CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp{
			Msg:   "NG",
			Error: err.Error(),
		})
		return
	}

	assignees := strings.Split(req.Assignees, ",")
	followers := strings.Split(req.Fllowers, ",")
	members := []*larktask.Member{}
	for _, assignee := range assignees {
		user, err := service.GetUserMsgFromUnionId(assignee)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error(),
			})
			return
		}

		members = append(members, larktask.NewMemberBuilder().
			Id(assignee).
			Type(`user`).
			Role(`assignee`).
			Name(*user.Name).
			Build())
	}
	for _, follower := range followers {
		user, err := service.GetUserMsgFromUnionId(follower)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error(),
			})
			return
		}
		members = append(members, larktask.NewMemberBuilder().
			Id(follower).
			Type(`user`).
			Role(`follower`).
			Name(*user.Name).
			Build())
	}

	t, err := service.CreateTask(req.Summary, req.Description, req.StartTimestamp, req.DueTimestamp, members, req.TasklistGuid, req.SectionGuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp{
			Msg:   "NG",
			Error: err.Error(),
		})
		return
	}
	attachmentTokens := strings.Split(req.AttachmentTokens, ",")
	for _, attachmentToken := range attachmentTokens {
		data, err := service.GetDownloadUrl([]string{attachmentToken})
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error()})
			return
		}

		if data[0].TmpDownloadUrl == nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: "未获取到下载链接"})
			return
		}
		file, filePath, err := utils.DownloadFileFromURL(*data[0].TmpDownloadUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error()})
			return
		}
		tenantAccessToken, err := token.GetTenantAccessToken(global.CONFIG.App.Id, global.CONFIG.App.Secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error()})
			return
		}
		if err = task.UploadAttachment(tenantAccessToken, *t.Task.Guid, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error()})
			return
		}
		utils.CleanupTmpFile(file)
	}
	c.JSON(200, model.SuccessResp{
		Msg: "success",
		Data: createTaskResponse{
			Guid:   *t.Task.Guid,
			Status: *t.Task.Status,
			Url:    *t.Task.Url,
		},
	})
}

type getTaskRequest struct {
	RecordIds string `json:"recordIds"`
	Guids     string `json:"guids"`
}

func GetTask(c *gin.Context) {
	var req getTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp{
			Msg:   "NG",
			Error: err.Error(),
		})
		return
	}
	recordId := strings.Split(req.RecordIds, ",")
	guid := strings.Split(req.Guids, ",")
	for i := range recordId {
		data, err := service.GetTask(guid[i])
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error(),
			})
			return
		}
		if err := service.UpdateRecord("V8AxbmAOXapXQesLIIFcJbkunae", "tbl17bp2dR5Renn5", recordId[i], map[string]any{
			"status": *data.Task.Status,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResp{
				Msg:   "NG",
				Error: err.Error(),
			})
			return
		}
	}
	c.JSON(200, model.SuccessResp{
		Msg: "OK",
	})
}

type deleteTaskRequest struct {
	RecordId string `json:"recordId"`
	Guid     string `json:"guid"`
}

// DeleteTask 删除任务(有bug，删除任务和删除记录无法保证原子性)
func DeleteTask(c *gin.Context) {
	var req deleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResp{
			Msg:   "NG",
			Error: err.Error(),
		})
		return
	}
	if err := service.DeleteTask(req.Guid); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp{
			Msg:   "NG",
			Error: err.Error(),
		})
		return
	}
	if err := service.BatchDeleteRecords("V8AxbmAOXapXQesLIIFcJbkunae", "tbl17bp2dR5Renn5", []string{req.RecordId}); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResp{
			Msg:   "NG",
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, model.SuccessResp{
		Msg: "OK",
	})
}
