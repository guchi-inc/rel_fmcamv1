// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/systems/sysutils"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	gqs "github.com/doug-martin/goqu/v9"
)

// 根据用户信息某个属性 登陆名查 doing
func (u *UserRepository) UserList(page, pageSize int, filters gqs.Ex, filter string) (*erps.UserListStruct, error) {

	var (
		total  int
		people = erps.UserListStruct{}
		users  = []erps.FmUser{}

		DBConn = databases.DBMysql
		Gq     = databases.Gq
	)

	//参数解读,生成sql语句
	//语句条件筛选
	fmt.Printf("Gq database:%#v is nil:%v\n", Gq, Gq == nil)

	totalSql, args, err := Gq.From("fm_user").Select(gqs.COUNT("*")).
		Where(filters).ToSQL()

	// 打印生成的 SQL
	logger.Printf("count all user filter:%v total sql %#v args:%v err:%v \n", filters, totalSql, args, err)
	err = DBConn.QueryRow(totalSql).Scan(&total)
	if err != nil {
		logger.Printf(" err:%#v totalSql:%v\n", err, totalSql)
		return &people, err
	}

	// aparamSql, args, err := Gq.From("fm_user").Where(filters).Order(gqs.I("id").Desc()).Order(gqs.I("username").Desc()).Limit(uint(pageSize)).Offset(uint(pageSize * (page - 1))).ToSQL() // .ScanStructs(&users)
	// aparamSql = strings.ReplaceAll(aparamSql, "\"", "'")
	// logger.Printf("all args %#v err:%#v aparamSql:%v\n", args, err, aparamSql)
	aparamSql := BaseUserSql + filter
	err = DBConn.Select(&users, aparamSql)
	if err != nil {
		logger.Printf(" err:%#v of aparamSql:%v\n", err, aparamSql)

		return &people, err
	}

	people.Data = users
	people.T = uint(total)
	people.S = uint(pageSize)
	if len(users) > 0 {
		people.Columns, _ = databases.FmGlobalMap("fm_user", nil)
	}
	return &people, nil
}

// 根据用户信息某个属性 比如tenant id 查
func (u *UserRepository) UserAccountInfo(accId int64, tenantId string) (*erps.FmUser, error) {

	var (
		people = erps.FmUser{}

		DBConn = databases.DBMysql
		Gq     = databases.Gq
	)

	//参数解读,生成sql语句
	//语句条件筛选
	Filter := " WHERE delete_flag != '1' "
	fmt.Printf("Gq database:%#v is nil:%v\n", Gq, Gq == nil)
	if accId != 0 {
		Filter += fmt.Sprintf(" AND id = '%v' ", accId)
	}
	if tenantId != "" {
		Filter += fmt.Sprintf(" AND (BIN_TO_UUID(tenant_id)) = '%v' ", uuid.MustParse(tenantId))
	}

	aparamSql := sysutils.BaseUserSql + Filter + " ORDER BY id DESC LIMIT 1 "
	err := DBConn.Get(&people, aparamSql)

	logger.Printf("all: %#v excel err:%#v aparamSql:%v\n", people, err, aparamSql)
	if err != nil {
		logger.Printf(" err:%#v of aparamSql:%v\n", err, aparamSql)

		return &people, err
	}

	return &people, nil
}

// 根据用户的 username 属性 查用户信息
func (u *UserRepository) UserInfo(pageSize int, usernames string) (*erps.UserListStruct, error) {

	var (
		total  int
		people = erps.UserListStruct{}
		users  = []erps.FmUser{}

		DBConn = databases.DBMysql
	)

	//参数解读,生成sql语句
	whereSql := fmt.Sprintf("WHERE username like '%v'", "%"+usernames+"%")
	//语句条件筛选
	totalSql := "SELECT COUNT(*) AS Total FROM fm_user " + whereSql

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logger.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &people, err
	}

	//最匹配的几个记录
	aparamSql := fmt.Sprintf("SELECT id,username,login_name,department,phonenum FROM fm_user %v ORDER BY username LIMIT %v", whereSql, pageSize)
	err = DBConn.Select(&users, aparamSql)
	logger.Printf("all peoples %#v err:%#v aparamSql:%v\n", people, err, aparamSql)
	people.Data = users
	people.T = uint(total)
	people.S = uint(len(users))

	return &people, nil
}

// 由请求token获取 用户的全部信息
func (u *UserRepository) GetUserByCtxToken(ctx *gin.Context) (*erps.FmUser, error) {

	var (
		DBConn = databases.DBMysql
	)
	// 查询单个数据 You can also get a single result, a la QueryRow

	//接收json参数
	//校验token信息
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		//apikey鉴权
		apiKey := ctx.GetHeader("X-API-KEY")
		logger.Printf("From X-API-KEY get api key:%v ", apiKey)
		if apiKey == "" {
			return nil, fmt.Errorf("Authorization missing X-API-KEY:%v AND Auth token", apiKey)
		}

		var keys []seclients.ApiKey
		baseSql := "SELECT * FROM ApiKeys WHERE enabled=1 AND expires_time > NOW()"
		err := DBConn.Select(&keys, baseSql)
		logger.Printf("baseSql:%v err:%v \n ", baseSql, err)
		if err != nil {
			return nil, fmt.Errorf("invalid API Key")
		}

		var uid int64
		for _, k := range keys {
			matched := bcrypt.CompareHashAndPassword([]byte(k.APIKey), []byte(apiKey))
			if matched == nil {
				uid = *k.UserID
			}
		}
		if uid == 0 {
			return nil, fmt.Errorf("invalid API Key:" + apiKey)
		}

		user, err := u.GetUserById(int(uid))
		if err != nil {
			logger.Printf("GetUserById:%v  error:%v", uid, err)
			return nil, err
		}

		if user.LoginName == "" {
			err = fmt.Errorf("user login request token:%v map error:%v", ctx.Request.Body, err)
			return nil, err
		}

		return user, nil
	}

	//token鉴权
	tokenString := authHeader[len(configs.AuthHeader):]
	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		err := fmt.Errorf("auth code:%v, err:%v", code.AuthorizationError, "非法token")
		return nil, err
	}

	//校验token信息,获取当前用户信息
	uid, _ := utils.GetIdFromToken(token)
	user, err := u.GetUserById(int(uid))
	if err != nil {
		logger.Printf("GetUserById:%v  error:%v", uid, err)
		return nil, err
	}

	if user.LoginName == "" {
		err = fmt.Errorf("user login request token:%v map error:%v", ctx.Request.Body, err)
		return nil, err
	}

	return user, nil
}

// 用户的全部信息
func (u *UserRepository) GetUserInfo(login_name string) (*erps.FmUser, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		usep = erps.FmUser{}

		DBConn  = databases.DBMysql
		baseSql = BaseUserSql
	)
	if login_name == "" {
		return nil, fmt.Errorf("login_name or id nil")
	}
	// fileds := strings.Join(models.TP_fm_user, ",")
	Fileter := fmt.Sprintf(" WHERE login_name = '%v' ", login_name)
	baseSql += Fileter
	err := DBConn.Get(&usep, baseSql)

	logger.Printf("from db get:%v user:%#v err:%v \n", login_name, baseSql, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil
}

// 操作工的登陆位置 信息
func (u *UserRepository) LocalHSET(login_name, remote_ip, local string) (*int64, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		DBRedis         = databases.DBS.DBRedis
		keyOfLoginCount = configs.GetDaysLoginKey(login_name)
		KeyOfMonth      = configs.GetMonthCountKey() //字典名称
	)

	rst := DBRedis.Client.HSet(configs.Ctx, configs.RedisKeyPrefixLocal, login_name, local+":"+remote_ip)
	sets, err := rst.Result()
	logger.Printf("redis db hset:%v user:%#v err:%v \n", sets, login_name, err)
	if err != nil {
		return nil, err
	}

	// 更新当天登陆次数, 当前日期
	rstSet, err := DBRedis.Client.IncrBy(configs.Ctx, keyOfLoginCount, 1).Result()
	logger.Printf("redis db rstSet count:%v user :%#v err:%v \n", rstSet, login_name, err)

	//当天第一次登陆，更新字典 对应位置
	if rstSet == 1 {
		rstMonthSet, err := DBRedis.Client.HIncrBy(configs.Ctx, KeyOfMonth, login_name, 1).Result()
		logger.Printf("redis db rstMonthSet count:%v user :%#v err:%v \n", rstMonthSet, login_name, err)

		if err != nil && rstMonthSet == 0 {
			rstMonthNew, _ := DBRedis.Client.HSetNX(configs.Ctx, KeyOfMonth, login_name, 1).Result()
			logger.Printf("new rstMonthSet count:%v user :%#v rstMonthNew:%v \n", KeyOfMonth, login_name, rstMonthNew)
		}
	}

	return &sets, nil
}

// 根据条件查询用户信息 ?
func (u *UserRepository) GetUserValidInfo(login_name, email, phonenum string) (*erps.FmUser, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		usep = erps.FmUser{}

		DBConn = databases.DBMysql
		Filter = " WHERE delete_flag != '1' "
	)
	if login_name == "" && email == "" && phonenum == "" {
		return nil, fmt.Errorf("login_name 和 email 参数为空")
	}
	if login_name != "" {
		Filter += fmt.Sprintf(" AND login_name = '%v' ", login_name)
	}

	if email != "" {
		Filter += fmt.Sprintf(" AND email = '%v' ", email)
	}

	if phonenum != "" {
		Filter += fmt.Sprintf(" AND phonenum = '%v' ", phonenum)
	}

	err := DBConn.Get(&usep, BaseUserSql+Filter)
	logger.Printf("from db get:%v user:%#v err:%v \n", login_name, usep, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil
}

// 更新单个数据
func (u *UserRepository) UpdateUserDeviceTime(userInfo *erps.FmUser, tenantId string) (any, error) {
	//  You can also get a single result, a la QueryRow

	var (
		DBConn = databases.DBMysql
	)

	//查待修改用户数据
	userOrigin, err := u.UserAccountInfo(userInfo.Id, tenantId)
	if err != nil {
		err = fmt.Errorf("数据不存在:%v", err)
		return nil, err
	}

	if userInfo.DeviceTime != "" {

		newTime, err := timeutil.ParseCSTInLocation(userInfo.DeviceTime)
		if err != nil {
			return nil, fmt.Errorf("时间格式RFC3339或年月日时分秒:%v", err)
		}
		userInfo.DeviceTime = timeutil.CSTLayoutString(&newTime)

		delTimeSql := fmt.Sprintf(`UPDATE fm_user SET device_time = '%v' WHERE id= '%v'`, userInfo.DeviceTime, userInfo.Id)
		rst, err := DBConn.Exec(delTimeSql)
		logger.Printf("update user:%v \n last time: %#v  err:%v \n", userInfo.Id, userOrigin.DeviceTime, err)
		if err != nil {
			return nil, err
		}
		return rst.RowsAffected()
	}
	return 0, fmt.Errorf("最近登陆时间无需更新")
}

// 更新单个用户数据
func (u *UserRepository) UpdateUser(userInfo *erps.FmUser, tenantId string) (any, error) {
	//  You can also get a single result, a la QueryRow

	var (
		DBConn = databases.DBMysql
	)

	//查待修改用户数据
	userOrigin, err := u.UserAccountInfo(userInfo.Id, tenantId)
	if err != nil {
		err = fmt.Errorf("数据不存在:%v", err)
		return nil, err
	}

	timeNow := time.Now().In(configs.Loc)
	tStr := timeutil.CSTLayoutString(&timeNow)
	if userInfo.DeletedAt != "" {

		if userInfo.LoginName != "admin" {
			newDelTime, err := timeutil.CheckDateEuq(userOrigin.DeletedAt, userInfo.DeletedAt)
			logger.Printf("from:%v orginDate:%v ifoDate:%#v Info.DeletedAt err:%v \n", newDelTime, userOrigin.DeletedAt, userInfo.DeletedAt, err)
			if err != nil {
				return nil, fmt.Errorf("时间格式RFC3339或年月日时分秒:%v", err)
			}
			if *newDelTime {
				newTime, err := timeutil.ParseCSTInLocation(userInfo.DeletedAt)
				if err != nil {
					return nil, fmt.Errorf("时间格式RFC3339或年月日时分秒:%v", err)
				}
				userInfo.DeletedAt = timeutil.CSTLayoutString(&newTime)
			}

			delTimeSql := fmt.Sprintf(`UPDATE fm_user SET deleted_time = '%v' WHERE id= '%v'`, userInfo.DeletedAt, userInfo.Id)
			rst, err := DBConn.Exec(delTimeSql)
			logger.Printf("update user:%v \n last time: %#v  err:%v \n", userInfo.Id, userOrigin.DeviceTime, err)
			if err != nil {
				return nil, err
			}
			return rst.RowsAffected()
		}
	}

	//优先更新删除状态，如果 DeleteFlag 状态为1，只更新删除标记
	if userInfo.DeleteFlag == "1" {

		upSql := fmt.Sprintf("UPDATE fm_user SET enabled = '%v', delete_flag = '%v'  WHERE id = '%v';",
			0, 1, userInfo.Id)

		rst, err := DBConn.Exec(upSql)
		logger.Printf("from db select:%#v  err:%v \n", rst, err)
		if err != nil {
			return nil, err
		}

		return rst.RowsAffected()
	}

	//更新账号的 短信验证码状态
	if (!userInfo.IsSms || userInfo.IsSms) && (userInfo.IsSms != userOrigin.IsSms) {
		upSql := fmt.Sprintf("UPDATE fm_user SET is_sms = %v  WHERE id = '%v';",
			userInfo.IsSms, userInfo.Id)

		rst, err := DBConn.Exec(upSql)
		logger.Printf("from db update:%#v  err:%v \n", upSql, err)
		if err != nil {
			return nil, err
		}

		return rst.RowsAffected()
	}

	//优先更新状态
	if userInfo.Enabled != nil && (*userOrigin.Enabled != *userInfo.Enabled) {

		var intEnabled int
		if !*userInfo.Enabled {
			intEnabled = 0
		} else {
			intEnabled = 1
		}
		upSql := fmt.Sprintf("UPDATE fm_user SET enabled = '%v', updated_time = '%v' WHERE id = '%v';",
			intEnabled, tStr, userInfo.Id)

		rst, err := DBConn.Exec(upSql)
		logger.Printf("from db select:%#v  err:%v \n", rst, err)
		if err != nil {
			return nil, err
		}

		return rst.RowsAffected()
	}

	//只允许 管理员 根据 id 更新
	upTotal := 0

	//更新其他信息
	if userInfo.LoginName != "" && (userOrigin.LoginName != userInfo.LoginName) {
		userOrigin.LoginName = userInfo.LoginName
		upTotal += 1
	}

	if userInfo.Username != "" && (userOrigin.Username != userInfo.Username) {
		userOrigin.Username = userInfo.Username
		upTotal += 1
	}

	if userInfo.Email != "" && (userOrigin.Email != userInfo.Email) {
		userOrigin.Email = userInfo.Email
		upTotal += 1
	}

	if userInfo.Phonenum != "" && (userOrigin.Phonenum != userInfo.Phonenum) {
		userOrigin.Phonenum = userInfo.Phonenum
		upTotal += 1
	}

	if userInfo.Position != "" && (userOrigin.Position != userInfo.Position) {
		userOrigin.Position = userInfo.Position
		upTotal += 1
	}

	if userInfo.Department != "" && (userOrigin.Department != userInfo.Department) {
		userOrigin.Department = userInfo.Department
		upTotal += 1
	}

	if userInfo.LeaderFlag != "" && (userOrigin.LeaderFlag != userInfo.LeaderFlag) {
		userOrigin.LeaderFlag = userInfo.LeaderFlag
		upTotal += 1
	}
	if userInfo.Ismanager != "" && (userOrigin.Ismanager != userInfo.Ismanager) {
		userOrigin.Ismanager = userInfo.Ismanager
		upTotal += 1
	}
	if userInfo.Local != "" && (userOrigin.Local != userInfo.Local) {
		userOrigin.Local = userInfo.Local
		upTotal += 1
	}

	if userInfo.LocalHost != "" && (userOrigin.LocalHost != userInfo.LocalHost) {
		userOrigin.LocalHost = userInfo.LocalHost
		upTotal += 1
	}

	if userInfo.Description != "" && (userOrigin.Description != userInfo.Description) {
		userOrigin.Description = userInfo.Description
		upTotal += 1
	}
	if userInfo.Ethnicity != "" && (userOrigin.Ethnicity != userInfo.Ethnicity) {
		userOrigin.Ethnicity = userInfo.Ethnicity
		upTotal += 1
	}
	if userInfo.Gender != "" && (userOrigin.Gender != userInfo.Gender) {
		userOrigin.Gender = userInfo.Gender
		upTotal += 1
	}
	if userInfo.MemberId != "" && (userOrigin.MemberId != userInfo.MemberId) {
		userOrigin.MemberId = userInfo.MemberId
		upTotal += 1
	}
	if userInfo.LeaderId != "" && (userOrigin.LeaderId != userInfo.LeaderId) {
		userOrigin.LeaderId = userInfo.LeaderId
		upTotal += 1
	}
	if userInfo.M2LocalHost != "" && (userOrigin.M2LocalHost != userInfo.M2LocalHost) {
		userOrigin.M2LocalHost = userInfo.M2LocalHost
		upTotal += 1
	}

	if upTotal >= 1 {
		//修改某些值的时候 更新该时间
		userOrigin.UpdatedAt = tStr

	} else {
		err = fmt.Errorf("数据无需修改。")
		return nil, err
	}

	personStructs := []erps.FmUser{
		*userOrigin,
	}

	rst, err := DBConn.NamedExec(`UPDATE fm_user SET login_name=:login_name, 
	username=:username,email=:email, phonenum=:phonenum,position=:position,
	description=:description, ethnicity=:ethnicity, gender=:gender,
	member_id=:member_id, ismanager=:ismanager, 
	department=:department,leader_flag=:leader_flag,m2_localhost=:m2_localhost,
	local=:local,localhost=:localhost
	WHERE id=:id`, personStructs)
	logger.Printf("update user:%v \n last time: %#v  err:%v \n", upTotal, userOrigin.DeviceTime, err)
	if err != nil {
		return nil, err
	}

	return rst.RowsAffected()
}

func (u *UserRepository) GetUserById(id int) (*erps.FmUser, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	// 返回用户信息
	var (
		usep = erps.FmUser{}

		DBConn  = databases.DBMysql
		baseSql = BaseUserSql
	)
	if id == 0 {
		return nil, fmt.Errorf("param id error")
	}
	// baseSql += fmt.Sprintf("WHERE id= '%v' ", id)
	err := DBConn.Get(&usep, baseSql+" WHERE id=?", id)

	logger.Printf("from db get user:%#v err:%v \n", usep, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil
}

// 用户信息，包括角色组和权限组
func (u *UserRepository) GetUserRolePowerById(id int) (*erps.FmUser, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	// 返回用户信息
	var (
		usep = erps.FmUser{}

		DBConn = databases.DBMysql
	)
	err := DBConn.Get(&usep, BaseUserSql+"  WHERE id=?", id)

	logger.Printf("from db get user:%#v err:%v \n", usep, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil
}

// 注册接口
func (u *UserRepository) Register(ctx *gin.Context, newOne *erps.GcParamUser) (*erps.GcParamUser, error) {

	var (
		defDepart = configs.DefDepart //默认部门
		// DefIsystem           = 2     //1 是，2 否
		// DefLeaderFlage       = 2     //1 是，2否
		// DefStatus            = 1     //1：正常，2：删除，3：封禁
		IsAllowDep = false
		DBConn     = databases.DBMysql
	)

	timeNow := time.Now().In(configs.Loc)
	tStr := timeutil.CSTLayoutString(&timeNow)

	if newOne.CreatedAt == "" {
		newOne.CreatedAt = tStr
	}

	//楼栋后缀
	var uu erps.FmUser
	err := DBConn.Get(&uu, fmt.Sprintf("SELECT * FROM fm_user WHERE  login_name = '%v' ", newOne.LoginName))
	logger.Printf("get user :%#v, with db rst err:%v\n", newOne.LoginName, err)

	//已经存在 姓名和地址 一样的信息
	if uu.Phonenum == newOne.Phonenum {
		return nil, fmt.Errorf("有:%v 相同用户:%v 手机号:%v 上次登记时间:%v", newOne.LoginName, uu.LoginName, uu.Phonenum, uu.CreatedAt)
	}

	//可用部门列表
	AllowDepartments, _ := sysutils.NewDepartmentUtil().GeDepartmentNameList(false)

	//检查是否在允许的部门内，默认有 生产部，检验部，管理部
	for _, dep_name := range AllowDepartments.Data {
		logger.Printf("match:gu.Department %v  dep_name.DepName: %v", newOne.Department, dep_name.DepName)
		if newOne.Department == dep_name.DepName {
			IsAllowDep = true
			break
		}
	}

	//如果不是允许的部门
	if !IsAllowDep {
		newOne.Department = defDepart
	}

	//先注册 酒店租户信息
	tenantInfo := &seclients.TenantFull{
		PhoneNum: newOne.Phonenum,
		Email:    newOne.Email,
		Contacts: newOne.Username,
		Creator:  newOne.LoginName,
		Supplier: newOne.Supplier}
	if newOne.Supplier == "" {
		tenantInfo.Supplier = "未填写企业公司名"
	}

	// 插入 SQL
	// 开始事务
	//鉴权
	user, err := u.GetUserByCtxToken(ctx)
	//只允许 管理员 根据 menu_number 更新 delete_flag, description, level
	if err != nil {
		logger.Printf("user_tenant:%v   err:%v\n", user, err)
	}

	//绑定用户 的 主账号
	if user != nil {
		newOne.LeaderId = user.LeaderId
	}

	errTx := DBConn.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %w", err)
	}

	//如果当前创建子账号的是付费账号 或内置租户的账号，则使用已有账号的 租户号，并关联子账号的 leader_id = 当前账号的租户号
	var newuuid string
	if user == nil {
		newuuid = uuid.New().String()
		fmt.Println("Tenant Binary:", newuuid, "new tx errTx:", errTx, "tx begin err:", err)

		newTentSql := fmt.Sprintf(`
			INSERT INTO fm_tenant (supplier,contacts,phone_num,tenant_id,creator) VALUES ('%v', '%v','%v', %v, '%v')
			`, tenantInfo.Supplier, tenantInfo.Contacts, tenantInfo.PhoneNum, fmt.Sprintf("(UUID_TO_BIN('%v'))", newuuid), tenantInfo.Creator)

		//写入 租户信息执行
		rst, err := tx.Exec(newTentSql)
		logger.Printf("new fm_tenant:%v insert:%#v, with db rst:%#v err:%v\n", newuuid, newTentSql, rst, err)
		if err != nil {
			logger.Println("插入失败: ", err)
			return nil, err
		}

		newTenantId, err := rst.LastInsertId()
		logger.Printf("new suppliy register:%#v, with new tenantid:%v err:%v\n", newOne, newTenantId, err)
		if err != nil {
			return nil, err
		}
		//0 没有上级账户
		newOne.LeaderId = 0

	} else {
		//付费会员 可创建用户
		if user.LeaderFlag == "4" {
			newuuid = user.TenantId
			newOne.LeaderId = user.Id
		} else {
			return nil, fmt.Errorf("%v:%v", code.ZhCNText[code.AuthorizedPowerError], user.LoginName)
		}

		newOne.Creator = user.LoginName
	}

	if newOne.Username == "" {
		//默认账号名
		newOne.Username = "员工"
	}
	if newOne.LeaderFlag == "" {
		//职员, 4付费会员账号,3顾客,2客户,1职员,0未定义
		newOne.LeaderFlag = "1"
	}

	//注册账号
	newOne.TenantId = newuuid
	personStructs := []erps.GcParamUser{
		*newOne,
	}

	rst, err := tx.NamedExec(`INSERT INTO fm_user (
		username,
		login_name,
		password,
		
		leader_flag,
		leader_id,
		position,

		department,
		email,
		phonenum)
        VALUES (:username, :login_name,:password,
			:leader_flag, :leader_id,:position,
			:department, :email,:phonenum)`, personStructs)

	logger.Printf("new user register:%#v, with db rst:%v err:%v\n", newOne, rst, err)
	if err != nil {
		// make sure err is a mysql.MySQLError.
		if errMySQL, ok := err.(*mysql.MySQLError); ok {
			switch errMySQL.Number {
			case 1062:
				// TODO handle Error 1062: Duplicate entry '%s' for key %d
				return nil, fmt.Errorf("已添加该用户:%v old:%v  手机号:%v", newOne.Username, uu.Username, newOne.Phonenum)
			}
		}
		return nil, err
	}

	lastId, err := rst.LastInsertId()
	logger.Printf("new user register:%#v, with db rst:%v \n", lastId, err)
	if err != nil {
		return nil, err
	}

	//更新uuid
	if newOne.TenantId != "" {
		upUUID := fmt.Sprintf(" UPDATE fm_user SET tenant_id = %v WHERE id = %v ", fmt.Sprintf("(UUID_TO_BIN('%v'))", newOne.TenantId), lastId)
		rst, errec := tx.Exec(upUUID)
		logger.Printf("new Profile:%#v err:%#v \n", rst, errec)
		if errec != nil {
			return nil, err
		}
	}

	//记录操作
	newSqlLog := fmt.Sprintf("INSERT INTO sql_logs (table_name,query,args,creator) VALUES ('%v', '%v', '%v','%v')", "fm_tenant:"+"fm_user", fmt.Sprintf("%v", newOne), newuuid, fmt.Sprintf("%v ", lastId))
	_, err = tx.Exec(newSqlLog)
	logger.Printf("new sql_logs:%#v, with db rst:%v \n", newSqlLog, err)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	logger.Printf("事务成功提交新:%v，租户uuid:%v 新增账号ID: %d\n", newOne.LoginName, newuuid, lastId)
	newOne.ID = lastId
	return newOne, err
}
