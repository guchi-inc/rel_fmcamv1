package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询 预警处理信息记录
func (that *FmClientApi) GetAlerts(c *gin.Context) {

	alert_level, _ := strconv.Atoi(c.DefaultQuery("alert_level", ""))
	ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	capture_log_id, _ := strconv.Atoi(c.DefaultQuery("capture_log_id", ""))
	device_id, _ := strconv.Atoi(c.DefaultQuery("device_id", ""))
	statusTr := c.DefaultQuery("status", "")

	tenantId := c.DefaultQuery("tenant_id", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished StartAt err:", StartAt, "EndAt:", EndAt, err)
		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	datas, err := SerEntORMApp.AlertsList(&Page, &PageSize, &ids, &device_id, &alert_level, &capture_log_id, c, statusTr, tenantId, StartAt, EndAt)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}

// 预警处理 更新
func (that *FmClientApi) UpdateAlerts(c *gin.Context) {

	var (
		PParm = seclients.Alerts{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.AlertUpdate(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
