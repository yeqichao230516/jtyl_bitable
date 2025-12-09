package service

import (
	"context"
	"fmt"
	"jtyl_bitable/global"

	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

func GetOpenIdsFromUnionIds(unionIds []string) ([]string, error) {
	openIds := make([]string, 0, len(unionIds))
	for _, uid := range unionIds {
		req := larkcontact.NewGetUserReqBuilder().
			UserId(uid).
			UserIdType(`union_id`).
			Build()

		resp, err := global.FEISHU.Contact.V3.User.Get(context.Background(), req)
		if err != nil {
			global.LOGGER.Error("获取用户ID失败", err)
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("获取用户ID失败, %s", resp.CodeError)
		}
		if resp.Data.User.UserId == nil {
			return nil, fmt.Errorf("union_id %s 未查询到对应 user_id", uid)
		}
		openIds = append(openIds, *resp.Data.User.OpenId)
	}
	return openIds, nil
}
func GetUserIdsFromUnionIds(unionIds []string) ([]string, error) {
	userIds := make([]string, 0, len(unionIds))
	for _, uid := range unionIds {
		req := larkcontact.NewGetUserReqBuilder().
			UserId(uid).
			UserIdType(`union_id`).
			Build()

		resp, err := global.FEISHU.Contact.V3.User.Get(context.Background(), req)
		if err != nil {
			global.LOGGER.Error("获取用户ID失败", err)
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("获取用户ID失败, %s", resp.CodeError)
		}
		if resp.Data.User.UserId == nil {
			return nil, fmt.Errorf("union_id %s 未查询到对应 user_id", uid)
		}
		userIds = append(userIds, *resp.Data.User.UserId)
	}
	return userIds, nil
}
