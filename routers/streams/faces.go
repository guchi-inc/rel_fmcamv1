package streams

import (
	"fmcam/common/databases"
	"fmcam/ctrlib/middleware"

	"github.com/gin-gonic/gin"
)

func (s *StreamRouter) InitFaceDevice(r *gin.RouterGroup) {

	//  路由
	fcr := r.Group("/fs")
	clientHand := ApiStreamCtl.ClientFaceApis
	fcr.GET("/profiles/type/list", clientHand.GetProfileType)                //人员类型列表
	fcr.POST("/profiles/type/new", clientHand.NewProfileType)                //人员类型新增
	fcr.POST("/profiles/type/update", clientHand.UpdateProfileType)          //人员类型修改
	fcr.POST("/profiles/group/new", clientHand.AlertGroupsNew)               //人员类型和组新增
	fcr.GET("/profiles/group/mapping", clientHand.GetGroupTypeMapping)       //人员类型和组映射记录
	fcr.GET("/profiles/group/info/list", clientHand.GetGroupTypeMappingInfo) //人员类型映射全信息

	fcr.GET("/profiles/list", clientHand.GetProfiles)       //人员记录
	fcr.POST("/profiles/update", clientHand.UpdateProfiles) //更新人员信息
	fcr.POST("/profiles/new", clientHand.NewProfiles)       //新增人员信息

	fcr.GET("/faces/list", clientHand.GetFaces)                  //人脸列表
	fcr.POST("/faces/update", clientHand.UpdateFace)             //人脸删除 ，用于 工作人员离职
	fcr.GET("/temporaryfaces/list", clientHand.GetTemporaryFace) //临时面部
	fcr.GET("/devices/list", clientHand.GetDevices)              //设备列表
	fcr.GET("/capturelogs/list", clientHand.GetCaptureLogs)      //抓拍记录
	fcr.GET("/gathers/list", clientHand.GetGatherLogs)           //采集 记录

	//预警数据信息列表
	fcr.GET("/alerts/list", clientHand.GetAlerts)
	fcr.POST("/alerts/update", clientHand.UpdateAlerts) //处理预警信息

	acr := r.Group("/alerts")
	acr.GET("/group/list", clientHand.GetAlertGroups)       //预警组信息
	acr.GET("/group/info", clientHand.GetAlertGroupsInfo)   //预警组信息
	acr.POST("/group/new", clientHand.AlertGroupsNew)       //预警组新增
	acr.POST("/group/update", clientHand.AlertGroupsUpdate) //预警组修改

	acr.GET("/level/list", clientHand.AlertLevelList)      //预警级别信息列表
	acr.GET("/level/info", clientHand.AlertLevelInfo)      //预警级别信息
	acr.POST("/level/new", clientHand.AlertLevelNew)       //预警级别信息新增
	acr.POST("/level/update", clientHand.AlertLevelUpdate) //预警级别信息修改

	//第三方调用鉴权
	thirdFace := r.Group("/third")
	thirdFace.Use(middleware.AuthorizeJWT())
	thirdFace.POST("/guest/checkout", clientHand.UpdateProfileCheckout)
	thirdFace.POST("/guest/checkin", clientHand.ProfileCheckin)
	//检查签名和时间戳
	thirdFace.POST("/yunlv/profile", databases.AuthMiddleware(), clientHand.YunLvRequest)

}
