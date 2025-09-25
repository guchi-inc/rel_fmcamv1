package users

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/sysutils"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (u *UserRepository) AuthProfileInfo(phonenum, login_name string) (users *erps.FmUser, err error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		usep = erps.FmUser{}

		DBConn  = databases.DBMysql
		baseSql = BaseUserSql
	)

	if login_name == "" && phonenum == "" {
		return nil, fmt.Errorf("账号参数错误")
	}

	Fileter := " WHERE id > 0 "
	if login_name != "" {
		Fileter += " AND login_name = " + "'" + login_name + "'"
	}
	if phonenum != "" {
		Fileter += " AND phonenum = " + "'" + phonenum + "'"
	}

	baseSql += Fileter
	err = DBConn.Get(&usep, baseSql)
	logger.Printf("from db get:%v user:%#v err:%v \n", login_name, baseSql, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil

}

// 注册接口
func (u *AccountsRes) NewAccounts(newOne *erps.GcParamUser) (*erps.GcParamUser, error) {

	var (
		defDepart = "社区" //默认部门
		// DefIsystem           = 2     //1 是，2 否
		// DefLeaderFlage       = 2     //1 是，2否
		// DefStatus            = 1     //1：正常，2：删除，3：封禁
		IsAllowDep = false

		DBConn = databases.DBMysql
	)

	timeNow := time.Now().In(configs.Loc)
	tStr := timeutil.CSTLayoutString(&timeNow)

	if newOne.CreatedAt == "" {
		newOne.CreatedAt = tStr
	}

	var uu erps.FmUser
	err := DBConn.Get(&uu, fmt.Sprintf(" %v WHERE  login_name = '%v' ", BaseUserSql, newOne.LoginName))
	logger.Printf("get user :%#v, with db rst err:%v\n", newOne.LoginName, err)

	//已经存在 姓名和地址 一样的信息
	if err == nil && uu.Username == newOne.Username && uu.LocalHost == newOne.LocalHost {
		return nil, fmt.Errorf("有:%v 相同用户:%v 地址:%v 上次登记时间:%v", newOne.Username, uu.Username, uu.LocalHost, uu.UpdatedAt)
	} else {
		logger.Printf("新增登录名。 已经存在负责人:%v 信息 %v", uu.Username, uu.LoginName)
		newOne.LoginName = newOne.LoginName + utils.GenStandString(2)
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

	//如果member id 为空
	// 更新 member_id 信息, 也就是员工号
	MonDay := configs.MonthDate[timeNow.Month().String()]
	if len(newOne.MemberId) <= 0 {
		if MonDay < 10 {
			newOne.MemberId = fmt.Sprintf("0%v%v%v", MonDay, timeNow.Day(), utils.GenerateRandomString(4))
		} else {
			newOne.MemberId = fmt.Sprintf("%v%v%v", MonDay, timeNow.Day(), utils.GenerateRandomString(3))
		}
	}

	personStructs := []erps.GcParamUser{
		*newOne,
	}

	rst, err := DBConn.NamedExec(`INSERT INTO fm_user (
		username,
		login_name,
		password,
		leader_flag,
		position,
		department,
		email,
		phonenum,
		ismanager,  
		isystem,

		enabled,
		is_sms,
		description, 
		ethnicity,
		gender,
		local,
		localhost,
		m2_localhost,
		member_id, 
		device_id,

		created_time)
        VALUES (:username, :login_name,:password,:leader_flag, :position,
			:department, :email,:phonenum,:ismanager, 2,
			1,:is_sms,:description,:ethnicity, :gender,:local,:localhost,:m2_localhost,
			:member_id, :device_id,:created_time)`, personStructs)

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
		rst, errec := DBConn.Exec(upUUID)
		logger.Printf("new Profile:%#v err:%#v \n", rst, errec)
		if errec != nil {
			return nil, err
		}
	}

	return newOne, err
}

// 根据租户号 tenant_id 查
func (u *UserRepository) UserAccountList(Page, PageSize int, TenantId, Phonenum string) (*seclients.FmAccountsList, error) {

	var (
		total  int
		people seclients.FmAccountsList
		users  []*genclients.FmUserAccount

		DBConn  = databases.DBMysql
		Gq      = databases.Gq
		baseSql = BaseUserSql
	)

	//参数解读,生成sql语句
	//语句条件筛选
	fmt.Printf("Gq database:%#v is nil:%v\n", Gq, Gq == nil)

	Filter := " WHERE delete_flag != '1'  "

	if TenantId != "" {
		Filter += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}

	if Phonenum != "" {
		Filter += fmt.Sprintf(" AND phonenum = '%v' ", Phonenum)
	}

	totalSql := "SELECT COUNT(*) AS 'total' FROM fm_user  " + Filter

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)

	logger.Printf("count all user err:%v total sql %#v \n", err, totalSql)
	if err != nil {
		logger.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &people, err
	}

	Limits := fmt.Sprintf(" ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filter
	baseSql += Limits

	err = DBConn.Select(&users, baseSql)
	logger.Printf("all args %#v err:%#v totalSql:%v\n", baseSql, err, totalSql)
	if err != nil {
		logger.Printf(" err:%#v of baseSql:%v\n", err, baseSql)
		return &people, err
	}

	people.Data = users
	people.Total = total
	people.PageSize = PageSize
	if len(users) > 0 {
		people.Columns, _ = databases.FmGlobalMap("fm_user", nil)
	}
	return &people, nil
}

// 用户的全部信息 包括租户信息
func (u *UserRepository) GetUserAccountInfo(login_name string) (*genclients.FmUserAccount, *seclients.TenantFull, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		usep = genclients.FmUserAccount{}

		DBConn  = databases.DBMysql
		baseSql = BaseUserSql
	)

	if login_name == "" {
		return nil, nil, fmt.Errorf("login_name or id nil")
	}
	// fileds := strings.Join(models.TP_fm_user, ",")
	Fileter := fmt.Sprintf(" WHERE login_name = '%v' ", login_name)
	baseSql += Fileter
	err := DBConn.Get(&usep, baseSql)

	logger.Printf("from db get:%v user:%#v err:%v \n", login_name, baseSql, err)
	if err != nil {
		return nil, nil, err
	}

	//租户信息
	customerInfo := sysutils.CustomerUtil{}
	if usep.TenantID.String() == configs.DefTenantId {
		return &usep, nil, nil
	}
	tenantInfo, err := customerInfo.GetCustomerInfo(0, "", "", usep.TenantID.String())
	if err != nil {
		return nil, nil, err
	}
	return &usep, tenantInfo, nil
}
