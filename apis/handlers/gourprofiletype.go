package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/grouprofiletypemapping"
	"fmcam/systems/genclients/predicate"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (that *FmClientApi) GetGroupTypeMapping(c *gin.Context) {

	group_id, _ := strconv.Atoi(c.DefaultQuery("group_id", ""))
	profile_type_id, _ := strconv.Atoi(c.DefaultQuery("profile_type_id", ""))
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.GrouProfileTypeMapping{}

	if group_id != 0 {
		conds = append(conds, grouprofiletypemapping.GroupIDEQ(int64(group_id)))
	}
	if profile_type_id != 0 {
		conds = append(conds, grouprofiletypemapping.ProfileTypeIDEQ(int64(profile_type_id)))
	}

	datas, err := SerEntORMApp.GroupProfileTypeMapList(Page, PageSize, c, conds)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

func (that *FmClientApi) GetGroupTypeMappingInfo(c *gin.Context) {

	group_id, _ := strconv.Atoi(c.DefaultQuery("group_id", ""))
	profile_type_id, _ := strconv.Atoi(c.DefaultQuery("profile_type_id", ""))
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	datas, err := SerEntORMApp.GroupProfileTypeInfoList(Page, PageSize, group_id, profile_type_id, c)
	apilog.Printf("query data:%#v gpid:%#v type id:%v, err:%v \n", datas, group_id, profile_type_id, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

func (that *FmClientApi) NewGroupTypeMappingInfo(c *gin.Context) {

	var PParm = seclients.ProfileTypeAlertInfoMap{}
	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		helpers.JSONs(c, code.ParamBindError, gin.H{"error": code.ZhCNText[code.ParamBindError] + err.Error()})
		return
	}

	datas, err := SerEntORMApp.GroupProfileTypeMapNew(c, &PParm)
	apilog.Printf("query data:%#v  err:%v \n", datas, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}
