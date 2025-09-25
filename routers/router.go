// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package routers

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/middleware"
	"fmcam/ctrlib/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Handlers = service.Handlers
	RS       = service.RS
	Router   = service.DefaultRouter
	Logger   = configs.NewLogger("routers-main") //configs.GLogs
)

// login服务 后台管理服务
func SetupRouters(sg *gin.RouterGroup) {

	// 路由注册
	router := Routers{}
	//总额大屏
	// router.Powers.InitHomeList(sg)

	// router.Maters.InitFileImg(sg)
	sg.Use(middleware.Cors())
	router.Users.InitUsers(sg)

	// 资源管理
	router.Recourse.InitResDevice(sg)

	//租户管理
	router.Recourse.InitCustomerRouter(sg)

}

// 操作页面 main 主路由注册
func InitRouter() http.Handler {

	// 捕获所有路由并返回index.html
	RS.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 404, "message": "404 not found"})
	})

	RS.Routers()
	SetupRouters(Router)

	return RS
}
