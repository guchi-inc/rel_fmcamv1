package routers

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/middleware"
	"fmcam/ctrlib/service"
	"fmcam/routers/streams"
	"fmcam/routers/sysrouter"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	ClientsApi = sysrouter.AlyRouter{}
	StreamsApi = streams.StreamRouter{}
)

// 客户端操作路由，包括mqtt
func ClientRouters(sg *gin.RouterGroup) {

	router := ClientsApi

	// 客户端操作
	router.InitAlysis(sg)

	sg.Use(middleware.Cors())
	//文件操作
	StreamsApi.InitFileObs(sg)

	//面部数据接口
	StreamsApi.InitFaceDevice(sg)
}

// 客户端服务注册
func InitClientRouter() http.Handler {

	MainEngine := service.MakeNewEngine()
	MainEngine.Use(gzip.Gzip(gzip.BestCompression))
	RS := service.NewGroupsServers(MainEngine)
	// 捕获所有路由并返回index.html
	RS.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404"})
	})
	//接口版本
	NewRouter := RS.Group(configs.ApiVersion)

	//注册客户端 接口到服务
	ClientRouters(NewRouter)

	return RS
}
