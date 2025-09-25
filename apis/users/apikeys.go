package users

import (
	"crypto/rand"
	"encoding/hex"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ===== 路由处理 =====
// ===== 工具函数 =====

func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashAPIKey(apiKey string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// 生成 API Key
func (that *UsersApi) GenerateAPIKeyHandler(c *gin.Context) {

	var (
		DBConn = databases.DBMysql
	)
	var req struct {
		LoginName   string `json:"login_name" db:"login_name"`
		PhoneNumber string `json:"phonenum" db:"phonenum"`
		Days        int    `json:"days" db:"days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Days <= 0 {
		apidebug.Printf("param bind:%#v days:%#v \n", err, req.Days)
		helpers.JSONs(c, code.ParamBindError, gin.H{"error": "bad request " + err.Error()})
		return
	}

	if req.Days < 7 {
		req.Days = 7
	}

	//创建账号信息
	var user erps.FmUser
	baseSql := fmt.Sprintf("SELECT id, login_name, leader_flag,phonenum,password FROM fm_user WHERE login_name='%v' AND phonenum='%v' ", req.LoginName, req.PhoneNumber)
	err := DBConn.Get(&user, baseSql)
	apidebug.Printf("get user:%v info:%#v LoginName:%#v \n", baseSql, err, req)

	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"data": nil, "message": "failed", "error": "user not found" + err.Error()})
		return
	}

	//必需是 付费主账号 leader_flag = 4
	if user.LeaderFlag != "4" {
		helpers.JSONs(c, code.AuthorizedPowerError, gin.H{"data": nil, "message": "failed", "error": code.ZhCNText[code.AuthorizedPowerError] + user.LeaderFlag})
		return
	}

	//查询是否已有
	var total int
	err = DBConn.Get(&total, "SELECT COUNT(id) FROM fm_api_keys WHERE  user_id=?", user.Id)
	if err == nil && total >= 3 {
		helpers.JSONs(c, code.TooManyRequests, gin.H{"data": nil, "message": "failed", "error": fmt.Sprintf("已创建%v keys: %v 个", user.LoginName, total)})
		return
	}

	var keys []seclients.ApiKey
	err = DBConn.Select(&keys, "SELECT id, enabled, expires_time, created_time FROM fm_api_keys WHERE enabled = '1' AND user_id=?", user.Id)
	if err == nil && len(keys) >= 1 {
		helpers.JSONs(c, code.TooManyRequests, gin.H{"data": nil, "message": "failed", "error": fmt.Sprintf("已存在%v可使用 keys: %v 个", user.LoginName, len(keys))})
		return
	}

	apiKeyPlain, _ := GenerateAPIKey()
	hashed, _ := HashAPIKey(apiKeyPlain)
	expiry := time.Now().Add(time.Duration(req.Days) * 24 * time.Hour)

	_, err = DBConn.Exec("INSERT INTO fm_api_keys (user_id, api_key, enabled, expires_time) VALUES (?, ?, '1', ?)",
		user.Id, hashed, expiry)
	apidebug.Printf("INSERT user fm_api_keys:%#v days:%#v \n", user.Id, err)

	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"data": nil, "message": "failed", "error": "failed to create api key"})
		return
	}

	helpers.JSONs(c, code.Success, gin.H{
		"message":      "仅分发一次APIKey !请谨慎保存",
		"data":         apiKeyPlain,
		"expires_time": expiry,
	})
}

// 管理 API Key 状态
func (that *UsersApi) UpdateAPIKeyHandler(c *gin.Context) {
	var (
		DBConn = databases.DBMysql
	)
	var req struct {
		LoginName string `json:"login_name" db:"login_name"`
		Phonenum  string `json:"phonenum" db:"phonenum"`

		DeleteFlag string `json:"delete_flag" db:"delete_flag"`
		Enabled    bool   `json:"enabled" db:"enabled"`

		KeyID int `json:"id" db:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.JSONs(c, code.ParamBindError, gin.H{"error": "bad request " + err.Error()})
		return
	}

	//更新的目标账号查询
	var keyuser erps.FmUser
	baseSql := fmt.Sprintf("SELECT id, login_name,phonenum,password  FROM fm_user WHERE login_name='%v' and phonenum='%v' ", req.LoginName, req.Phonenum)
	err := DBConn.Get(&keyuser, baseSql)
	apidebug.Printf("get apikey  user nil?:%v err:%v \n", keyuser.Id != 0, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"data": nil, "message": "failed", "error": "user not found" + err.Error()})
		return
	}

	//当前登陆账号鉴权
	userNow, err := UserService.GetUserByCtxToken(c)
	apidebug.Printf("get apikey  user nil?:%v err:%v \n", userNow != nil, err)
	if err != nil {
		helpers.JSONs(c, code.AuthorizationError, gin.H{"data": nil, "message": "API Key delete failed", "error": "failed to revoke api key"})
		return
	}

	//删除用户账号自己的token
	if req.DeleteFlag == "1" {
		if userNow.Id != keyuser.Id {
			helpers.JSONs(c, code.AuthorizedPowerError, gin.H{"data": nil, "message": "API Key delete failed", "error": "只能删除自己的apikey"})
			return
		}

		delSql := fmt.Sprintf(" DELETE FROM fm_api_keys WHERE id='%v' AND user_id='%v' ", req.KeyID, keyuser.Id)
		rst, err := DBConn.Exec(delSql)
		apidebug.Printf("delete apikey rst nil?:%v err:%v \n", rst != nil, err)
		if err != nil {
			helpers.JSONs(c, code.NullData, gin.H{"data": nil, "message": "API Key delete failed", "error": "failed to revoke api key"})
			return
		}
	}

	baseSql = fmt.Sprintf(" UPDATE fm_api_keys SET enabled=%v WHERE id=%v AND user_id='%v' ", req.Enabled, req.KeyID, keyuser.Id)
	rst, err := DBConn.Exec(baseSql)
	apidebug.Printf("update baseSql:%v apikey rst nil?:%v err:%v \n", baseSql, rst != nil, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"data": nil, "message": "API Key update failed", "error": "failed to revoke api key"})
		return
	}
	rows, _ := rst.RowsAffected()
	msg := "API Key update done"
	if rows == 0 {
		msg = " 未找到API Key 更新"
	}
	helpers.JSONs(c, code.Success, gin.H{"data": rows, "message": msg})

}

// 查询用户所有 API Keys
func (that *UsersApi) ListAPIKeysHandler(ctx *gin.Context) {

	PageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	loginName := ctx.Query("login_name")
	Phonenum := ctx.Query("phonenum")

	userNow, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, gin.H{"data": nil, "message": "API Key delete failed", "error": "failed to revoke api key"})
		return
	}

	datas, err := UserService.QueryApiKeysList(Page, PageSize, loginName, Phonenum, userNow.TenantId)
	apidebug.Printf("list apikey nil?:%v err:%v \n", datas != nil, err)
	if err != nil {
		helpers.JSONs(ctx, code.NullData, gin.H{"code": code.NullData, "data": datas, "message": "failed", "error": code.ZhCNText[code.NullData] + err.Error()})
		return
	}
	helpers.JSONs(ctx, code.Success, datas)
}
