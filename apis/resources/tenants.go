package resources

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/predicate"
	"fmcam/systems/genclients/tenants"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 租户信息
func (that *TenantsApi) TenantsList(ctx *gin.Context) {

	//Page, PageSize  supplier, province, city, area, street, fullAddress, TenantId

	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	SupplyName := ctx.DefaultQuery("supplier", "")
	Telephone := ctx.DefaultQuery("telephone", "")
	Province := ctx.DefaultQuery("province", "")
	City := ctx.DefaultQuery("city", "")
	Area := ctx.DefaultQuery("area", "")
	Street := ctx.DefaultQuery("street", "")
	Address := ctx.DefaultQuery("address", "")
	TenantId := ctx.DefaultQuery("tenant_id", "")

	StartAt := ctx.DefaultQuery("start_time", "") //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")     //结束时间

	fmt.Printf("reveice TenantId:%v Name:%v Codes:%v TableName:%v, DataType:%v Area:%v, Street:%v,FullAddress:%v, PageSize:%v, Page:%v StartAt:%v EndAt:%v \n",
		TenantId, SupplyName, Telephone, Province, City, Area, Street, Address, PageSize, Page, StartAt, EndAt)

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 10 {
		PageSize = 10
	}

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	var conds = []predicate.Tenants{}

	//  PageSize, Page, StartAt, EndAt
	if TenantId != "" {
		uids, err := uuid.Parse(TenantId)
		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
		conds = append(conds, tenants.TenantID(uids))
	}

	if SupplyName != "" {
		conds = append(conds, tenants.SupplierContains(SupplyName))
	}
	if Telephone != "" {
		conds = append(conds, tenants.TelephoneContains(Telephone))
	}
	if Province != "" {
		conds = append(conds, tenants.ProvinceContains(Province))
	}
	if City != "" {
		conds = append(conds, tenants.CityContains(City))
	}
	if Area != "" {
		conds = append(conds, tenants.AreaContains(Area))
	}
	if Street != "" {
		conds = append(conds, tenants.StreetContains(Street))
	}
	if Address != "" {
		conds = append(conds, tenants.AddressContains(Address))
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apilog.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, tenants.CreatedTimeGT(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apilog.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, tenants.CreatedTimeLT(EndAtDate))
	}

	result, err := SerEntORMApp.TenantList(Page, PageSize, ctx, conds)
	apilog.Printf("ent conds:%#v to list:%#v err:%#v \n", conds, result.Data, err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	helpers.JSONs(ctx, code.Success, result)

}

// 更新租户信息
func (that *TenantsApi) UpdateTenant(c *gin.Context) {

	var (
		PParm = seclients.TenantFull{}
	)

	//鉴权
	user, err := UserService.GetUserByCtxToken(c)

	apilog.Printf("receive tenant edit by user:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(c, code.AuthorizationError, err)
		return
	}

	err = c.ShouldBindJSON(&PParm)
	apilog.Printf("receive PParm:%#v \n err:%#v \n", PParm, err)

	if err != nil {
		err = fmt.Errorf("%v:更新失败:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	PParm.Creator = user.LoginName

	successId, err := SerEntORMApp.TenantInfoUpdate(c, &PParm)
	if err != nil {
		err = fmt.Errorf("%v:更新失败:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": successId, "message": "success"})
}
