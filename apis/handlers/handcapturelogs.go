package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询抓拍记录
func (that *FmClientApi) GetCaptureLogs(c *gin.Context) {

	typeIds, _ := strconv.Atoi(c.DefaultQuery("type_id", "0")) // 人员类型id
	deviceids, _ := strconv.Atoi(c.DefaultQuery("device_id", "0"))
	tenantId := c.DefaultQuery("tenant_id", "")
	profileId, _ := strconv.Atoi(c.DefaultQuery("matched_profile_id", ""))
	FuncType := c.DefaultQuery("func_type", "") //0 采集 1 抓拍
	status := c.DefaultQuery("status", "")      //处理状态 0，1，2
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	devid := int64(deviceids)
	profId := int64(profileId)
	typeId := int64(typeIds)

	var funcT int8
	if FuncType == "0" {
		funcT = 0
	} else {
		funcT = 1
	}

	PParm := seclients.CaptureLogInfo{FuncType: &funcT}
	if devid != 0 {
		PParm.DeviceID = &devid
	}
	if profId != 0 {
		PParm.MatchedProfileID = &profId
	}
	if typeId != 0 {
		PParm.TypeId = &typeId
	}
	if tenantId != "" {
		PParm.TenantID = &tenantId
	}
	if status != "" {
		newSta, err := strconv.Atoi(status)
		apilog.Printf("query status:%#v err:%#v \n", status, err)
		PParm.Status = &newSta
	}

	//排除采集查询条件
	tagFilter := seclients.FieldTags{SortTagFilter: "gat"}
	datas, errcap := SerEntORMApp.CaptureLogsList(int64(Page), int64(PageSize), StartAt, EndAt, c, &PParm, &tagFilter)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if errcap != nil {
		helpers.JSONs(c, code.NullData, errcap)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}

// 查询 采集记录
func (that *FmClientApi) GetGatherLogs(c *gin.Context) {

	profileTypeID, _ := strconv.Atoi(c.DefaultQuery("type_id", "0"))
	faceID, _ := strconv.Atoi(c.DefaultQuery("face_id", "0"))
	profileId, _ := strconv.Atoi(c.DefaultQuery("matched_profile_id", ""))
	deviceids, _ := strconv.Atoi(c.DefaultQuery("device_id", "0"))
	tenantId := c.DefaultQuery("tenant_id", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	RoomId := c.DefaultQuery("room_id", "")     //房间号
	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	devid := int64(deviceids)
	profId := int64(profileId)
	typeId := int64(profileTypeID)

	var funcT int8 = 0

	PParm := seclients.CaptureLogInfo{FuncType: &funcT}
	if devid != 0 {
		PParm.DeviceID = &devid
	}
	if profId != 0 {
		PParm.MatchedProfileID = &profId
	}
	if RoomId != "" {
		PParm.RoomId = &RoomId
	}
	if typeId != 0 {
		PParm.TypeId = &typeId
	}
	if tenantId != "" {
		PParm.TenantID = &tenantId
	}
	if faceID != 0 {
		FaceID := int64(faceID)
		PParm.FaceId = &FaceID
	}

	//排除抓拍字段
	tagFilter := seclients.FieldTags{SortTagFilter: "cap"}
	datas, err := SerEntORMApp.CaptureLogsList(int64(Page), int64(PageSize), StartAt, EndAt, c, &PParm, &tagFilter)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}
