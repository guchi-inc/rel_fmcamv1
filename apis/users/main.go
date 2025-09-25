// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/systems"
	"fmcam/systems/sysutils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type ApiGroup struct {
	UsersApi
}

var (
	UserUtil      = sysutils.NewUserUtil()
	DepUtils      = systems.DepUtils //使用系统层的 部门工具
	CustomerUtils = systems.CustomerUtils

	UserService = systems.ServiceGroupApp.SysUserGroup
	logger      = configs.NewLogger("apis-users")
	apidebug    = log.New(os.Stdout, "DEBUG -", 13)
)

type Handler interface {
	w()

	//服务检测
	Ping(ctx *gin.Context)

	//返回视频信息
	List(ctx *gin.Context)
}

func (h *UsersApi) w() {}

func NewUserApis() *ApiGroup {
	return &ApiGroup{
		UsersApi: UsersApi{},
	}
}
