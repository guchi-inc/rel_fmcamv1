// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

// 刷新token
func (uc *UsersApi) TokenUser(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	tokenString := authHeader[len(configs.AuthHeader):]

	var token *jwt.Token
	token, err := utils.ValidateToken(tokenString)

	if err != nil {
		fmt.Println("token:", token, "token str", tokenString, err.Error())
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": code.AuthorizationError, "status": "failed", "message": "Non-Valid Token"})
		return
	}

	//校验token信息
	uid, iatInt, err := utils.GetInfoFromToken(token, "iat")
	if err != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return
	}

	// Logger.Infof("token reflush:%v,user_id:%v \n", token, user_id)
	total := int(configs.AccessExt) * 3600
	restTotal := int(time.Now().Unix() - iatInt) //已逝去的 多少秒,当前时间 减去 次token 发布时间
	expRedis := total - restTotal                //剩余多少秒

	//旧有token失效
	//token 与 账号 id关联 过期时间
	newKeys := fmt.Sprintf("%v:%v", configs.LogoutIatRedisTag, uid)
	newKeys += tokenString
	rst := databases.DBS.DBRedis.Client.SetEX(configs.Ctx, newKeys, time.Duration(expRedis)*time.Second, time.Duration(expRedis)*time.Second)
	rstRest, err := rst.Result()
	if err != nil {
		newErr := fmt.Errorf("rst:%v Set Token exp Err:%v", rstRest, code.ZhCNText[code.AuthorizedDetailError])
		helpers.JSONs(ctx, code.AuthorizedDetailError, newErr)
		return
	}

	newtokens := utils.GenerateToken(uint(uid))
	helpers.JSONs(ctx, http.StatusOK, gin.H{"data": newtokens})

}
