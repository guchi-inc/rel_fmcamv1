package sysrouter

import (
	"fmcam/ctrlib/middleware"

	"github.com/gin-gonic/gin"
)

// 用户资源查询
func (s *AlyRouter) InitAlysis(r *gin.RouterGroup) {

	alyCtl := ApiClientsCtl

	r.GET("/ok", alyCtl.Pok) //当前服务测试信息
	// 路由
	alys := r.Group("/aly")
	// alys.GET("/warnings", alyCtl.MorkHomeLastWarns)       //模拟
	// alys.GET("/week/captures", alyCtl.MorkCaptureDatas)   //模拟
	// alys.GET("/week/warns", alyCtl.MorkWarnDatas)         //模拟
	// alys.GET("/week/gathers", alyCtl.MorkHomeGatherDatas) //模拟

	alys.GET("/warning", alyCtl.DBHomeLastWarns)
	alys.GET("/week/warn", alyCtl.DBWarnDatas)
	//采集统计
	alys.GET("/week/gather", alyCtl.DBHomeGatherDatas)
	//抓拍统计
	alys.GET("/week/capture", alyCtl.DBCaptureDatas)
	//采集统计图
	alys.GET("/accrued/gather", alyCtl.DBGatherDatas)
	//轨迹查询
	alys.GET("/history/capture/list", alyCtl.CaptionHistory)

	// WebSocket 连接路由
	// 接收消息并广播
	r.Use(middleware.AuthorizeJWT())
	r.POST("/alert/news", alyCtl.PostAlertNews)
	stream := r.Group("/ws")
	stream.GET("/alert", alyCtl.AlertWsHandler)

}
