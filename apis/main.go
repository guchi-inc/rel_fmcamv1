// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package apis

import (
	"fmcam/apis/alysis"
	"fmcam/apis/handlers"
	"fmcam/apis/resources"
	"fmcam/apis/streams"
	"fmcam/apis/users"
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/sysaly"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiClientsCtl struct {
	DevStreamApis  *streams.StreamGroup
	ClientAlyApis  *alysis.AlysisApi
	ClientFaceApis *handlers.FmClientApi
}

type ApiGroupCtl struct {
	UserApis *users.ApiGroup

	LimitsApis         *resources.LimitsApi
	GAddrsApis         *handlers.GovAddrApi
	CustomerTenantApis *resources.TenantsApi
}

var (
	Logger      = configs.NewLogger("apis")
	ApiGroupApp = NewApiGroupCtl()
	//客户端操作路由 2818
	ApiClientApp = NewApiClientsCtl()
)

func NewApiGroupCtl() *ApiGroupCtl {
	users := users.NewUserApis()
	return &ApiGroupCtl{
		UserApis: users,

		LimitsApis:         &resources.LimitsApi{},
		GAddrsApis:         &handlers.GovApis,
		CustomerTenantApis: &resources.TenantsApi{},
	}
}

func NewApiClientsCtl() *ApiClientsCtl {
	streams := streams.BuildFacer()
	alysis := &alysis.AlysisApi{}
	fmclients := &handlers.FmClientApi{}
	return &ApiClientsCtl{
		DevStreamApis:  streams,
		ClientAlyApis:  alysis,
		ClientFaceApis: fmclients,
	}
}

// 接口问题跟踪
type TraceProblem struct{}

// 新增验证记录
func (TG *TraceProblem) OperateLog(ctx *gin.Context) {

	var NodeIds = sysaly.TraceInfo{}
	err := ctx.ShouldBindJSON(&NodeIds)
	if err != nil {

		Logger.Infof(" OperateLog datas of param Err:%v", err)
		helpers.JSONs(ctx, http.StatusBadRequest, err)
		return
	}

	if NodeIds.Data == "" && NodeIds.Brower == "" {
		newErr := fmt.Errorf(" OperateLog datas of param Data or Brower Err:%v", NodeIds)
		Logger.Infof("%v", newErr)
		helpers.JSONs(ctx, http.StatusBadRequest, newErr)
		return
	}
	Logger.Infof("NodeIds:%#v \n", NodeIds)

	helpers.JSONs(ctx, http.StatusOK, gin.H{"message": "ok", "data": "", "code": code.Success})
}
