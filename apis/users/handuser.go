package users

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/fmuseraccount"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 根据手机号 查询账号列表
func (that *UsersApi) GetAccountsByPhone(ctx *gin.Context) {

	phoneNumber := ctx.DefaultQuery("phonenum", "")

	//账号列表 属于全局管理 可以使用租户号筛选
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if phoneNumber != "" {
		userList, err := UserService.UserAccountList(page, PageSize, "", phoneNumber)
		apidebug.Printf("profile retrieved:%#v phoneNumber:%#v err:%v\n", userList, phoneNumber, err)
		userList.Columns = nil
		helpers.JSONs(ctx, code.Success, userList)
		return
	}

	helpers.JSONs(ctx, code.Success, seclients.FmAccountsList{})

}

// 查询账号列表
func (that *UsersApi) GetAccountsList(ctx *gin.Context) {

	var (
		client = databases.EntClient
		Ctx    = configs.Ctx
	)

	name := ctx.DefaultQuery("username", "")
	LoginName := ctx.DefaultQuery("login_name", "")
	phoneNumber := ctx.DefaultQuery("phonenum", "")
	Email := ctx.DefaultQuery("email", "")         //邮箱
	Dep := ctx.DefaultQuery("department", "")      //部门
	Local := ctx.DefaultQuery("local", "")         //楼栋
	Localhost := ctx.DefaultQuery("localhost", "") //房号

	//账号列表 属于全局管理 可以使用租户号筛选
	TenantId := ctx.DefaultQuery("tenant_id", "")
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	StartAt := ctx.DefaultQuery("start_time", "") //创建开始时间
	EndAt := ctx.DefaultQuery("end_time", "")     //查询结束时间

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	offset := (Page - 1) * PageSize
	var conds = []predicate.FmUserAccount{}

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apidebug.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}

		startTime, err := timeutil.ParseCSTInLocation(StartAt)
		if err == nil {

			conds = append(conds, fmuseraccount.CreatedTimeGT(startTime))
		}
		endTime, err := timeutil.ParseCSTInLocation(EndAt)
		if err == nil {
			conds = append(conds, fmuseraccount.CreatedTimeLT(endTime))
		}
	}

	if name != "" {
		conds = append(conds, fmuseraccount.UsernameContains(name))
	}
	if Email != "" {
		conds = append(conds, fmuseraccount.Email(Email))
	}
	if Dep != "" {
		conds = append(conds, fmuseraccount.Department(Dep))
	}
	if Local != "" {
		conds = append(conds, fmuseraccount.Local(Local))
	}
	if Localhost != "" {
		conds = append(conds, fmuseraccount.Localhost(Localhost))
	}
	if LoginName != "" {
		conds = append(conds, fmuseraccount.LoginName(LoginName))
	}
	if phoneNumber != "" {
		conds = append(conds, fmuseraccount.Phonenum(phoneNumber))
	}

	//处理租户号筛选
	if TenantId != "" {
		userList, err := UserService.UserAccountList(Page, PageSize, TenantId, "")
		apidebug.Printf("profile retrieved:%#v tenant id:%#v err:%v\n", userList, TenantId, err)

		helpers.JSONs(ctx, code.Success, userList)
		return
	}

	var result seclients.FmAccountsList

	list, err := client.Debug().FmUserAccount.Query().
		Where(conds...).
		Limit(PageSize).
		Offset(offset).All(Ctx)

	apidebug.Printf("ent err:%v, name:%v typeId:%v, TenantId:%v conds:%#v \n", err, name, LoginName, TenantId, conds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	total, err := client.Debug().FmUserAccount.Query().Where(conds...).Count(Ctx)
	apidebug.Printf("ent err:%v, name:%v typeId:%v, TenantId:%v conds:%#v \n", err, name, LoginName, TenantId, conds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	apidebug.Println("profile retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("fm_user", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize

		helpers.JSONs(ctx, code.Success, result)
	} else {
		helpers.JSONs(ctx, code.Success, seclients.FmAccountsList{})
	}
}

// 更新某个 用户账户 信息
func (ap *UsersApi) UpdateUserAccount(ctx *gin.Context) {

	var (
		valueUser = &erps.FmUser{}
		err       error
	)

	err = ctx.ShouldBindJSON(valueUser)
	fmt.Printf("receive params rst:%#v, err:%v\n", valueUser, err)
	if valueUser == nil || err != nil {
		err = fmt.Errorf("request body datas:%v err:%v", valueUser, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if valueUser.Id == 0 {
		err = fmt.Errorf(":%v :%v", code.ZhCNText[code.ParamError], err)
		helpers.JSONs(ctx, code.ParamError, err)
		return
	}

	//鉴权
	user, err := UserService.GetUserByCtxToken(ctx)

	//只允许 管理员 根据 menu_number 更新 delete_flag, description, level
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	//执行更新
	var userid any
	var errExec error
	if user.Isystem != "1" {
		//非管理只能更新自己的
		if user.Ismanager == "2" {
			if valueUser.Id != user.Id {
				err = fmt.Errorf("%v:仅能修改本人信息", code.ZhCNText[code.AuthorizedPowerError])
				helpers.JSONs(ctx, code.AuthorizedPowerError, err)
				return
			}
		}
		//管理员可以更新自己租户下的账号
		userid, errExec = UserService.UpdateUser(valueUser, user.TenantId)
	} else {
		//查全部
		userid, errExec = UserService.UpdateUser(valueUser, "")
	}

	if errExec != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], errExec.Error())
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": userid})
}
