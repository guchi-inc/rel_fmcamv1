package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
未使用，不推荐绑定
*/
// 返回 省份资源信息
func (that *GovAddrApi) GAddrsProvinceList(ctx *gin.Context) {

	var (
		err error
	)

	GAddrsId := ctx.DefaultQuery("code", "")
	ProName := ctx.DefaultQuery("name", "")

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

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	if user.Ismanager != "2" || user.Isystem == "1" {
		datas, err := SysGovAddGroup.GAddrProvinceList(Page, PageSize, ProName, GAddrsId)
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

// 返回 城市资源信息
func (that *GovAddrApi) GAddrsCityList(ctx *gin.Context) {

	var (
		err error
	)

	GAddrsId := ctx.DefaultQuery("code", "")
	ProName := ctx.DefaultQuery("name", "")
	GProvId := ctx.DefaultQuery("province_code", "")

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

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	if user.Ismanager != "2" || user.Isystem == "1" {
		datas, err := SysGovAddGroup.GAddrCityList(Page, PageSize, ProName, GAddrsId, GProvId)
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

// 返回 区县资源信息
func (that *GovAddrApi) GAddrsAreaList(ctx *gin.Context) {

	var (
		err error
	)

	GAddrsId := ctx.DefaultQuery("code", "")
	ProName := ctx.DefaultQuery("name", "")
	GProvId := ctx.DefaultQuery("province_code", "")
	GCityId := ctx.DefaultQuery("city_code", "")

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

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	if user.Ismanager != "2" || user.Isystem == "1" {
		datas, err := SysGovAddGroup.GAddrAreaList(Page, PageSize, ProName, GAddrsId, GProvId, GCityId)
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

// 返回 街道资源信息
func (that *GovAddrApi) GAddrsStreetList(ctx *gin.Context) {

	var (
		err error
	)

	GAddrsId := ctx.DefaultQuery("code", "")
	ProName := ctx.DefaultQuery("name", "")
	GProvId := ctx.DefaultQuery("province_code", "")
	GCityId := ctx.DefaultQuery("city_code", "")
	GAreaId := ctx.DefaultQuery("area_code", "")

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

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	if user.Ismanager != "2" || user.Isystem == "1" {
		datas, err := SysGovAddGroup.GAddrStreetList(Page, PageSize, ProName, GAddrsId, GProvId, GCityId, GAreaId)
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
