// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 返回全部后台管理角色信息
func (that *UsersApi) RolesList(ctx *gin.Context) {

	var (
		err  error
		Name string
	)

	RoleName := ctx.DefaultQuery("role_name", "")
	TenantId := ctx.DefaultQuery("tenant_id", "")

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
		apidebug.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	user, err := UserService.GetUserByCtxToken(ctx)

	apidebug.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	logger.Infof("user sign info:%#v, default en:%v", user, user.Email)
	dbUser, err := UserService.GetUserInfo(user.LoginName)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}

	apidebug.Printf("db user password:%v, user password en:%v", dbUser.Password, user.Password)

	//非管理员权限不可使用后台管理
	if dbUser.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name = user.LoginName
	if dbUser.Ismanager != "2" || dbUser.Isystem == "1" {
		//内置系统账号 可以查询任意租户
		if dbUser.Isystem != "1" {
			TenantId = dbUser.TenantId
		}
		datas, err := DepUtils.GetRolesList(Page, PageSize, Name, RoleName, TenantId, StartAt, EndAt)
		if err != nil {
			logger.Infof("Get datas of depment Err:%v", err)
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

// POST 新增角色，仅授权者可执行
func (that *UsersApi) RoleInsert(ctx *gin.Context) {

	var (
		valueStructs = erps.GcRoleParam{}
		err          error
	)

	err = ctx.ShouldBindJSON(&valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v", valueStructs, err)
	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		logger.Infow("%v:%v", valueStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return

	}
	//参数判空
	if valueStructs.RoleName == "" || valueStructs.Description == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		logger.Infow("%v:%v", valueStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	user, err := UserService.GetUserByCtxToken(ctx)
	apidebug.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//非管理员 非系统内置账号
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name := user.LoginName
	//主账号的租户号
	valueStructs.TenantId = user.TenantId
	//doing
	datas, err := DepUtils.InsertRole(Name, &valueStructs)
	fmt.Printf("role update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		logger.Infow("%v:%v", valueStructs, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// POST 修改后台管理菜单信息，仅admin可执行
func (that *UsersApi) RoleUpdate(ctx *gin.Context) {

	var (
		valueStructs erps.GcRole
		err          error
	)

	//校验token信息

	err = ctx.ShouldBindJSON(&valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v", valueStructs, err)

	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		logger.Infow("%v:%v", valueStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)
	apidebug.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//非管理员 非系统内置账号
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//执行更新
	datas, err := DepUtils.PostRoleUpdate(&valueStructs)
	fmt.Printf("role update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AdminModifyPersonalInfoError], user.LoginName)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// 返回全部后台管理【角色/用户】信息
func (that *UsersApi) RoleUserList(ctx *gin.Context) {

	var (
		err  error
		Name string
	)

	login_name := ctx.DefaultQuery("login_name", "") //登录账号

	role_name := ctx.DefaultQuery("role_name", "") //角色名
	role_type := ctx.DefaultQuery("type", "")      //角色类型名
	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}
	user, err := UserService.GetUserByCtxToken(ctx)

	fmt.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	logger.Infof("user sign info:%#v, default en:%v", user, user.Email)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name = login_name

	datas, err := DepUtils.GcRoleUserList(Page, PageSize, role_type, Name, role_name)
	if err != nil {

		logger.Infof("Get datas of depment Err:%v", err)
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, datas)

}

// POST 修改后台管理菜单信息，批量更新，用于脚本执行
func (that *UsersApi) RoleUserOneUpdate(ctx *gin.Context) {

	var (
		valueStructs erps.GcRoleUser
		err          error
	)

	//校验token信息,获取登陆账号信息
	user, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		err = fmt.Errorf("%v:%v err:%v", code.ZhCNText[code.AuthorizationError], valueStructs, err)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	err = ctx.ShouldBindJSON(&valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v", valueStructs, err)
	if valueStructs.RoleName == "" || err != nil {
		err = fmt.Errorf("%v:%v err:%v", code.ZhCNText[code.ParamBindError], valueStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name := user.LoginName
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizedPowerError], Name)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	var UpdateDeleteFlag = ""
	if valueStructs.DeleteFlag != "" {
		UpdateDeleteFlag = valueStructs.DeleteFlag
	}
	var TenantId string
	if user.TenantId != configs.DefTenantId {
		TenantId = user.TenantId
	}
	datas, err := DepUtils.UpdateOneRoleUsers(&valueStructs, UpdateDeleteFlag, TenantId)
	fmt.Printf("role user update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// 新增角色和用户关系
func (that *UsersApi) RoleUserInsert(ctx *gin.Context) {

	var (
		valueStructs erps.GcRoleUser
		err          error
	)

	//校验token信息
	err = ctx.ShouldBindJSON(&valueStructs)
	apidebug.Printf("receive params rst:%#v, err:%v", valueStructs, err)
	if err != nil || valueStructs.RoleID == 0 || valueStructs.UserID == 0 {
		err = fmt.Errorf("%v: err:%v", code.ZhCNText[code.ParamBindError], err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	//校验token信息,获取登陆账号信息
	user, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		err = fmt.Errorf("%v:%v  ", code.ZhCNText[code.AuthorizationError], err)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name := user.LoginName
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v  ", code.ZhCNText[code.AuthorizedPowerError], Name)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	datas, err := DepUtils.InsertRoleUsers(Name, &valueStructs)
	apidebug.Printf("fm_role_user insert rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyDataError], err)
		helpers.JSONs(ctx, code.AdminModifyDataError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}
