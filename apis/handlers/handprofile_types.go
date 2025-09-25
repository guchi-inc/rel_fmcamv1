package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/predicate"
	"fmcam/systems/genclients/profiletype"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询人员类型
func (that *FmClientApi) GetProfileType(c *gin.Context) {

	name := c.DefaultQuery("type_name", "")
	tenantId := c.DefaultQuery("tenant_id", "")
	typeGroupId, _ := strconv.Atoi(c.DefaultQuery("type_group_id", "0"))
	typeCode := c.DefaultQuery("type_code", "")
	enabled := c.DefaultQuery("enabled", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.ProfileType{}

	if name != "" {
		conds = append(conds, profiletype.TypeNameContains(name))
	}
	if typeCode != "" {
		conds = append(conds, profiletype.TypeCodeContains(typeCode))
	}

	var enabledNil *int
	if enabled != "" {
		valEnabled, err := strconv.Atoi(enabled)
		if err != nil {
			helpers.JSONs(c, code.NullData, err)
			return
		}
		if valEnabled == 0 {
			conds = append(conds, profiletype.EnabledEQ(false))
		} else {
			conds = append(conds, profiletype.EnabledEQ(true))
		}
		enabledNil = &valEnabled
	} else {
		enabledNil = nil
	}

	datas, err := SerEntORMApp.ProfileTypesList(Page, PageSize, typeGroupId, enabledNil, c, conds, name, typeCode, tenantId)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 修改人员类型
func (that *FmClientApi) UpdateProfileType(ctx *gin.Context) {

	var (
		err   error
		PParm = seclients.ProfileType{}
	)

	err = ctx.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if PParm.ID == 0 {
		err = fmt.Errorf("%v:param:%v ", code.ZhCNText[code.ParamError], PParm)
		helpers.JSONs(ctx, code.ParamError, err)
		return
	}

	//执行
	datas, err := SerEntORMApp.ProfileTypeUpdate(ctx, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(ctx, code.NullData, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": datas, "message": code.ZhCNText[code.Success], "error": err})

}

// 新增人员类型信息
func (that *FmClientApi) NewProfileType(ctx *gin.Context) {

	var (
		err   error
		PParm = seclients.ParamProfileTypeLists{}
	)

	err = ctx.ShouldBindJSON(&PParm)
	if err != nil || len(PParm.NewTypeList) == 0 {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//单个更新
	var datas []int64
	if len(PParm.NewTypeList) == 1 {
		newTypeTenant := &PParm.NewTypeList[0]
		newTypeTenant.TenantID = PParm.TenantID
		data, err := SerEntORMApp.ProfileTypeNew(ctx, newTypeTenant)
		apilog.Printf("query data:%#v conds:%#v \n", data, PParm)
		if err != nil {
			helpers.JSONs(ctx, code.NullData, err)
			return
		}
		datas = append(datas, data)
		//关联表map中间表信息
		_, err = SerEntORMApp.GroupProfileTypeMapNew(ctx, &seclients.ProfileTypeAlertInfoMap{GroupID: PParm.ID,
			ProfileTypeID: data, TenantID: &PParm.TenantID})
		apilog.Printf("new prfile group map:%#v type id:%#v err:%#v \n", PParm.ID, data, err)

		if err != nil {
			apilog.Printf("new map type:%#v err:%#v \n", PParm.ID, err)
			helpers.JSONs(ctx, code.Success, gin.H{"data": datas, "message": code.ZhCNText[code.Success], "error": err})
			return
		}
		helpers.JSONs(ctx, code.Success, gin.H{"data": datas, "message": code.ZhCNText[code.Success], "error": err})
		return
	} else {
		//批量更新
		datas, err = SerEntORMApp.InsertProfileTypeList(ctx, &PParm)
		if err != nil {
			helpers.JSONs(ctx, code.NullData, err)
			return
		}
		helpers.JSONs(ctx, code.Success, gin.H{"data": datas, "message": code.ZhCNText[code.Success], "error": err})

	}

}
