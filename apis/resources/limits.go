package resources

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 返回全部后台管理 资源组信息
func (that *LimitsApi) LimitsList(ctx *gin.Context) {

	var (
		err  error
		Name string
	)

	LimitsId := ctx.DefaultQuery("limits_id", "")
	TenantId := ctx.DefaultQuery("tenant_id", "")

	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	user, err := UserService.GetUserByCtxToken(ctx)

	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	apilog.Printf("user sign info:%#v, default en:%v", user, user.Email)
	dbUser, err := UserService.GetUserInfo(user.LoginName)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return

	}

	apilog.Printf("db user password:%v, user password en:%v", dbUser.Password, user.Password)

	//非管理员权限不可使用后台管理
	if dbUser.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name = user.LoginName
	if dbUser.Ismanager != "2" || dbUser.Isystem == "1" {
		datas, err := LimitsServiceGroup.GetLimitsList(Page, PageSize, Name, LimitsId, TenantId)
		if err != nil {
			apilog.Printf("Get datas of depment Err:%v", err)
			err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
			helpers.JSONs(ctx, code.AdminListError, err)
			return
		}
		helpers.JSONs(ctx, code.Success, datas)
	} else {
		//其他管理员无权限
		helpers.JSONs(ctx, code.AuthorizedPowerError, code.ZhCNText[code.AuthorizedPowerError])
		return
	}
}

// POST 修改后台管理菜单信息，仅admin可执行
func (that *LimitsApi) LimitsUpdate(ctx *gin.Context) {

	var (
		vStructs erps.GcLimits
		err      error
	)

	//校验token信息
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	err = ctx.ShouldBindJSON(&vStructs)
	apilog.Printf("receive params rst:%#v, err:%v", vStructs, err)

	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		apilog.Printf("%v:%v", vStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)

	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表,三者满足一个即可
	Name := user.LoginName
	if Name != "admin" && user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizedPowerError], Name)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//执行更新
	datas, err := LimitsServiceGroup.PostLimitsUpdate(&vStructs)
	apilog.Printf("role update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AdminModifyPersonalInfoError], Name)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// POST 新增角色，仅授权者可执行
func (that *LimitsApi) LimitsInsert(ctx *gin.Context) {

	var (
		vLimits = erps.GcLimits{}
		err     error
	)

	//校验token信息
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	err = ctx.ShouldBindJSON(&vLimits)
	fmt.Printf("receive params rst:%#v, err:%v", vLimits, err)

	//参数判空
	if vLimits.LimitsId == "" || vLimits.Description == "" || vLimits.LimitsName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		apilog.Printf("%v:%v", vLimits, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		apilog.Printf("%v:%v", vLimits, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)

	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name := user.LoginName
	//doing
	datas, err := LimitsServiceGroup.LimitsNew(Name, &vLimits)
	fmt.Printf("role update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		apilog.Printf("%v:%v", vLimits, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}
