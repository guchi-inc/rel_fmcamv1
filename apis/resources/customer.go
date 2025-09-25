package resources

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TenantsApi struct{}

// 租户信息
func (that *TenantsApi) GetTenantsList(ctx *gin.Context) {

	var (
		err error
	)

	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	CustomName := ctx.DefaultQuery("supplier", "")
	Telephone := ctx.DefaultQuery("telephone", "")
	Contacts := ctx.DefaultQuery("contacts", "")
	FullAddress := ctx.DefaultQuery("full_address", "")
	TenantId := ctx.DefaultQuery("tenant_id", "")

	StartAt := ctx.DefaultQuery("start_time", "") //开始时间
	EndAt := ctx.DefaultQuery("end_time", "")     //结束时间

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}
	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
	}

	user, err := UserService.GetUserByCtxToken(ctx)

	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	apilog.Printf("db user password:%v, user password en:%v", user.Password, user.Password)

	//非管理员权限不可使用后台管理
	if user.Ismanager == "2" {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	if user.Ismanager != "2" || user.Isystem == "1" {
		datas, err := SysCustomesGroup.GetCustomerList(Page, PageSize, CustomName, Telephone, Contacts, FullAddress, StartAt, EndAt, TenantId)
		if err != nil {
			apilog.Printf("Get datas of depment Err:%v", err)
			err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], err)
			helpers.JSONs(ctx, code.AdminListError, err)
			return
		}
		helpers.JSONs(ctx, code.Success, datas)
	} else {
		//其他管理员无权限
		helpers.JSONs(ctx, code.AuthorizedPowerError, code.ZhCNText[code.AuthorizedPowerError])
		return
	}
}

// 新增租户 信息
func (that *TenantsApi) NewTenantsInfo(ctx *gin.Context) {

	var (
		vTenant = seclients.TenantFull{}
		err     error
	)

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)

	apilog.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v err:%v", code.ZhCNText[code.AuthorizationError], user, err)
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}

	err = ctx.ShouldBindJSON(&vTenant)
	fmt.Printf("receive params rst:%#v, err:%v", vTenant, err)

	//参数判空
	if vTenant.Supplier == "" || vTenant.Contacts == "" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		apilog.Printf("%v:%v", vTenant, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
		apilog.Printf("%v:%v", vTenant, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//管理员需要前端调用menu list 接口获取菜单功能列表
	datas, err := SysCustomesGroup.NewCustomer(user.LoginName, &vTenant)
	fmt.Printf("role update rst:%v, err:%v", datas, err)
	if err != nil {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
		apilog.Printf("%v:%v", vTenant, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
}

// // 更新 租户 信息
// func (that *TenantsApi) UpdateTenantsInfo(ctx *gin.Context) {

// 	var (
// 		vTenant = erps.FmTenant{}
// 		err     error
// 	)

// 	//鉴权
// 	user, err := UserService.GetUserByCtxToken(ctx)

// 	apilog.Printf("receive params:%v", user)
// 	if err != nil || user.LoginName == "" {
// 		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
// 		helpers.JSONs(ctx, code.AuthorizationError, err)
// 		return
// 	}

// 	err = ctx.ShouldBindJSON(&vTenant)
// 	fmt.Printf("receive params rst:%#v, err:%v", vTenant, err)

// 	//参数判空
// 	if vTenant.TenantId == "" || vTenant.Supplier == "" || vTenant.Address == "" {
// 		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
// 		apilog.Printf("%v:%v", vTenant, err)
// 		helpers.JSONs(ctx, code.ParamBindError, err)
// 		return
// 	}

// 	if err != nil {
// 		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.ParamBindError], err)
// 		apilog.Printf("%v:%v", vTenant, err)
// 		helpers.JSONs(ctx, code.ParamBindError, err)
// 		return
// 	}

// 	//管理员需要前端调用menu list 接口获取菜单功能列表
// 	datas, err := SysCustomesGroup.CustomerUpdate(user.LoginName, &vTenant)
// 	fmt.Printf("role update rst:%v, err:%v", datas, err)
// 	if err != nil {
// 		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AdminModifyPersonalInfoError], err)
// 		apilog.Printf("%v:%v", vTenant, err)
// 		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
// 		return
// 	}
// 	helpers.JSONs(ctx, code.Success, gin.H{"data": datas})
// }
