package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/govcity"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (that *GovAddrApi) GetCitys(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	codes := c.DefaultQuery("code", "")
	provinceCode := c.DefaultQuery("province_code", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.GovCity{govcity.DeleteFlagNEQ("1")}

	if name != "" {
		conds = append(conds, govcity.NameContains(name))
	}
	if codes != "" {
		conds = append(conds, govcity.CodeContains(codes))
	}
	if provinceCode != "" {
		conds = append(conds, govcity.ProvinceCodeContains(provinceCode))
	}

	datas, err := SerEntORMApp.CityList(Page, PageSize, c, conds)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 更新城市
func (that *GovAddrApi) UpdateCity(c *gin.Context) {

	var (
		PParm = genclients.GovCity{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.CityUpdate(c, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 新增城市
func (that *GovAddrApi) NewCity(c *gin.Context) {

	var (
		PParm = genclients.GovCity{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.CityNew(c, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
