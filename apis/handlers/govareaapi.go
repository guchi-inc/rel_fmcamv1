package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/govarea"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (that *GovAddrApi) GetAreas(c *gin.Context) {

	name := c.Query("name")
	codes := c.Query("code")
	provinceCode := c.DefaultQuery("province_code", "")
	cityCode := c.DefaultQuery("city_code", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.GovArea{govarea.DeleteFlagNEQ("1")}

	if name != "" {
		conds = append(conds, govarea.NameContains(name))
	}
	if codes != "" {
		conds = append(conds, govarea.CodeContains(codes))
	}
	if provinceCode != "" {
		conds = append(conds, govarea.ProvinceCodeContains(provinceCode))
	}
	if cityCode != "" {
		conds = append(conds, govarea.CityCodeContains(cityCode))
	}

	datas, err := SerEntORMApp.AreasList(Page, PageSize, c, conds)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}

// 更新区县
func (that *GovAddrApi) UpdateArea(c *gin.Context) {

	var (
		PParm = genclients.GovArea{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.AreaUpdate(c, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 新增区县
func (that *GovAddrApi) NewArea(c *gin.Context) {

	var (
		PParm = genclients.GovArea{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.AreaNew(c, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
