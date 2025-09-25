// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 返回全部后台管理菜单信息
func (that *UsersApi) MenuLists(ctx *gin.Context) {

	var (
		err  error
		Name string
	)
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	user, err := UserService.GetUserByCtxToken(ctx)
	fmt.Printf("receive params:%v\n", user)
	if err != nil || user == nil || user.LoginName == "" {
		err = fmt.Errorf("%v:%v:%v", code.ZhCNText[code.AuthorizedDeleteError], user, err)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AuthorizedDeleteError, err)
		return
	}

	fmt.Printf("db user password  user password en:%v\n", user.Password)

	//非管理员权限不可使用后台管理
	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name = user.LoginName

	tenantId := user.TenantId
	if user.Ismanager != "2" && user.Isystem == "1" {
		datas, err := DepUtils.GeMenuFullList(0, 0, Name, "", "")
		if err != nil {
			err = fmt.Errorf("%v:%v:%v", code.ZhCNText[code.AdminListError], user, err)
			logger.Infof("get menu:%v by user info:%#v, AdminListError err:%v\n", user, user, err)
			helpers.JSONs(ctx, code.AdminListError, err)
			return
		}
		helpers.JSONs(ctx, http.StatusOK, datas)
	} else {
		//其他管理员 根据自己角色组权限  获取对应的菜单列表
		datas, err := DepUtils.GeMenuLimitsList(Name, tenantId)
		if err != nil {
			err = fmt.Errorf("%v:%v:%v", code.ZhCNText[code.AdminListError], user, err)
			logger.Infof("get menu:%v by user info:%#v, AdminListError err:%v\n", user, user, err)
			helpers.JSONs(ctx, code.AdminListError, err)
			return
		}
		helpers.JSONs(ctx, code.Success, datas)
	}

}

// 返回 可管理功能信息列表
func (that *UsersApi) GetMenuList(ctx *gin.Context) {

	menuNumber := ctx.DefaultQuery("menu_number", "")  //功能编号
	description := ctx.DefaultQuery("description", "") //功能描述名
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

	//非管理员 或内部账户 权限不可使用后台管理
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	dbUser, err := UserService.GetUserInfo(user.LoginName)
	if dbUser == nil || err != nil {
		err = fmt.Errorf("%v:%v:%v", code.ZhCNText[code.AuthorizedDeleteError], dbUser, err)
		logger.Infof("get dbUser:%v by user info:%#v, AuthorizedDeleteError err:%v\n", dbUser, user, err)
		helpers.JSONs(ctx, code.AuthorizedDeleteError, err)
		return
	}

	fmt.Printf("db user password:%v, user password en:%v\n", dbUser.Password, user.Password)

	//非管理员权限不可使用后台管理
	//管理员需要前端调用menu list 接口获取菜单功能列表
	datas, err := DepUtils.GeMenuFullList(Page, PageSize, user.LoginName, menuNumber, description)
	if err != nil {
		err = fmt.Errorf("%v:%v:%v", code.ZhCNText[code.AdminListError], dbUser, err)
		logger.Infof("get menu:%v by user info:%#v, AdminListError err:%v\n", dbUser, user, err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, http.StatusOK, datas)

}

// POST 修改后台管理菜单信息，仅admin可执行
func (that *UsersApi) MenusUpdate(ctx *gin.Context) {

	var (
		valueStructs []*erps.GcMenu
		err          error
	)

	//校验token信息
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	err = ctx.ShouldBindJSON(&valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v\n", valueStructs, err)
	if len(valueStructs) == 0 || err != nil {
		err = fmt.Errorf("%v err:%v", code.ZhCNText[code.ParamBindError], err)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//鉴权
	tokenString := authHeader[len(configs.AuthHeader):]

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		err = fmt.Errorf("%v :%v err:%v", code.ZhCNText[code.AuthorizedDetailError], tokenString, err)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AuthorizedDetailError, err)
		return
	}

	//校验token信息,获取用户信息
	uid, _ := utils.GetIdFromToken(token)

	user, err := UserService.GetUserById(int(uid))
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizationError], uid)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}
	//管理员需要前端调用menu list 接口获取菜单功能列表  全部满足条件
	Name := user.LoginName
	if Name != "admin" || user.Ismanager == "2" || user.Isystem != "1" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], Name)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	datas, err := DepUtils.GeMenuUpdate(Name, valueStructs)
	fmt.Printf("menu update rst:%v, err:%v\n", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], Name)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}
