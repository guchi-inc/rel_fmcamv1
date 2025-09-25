package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmcam/systems/genclients/province"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (that *GovAddrApi) GetProvinces(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	codes := c.DefaultQuery("code", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.Province{province.DeleteFlagNEQ("1")}

	if name != "" {
		conds = append(conds, province.NameContains(name))
	}
	if codes != "" {
		conds = append(conds, province.CodeContains(codes))
	}

	datas, err := SerEntORMApp.ProvinceList(Page, PageSize, c, conds)

	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 省份更新
func (that *GovAddrApi) UpdateProvince(c *gin.Context) {

	var (
		PParm = genclients.Province{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.ProvenceUpdate(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 新增省份
func (that *GovAddrApi) NewProvince(c *gin.Context) {

	var (
		PParm = genclients.Province{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.ProvenceNew(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
