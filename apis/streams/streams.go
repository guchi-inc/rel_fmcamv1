package streams

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/iotmodel"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceTask struct{}

// 流组设备查询
func (st *DeviceTask) StreamQuery(ctx *gin.Context) {

	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	PageId := ctx.DefaultQuery("page_id", "")
	RtspUrl := ctx.DefaultQuery("url", "")

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 10 {
		PageSize = 10
	}

	user, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, gin.H{"error": err, "message": code.ZhCNText[code.AuthorizationError]})
		return
	}

	datas, err := SysStreamGroup.StreamSelect(Page, PageSize, user.TenantId, PageId, RtspUrl)
	apilog.Printf("param id:%v url:%v  datas:%v, err:%v \n", PageId, RtspUrl, datas, err)

	if err != nil {
		helpers.JSONs(ctx, code.NullData, gin.H{"error": err, "message": code.ZhCNText[code.NullData]})
		return
	}

	helpers.JSONs(ctx, code.Success, datas)
}

// 流组设备新增
func (st *DeviceTask) StreamNew(ctx *gin.Context) {

	paramsDevice := iotmodel.FmStreamTasks{}

	err := ctx.ShouldBindJSON(&paramsDevice)
	if err != nil {
		helpers.JSONs(ctx, code.ParamBindError, gin.H{"message": code.ZhCNText[code.ParamBindError], "error": err})
		return
	}

	data, err := SysStreamGroup.StreamDeviceNew(&paramsDevice)
	if err != nil {
		helpers.JSONs(ctx, code.AdminModifyDataError, gin.H{"message": code.ZhCNText[code.AdminModifyDataError], "error": err})
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"message": "ok", "data": data})
}

// 流组设备修改
func (st *DeviceTask) StreamUpdate(ctx *gin.Context) {

	paramsDevice := iotmodel.FmStreamTasks{}

	err := ctx.ShouldBindJSON(&paramsDevice)
	if err != nil {
		helpers.JSONs(ctx, code.ParamBindError, gin.H{"message": code.ZhCNText[code.ParamBindError], "error": err})
		return
	}

	if paramsDevice.ID <= 0 {
		helpers.JSONs(ctx, code.ParamBindError, gin.H{"message": code.ZhCNText[code.ParamBindError], "error": "数据唯一id 错误"})
		return
	}
	data, err := SysStreamGroup.StreamDeviceUpdate(&paramsDevice)
	if err != nil {
		helpers.JSONs(ctx, code.AdminModifyDataError, gin.H{"message": code.ZhCNText[code.AdminModifyDataError], "error": err})
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"message": "ok", "data": data})
}

// 校验查询
func (st *DeviceTask) CheckLogSelect(ctx *gin.Context) {
	helpers.JSONs(ctx, code.Success, gin.H{"message": "ok"})

}

// 校验新增
func (st *DeviceTask) CheckLogNew(ctx *gin.Context) {
	helpers.JSONs(ctx, code.Success, gin.H{"message": "ok"})
}

// 校验修改
func (st *DeviceTask) CheckLogUpdate(ctx *gin.Context) {
	helpers.JSONs(ctx, code.Success, gin.H{"message": "ok"})

}
