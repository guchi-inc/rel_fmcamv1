package users

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	GenerateAPIKey = utils.GenerateAPIKey
	HashAPIKey     = utils.HashAPIKey
)

func (u *UserRepository) EntApiKeysList(Page, PageSize int, Ctx *gin.Context, conds []predicate.Apikeys) (*seclients.GenApiKeyPages, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.GenApiKeyPages
	offset := (Page - 1) * PageSize

	total, err := client.Debug().Apikeys.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().Apikeys.Query().
		Where(conds...).
		Limit(PageSize).
		Offset(offset).Select("id", "user_id", "tenant_id", "api_key", "key_name", "usage_count", "expires_time", "enabled", "created_time", "updated_time", "last_used_time", "type"). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//.All(Ctx)
	list := result.Data
	logger.Println("Apikeys retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("ApiKeys", nil)
		cm = append(cm, erps.GcDesc{Name: "login_name", CName: "登陆名", IsSearchable: true})
		cm = append(cm, erps.GcDesc{Name: "phonenum", CName: "手机号", IsSearchable: true})

		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

func (u *UserRepository) EntApiKeysUpdate(c *gin.Context, PParm *seclients.GenApikeys) (int, error) {

	var (
		client = databases.EntClient
	)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)
	originPm, err := client.Apikeys.Get(Ctx, PParm.ID)

	logger.Printf("originPm:%#v receive params:%#v", originPm, PParm)
	if err != nil {
		return 0, err
	}

	update := client.Apikeys.UpdateOneID(PParm.ID)
	//	// Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag == "1" {
		delOne := client.Apikeys.DeleteOneID(PParm.ID)
		err := delOne.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

	var changed = false
	if PParm.Enabled != nil && (*PParm.Enabled != originPm.Enabled) {
		update = update.SetEnabled(*PParm.Enabled)
		changed = true
	}

	if PParm.KeyName != "" && PParm.KeyName != originPm.KeyName {
		update = update.SetKeyName(PParm.KeyName)
		changed = true
	}
	if PParm.Type != nil && *PParm.Type != originPm.Type {
		update = update.SetType(*PParm.Type)
		changed = true
	}

	if PParm.LastUsedTime != originPm.LastUsedTime {
		tstr := timeutil.CSTLayoutString(&PParm.LastUsedTime)
		//时间不得早于2006
		err := timeutil.DateCheckStart(tstr)
		logger.Printf("last used time:%v format:%#v \n", tstr, err)
		if err == nil {
			update = update.SetLastUsedTime(PParm.LastUsedTime)
			changed = true
		}

	}

	if changed {
		//鉴权

		//更新执行
		err = update.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return 0, fmt.Errorf("nothing to do")

}

func (u *UserRepository) ApiKeysNew(PParm *seclients.ApiKey) (string, error) {

	var (
		DBConn = databases.DBMysql
		users  *erps.FmUser
		err    error
	)

	if PParm.UserID != nil {
		users, err = u.GetUserById(int(*PParm.UserID))
		if err != nil {
			return "", fmt.Errorf("用户信息错误%v", PParm.UserID)
		}
	}

	//创建账号信息
	logger.Printf("get days :%v LoginName:%#v \n", PParm.Days, users)

	//查询是否已有
	var total int
	err = DBConn.Get(&total, "SELECT COUNT(id) FROM ApiKeys WHERE user_id=?", PParm.UserID)
	if err == nil && total >= 3 {
		return "", fmt.Errorf("已创建%v keys: %v 个", users, total)
	}

	var enabledkeys int
	err = DBConn.Get(&enabledkeys, "SELECT Count(id) FROM ApiKeys WHERE enabled = '1' AND user_id=?", PParm.UserID)
	if err == nil && enabledkeys >= 1 {
		return "", fmt.Errorf("已存在%v可使用 keys: %v 个", users, enabledkeys)
	}

	apiKeyPlain, _ := GenerateAPIKey()
	hashed, _ := HashAPIKey(apiKeyPlain)
	expiry := time.Now().Local().In(configs.Loc).Add(time.Duration(*PParm.Days) * 24 * time.Hour)

	//默认使用场景
	defType := 0 //web
	if PParm.Type != nil {
		defType = *PParm.Type //内部
	}

	if PParm.KeyName == nil {
		newName := utils.GenStandString(3) + *PParm.LoginName
		PParm.KeyName = &newName
	}
	newSql := fmt.Sprintf(`INSERT INTO ApiKeys (user_id, tenant_id, api_key, key_name, enabled, expires_time,type) 
	VALUES ('%v',%v ,'%v','%v',1, '%v', '%v')`, users.Id, fmt.Sprintf("(UUID_TO_BIN('%v'))", users.TenantId), hashed, *PParm.KeyName, timeutil.CSTLayoutString(&expiry), defType)
	_, err = DBConn.Exec(newSql)
	logger.Printf("INSERT user ApiKeys:%#v days:%#v \n", newSql, err)

	if err != nil {
		return "", fmt.Errorf("failed to create api key :%v", err)
	}

	return apiKeyPlain, nil

}

func (u *UserRepository) QueryApiKeysList(Page, PageSize int, loginName, Phonenum, TenantId string) (*seclients.ApiKeyList, error) {
	var (
		DBConn  = databases.DBMysql
		RetKeys = seclients.ApiKeyList{}
	)

	var filter = " WHERE a.id > 0"

	if loginName != "" || Phonenum != "" {
		filter += fmt.Sprintf(" AND (b.login_name= '%v' OR b.phonenum = '%v') ", loginName, Phonenum)
	}

	//非默认租户
	UserTenantFilter := ""
	if TenantId != configs.DefTenantId {
		UserTenantFilter += fmt.Sprintf(" AND BIN_TO_UUID(b.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	//总数
	totalSql := fmt.Sprintf(`SELECT COUNT(a.id) FROM fm_api_keys a 
		LEFT JOIN fm_user b ON b.id = a.user_id
		 %v  `, filter)
	var total int
	err := DBConn.Get(&total, totalSql)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys" + err.Error())
	}

	//数据
	var keys []seclients.ApiKey
	Limits := fmt.Sprintf(" ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql := fmt.Sprintf(`SELECT a.id, a.enabled, a.expires_time,b.login_name,b.phonenum, a.created_time FROM fm_api_keys a 
	LEFT JOIN fm_user b ON b.id = a.user_id %v
	 %v  `, UserTenantFilter, filter)

	baseSql += Limits
	err = DBConn.Select(&keys, baseSql)
	logger.Printf("get keys:%#v baseSql:%#v err:%#v \n", keys, baseSql, err)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys" + err.Error())
	}

	RetKeys.Data = keys
	RetKeys.Total = total
	RetKeys.Size = PageSize
	if len(keys) > 0 {
		RetKeys.Columns = []erps.GcDesc{
			{Name: "id", CName: "令牌号", IsVisible: false},
			{Name: "enabled", CName: "状态", IsVisible: true, IsEditable: true},
			{Name: "login_name", CName: "登陆名", IsVisible: false, IsSearchable: true},
			{Name: "phonenum", CName: "手机号", IsVisible: false, IsSearchable: true},

			{Name: "expires_time", CName: "到期", IsVisible: true, IsEditable: false},
			{Name: "created_time", CName: "生成日期", IsVisible: true, IsEditable: false},
		}
	}

	return &RetKeys, nil
}
