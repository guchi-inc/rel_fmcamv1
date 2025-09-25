// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package alysis

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlysisApi struct{}

var (
	dataTypes = configs.GraphClientType
)

// 返回抓拍统计类型 数据
func (that *AlysisApi) DBCaptureDatas(ctx *gin.Context) {

	var (
		dataType = "capture"
	)

	StartAt := ctx.DefaultQuery("start_time", "")   //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")       //结束时间
	SubType := ctx.DefaultQuery("sub_type", "")     //子图类型名
	RangeType := ctx.DefaultQuery("range_type", "") //时间范围类型 date 天，week 周 month 月
	Device := ctx.DefaultQuery("device", "")        //设备唯一标记

	TwoModel := ctx.DefaultQuery("two_model", "") //不传该参数 则不是二维模式： 设备 + 时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		logapi.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)
		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	//鉴权
	tenUser, err := UserUtil.GetUserByCtx(ctx)
	if err != nil {
		helpers.JSONs(ctx, http.StatusUnauthorized, err)
		return
	}

	//二维模式
	isTwo := false
	if TwoModel != "" {
		isTwo = true
	}

	//查当前租户相关信息
	alt, errcap := DBs1Datas(dataType, tenUser.TenantId, StartAt, EndAt, SubType, RangeType, &Device, &isTwo)
	if errcap != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, errcap)
		return
	}

	//它向响应正文添加填充，以从驻留在与客户端不同的域中的服务器请求数据。
	//它填充请求头 "application/javascript"
	//ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*") //导致跨域问题
	logapi.Printf("DBCaptureDatas:%#v err:%#v \n", alt, err)
	if alt != nil {
		alt.DataType = dataType
	}
	helpers.JSONs(ctx, http.StatusOK, alt)
}

// 返回预警数据 信息
func (that *AlysisApi) DBWarnDatas(ctx *gin.Context) {
	var (
		dataType = "warn"
		warnsub  = configs.SubType
	)
	StartAt := ctx.DefaultQuery("start_time", "")   //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")       //结束时间
	SubType := ctx.DefaultQuery("sub_type", "")     //子图类型名 如 时间 date，位置 position，累计 accrued
	RangeType := ctx.DefaultQuery("range_type", "") //子图时间范围类型名 天 date 周 week 月 month
	Device := ctx.DefaultQuery("device", "")        //设备唯一标记, 有多个则使用逗号分隔
	TwoModel := ctx.DefaultQuery("two_model", "")   //不传该参数 则不是二维模式： 设备 + 时间

	logapi.Println("select SubType :", SubType, warnsub[SubType])
	if SubType != "" && !warnsub[SubType] {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("不支持的子图类型"))
		return
	}
	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		logapi.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)
		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	tenUser, err := UserUtil.GetUserByCtx(ctx)
	if err != nil {
		helpers.JSONs(ctx, http.StatusUnauthorized, err)
		return
	}

	//二维模式
	isTwo := false
	if TwoModel != "" {
		//设备id 不能为空
		if Device == "" {
			helpers.JSONs(ctx, code.ParamError, fmt.Errorf("多维统计Device不能为空:%v", Device))
			return
		}
		isTwo = true
	}

	alt, errwarn := DBs1Datas(dataType, tenUser.TenantId, StartAt, EndAt, SubType, RangeType, &Device, &isTwo)
	if errwarn != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, errwarn)
		return
	}
	if alt != nil {
		alt.DataType = dataType
	}
	ctx.JSONP(http.StatusOK, alt)
}

// 采集数据 信息
func (that *AlysisApi) DBHomeGatherDatas(ctx *gin.Context) {

	var (
		dataType = "gather"
	)

	StartAt := ctx.DefaultQuery("start_time", "")   //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")       //结束时间
	SubType := ctx.DefaultQuery("sub_type", "")     //子图类型名
	RangeType := ctx.DefaultQuery("range_type", "") //时间范围类型 date 天，week 周 month 月
	Device := ctx.DefaultQuery("device", "")        //设备唯一标记

	TwoModel := ctx.DefaultQuery("two_model", "") //不传该参数 则不是二维模式： 设备 + 时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		logapi.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)
		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	//它向响应正文添加填充，以从驻留在与客户端不同的域中的服务器请求数据。
	//它填充请求头 "application/javascript"
	//ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*") //导致跨域问题
	tenUser, err := UserUtil.GetUserByCtx(ctx)
	if err != nil {
		helpers.JSONs(ctx, http.StatusUnauthorized, err)
		return
	}

	//二维模式
	isTwo := false
	if TwoModel != "" {
		isTwo = true
	}

	alt, errgather := DBs1Datas(dataType, tenUser.TenantId, StartAt, EndAt, SubType, RangeType, &Device, &isTwo)
	logapi.Println("select alt:", alt != nil, " finished total err:", StartAt, "Fileter:", EndAt, err)
	if errgather != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return
	}
	if alt != nil {
		alt.DataType = dataType
	}
	helpers.JSONs(ctx, http.StatusOK, alt)

}

// 采集数据 信息
func (that *AlysisApi) DBGatherDatas(ctx *gin.Context) {

	var (
		dataType = "accrued_gather"
	)

	StartAt := ctx.DefaultQuery("start_time", "")   //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")       //结束时间
	SubType := ctx.DefaultQuery("sub_type", "")     //子图类型名
	RangeType := ctx.DefaultQuery("range_type", "") //时间范围类型 date 天，week 周 month 月
	Types := ctx.DefaultQuery("types", "")          //人员类型id : 18,55,65

	TwoModel := ctx.DefaultQuery("two_model", "") //不传该参数 则不是二维模式： 设备 + 时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		logapi.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)
		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	//它向响应正文添加填充，以从驻留在与客户端不同的域中的服务器请求数据。
	//它填充请求头 "application/javascript"
	//ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*") //导致跨域问题
	tenUser, err := UserUtil.GetUserByCtx(ctx)
	if err != nil {
		helpers.JSONs(ctx, http.StatusUnauthorized, err)
		return
	}

	//二维模式
	isTwo := false
	if TwoModel != "" {
		isTwo = true
	}

	alt, errAccruedgather := DBs1Datas(dataType, tenUser.TenantId, StartAt, EndAt, SubType, RangeType, &Types, &isTwo)
	logapi.Println("select errAccruedgather:", alt != nil, " finished total err:", StartAt, "Fileter:", EndAt, err)
	if errAccruedgather != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return
	}
	if alt != nil {
		alt.DataType = dataType
	}
	helpers.JSONs(ctx, http.StatusOK, alt)

}

// 最近预警数据 信息
func (that *AlysisApi) DBHomeLastWarns(ctx *gin.Context) {

	var (
		dataType = "warning"
	)

	StartAt := ctx.DefaultQuery("start_time", "")   //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")       //结束时间
	SubType := ctx.DefaultQuery("sub_type", "")     //子图类型名
	RangeType := ctx.DefaultQuery("range_type", "") //时间范围类型 date 天，week 周 month 月
	Device := ctx.DefaultQuery("device", "")        //设备唯一标记列表

	TwoModel := ctx.DefaultQuery("two_model", "") //不传该参数 则不是二维模式： 设备 + 时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		logapi.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)
		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	//鉴权
	tenUser, err := UserUtil.GetUserByCtx(ctx)
	if err != nil {
		helpers.JSONs(ctx, http.StatusUnauthorized, err)
		return
	}

	//二维模式
	isTwo := false
	if TwoModel != "" {
		isTwo = true
	}

	//查当前租户相关信息
	alt, errlastwarning := DBs1Datas(dataType, tenUser.TenantId, StartAt, EndAt, SubType, RangeType, &Device, &isTwo)
	if errlastwarning != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return
	}
	//它向响应正文添加填充，以从驻留在与客户端不同的域中的服务器请求数据。
	//它填充请求头 "application/javascript"
	//ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*") //导致跨域问题
	alt.DataType = dataType
	helpers.JSONs(ctx, http.StatusOK, alt)

}

// 心跳接口
func (that *AlysisApi) Pok(ctx *gin.Context) {

	// 服务鉴定
	helpers.JSONs(ctx, http.StatusOK, gin.H{"message": "user service Success", "data": 1, "status": "OK.", "code": code.Success})
}
