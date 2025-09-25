package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/fmalertdefinition"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 预警组 查询预警组 信息
func (that *FmClientApi) GetAlertGroupsInfo(c *gin.Context) {

	ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	groupName := c.DefaultQuery("group_name", "")

	//查询预警组信息
	datas, err := SerEntORMApp.GroupAlertsInfo(c, ids, groupName)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})
}

// 预警组 查询预警组 列表
func (that *FmClientApi) GetAlertGroups(c *gin.Context) {

	group_name := c.DefaultQuery("group_name", "")
	ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	groupType, _ := strconv.Atoi(c.DefaultQuery("group_type", "0"))

	customization := c.DefaultQuery("customization", "")
	tenantId := c.DefaultQuery("tenant_id", "")

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

	datas, err := SerEntORMApp.GroupAlertsList(Page, PageSize, ids, groupType, c, group_name, customization, tenantId, StartAt, EndAt)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 预警组 新增信息
func (that *FmClientApi) AlertGroupsNew(c *gin.Context) {

	params := seclients.FmAlertGroup{}
	err := c.ShouldBind(&params)
	if err != nil {
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	if params.GroupName == "" {
		helpers.JSONs(c, code.ParamError, err)
		return
	}

	user, err := UserService.GetUserByCtxToken(c)
	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(c, code.AuthorizationError, err)
		return
	}

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//新增预警组信息
	params.Creator = user.LoginName
	datas, err := SerEntORMApp.NewGroupAlertsInfo(c, &params)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}

	var datasAll = []int64{datas}
	if params.GroupType == 1 {
		//组相关预警等级 新增
		var AlertGPLevels = seclients.ParamAlertList{ID: datas, NewLevelList: params.NewLevelList}
		datasLevel, err := SerEntORMApp.LevelAlertNewList(c, &AlertGPLevels)
		apilog.Printf("query data:%#v err:%#v \n", datasLevel, err)
		if err != nil {
			helpers.JSONs(c, code.NullData, err)
			return
		}
		datasAll = append(datasAll, datasLevel...)
	} else if params.GroupType == 2 {
		//人员类型组相关新增
		if len(params.NewProfileTypeList) > 0 {
			datasTypes, err := SerEntORMApp.InsertProfileTypeList(c, &seclients.ParamProfileTypeLists{
				ID:          datas,
				TenantID:    params.TenantID,
				NewTypeList: params.NewProfileTypeList})
			if err != nil {
				helpers.JSONs(c, code.NullData, err)
				return
			}
			datasAll = append(datasAll, datasTypes...)
		}
	} else {
		helpers.JSONs(c, code.NullData, fmt.Errorf("组类型只能是1预警或2人员类型组"))
		return
	}

	helpers.JSONs(c, code.Success, gin.H{"data": datasAll, "message": "success"})

}

// 预警组 修改预警组信息
func (that *FmClientApi) AlertGroupsUpdate(c *gin.Context) {

	params := seclients.FmAlertGroup{}
	err := c.ShouldBind(&params)
	if err != nil {
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	if params.ID == 0 {
		helpers.JSONs(c, code.ParamError, gin.H{"error": fmt.Errorf("%v:%v", params.ID, err)})
		return
	}

	//修改鉴权
	user, err := UserService.GetUserByCtxToken(c)
	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(c, code.AuthorizationError, gin.H{"error": err})
		return
	}

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//修改预警组信息
	params.Creator = user.LoginName
	datas, err := SerEntORMApp.GroupAlertsInfoUpdate(c, &params)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"error": err})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})
}

// 查询某个预警级别信息
func (that *FmClientApi) AlertLevelInfo(c *gin.Context) {
	ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	level, _ := strconv.Atoi(c.DefaultQuery("level", "0"))

	alert_group_id, _ := strconv.Atoi(c.DefaultQuery("alert_group_id", "0"))
	profile_type_id, _ := strconv.Atoi(c.DefaultQuery("profile_type_id", "0"))

	var conds = []predicate.FmAlertDefinition{}

	if ids != 0 {
		conds = append(conds, fmalertdefinition.IDEQ(int64(ids)))
	}
	if alert_group_id != 0 {
		conds = append(conds, fmalertdefinition.AlertGroupIDEQ(int64(alert_group_id)))
	}
	if profile_type_id != 0 {
		conds = append(conds, fmalertdefinition.ProfileTypeIDEQ(int64(profile_type_id)))
	}
	if level != 0 {
		conds = append(conds, fmalertdefinition.LevelEQ(level))
	}

	//查询预警组信息
	datas, err := SerEntORMApp.LevelAlertInfo(c, conds)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)

	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"error": err})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})
}

// 查询预警级别列表信息
func (that *FmClientApi) AlertLevelList(c *gin.Context) {
	ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	level, _ := strconv.Atoi(c.DefaultQuery("level", "0"))

	alert_group_id, _ := strconv.Atoi(c.DefaultQuery("alert_group_id", "0"))
	profile_type_id, _ := strconv.Atoi(c.DefaultQuery("profile_type_id", "0"))
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.FmAlertDefinition{}

	if ids != 0 {
		conds = append(conds, fmalertdefinition.IDEQ(int64(ids)))
	}
	if alert_group_id != 0 {
		conds = append(conds, fmalertdefinition.AlertGroupIDEQ(int64(alert_group_id)))
	}
	if profile_type_id != 0 {
		conds = append(conds, fmalertdefinition.ProfileTypeIDEQ(int64(profile_type_id)))
	}
	if level != 0 {
		conds = append(conds, fmalertdefinition.LevelEQ(level))
	}

	//查询预警组信息
	datas, err := SerEntORMApp.LevelAlertList(Page, PageSize, c, conds)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 新增预警级别信息
func (that *FmClientApi) AlertLevelNew(c *gin.Context) {

	params := seclients.ParamAlertList{} //genclients.FmAlertDefinition{}
	err := c.ShouldBind(&params)
	if err != nil || len(params.NewLevelList) == 0 {
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	//新增等级与 组一起操作
	if len(params.NewLevelList) >= 1 {
		helpers.JSONs(c, code.Success, gin.H{"data": len(params.NewLevelList)})
		return
	}
	var datas []int64
	if len(params.NewLevelList) == 1 {
		newLevel := params.NewLevelList[0]
		if newLevel.Level == 0 || newLevel.AlertGroupID == 0 {
			helpers.JSONs(c, code.ParamError, code.ZhCNText[code.ParamError])
			return
		}

		data, err := SerEntORMApp.LevelAlertNew(c, &newLevel)
		apilog.Printf("query data:%#v err:%#v \n", datas, err)
		if err != nil {
			helpers.JSONs(c, code.NullData, err)
			return
		}
		datas = append(datas, data)
		helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})

	}

	//批量新增
	datas, err = SerEntORMApp.LevelAlertNewList(c, &params)
	apilog.Printf("new data:%#v err:%#v \n", datas, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})
}

// 修改预警 级别信息
func (that *FmClientApi) AlertLevelUpdate(c *gin.Context) {
	params := seclients.ParamAlertList{} //genclients.FmAlertDefinition{}
	err := c.ShouldBind(&params)
	if err != nil {
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	if len(params.NewLevelList) == 0 {
		helpers.JSONs(c, code.ParamError, code.ZhCNText[code.ParamError])
		return
	}

	datas, err := SerEntORMApp.LevelAlertListUpdate(c, &params)
	apilog.Printf("query data:%#v err:%#v \n", datas, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})
}
