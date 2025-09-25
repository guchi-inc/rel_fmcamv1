package resources

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 返回全部后台管理【角色/用户】信息
func (that *LimitsApi) RoleLimitsList(ctx *gin.Context) {

	var (
		err error
	)

	btn_id, _ := strconv.Atoi(ctx.DefaultQuery("btn_id", "0")) //功能号

	tenant_id := ctx.DefaultQuery("tenant_id", "") //租户号

	role_name := ctx.DefaultQuery("role_name", "") //角色名
	username := ctx.DefaultQuery("username", "")   //用户名
	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	StartAt := ctx.DefaultQuery("start_time", "") //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")     //结束时间

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}
	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	user, err := UserService.GetUserByCtxToken(ctx)

	fmt.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	apilog.Printf("user sign info:%#v, default en:%v", user, user.Email)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	datas, err := LimitsServiceGroup.GetRoleLimitsList(Page, PageSize, btn_id, role_name, username, StartAt, EndAt, tenant_id)
	if err != nil {

		apilog.Printf("Get datas of depment Err:%v", err)
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, datas)

}

// 修改后台 角色功能，单个删除
func (that *LimitsApi) RoleLimitsOneUpdate(ctx *gin.Context) {

	var (
		valueStructs erps.GcRoleLimits
		err          error
	)

	err = ctx.ShouldBindJSON(&valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v", valueStructs, err)
	if err != nil {
		err = fmt.Errorf("%v:%v err:%v", code.ZhCNText[code.ParamBindError], valueStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)

	fmt.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	apilog.Printf("user sign info:%#v, default en:%v", user, user.Email)

	datas, err := LimitsServiceGroup.PostLimitRolesUpdate(&valueStructs)
	fmt.Printf("role user update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// 新增角色和资源关系
func (that *LimitsApi) RoleLimitsInsert(ctx *gin.Context) {

	var (
		valueStructLists erps.FmRoleLimitNews
		err              error
	)

	//校验token信息 只能通过 jwt 认证
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	err = ctx.ShouldBindJSON(&valueStructLists)
	fmt.Printf("receive params rst:%#v, err:%v", valueStructLists, err)
	if err != nil || valueStructLists.RoleName == "" || valueStructLists.BtnIdList == nil {
		err = fmt.Errorf("%v: 角色 功能号 资源组号缺一不可:%v", code.ZhCNText[code.ParamBindError], err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)
	fmt.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//管理员  内部账号和管理员属性缺一不可
	Name := user.LoginName
	if user.Ismanager == "2" || user.Isystem != "1" {
		err = fmt.Errorf("%v:%v  ", code.ZhCNText[code.AuthorizedPowerError], Name)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	datas, err := LimitsServiceGroup.InsertRoleLimits(Name, &valueStructLists)
	fmt.Printf("fm_role_limits insert rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// 查已有角色和资源关系
func (that *LimitsApi) RoleLimitsListHistory(ctx *gin.Context) {

	var (
		err error
	)

	//校验token信息
	role_name := ctx.DefaultQuery("role_name", "") //角色名

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)

	fmt.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//管理员  内部账号和管理员属性即可
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v  ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	datas, err := LimitsServiceGroup.GetRoleBtnMenus(role_name)
	fmt.Printf("fm_menu list rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, datas)
}

// 删除 角色和资源关系
func (that *LimitsApi) DelRoleLimits(ctx *gin.Context) {

	var (
		valueStructs erps.GcRoleLimits
		err          error
	)

	//校验信息
	err = ctx.ShouldBindJSON(&valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v", valueStructs, err)
	if err != nil || valueStructs.RoleName == "" || valueStructs.BtnId == 0 {
		err = fmt.Errorf("%v: err:%v", code.ZhCNText[code.ParamBindError], err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)
	fmt.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v  ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	datas, err := LimitsServiceGroup.DelRoleLimits(valueStructs.RoleName, valueStructs.BtnId, valueStructs.ID)
	fmt.Printf("fm_role_limits insert rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}
