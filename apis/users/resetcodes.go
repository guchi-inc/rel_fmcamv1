package users

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmcam/models/erps"

	"fmt"
	"time"

	"fmcam/ctrlib/utils/resetcode"
	"fmcam/ctrlib/utils/timeutil"

	"github.com/gin-gonic/gin"
)

/*
r := gin.Default()
r.POST("/send-code", SendCodeHandler)
r.POST("/reset-password", ResetPasswordWithCodeHandler)
r.Run(":8080")
*/

func (that *UsersApi) EmailCodeSend(ctx *gin.Context) {

	var (
		DBRedis = databases.DBS.DBRedis.Client
	)

	email := ctx.DefaultQuery("email", "")
	phonenum := ctx.DefaultQuery("phonenum", "")
	regnew := ctx.DefaultQuery("regnew", "0") // 1为新注册

	if phonenum != "" && email == "" {
		if regnew == "0" {
			userinfo, err := UserService.GetUserValidInfo("", "", phonenum)
			if err != nil {
				helpers.JSONs(ctx, code.NullData, gin.H{"error": "该手机未注册:" + phonenum})
				return
			}

			email = userinfo.Email
		}

		if email == "" {
			email = phonenum
		}
	}

	storedTime := DBRedis.TTL(ctx, configs.DefResetKey+email)
	stilTime, err := storedTime.Result()
	apidebug.Printf("to email:%v, stilTime:%v err:%#v \n", configs.DefResetKey+":"+email, stilTime, err)
	if err == nil {
		fmt.Println("获取该键:", stilTime)
		intTime := stilTime.Seconds()
		fmt.Printf("剩余有效时间：%v 秒\n", intTime)
		if intTime > 300 {
			helpers.JSONs(ctx, code.ServerError, gin.H{"error": "获取验证码频率180秒一次：" + stilTime.String()})
			return
		}
	}

	if stilTime == -1 {
		fmt.Println("该键没有设置过期时间")
		helpers.JSONs(ctx, code.ServerError, gin.H{"error": "获取验证码频率180秒一次"})
		return
	}

	if stilTime == -2 {
		fmt.Println("该键不存在")
	}

	newcode := resetcode.GenerateCode()

	// 保存到 Redis，有效期 10 分钟
	err = DBRedis.Set(ctx, configs.DefResetKey+email, newcode, 10*time.Minute).Err()
	if err != nil {
		helpers.JSONs(ctx, code.ServerError, gin.H{"error": "保存验证码失败"})
		return
	}

	//TODO 发送手机验证码
	if phonenum != "" {

		err = resetcode.AliSendSMSCode(phonenum, newcode)
		if err != nil {
			helpers.JSONs(ctx, code.ServerError, gin.H{"error": "发送验证码失败" + err.Error()})
			return
		}
	} else {
		err = resetcode.SendApiVerifyCodeWithEmail(email, newcode)
		apidebug.Printf("to email:%v, newcode:%#v, err:%#v \n", email, newcode, err)
		if err != nil {
			helpers.JSONs(ctx, code.ServerError, gin.H{"error": "发送验证码失败" + err.Error()})
			return
		}
	}

	tn := time.Now().Local().In(configs.Loc)
	helpers.JSONs(ctx, code.Success, gin.H{"message": "验证码已发送", "data": timeutil.CSTLayoutString(&tn)})

}

// 跟验证码 更新密码
func (that *UsersApi) ResetPasswdWithEmail(ctx *gin.Context) {

	var (
		msgBack  = erps.MsgReback{}
		userinfo *erps.FmUser
		err      error
	)

	err = ctx.ShouldBind(&msgBack)
	apidebug.Printf("to err:%v, msgBack param:%#v,  \n", err, msgBack)

	var (
		DBRedis = databases.DBS.DBRedis.Client
	)

	if msgBack.Email == "" && msgBack.PhoneNum == "" {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("Email和PhoneNum%参数错误:%v", err))
		return
	}

	if msgBack.Email == "" {
		userinfo, err = UserService.GetUserValidInfo("", "", msgBack.PhoneNum)
		if err != nil {
			helpers.JSONs(ctx, code.NullData, gin.H{"error": "该手机未注册:" + msgBack.PhoneNum})
			return
		}
		msgBack.Email = userinfo.Email
	} else {
		//校验 id 并 限制只能更新自己的信息
		userinfo, err = UserService.GetUserValidInfo("", msgBack.Email, "")
		apidebug.Printf("user info:%v err:%v \n", userinfo, err)
		if err != nil || userinfo == nil {
			helpers.JSONs(ctx, code.NullData, fmt.Errorf("该用户不存在"))
			return
		}
	}

	// 直接从 Redis 获取验证码
	storedCode, err := DBRedis.Get(ctx, configs.DefResetKey+msgBack.Email).Result()
	apidebug.Printf("to email:%v, storedCode:%#v, err:%#v \n", msgBack, storedCode, err)

	if err != nil || storedCode != msgBack.Codes {
		helpers.JSONs(ctx, code.ParamError, fmt.Errorf("验证码无效或已过期"))
		return
	}

	// 更新数据库中用户的密码（密码加密处理）
	//密码加密存储, id不变
	if userinfo.Password != "" {
		//密码一致无需更改
		newPassSame := utils.BcryptCheck(msgBack.Password, userinfo.Password)
		apidebug.Printf("compare password:%v, %v, result:%v", msgBack.Password, userinfo.Password, newPassSame)
		if newPassSame {
			err := fmt.Errorf("新旧密码一致")
			helpers.JSONs(ctx, code.ParamError, err)
			// 删除验证码
			DBRedis.Del(ctx, configs.DefResetKey+msgBack.Email)
			return
		}
		userinfo.Password = utils.HashPassword(&msgBack.Password)
	}

	//执行更新
	userid, err := UserService.BussUpdateUser(userinfo)
	apidebug.Printf("to BussUpdateUser userid:%v, userinfo:%#v, err:%#v \n", userid, userinfo, err)

	if err != nil {
		err := fmt.Errorf("登陆密码修改为敏感操作，失败请找管理员协助 %v:您的访问地址:%v:  ", code.ZhCNText[code.AdminModifyPersonalInfoError], ctx.ClientIP())
		helpers.JSONs(ctx, code.AdminModifyPersonalInfoError, err)
		return
	}

	// 删除验证码
	DBRedis.Del(ctx, configs.DefResetKey+msgBack.Email)
	helpers.JSONs(ctx, code.Success, gin.H{"data": userid})

}
