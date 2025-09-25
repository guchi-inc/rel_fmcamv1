package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/fmdedicatedservices"
	"fmcam/systems/genclients/fmdemands"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 新增需求留言
func (that *GovAddrApi) NewDemands(c *gin.Context) {

	var (
		PParm = genclients.FmDemands{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.DemandsNew(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 需求留言列表
func (that *GovAddrApi) GetDemands(c *gin.Context) {

	username := c.DefaultQuery("username", "") //己方客服姓名
	supplier := c.DefaultQuery("supplier", "") //服务公司名
	phonenum := c.DefaultQuery("phonenum", "") //己方客服电话
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.FmDemands{}

	if username != "" {
		conds = append(conds, fmdemands.UsernameContains(username))
	}
	if phonenum != "" {
		conds = append(conds, fmdemands.PhonenumEQ(phonenum))
	}
	if supplier != "" {
		conds = append(conds, fmdemands.SupplierEQ(supplier))
	}

	datas, err := SerEntORMApp.Demandsist(Page, PageSize, c, conds)

	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 查询客服人员
func (that *GovAddrApi) GetDedicatedService(c *gin.Context) {

	contacts := c.DefaultQuery("contacts", "")  //服务客服姓名
	supplier := c.DefaultQuery("supplier", "")  //酒店公司名
	phonenum := c.DefaultQuery("phonenum", "")  //服务客服电话
	TenantId := c.DefaultQuery("tenant_id", "") //酒店 号

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.FmDedicatedServices{}

	if contacts != "" {
		conds = append(conds, fmdedicatedservices.ContactsEQ(contacts))
	}
	if phonenum != "" {
		conds = append(conds, fmdedicatedservices.PhonenumEQ(phonenum))
	}
	if supplier != "" {
		conds = append(conds, fmdedicatedservices.SupplierEQ(supplier))
	}

	datas, err := SerEntORMApp.DedicateServiceList(Page, PageSize, c, conds, TenantId)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 客服人员更新
func (that *GovAddrApi) UpdateDedicatedService(c *gin.Context) {

	var (
		PParm = seclients.DedicateServiceInfo{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.DedicateUpdate(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}

// 新增客服人员信息
func (that *GovAddrApi) NewDedicatedService(c *gin.Context) {

	var (
		PParm = genclients.FmDedicatedServices{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.DedicateNew(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
