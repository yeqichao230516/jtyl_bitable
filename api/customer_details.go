package api

import (
	"fmt"
	"jtyl_bitable/global"
	"jtyl_bitable/model"
	"jtyl_bitable/service"
	"jtyl_bitable/utils"
	"sync"

	"github.com/gin-gonic/gin"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func PostCustomerDetails(c *gin.Context) {
	var req model.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.Resp{
			Code: 1,
			Msg:  "请求参数错误",
		})
		return
	}
	companyName := req.Data.(map[string]any)["商业公司"].(string)
	appToken := req.Data.(map[string]any)["appToken"].(string)
	salesTableId := req.Data.(map[string]any)["销售统计"].(string)
	shipmentDetailTableId := req.Data.(map[string]any)["发货明细"].(string)
	companyDetailTableId := req.Data.(map[string]any)["公司明细"].(string)

	if len(global.RECORDS_ID) > 0 {
		err := service.BatchDeleteRecords(appToken, companyDetailTableId, global.RECORDS_ID)
		if err != nil {
			c.JSON(400, model.Resp{
				Code: 400,
				Msg:  "删除失败",
			})
			return
		}
		global.RECORDS_ID = []string{}
	}

	var (
		records = []*larkbitable.AppTableRecord{}
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		filter_sales := larkbitable.NewFilterInfoBuilder().
			Conjunction(`and`).
			Conditions([]*larkbitable.Condition{
				larkbitable.NewConditionBuilder().
					FieldName(`商业公司`).
					Operator(`is`).
					Value([]string{companyName}).
					Build(),
			}).
			Build()
		body_sales := larkbitable.NewSearchAppTableRecordReqBodyBuilder().
			FieldNames([]string{`日期`, `序号`, `产品代码`, `产品名称`, `产品型号`, `产品规格`, `销售数量`, `商业单价`, `商业总额`, `诊所名称`, `金塔单价`, `金塔总额`, `开票备注`}).
			Filter(filter_sales).
			Build()
		salesRecords := service.SearchRecords(appToken, salesTableId, 500, body_sales)

		var salesToAdd []*larkbitable.AppTableRecord
		for _, r := range salesRecords {
			if r == nil {
				continue
			}
			record := larkbitable.NewAppTableRecordBuilder().
				Fields(map[string]any{
					`商业公司`: companyName,
					`序号`:   utils.GetNestedString(r, `序号`, `text`),
					`日期`:   utils.GetNestedFloat64(r, `日期`),
					`产品代码`: utils.GetNestedString(r, `产品代码`, `text`),
					`产品名称`: utils.GetNestedString(r, `产品名称`, `text`),
					`产品型号`: utils.GetNestedString(r, `产品型号`, `text`),
					`产品规格`: utils.GetNestedString(r, `产品规格`, `text`),
					`产品数量`: -utils.GetNestedFloat64(r, `销售数量`),
					`商业单价`: utils.GetNestedFloat64(r, `商业单价`),
					`商业总额`: utils.GetNestedFloat64(r, `商业总额`),
					`诊所名称`: utils.GetNestedString(r, `诊所名称`, `text`),
					`金塔单价`: utils.GetNestedFloat64(r, `金塔单价`),
					`金塔总额`: utils.GetNestedFloat64(r, `金塔总额`),
					`开票备注`: utils.GetNestedString(r, `开票备注`, `text`),
				}).
				Build()
			if record != nil {
				salesToAdd = append(salesToAdd, record)
			}
		}
		fmt.Println("salesToAdd:", salesToAdd)
		mu.Lock()
		records = append(records, salesToAdd...)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		filter_shipment := larkbitable.NewFilterInfoBuilder().
			Conjunction(`and`).
			Conditions([]*larkbitable.Condition{
				larkbitable.NewConditionBuilder().
					FieldName(`商业公司`).
					Operator(`is`).
					Value([]string{companyName}).
					Build(),
			}).
			Build()
		body_shipment := larkbitable.NewSearchAppTableRecordReqBodyBuilder().
			FieldNames([]string{`日期`, `序号`, `出库单号`, `产品代码`, `产品名称`, `产品型号`, `产品规格`, `产品数量`, `金塔单价`, `金塔总额`, `收款金额`, `应收余额`, `应收编码`}).
			Filter(filter_shipment).
			Build()
		shipmentRecords := service.SearchRecords(appToken, shipmentDetailTableId, 500, body_shipment)

		var shipmentToAdd []*larkbitable.AppTableRecord
		for _, r := range shipmentRecords {
			if r == nil {
				continue
			}
			record := larkbitable.NewAppTableRecordBuilder().
				Fields(map[string]any{
					`商业公司`: companyName,
					`序号`:   utils.GetNestedString(r, `序号`, `text`),
					`日期`:   utils.GetNestedFloat64(r, `日期`),
					`出库单号`: utils.GetNestedString(r, `出库单号`, `text`),
					`产品代码`: utils.GetNestedString(r, `产品代码`, `text`),
					`产品名称`: utils.GetNestedString(r, `产品名称`, `text`),
					`产品型号`: utils.GetNestedString(r, `产品型号`, `text`),
					`产品规格`: utils.GetNestedString(r, `产品规格`, `text`),
					`产品数量`: utils.GetNestedFloat64(r, `产品数量`),
					`金塔单价`: utils.GetNestedFloat64(r, `金塔单价`),
					`金塔总额`: utils.GetNestedFloat64(r, `金塔总额`),
					`收款金额`: utils.GetNestedFloat64(r, `收款金额`),
					`应收余额`: utils.GetNestedFloat64(r, `应收余额`),
					`应收编码`: utils.GetNestedString(r, `应收编码`, `text`),
				}).
				Build()
			if record != nil {
				shipmentToAdd = append(shipmentToAdd, record)
			}
		}

		fmt.Println("shipmentToAdd:", shipmentToAdd)
		mu.Lock()
		records = append(records, shipmentToAdd...)
		mu.Unlock()
	}()
	wg.Wait()

	createdIds := service.BatchCreateRecords(appToken, companyDetailTableId, records)

	global.RECORDS_ID = append(global.RECORDS_ID, createdIds...)

	c.JSON(200, model.Resp{
		Code: 0,
		Msg:  "success",
	})
}
