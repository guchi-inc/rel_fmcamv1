// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 新增用户
func (uc *UsersApi) Register(ctx *gin.Context) {

	// utils.RateLimit(ctx)
	var userRegisters erps.GcParamUser
	//无需解密 解密 Json
	if err := ctx.ShouldBindJSON(&userRegisters); err != nil {
		apidebug.Printf("should bind json to user:%v err:%v\n", userRegisters, err)
		//参数信息错误: username必须要
		helpers.JSONs(ctx, code.ParamBindError, fmt.Errorf("%v: %v", code.ZhCNText[code.ParamBindError], userRegisters))
		return
	}

	apidebug.Printf("receive params:%#v,  ", userRegisters)
	if userRegisters.Phonenum == "" {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v:Phonenum %v", code.ZhCNText[code.ParamError], userRegisters.Phonenum))
		return
	}

	if userRegisters.Password == "" {
		userRegisters.Password = configs.DefPassWord
	}

	apidebug.Printf("users.Password:%v \n", userRegisters.Password)
	userRegisters.Password = utils.HashPassword(&userRegisters.Password)
	apidebug.Printf("hash password:%v \n", userRegisters.Password)

	//如果登陆名没有设置，则默认使用 用户名作为登陆名
	if userRegisters.LoginName == "" {
		if userRegisters.Username == "" {
			userRegisters.LoginName = "es_" + utils.GenStandString(5)
		} else {
			newLogin := utils.PinMain(userRegisters.Username)
			userRegisters.LoginName = newLogin //users.Username
		}

	} else {
		userRegisters.LoginName = strings.ToLower(userRegisters.LoginName)
	}

	if userRegisters.Email != "" {
		if !strings.Contains(userRegisters.Email, "@") {
			helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v:%v", code.ZhCNText[code.ParamError], userRegisters.Email))
			return
		}
	} else {
		//强制填一个 模拟邮箱，因为需要用来找回密码
		userRegisters.Email = fmt.Sprintf("es_%v@demo.com", utils.GenStandString(5))
	}

	//电话不能为空
	if userRegisters.Phonenum == "" {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v:%v", code.ZhCNText[code.ParamError], userRegisters.Phonenum))
		return
	}

	if userRegisters.Ismanager == "" {
		userRegisters.Ismanager = "0"
	}

	if userRegisters.Isystem == "" {
		userRegisters.Isystem = "0"
	}

	newUser, err := UserService.Register(ctx, &userRegisters)
	if err != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		logger.Infof("Insert datas of users:%v err:%v", userRegisters, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": map[string]any{"id": newUser.ID, "member_id": newUser.MemberId, "login_name": newUser.LoginName}})
}

// 新增 客户端内部用户
func (uc *UsersApi) RegisterAccount(ctx *gin.Context) {

	// utils.RateLimit(ctx)
	var userRegisters erps.GcParamUser
	//无需解密 解密 Json
	if err := ctx.ShouldBindJSON(&userRegisters); err != nil {
		apidebug.Printf("should bind json to user:%v err:%v\n", userRegisters, err)
		//参数信息错误: username必须要
		helpers.JSONs(ctx, code.ParamBindError, fmt.Errorf("%v: %v", code.ZhCNText[code.ParamBindError], userRegisters))
		return
	}

	apidebug.Printf("receive params:%#v,  ", userRegisters)
	if userRegisters.Phonenum == "" {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v:Phonenum %v", code.ZhCNText[code.ParamError], userRegisters.Phonenum))
		return
	}

	if userRegisters.Password == "" {
		userRegisters.Password = configs.DefPassWord
	}

	apidebug.Printf("users.Password:%v \n", userRegisters.Password)
	userRegisters.Password = utils.HashPassword(&userRegisters.Password)
	apidebug.Printf("hash password:%v \n", userRegisters.Password)

	//如果登陆名没有设置，则默认使用 用户名作为登陆名
	if userRegisters.LoginName == "" {
		if userRegisters.Username == "" {
			userRegisters.LoginName = "es_" + utils.GenStandString(5)
		} else {
			newLogin := utils.PinMain(userRegisters.Username)
			userRegisters.LoginName = newLogin //users.Username
		}

	} else {
		userRegisters.LoginName = strings.ToLower(userRegisters.LoginName)
	}

	if userRegisters.Email != "" {
		if !strings.Contains(userRegisters.Email, "@") {
			helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v:%v", code.ZhCNText[code.ParamError], userRegisters.Email))
			return
		}
	} else {
		//强制填一个 模拟邮箱，因为需要用来找回密码
		userRegisters.Email = fmt.Sprintf("es_%v@demo.com", utils.GenStandString(5))
	}

	//电话不能为空
	if userRegisters.Phonenum == "" {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v:%v", code.ZhCNText[code.ParamError], userRegisters.Phonenum))
		return
	}

	if userRegisters.Ismanager == "" {
		userRegisters.Ismanager = "0"
	}

	if userRegisters.Isystem == "" {
		userRegisters.Isystem = "0"
	}

	if userRegisters.LeaderFlag != "" {
		userRegisters.Position = configs.LeaderFlagMap[userRegisters.LeaderFlag]
	}

	newUser, errInnerAccount := UserService.Register(ctx, &userRegisters)
	if errInnerAccount != nil {
		err := fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], errInnerAccount)
		logger.Infof("Insert datas of users:%v err:%v", userRegisters, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": map[string]any{"id": newUser.ID, "member_id": newUser.MemberId, "login_name": newUser.LoginName}})
}

// 管理员重置密码
func (uc *UsersApi) RestorePasswd(ctx *gin.Context) {

	// utils.RateLimit(ctx)
	var userReset erps.FmUser
	//无需解密 解密 Json
	if err := ctx.ShouldBindJSON(&userReset); err != nil {
		apidebug.Printf("should bind json to user:%v err:%v\n", userReset, err)
		//参数信息错误: username必须要
		helpers.JSONs(ctx, code.ParamBindError, fmt.Errorf("%v: %v", code.ZhCNText[code.ParamBindError], userReset))
		return
	}

	if userReset.Id == 0 {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v: %v", code.ZhCNText[code.ParamError], userReset))
		return
	}

	userReset.Password = configs.DefPassWord

	apidebug.Printf("users.Password:%v \n", userReset.Password)
	userReset.Password = utils.HashPassword(&userReset.Password)
	apidebug.Printf("hash password:%v \n", userReset.Password)

	//校验 id 并 限制, 非管理员（租户需要是付费账号才能操作） 只能更新自己的信息
	user, err := UserService.GetUserByCtxToken(ctx)
	apidebug.Printf("user info:%v err:%v \n", user, err)
	if err != nil || user == nil {
		helpers.JSONs(ctx, code.AuthorizationError, fmt.Errorf("该次调用鉴权错误%v", err))
		return
	}

	//当前账号需要是内置管理员 is_manager = 3（付费主账号）才可以使用重置密码 功能重置，否则只能修改自己的
	if user.Ismanager != "3" {
		if int(user.Id) != int(userReset.Id) {
			err := fmt.Errorf("只能修改自己的信息 %v:%v", code.ZhCNText[code.AuthorizedPowerError], user.Id == userReset.Id)
			helpers.JSONs(ctx, code.AuthorizedPowerError, err)
			return
		}
	}

	//非系统内置账号，只能更新 本租户内的 账号
	if user.Isystem != "1" {
		//当前账号 租户号与 目标账号的租户号不同
		if user.TenantId != userReset.TenantId {
			err := fmt.Errorf("只能修改本店的账号 %v:%v", code.ZhCNText[code.AuthorizedPowerError], user.TenantId == userReset.TenantId)
			helpers.JSONs(ctx, code.AuthorizedPowerError, err)
			return
		}
	}

	//执行重置
	userid, err := UserService.BussUpdateUser(&userReset)
	if err != nil {
		err := fmt.Errorf("登陆密码修改为敏感操作，失败请找管理员协助 %v:%v: err:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], userReset.Id, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": userid})

}
