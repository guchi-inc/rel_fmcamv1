// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 返回部门全部信息
func (that *UsersApi) GetDepments(ctx *gin.Context) {

	Name := ctx.DefaultQuery("dep_name", "") //部门名

	fmt.Printf("UserApi GetDepments dep_name :%v \n", Name)

	if Name == "" {
		datas, err := DepUtils.GeDepartmentNameList(true)
		if err != nil {
			logger.Infof("Get datas of depment Err:%v", err)
			err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
			helpers.JSONs(ctx, code.AdminListError, err)
			return
		}
		helpers.JSONs(ctx, code.Success, datas)
	} else {
		datas, err := DepUtils.GeDepartmentByName(Name)
		if err != nil {
			logger.Infof("Get datas of depment Err:%v", err)
			err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
			helpers.JSONs(ctx, code.AdminListError, err)
			return
		}
		helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
	}

}

// 新增部门信息
func (that *UsersApi) NewDepments(ctx *gin.Context) {

	var (
		valueStructs = &erps.GcDepartmentStruct{}
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
	if valueStructs == nil || err != nil {
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

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name := user.LoginName
	if Name != "admin" || (user.Ismanager == "2" || user.Isystem != "1") {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], Name)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	valueStructs.Creator = Name
	datas, err := DepUtils.InsertDepartment(valueStructs)
	fmt.Printf("menu update rst:%v, err:%v\n", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], Name)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})

}

// 更新部门信息
func (that *UsersApi) UpdateDepments(ctx *gin.Context) {

	var (
		valueStructs = &erps.GcDepartmentStruct{}
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
	if valueStructs == nil || err != nil {
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

	//管理员需要前端调用menu list 接口获取菜单功能列表
	Name := user.LoginName
	if Name != "admin" || (user.Ismanager == "2" || user.Isystem != "1") {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], Name)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	valueStructs.Creator = Name
	datas, err := DepUtils.UpdateDepartment(valueStructs)
	fmt.Printf("menu update rst:%v, err:%v\n", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], Name)
		logger.Infof("%v", err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})

}

//客户类型
