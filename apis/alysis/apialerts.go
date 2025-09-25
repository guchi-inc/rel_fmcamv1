package alysis

import (
	"encoding/json"
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/capturelogs"
	"fmcam/systems/genclients/devices"
	"fmcam/systems/genclients/fmalertdefinition"
	"fmcam/systems/genclients/predicate"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 全局 WebSocket 客户端池
var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// 预警通知
func (that *AlysisApi) AlertWsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logapi.Println("WebSocket 升级失败:", err, c.ClientIP())
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	logapi.Println("新客户端已连接 WebSocket:", c.ClientIP())

	// 保持连接，直到客户端断开
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			logapi.Println("客户端断开:", err, c.ClientIP())
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
	}
}

// 预警接收
func (that *AlysisApi) PostAlertNews(c *gin.Context) {
	var news seclients.AlertNews
	if err := c.ShouldBindJSON(&news); err != nil {
		logapi.Println("客户端参数错误:", err, c.ClientIP())
		helpers.JSONs(c, code.ParamBindError, gin.H{"error": "JSON 解析失败"})
		return
	}

	if news.CaptureLogID != 0 {

		var conds = []predicate.CaptureLogs{}
		conds = append(conds, capturelogs.IDEQ(news.CaptureLogID))
		//图像信息
		capLog, err := SerEntORMApp.CaptureLogInfoById(c, conds)
		logapi.Println("预警抓拍id:", news.CaptureLogID, err)

		if err != nil {
			helpers.JSONs(c, code.ParamBindError, gin.H{"error": "capture log 查询失败" + err.Error()})
			return
		}

		news.CaptureImageUrl = capLog.CaptureImageURL
		news.MatchedProfileId = &capLog.MatchedProfileID
		news.TenantID = capLog.TenantID.String()

		// 预警动作查询
		var condDefs = []predicate.FmAlertDefinition{}
		condDefs = append(condDefs, fmalertdefinition.ProfileTypeIDEQ(*news.ProfileTypeId))
		alertDef, err := SerEntORMApp.LevelAlertInfo(c, condDefs)
		if err != nil {
			helpers.JSONs(c, code.ParamBindError, gin.H{"error": "预警动作查询失败" + err.Error()})
			return
		}

		news.Action = &alertDef.Action
		news.AlertGroupId = &alertDef.AlertGroupID
		news.AlarmSound = &alertDef.AlarmSound
		news.Level = &alertDef.Level
		to := time.Now().Local().In(configs.Loc)
		news.CreatedAt = &to
		//位置信息
		devconds := []predicate.Devices{}
		devconds = append(devconds, devices.IDEQ(capLog.DeviceID))
		devices, err := SerEntORMApp.DevicesList(1, 1, c, devconds, "")
		logapi.Println("预警抓拍id:", news.CaptureLogID, err)
		if err != nil || len(devices.Data) == 0 {
			helpers.JSONs(c, code.ParamBindError, gin.H{"error": "capture log 查询失败" + err.Error()})
			return
		}
		news.Location = &devices.Data[0].Location
		news.DeviceID = capLog.DeviceID
		logapi.Printf("预警弹窗news:%#v \n err:%#v\n", news, err)
		// 广播消息
		mutex.Lock()
		go that.broadcast(&news)
		mutex.Unlock()
	}

	helpers.JSONs(c, code.Success, gin.H{"data": news, "message": "消息已推送"})

}

// 广播函数：推送消息给所有客户端
func (that *AlysisApi) broadcast(news *seclients.AlertNews) {
	mutex.Lock()
	defer mutex.Unlock()

	data, _ := json.Marshal(news)

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			logapi.Println("发送失败，移除客户端:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
