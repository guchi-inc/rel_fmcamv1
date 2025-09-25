// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"encoding/json"
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	gqs "github.com/doug-martin/goqu/v9"
)

type UsersApi struct{}

// 心跳接口
func (that *UsersApi) Pok(ctx *gin.Context) {

	// 服务鉴定
	helpers.JSONs(ctx, http.StatusOK, gin.H{"message": "Success", "data": 1, "status": "OK.", "code": code.Success})
}

func (that *UsersApi) Ping(ctx *gin.Context) {

	// 由token 获取用户信息
	gu, err := UserUtil.GetUserByCtx(ctx)
	if err != nil {
		helpers.JSONs(ctx, http.StatusOK, gin.H{"message": fmt.Sprintf("user does not exit %#v", err), "status": "ERROR.", "code": 401})
		return
	}
	helpers.JSONs(ctx, http.StatusOK, gin.H{"message": "Success", "data": gu, "status": "OK.", "code": 200})
}

// 返回users信息  当前用户租户号 之下的全部
func (that *UsersApi) GetUsersList(ctx *gin.Context) {

	userInfo, err := UserService.GetUserByCtxToken(ctx)
	apidebug.Printf("reveice userlogin:%#v,errt:%v \n", userInfo, err)

	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return
	}
	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "5"))
	Name := ctx.DefaultQuery("username", "")        //username 用户名
	LoginName := ctx.DefaultQuery("login_name", "") //登录账号
	Email := ctx.DefaultQuery("email", "")          //邮箱
	PhoneNum := ctx.DefaultQuery("phonenum", "")    //电话号
	Dep := ctx.DefaultQuery("department", "")       //部门
	Local := ctx.DefaultQuery("local", "")          //楼栋
	Localhost := ctx.DefaultQuery("localhost", "")  //房号
	StartAt := ctx.DefaultQuery("start_time", "")   //入住开始时间
	EndAt := ctx.DefaultQuery("end_time", "")       //入住结束时间

	apidebug.Printf("reveice param:%v, login name:%v, email:%v, phonenum:%v, Dep:%v Local:%v Localhost:%v \n", Name, LoginName, Email, PhoneNum, Dep, Local, Localhost)
	var gex = gqs.Ex{}
	var Filter string = " WHERE delete_flag != '1' "
	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apidebug.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}

		gex["created_time"] = gqs.Op{"between": gqs.Range(StartAt, EndAt)}
		Filter += fmt.Sprintf(" AND (created_time between '%v' AND '%v') ", StartAt, EndAt)
	}
	// var userGc = erps.GcUser{}
	if Email != "" && strings.Count(Email, "@") == 1 {
		gex["email"] = Email
		Filter += fmt.Sprintf(" AND email = '%v' ", Email)
	}

	if len(Name) >= 1 && len(Name) <= 100 {
		gex["username"] = gqs.Op{"like": fmt.Sprintf("%v", Name+"%")}
		Filter += fmt.Sprintf(" AND username like '%v' ", Name+"%")

	}

	if len(LoginName) >= 2 && len(LoginName) <= 100 {
		fmt.Println("LoginName:", len(LoginName))
		gex["login_name"] = gqs.Op{"like": fmt.Sprintf("%v", LoginName+"%")}
		Filter += fmt.Sprintf(" AND login_name = '%v' ", LoginName)
	}

	if len(PhoneNum) >= 2 && len(PhoneNum) <= 20 {
		gex["phonenum"] = gqs.Op{"like": fmt.Sprintf("%v", PhoneNum+"%")}
		Filter += fmt.Sprintf(" AND phonenum = '%v' ", PhoneNum)
	}

	if Dep != "" {
		gex["department"] = gqs.Op{"like": fmt.Sprintf("%v", Dep+"%")}
		Filter += fmt.Sprintf(" AND department = '%v' ", Dep)
	}

	if Local != "" {
		gex["local"] = gqs.Op{"like": fmt.Sprintf("%v", Local+"%")}
		Filter += fmt.Sprintf(" AND local = '%v' ", Local)
	}

	if Localhost != "" {
		gex["localhost"] = gqs.Op{"like": fmt.Sprintf("%v", "%"+Localhost+"%")}
		Filter += fmt.Sprintf(" AND localhost like '%v' ", Localhost+"%")
	}

	gex["id"] = gqs.Op{"Neq": 120} //排除admin
	apidebug.Printf("search user send gex:%v\n", gex)

	//只查有效的 当前账号 租户信息相关的用户账号
	tenantInfo, err := CustomerUtils.GetCustomerInfo(0, "", "", userInfo.TenantId)
	if err == nil || tenantInfo != nil {
		apidebug.Printf("tenant info:%#v , err:%#v \n", tenantInfo, err)
		Filter += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id) = '%v' ", uuid.MustParse(tenantInfo.TenantId))
	}

	datas, err := UserService.UserList(Page, PageSize, gex, Filter)
	if err != nil {
		logger.Infof("%v Err:%v", code.ZhCNText[code.AdminListError], err)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}

	helpers.JSONs(ctx, http.StatusOK, datas)
}

// 获取用户信息 包括租户信息
func (that *UsersApi) GetAccountInfo(ctx *gin.Context) {

	users, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	LoginName := users.LoginName
	logger.Infof("user sign info:%#v, default en:%v\n", LoginName)
	accountInfo, tenantInfo, err := UserService.GetUserAccountInfo(LoginName)
	if err != nil {
		helpers.JSONs(ctx, code.AdminDetailError, code.ZhCNText[code.AdminDetailError]+err.Error())
		return
	}

	//状态判断.1：正常，2：删除，3：封禁
	if !accountInfo.Enabled {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminListError], accountInfo.LoginName)
		helpers.JSONs(ctx, code.AdminListError, err)
		return
	}
	apidebug.Printf("db user password:%v, user password en:%v\n", accountInfo.Password, LoginName)

	if tenantInfo == nil {
		helpers.JSONs(ctx, code.Success, gin.H{"data": seclients.FmAccountsInfo{Account: accountInfo, TenantInfo: nil}})
		return
	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": seclients.FmAccountsInfo{Account: accountInfo, TenantInfo: tenantInfo}})
}

// 用户json登录 支持body 加密
func (uc *UsersApi) SignInUser(ctx *gin.Context) {

	/*
		登陆时如果该用户token错误，返回失败
		该用户token正确时
			1  如果为管理员，获取管理员菜单列表
				如果管理员没有任何权限，返回空菜单列表
				获取权限分组信息并获得访问控制列表
			2  如果为普通用户，返回成功
	*/

	//解密并绑定登陆参数
	users, err := uc.signJsonParam(ctx)
	if err != nil {
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if users == nil || (users.LoginName == "" && users.Phonenum == "") {

		paramInfo := fmt.Errorf("receive user:%v ", users)
		err := fmt.Errorf("%v:%v", code.ZhCNText[code.ParamBindError], paramInfo)
		logger.Infow("%v", err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	dbUser, err := UserService.GetUserValidInfo(users.LoginName, "", users.Phonenum)
	logger.Infof("user sign info:%#v, find:%#v client:%v\n", users, dbUser, ctx.ClientIP())
	if err != nil {
		helpers.JSONs(ctx, code.AdminDetailError, code.ZhCNText[code.AdminDetailError]+err.Error())
		return
	}

	//状态判断.1：正常，0 无效
	if dbUser.Enabled != nil && !*dbUser.Enabled {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminLoginError], dbUser.LoginName)
		helpers.JSONs(ctx, code.AdminLoginError, err)
		return
	}

	// 以付费账号为准 是否超时到期
	if dbUser.LeaderId != "" && dbUser.LeaderId != "0" {
		leaderId, _ := strconv.Atoi(dbUser.LeaderId)
		leaderInfo, err := UserService.GetUserById(leaderId)
		if err != nil {
			err = fmt.Errorf("%v:%v", code.ZhCNText[code.NullData], err)
			helpers.JSONs(ctx, code.NullData, err)
			return
		}
		dbUser.DeletedAt = leaderInfo.DeletedAt
	}

	//账户过期超时 无法登陆
	loginNow := time.Now().Local().In(configs.Loc)
	loginStr := timeutil.CSTLayoutString(&loginNow)
	_, _, err = timeutil.CheckDateRange(loginStr, dbUser.DeletedAt)
	apidebug.Printf("db user password:%v, user password en:%v login delete date:%v\n", dbUser.Password, users.Password, err)
	if err != nil {
		helpers.JSONs(ctx, code.RBACError, code.ZhCNText[code.RBACError]+"此账户已经过期:"+dbUser.DeletedAt)
		return
	}

	//判定和注册位置前端, 当用户退出时清理该信息, 如果是后台登陆，则不判断
	authPosition := ctx.GetHeader("local")
	if authPosition != "" {
		clientIp := ctx.ClientIP()
		rst, err := UserService.LocalHSET(users.LoginName, clientIp, authPosition)
		apidebug.Printf("\nuser:%#v login local code is:%#v, client ip:%v redis cache hset:%#v err:%v\n", dbUser, authPosition, clientIp, rst, err)
		if err != nil {
			helpers.JSONs(ctx, code.CacheSetError, err)
			return
		}
	}

	//组装返回信息
	token := utils.GenerateToken(uint(dbUser.Id))
	LoginData := gin.H{"data": token,
		"is_manager":  fmt.Sprintf("%v", dbUser.Ismanager),
		"is_system":   fmt.Sprintf("%v", dbUser.Isystem),
		"code":        code.Success,
		"menu":        false,
		"login_name":  dbUser.LoginName,
		"username":    dbUser.Username,
		"tenant_id":   dbUser.TenantId,
		"device_time": dbUser.DeviceTime,
		"phonenum":    dbUser.Phonenum,
		"id":          dbUser.Id,
		"message":     "success"}

	//添加管理员权限菜单内容,需要前端调用menu list 接口获取菜单功能列表
	if dbUser.Ismanager != "2" {
		LoginData["menu"] = true
	}

	//ComparePassword(dbUser.Password, user.Password)
	if isTrue := utils.BcryptCheck(users.Password, dbUser.Password); isTrue {
		logger.Infof("users before:%v", dbUser.Id)
		apidebug.Printf("db user password:%v, user password en:%v isTrue:%v \n", dbUser.Password, users.Password, isTrue)

		//更新登陆时间
		tn := time.Now().Local().In(configs.Loc)
		dbUser.DeviceTime = timeutil.CSTLayoutString(&tn)
		_, err = UserService.UpdateUserDeviceTime(dbUser, "")
		apidebug.Printf("\nuser:%#v  client ip:%v   err:%v\n", dbUser, ctx.ClientIP(), err)

		ctx.JSON(http.StatusOK, LoginData)
		return
	}

	helpers.JSONs(ctx, code.AuthorizedPasswordError, gin.H{"message": code.ZhCNText[code.AuthorizedPasswordError]})

}

// 仅用于登陆 解码
func (uc *UsersApi) signJsonParam(ctx *gin.Context) (*erps.FmUser, error) {

	var (
		user erps.FmUser
	)

	//接收json参数
	var body []byte
	body, err := io.ReadAll(ctx.Request.Body)
	apidebug.Printf("before dencry byte:%v", string(body))
	if err != nil {
		logger.Errorf("ioutil.ReadAll:%v failed:%v", body, err)
		return nil, err
	}

	//解密
	if configs.PostEncrypt {
		// body = utils.XDencryptByte(body, []byte(configs.AccessKey))
		strBody := utils.XDencryptString(string(body), configs.AccessKey)
		if strBody == "" {
			apidebug.Printf("dencry str:%v strBody:%v", string(body), string(strBody))
			err = fmt.Errorf("after dencry:%v map error:%v", string(strBody), err)
			return nil, err
		}
		body = []byte(strBody)

		apidebug.Println("dencry byte str:", string(body))
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		logger.Errorf("after dencry byte:%v to map error:%v", string(body), err)
		return nil, err
	}

	if user.LoginName != "" || user.Phonenum != "" {
		return &user, nil
	}

	//无需解密 解密 Json
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Printf("should bind json to user:%v err:%v\n", user, err)

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "check ContentType should be application/json", "error": err.Error()})
		return nil, err
	}
	logger.Infof("receive json params:%v", user)

	return &user, nil

}

// 用户登出
func (uc *UsersApi) LogOut(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) < len(configs.AuthHeader) {
		helpers.JSONs(ctx, code.AuthorizedDetailError, fmt.Errorf("鉴权错误:%v", authHeader))
		return
	}
	tokenString := authHeader[len(configs.AuthHeader):]

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizedDetailError, err)
		return
	}

	//校验token信息
	uid, iatInt, _ := utils.GetInfoFromToken(token, "iat")

	//校验token信息

	//token解析失败，或者发行时间大于当前时间，返回报错
	if err != nil || int64(iatInt) > time.Now().Unix() {
		//如果token已经失效，不记录到redis缓存处理
		logger.Infof("token:%v err:%#v, info:%#v iatInt:%v timeStuff:%v", uid, err, iatInt, int64(iatInt), time.Now().Unix())
		ferr := fmt.Errorf("%v", code.ZhCNText[code.AuthorizedDetailError])
		helpers.JSONs(ctx, code.AuthorizedDetailError, ferr)
		return
	} else {
		//token 还有效，在缓存中处理，
		//其他接口的token判定，需要先查看是否存在缓存中。
		var newErr error

		//设置边界，大于 该 timStuff的token才被允许进行校验。
		total := int(configs.AccessExt) * 3600
		restTotal := int(time.Now().Unix() - iatInt) //已逝去的 多少秒,当前时间 减去 次token 发布时间

		expRedis := total - restTotal //剩余多少秒
		fmt.Printf("total:%v, used Total:%v remain expRedis:%v \n", total, restTotal, expRedis)
		newKeys := fmt.Sprintf("%v:%v", configs.LogoutIatRedisTag, uid)
		rst := databases.DBS.DBRedis.Client.SetEX(configs.Ctx, newKeys, time.Now().Unix(), time.Duration(expRedis)*time.Second)
		rstRest, err := rst.Result()
		if err != nil {
			newErr = fmt.Errorf("rst:%v Set Token Err:%v", rstRest, code.ZhCNText[code.AuthorizedDetailError])
			helpers.JSONs(ctx, code.AuthorizedDetailError, newErr)
			return
		}

		//token 与 账号 id关联 过期时间
		newKeys += tokenString
		rst = databases.DBS.DBRedis.Client.SetEX(configs.Ctx, newKeys, time.Duration(expRedis)*time.Second, time.Duration(expRedis)*time.Second)
		rstRest, err = rst.Result()
		if err != nil {
			newErr = fmt.Errorf("rst:%v Set Token exp Err:%v", rstRest, code.ZhCNText[code.AuthorizedDetailError])
			helpers.JSONs(ctx, code.AuthorizedDetailError, newErr)
			return
		}

	}
	helpers.JSONs(ctx, code.Success, gin.H{"data": err})
}

// 更新某个 用户 信息
func (ap *UsersApi) UpdateUser(ctx *gin.Context) {

	var (
		valueStructs = &erps.FmUser{}
		err          error
	)

	err = ctx.ShouldBindJSON(valueStructs)
	fmt.Printf("receive params rst:%#v, err:%v\n", valueStructs, err)
	if valueStructs == nil || err != nil {
		err = fmt.Errorf("request body datas:%v err:%v", valueStructs, err)
		helpers.JSONs(ctx, code.ParamBindError, err)
		return
	}

	if valueStructs.Id == 0 {
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

	//当前管理员需要前端调用 menu list 接口获取菜单功能列表
	//根据用户属性判定，如果仅仅是管理员 不是系统分配账号 可以修改自己 租户内的
	Name := user.LoginName
	if user.Ismanager == "2" && user.Isystem != "1" {
		err = fmt.Errorf("%v:%v ", code.ZhCNText[code.AuthorizedPowerError], Name)
		helpers.JSONs(ctx, code.AuthorizedPowerError, err)
		return
	}

	//执行更新
	var userid any
	var errExec error
	if user.Isystem != "1" {
		//仅自己租户下的账号
		userid, errExec = UserService.UpdateUser(valueStructs, user.TenantId)
	} else {
		//查全部
		userid, errExec = UserService.UpdateUser(valueStructs, "")
	}

	if errExec != nil {
		err = fmt.Errorf("%v:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], errExec.Error())
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": userid})
}

// 用户修改自己的密码
func (BA *UsersApi) UpdatePassword(ctx *gin.Context) {

	var userinfo erps.FmUser
	err := ctx.ShouldBindJSON(&userinfo)
	apidebug.Printf("user update :%v err:%#v   \n", userinfo, err)

	if err != nil { //|| userinfo.Id == 0
		err := fmt.Errorf("信息不可为空 %v:%v", code.ZhCNText[code.ParamError], userinfo.LocalHost)
		helpers.JSONs(ctx, code.ParamError, err)
		return
	}

	//根据login_name查id
	var userPass = &erps.FmUser{}
	if userinfo.LoginName != "" && userinfo.Id == 0 {
		userPass, err = UserService.GetUserInfo(userinfo.LoginName)
		apidebug.Printf("userPass info by loginname:%v err:%v \n", userPass, err)

		if err != nil || userPass == nil {
			helpers.JSONs(ctx, code.NullData, fmt.Errorf("查询%v错误%v", userinfo.LoginName, err))
			return
		}
	} else {
		if userinfo.Id != 0 {
			userPass, err = UserService.GetUserById(int(userinfo.Id))
			apidebug.Printf("userPass info by id:%v err:%v \n", userPass, err)

			if err != nil || userPass == nil {
				helpers.JSONs(ctx, code.NullData, fmt.Errorf("查询%v错误%v", userinfo.Id, err))
				return
			}
		} else {
			helpers.JSONs(ctx, code.NullData, fmt.Errorf("id查询%v错误%v", userinfo.Id, err))
			return
		}
	}

	//校验 id 并 限制只能更新自己的信息
	user, err := UserService.GetUserByCtxToken(ctx)
	apidebug.Printf("user info:%v err:%v \n", user, err)
	if err != nil || user == nil {
		helpers.JSONs(ctx, code.AuthorizationError, fmt.Errorf("该次调用鉴权错误%v", err))
		return
	}

	if int(user.Id) != int(userPass.Id) {
		err := fmt.Errorf("只能修改自己的信息 %v:%v", code.ZhCNText[code.ParamError], user.Id == userPass.Id)
		helpers.JSONs(ctx, code.ParamError, err)
		return
	}

	//密码加密存储, id不变
	userinfo.Id = userPass.Id
	if userinfo.Password != "" {
		//密码一致无需更改
		newPassSame := utils.BcryptCheck(userinfo.Password, userPass.Password)
		apidebug.Printf("compare password:%v, %v, result:%v", userPass.Password, userinfo.Password, newPassSame)
		if newPassSame {
			err := fmt.Errorf("新旧密码一致")
			helpers.JSONs(ctx, code.ParamError, err)
			return
		}
		userinfo.Password = utils.HashPassword(&userinfo.Password)

	}

	//执行更新
	userid, err := UserService.BussUpdateUser(&userinfo)
	if err != nil {
		err := fmt.Errorf("登陆密码修改为敏感操作，失败请找管理员协助 %v:%v: err:%v", code.ZhCNText[code.AdminModifyPersonalInfoError], userinfo.LocalHost, err)
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	helpers.JSONs(ctx, code.Success, gin.H{"data": userid})
}

// 客户端用户 刷新token 对比新旧密码 ，更新新密码
func (uc *UsersApi) BussTokenUser(ctx *gin.Context) {

	var (
		user erps.FmUser
	)
	apidebug.Println("user init:", user)

	// 读取数据
	var body []byte
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		helpers.JSONs(ctx, code.ParamBindError, fmt.Errorf("read body err:%v", err))
		return
	}
	//解密
	if configs.PostEncrypt {
		// body = utils.XDencryptByte(body, []byte(configs.AccessKey))
		strBody := utils.XDencryptString(string(body), configs.AccessKey)
		if strBody == "" {
			apidebug.Printf("dencry str:%v error:%v", string(body), string(strBody))
			return
		}
		body := []byte(strBody)

		apidebug.Println("dencry byte str:", string(body))
		err = json.Unmarshal(body, &user)
		if err != nil {
			err := fmt.Errorf("after dencry byte:%v to map error:%v", string(body), err)
			helpers.JSONs(ctx, http.StatusBadRequest, err)
			return
		}
	}

	apidebug.Printf("user unmarshal:%#v", user)

	if user.LoginName == "" {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("%v", code.ZhCNText[code.ParamError]))
		return
	}

	dbUser, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, err)
		return

	}

	//比对提供的旧密码
	if isTrue := utils.BcryptCheck(user.Password, dbUser.Password); !isTrue {
		logger.Infof("user before:%v", dbUser.Id)
		helpers.JSONs(ctx, code.AuthorizedPasswordError, fmt.Errorf("旧%v%v", code.ZhCNText[code.AuthorizedPasswordError], isTrue))
		return
	}

	//校验现有token
	authHeader := ctx.GetHeader("Authorization")
	tokenString := authHeader[len(configs.AuthHeader):]
	var token *jwt.Token
	token, err = utils.ValidateToken(tokenString)

	if err != nil {
		fmt.Println("token:", token, "token str", tokenString, err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": code.AuthorizationError, "status": "failed", "message": "Non-Valid Token"})
		return
	}

	uid, err := utils.GetIdFromToken(token)
	if err != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return
	}

	//客户端新的token
	tokens := utils.GenerateToken(uint(uid))
	logger.Infof("user new tokens:%v user id:%v", tokens, uid)

	helpers.JSONs(ctx, http.StatusOK, gin.H{"data": tokens})

}
