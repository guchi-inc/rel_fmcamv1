package resources

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/fmpmsapi"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// pms api 管理列表
func (that *TenantsApi) GetPMSApiList(c *gin.Context) {

	name := c.DefaultQuery("pms_name", "")
	contact := c.DefaultQuery("contact", "")
	phonenum := c.DefaultQuery("phonenum", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.FMPMSApi{fmpmsapi.DeleteFlagNEQ("1")}

	if name != "" {
		conds = append(conds, fmpmsapi.PmsNameContains(name))
	}
	if phonenum != "" {
		conds = append(conds, fmpmsapi.PhonenumContains(phonenum))
	}
	if contact != "" {
		conds = append(conds, fmpmsapi.ContactContains(contact))
	}

	datas, err := SerEntORMApp.FMPMSAPIList(Page, PageSize, c, conds)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 更新 pms api
func (that *TenantsApi) UpdatePMSApi(c *gin.Context) {

	var (
		PParm = seclients.FMPMSApi{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.FMPMSApiUpdate(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 新增pms api
func (that *TenantsApi) NewPMSApi(c *gin.Context) {

	var (
		PParm = seclients.FMPMSApi{}
	)

	err := c.ShouldBindJSON(&PParm)
	apilog.Printf("new pms:%#v api bindjson err:%#v \n", PParm, err)
	if err != nil {
		err = fmt.Errorf("%v:pmsapi:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.FMPMSApiNew(c, &PParm)
	apilog.Printf("query data:%#v \n PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
