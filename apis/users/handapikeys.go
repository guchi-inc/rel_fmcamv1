package users

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/apikeys"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (that *UsersApi) GetApiKeys(c *gin.Context) {

	LoginName := c.DefaultQuery("login_name", "")
	Phonenum := c.DefaultQuery("phonenum", "")
	keyName := c.DefaultQuery("key_name", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.Apikeys{}

	if LoginName != "" && Phonenum != "" {
		us, err := UserService.GetUserValidInfo(LoginName, "", Phonenum)
		if err != nil {
			helpers.JSONs(c, code.NullData, err)
			return
		}
		conds = append(conds, apikeys.UserID(us.Id))
	}

	if keyName != "" {
		conds = append(conds, apikeys.KeyNameEQ(keyName))
	}

	datas, err := UserService.EntApiKeysList(Page, PageSize, c, conds)
	apidebug.Printf("list apikey nil?:%v err:%v \n", datas != nil, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"code": code.NullData, "data": datas, "message": "failed", "error": code.ZhCNText[code.NullData] + err.Error()})
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 更新 delete, SetEnabled, lastused,KeyName
func (that *UsersApi) UpdateApiKeys(c *gin.Context) {

	var (
		PParm = seclients.GenApikeys{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := UserService.EntApiKeysUpdate(c, &PParm)
	apidebug.Printf("list apikey nil?:%v err:%v \n", datas != 0, err)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"code": code.NullData, "data": datas, "message": "failed", "error": code.ZhCNText[code.NullData] + err.Error()})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas, "message": "success"})

}

// 新增城市
func (that *UsersApi) NewApiKeys(c *gin.Context) {

	var (
		PParm = seclients.ApiKey{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	userNow, err := UserService.GetUserByCtxToken(c)
	if err != nil {
		helpers.JSONs(c, code.AuthorizationError, gin.H{"data": nil, "message": "API Key new failed", "error": "failed to revoke api key"})
		return
	}

	//必需是 付费主账号 leader_flag = 4
	if userNow.LeaderFlag != "4" && userNow.Isystem != "1" {
		helpers.JSONs(c, code.AuthorizedPowerError, gin.H{"data": nil, "message": "API Key new failed", "error": code.ZhCNText[code.AuthorizedPowerError] + "当前账号必需为 付费账号或内置账号"})
		return
	}

	//查询目标用户的信息 赋予id
	distUser, err := UserService.GetUserValidInfo(*PParm.LoginName, "", *PParm.Phonenum)
	if err != nil {
		helpers.JSONs(c, code.AuthorizationError, gin.H{"data": nil, "message": "API Key new failed", "error": "该用户不存在" + *PParm.Phonenum + ":" + *PParm.LoginName})
		return
	}

	//赋值
	PParm.UserID = &distUser.Id
	datas, err := UserService.ApiKeysNew(&PParm)
	apidebug.Printf("query data:%#v conds:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
