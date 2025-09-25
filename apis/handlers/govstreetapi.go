package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/govstreet"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (that *GovAddrApi) GetStreets(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	codes := c.DefaultQuery("code", "")
	provinceCode := c.DefaultQuery("province_code", "")
	cityCode := c.DefaultQuery("city_code", "")
	areaCode := c.DefaultQuery("area_code", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.GovStreet{govstreet.DeleteFlagNEQ("1")}

	if name != "" {
		conds = append(conds, govstreet.NameContains(name))
	}
	if codes != "" {
		conds = append(conds, govstreet.CodeContains(codes))
	}
	if provinceCode != "" {
		conds = append(conds, govstreet.ProvinceCodeContains(provinceCode))
	}
	if cityCode != "" {
		conds = append(conds, govstreet.CityCodeContains(cityCode))
	}

	if areaCode != "" {
		conds = append(conds, govstreet.AreaCodeContains(areaCode))
	}

	datas, err := SerEntORMApp.StreetsList(Page, PageSize, c, conds)

	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}

// 更新街道乡镇
func (that *GovAddrApi) UpdateStreet(c *gin.Context) {

	var (
		PParm = genclients.GovStreet{}
		err   error
	)

	err = c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:GovStreet:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.StreetUpdate(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 新增街道
func (that *GovAddrApi) NewStreet(c *gin.Context) {

	var (
		PParm = genclients.GovStreet{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:street:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	apilog.Println("receive street param:", PParm)

	datas, err := SerEntORMApp.StreetNew(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
