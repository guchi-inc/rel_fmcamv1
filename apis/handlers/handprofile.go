package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/models/sysaly"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmcam/systems/genclients/profiles"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 查询人员列表
func (that *FmClientApi) GetProfiles(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	roomId := c.DefaultQuery("room_id", "")
	typeId, _ := strconv.Atoi(c.DefaultQuery("type_id", ""))
	cardCode := c.DefaultQuery("id_card_number", "")
	phoneNumber := c.DefaultQuery("phone_number", "")
	tenantId := c.DefaultQuery("tenant_id", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	PParm := genclients.Profiles{}
	if tenantId != "" {
		PParm.TenantID = uuid.MustParse(tenantId)
	}

	var conds = []predicate.Profiles{}

	if name != "" {
		conds = append(conds, profiles.NameContains(name))
		PParm.Name = &name
	}
	if cardCode != "" {
		conds = append(conds, profiles.IDCardNumberEQ(cardCode))
		PParm.IDCardNumber = &cardCode
	}
	if typeId != 0 {
		newTypeId := int64(typeId)
		conds = append(conds, profiles.TypeIDEQ(newTypeId))
		PParm.TypeID = &newTypeId
	}
	if phoneNumber != "" {
		conds = append(conds, profiles.PhoneNumberContains(phoneNumber))
		PParm.PhoneNumber = &phoneNumber
	}
	if roomId != "" {
		conds = append(conds, profiles.RoomIDEQ(roomId))
		PParm.RoomID = &roomId
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apilog.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, profiles.CreatedTimeGT(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apilog.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, profiles.CreatedTimeLT(EndAtDate))
	}

	datas, err := SerEntORMApp.ProfilesList(Page, PageSize, c, conds, StartAt, EndAt, &PParm)
	apilog.Printf("query data err:%#v conds:%#v \n", err, conds)

	if err != nil {
		err = fmt.Errorf("查询失败:%v:%v ", code.ZhCNText[code.NullData], err)
		helpers.JSONs(c, code.NullData, gin.H{"message": code.ZhCNText[code.NullData], "error": err})
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 新增人员信息的 类型
func (that *FmClientApi) NewProfiles(ctx *gin.Context) {

	var (
		err   error
		PParm = seclients.ProfileInfos{}
	)

	err = ctx.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	//执行
	datas, err := SerEntORMApp.ProfileNew(ctx, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(ctx, code.NullData, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": datas, "message": code.ZhCNText[code.Success], "error": err})

}

// 用户修改人员信息的 类型
// 人员信息的类型修改
func (that *FmClientApi) UpdateProfiles(ctx *gin.Context) {

	var (
		err   error
		PParm = seclients.ProfileInfos{}
	)

	err = ctx.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if PParm.ID == 0 {
		err = fmt.Errorf("%v:param:%v ", code.ZhCNText[code.ParamError], PParm)
		helpers.JSONs(ctx, code.ParamError, err)
		return
	}

	//执行
	datas, err := SerEntORMApp.ProfileUpdate(ctx, &PParm)
	apilog.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(ctx, code.NullData, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": datas, "message": code.ZhCNText[code.Success], "error": err})

}

// 人员离店
func (that *FmClientApi) UpdateProfileCheckout(ctx *gin.Context) {

	var (
		Authed bool = false
		err    error
		dbUser *erps.FmUser
	)

	/*
		登陆解密用户传入数据， id_card_number  必填，否则返回失败
		然后鉴权

		校验token

			该登陆token正确时
				进入checkout流程
			token有，但是校验失败
				记录错误信息
				尝试校验账号信息

		校验账号密码信息
			1 校验账密成功
				进入checkout流程
			2  校验账密失败
				记录错误信息

	*/

	//解密并绑定登陆参数
	checkOutProfile, err := SysCheckOutIn.ProfileJsonParam(ctx)
	apilog.Printf("profile login:%#v err:%#v\n", checkOutProfile, err)
	if err != nil {
		err = fmt.Errorf("参数读取失败:%v:%v ", code.ZhCNText[code.ParamBindError], err)
		helpers.JSONs(ctx, code.ParamBindError, gin.H{"message": code.ZhCNText[code.ParamBindError], "error": err})
		return
	}

	//必填校验
	if checkOutProfile != nil && checkOutProfile.IdCardGuest == nil {
		helpers.JSONs(ctx, code.ParamError, gin.H{"message": code.ZhCNText[code.ParamError]})
		return
	}

	//token鉴权
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		dbUser, err = UserService.GetUserByCtxToken(ctx)
		apilog.Printf("receive params:%v", dbUser)
		if err != nil {
			err = fmt.Errorf("token错误:%v:%v ", code.ZhCNText[code.AuthorizationError], err)
			helpers.JSONs(ctx, code.AuthorizationError, gin.H{"message": code.ZhCNText[code.AuthorizationError], "error": err})
			return
		} else {
			//token 校验成功
			Authed = true
		}
	}

	//校验用户信息
	if !Authed {
		//接收的 参数解析失败
		if checkOutProfile == nil {
			helpers.JSONs(ctx, code.AuthorizedListError, gin.H{"message": code.ZhCNText[code.AuthorizedListError], "error": err})
			return
		}

		//如果没有通过token 则校验用户信息
		if checkOutProfile.LoginPassword == nil || (checkOutProfile.LoginName == nil && checkOutProfile.LoginPhonenum == nil) {
			err = fmt.Errorf("手机号或邮箱或账户密码错误:%v:%v ", code.ZhCNText[code.ParamError], err)
			helpers.JSONs(ctx, code.AuthorizedDetailError, gin.H{"message": code.ZhCNText[code.AuthorizedDetailError], "error": err})
			return
		}

		dbUser, err = UserService.AuthProfileInfo(*checkOutProfile.LoginPhonenum, *checkOutProfile.LoginName)
		apilog.Printf("user sign info:%#v, find:%#v client:%v\n", checkOutProfile, dbUser, ctx.ClientIP())
		if err != nil {
			err = fmt.Errorf("门店账号不存在:%v:%v ", code.ZhCNText[code.AuthorizedListError], err)
			helpers.JSONs(ctx, code.AuthorizedDetailError, gin.H{"message": code.ZhCNText[code.AuthorizedDetailError], "error": err})
			return
		}

		//获取用户信息
		//ComparePassword(dbUser.Password, user.Password)
		if isTrue := utils.BcryptCheck(*checkOutProfile.LoginPassword, dbUser.Password); !isTrue {
			apilog.Printf("user before:%v", dbUser.Id)
			err = fmt.Errorf("密码错误:%v:%v ", code.ZhCNText[code.AuthorizedPasswordError], err)
			helpers.JSONs(ctx, code.AuthorizedPasswordError, gin.H{"message": code.ZhCNText[code.AuthorizedPasswordError], "error": err})
			return
		} else {
			Authed = true
		}
	}

	//鉴权成功,更新或创建一个离店记录
	if Authed {
		rst, err := SysCheckOutIn.ProfileCheckout(ctx, dbUser, checkOutProfile)
		if err != nil {
			helpers.JSONs(ctx, code.AdminModifyDataError, gin.H{"message": code.ZhCNText[code.AdminModifyDataError], "error": err})
			return
		}
		helpers.JSONs(ctx, code.Success, gin.H{"data": rst, "message": "success"})
		return
	}
	//无法获取鉴权信息
	helpers.JSONs(ctx, code.AuthorizedDetailError, gin.H{"message": code.ZhCNText[code.AuthorizedDetailError], "error": err})

}

// 人员住店 checkin
func (that *FmClientApi) ProfileCheckin(ctx *gin.Context) {

	var (
		Authed bool = false
		err    error
		dbUser *erps.FmUser
	)

	/*
		登陆解密用户传入数据， id_card_number  必填，否则返回失败
		然后鉴权

		校验token

			该登陆token正确时
				进入checkout流程
			token有，但是校验失败
				记录错误信息
				尝试校验账号信息

		校验账号密码信息
			1 校验账密成功
				进入checkout流程
			2  校验账密失败
				记录错误信息
	*/

	//解密并绑定登陆参数
	checkInProfile, err := SysCheckOutIn.ProfileJsonParam(ctx)
	apilog.Printf("profile login:%#v err:%#v\n", checkInProfile, err)
	if err != nil {
		err = fmt.Errorf("参数读取失败:%v:%v ", code.ZhCNText[code.ParamBindError], err)
		helpers.JSONs(ctx, code.ParamBindError, gin.H{"message": code.ZhCNText[code.ParamBindError], "error": err})
		return
	}

	//必填校验
	if checkInProfile != nil && checkInProfile.IdCardGuest == nil {
		helpers.JSONs(ctx, code.ParamError, gin.H{"message": code.ZhCNText[code.ParamError]})
		return
	}

	//token或 x-apikey鉴权
	dbUser, err = UserService.GetUserByCtxToken(ctx)
	apilog.Printf("this dbUser:%v", dbUser)
	if err != nil {
		err = fmt.Errorf("token错误:%v:%v ", code.ZhCNText[code.AuthorizationError], err)
		helpers.JSONs(ctx, code.AuthorizationError, gin.H{"message": code.ZhCNText[code.AuthorizationError], "error": err})
		return
	} else {
		//token 校验成功
		Authed = true
	}

	//校验用户信息
	if !Authed {
		//接收的 参数解析失败
		if checkInProfile == nil {
			helpers.JSONs(ctx, code.AuthorizedListError, gin.H{"message": code.ZhCNText[code.AuthorizedListError], "error": err})
			return
		}

		//如果没有通过token 则校验用户信息
		if checkInProfile.LoginPassword == nil || (checkInProfile.LoginName == nil && checkInProfile.LoginPhonenum == nil) {
			err = fmt.Errorf("手机号或邮箱或账户密码错误:%v:%v ", code.ZhCNText[code.ParamError], err)
			helpers.JSONs(ctx, code.AuthorizedDetailError, gin.H{"message": code.ZhCNText[code.AuthorizedDetailError], "error": err})
			return
		}

		dbUser, err = UserService.AuthProfileInfo(*checkInProfile.LoginPhonenum, *checkInProfile.LoginName)
		apilog.Printf("user sign info:%#v, find:%#v client:%v\n", checkInProfile, dbUser, ctx.ClientIP())
		if err != nil {
			err = fmt.Errorf("门店账号不存在:%v:%v ", code.ZhCNText[code.AuthorizedListError], err)
			helpers.JSONs(ctx, code.AuthorizedDetailError, gin.H{"message": code.ZhCNText[code.AuthorizedDetailError], "error": err})
			return
		}

		//获取用户信息
		//ComparePassword(dbUser.Password, user.Password)
		if isTrue := utils.BcryptCheck(*checkInProfile.LoginPassword, dbUser.Password); !isTrue {
			apilog.Printf("user before:%v", dbUser.Id)
			err = fmt.Errorf("密码错误:%v:%v ", code.ZhCNText[code.AuthorizedPasswordError], err)
			helpers.JSONs(ctx, code.AuthorizedPasswordError, gin.H{"message": code.ZhCNText[code.AuthorizedPasswordError], "error": err})
			return
		} else {
			Authed = true
		}
	}

	//鉴权成功,更新或创建一个 入店 记录
	if Authed {
		rst, err := SysCheckOutIn.ProfileCheckin(ctx, dbUser, checkInProfile)
		if err != nil {
			helpers.JSONs(ctx, code.AdminModifyDataError, gin.H{"message": code.ZhCNText[code.AdminModifyDataError], "error": err})
			return
		}
		helpers.JSONs(ctx, code.Success, gin.H{"data": rst, "message": "success"})
		return
	}
	//无法获取鉴权信息
	helpers.JSONs(ctx, code.AuthorizedDetailError, gin.H{"message": code.ZhCNText[code.AuthorizedDetailError], "error": err})

}

// 绿云 人员状态检测
func (that *FmClientApi) YunLvRequest(c *gin.Context) {

	var (
		req = sysaly.YunLvPostData{}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		apilog.Printf("req bind:%#v \n", err)
		c.JSON(http.StatusBadRequest, gin.H{"rescode": 400, "resDesc": "Invalid request"})
		return
	}

	if req.IdNo == "" && req.MasterId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"rescode": 400, "resDesc": "idNo  or masterid request can not be empty"})
		return
	}

	//"rescode":"200","resDesc":"访问成功"
	err := SysCheckOutIn.YunLvPostDeal(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"rescode": 400, "resDesc": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rescode": 200, "resDesc": "OK", "serverTime": time.Now().UnixMilli()})
}
