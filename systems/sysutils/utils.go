// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package sysutils

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils"
	"fmcam/models/erps"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	UsersData   = &erps.FmUser{}
	BaseUserSql = `SELECT id,
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
	leader_id,
	device_time,
	BIN_TO_UUID(tenant_id) as tenant_id,
	created_time,
	updated_time,
	deleted_time FROM fm_user `
)

// 获取后台管理 [菜单] 列表
func (ut *DepartmentUtil) GetRolesList(Page, PageSize int, LogiName, RoleName, TenantId, StartAt, EndAt string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT  id, role_name, type, price_limit, description,enabled,sort,BIN_TO_UUID(tenant_id) as tenant_id,created_time,delete_flag FROM `fm_role` "

		total  int
		roles  = erps.GcRoleList{}
		DBConn = databases.DBMysql
	)

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_role`  "

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)
		return &roles, err
	}

	//条件
	Filters := " WHERE delete_flag != '1' "
	if RoleName != "" {
		Filters += " AND role_name LIKE '%" + RoleName + "%' "
	}
	if TenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}

	if StartAt != "" && EndAt != "" {
		Filters += " AND (created_time between '" + StartAt + "' AND ' " + EndAt + "')"

	}

	Limits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	//admin 查询全部角色数据，具有编辑权限
	if LogiName == "admin" {
		baseSql = "SELECT  id, role_name, type, price_limit,description,enabled,sort,BIN_TO_UUID(tenant_id) as tenant_id,delete_flag FROM  `fm_role`  "
		baseSql += Filters
		baseSql += Limits
		err = DBConn.Select(&roles.Data, baseSql)
		logsys.Printf("from db admin baseSql:%#v  err:%v \n", baseSql, err)
		if err != nil {
			return nil, err
		}

		roles.Total = uint(total)
		roles.Size = uint(PageSize)
		if len(roles.Data) > 0 {
			roles.Columns = FiledOper.FmRolesMap()
		}
		return &roles, nil
	}

	baseSql += Filters
	baseSql += Limits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&roles.Data, baseSql)
	logsys.Printf("from role select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	roles.Total = uint(total)
	roles.Size = uint(PageSize)
	if len(roles.Data) > 0 {
		roles.Columns = FiledOper.FmRolesMap()
	}
	return &roles, nil
}

// 用户账号id 和 角色id是否同一个 租户
func (ut *DepartmentUtil) SamceTenantUserAndRole(userId, roleId int) bool {

	logsys.Printf("userAccount.userId %v RoleList.roleId %v \n", userId, roleId)

	userAccount, err := UserUtils.GetUserById(fmt.Sprintf("%v", userId))
	if err != nil {
		return false
	}

	RoleList, err := ut.GetRoleInfoById(1, 10, roleId, "")
	if err != nil || userAccount == nil {
		return false
	}

	if len(RoleList.Data) == 0 {
		return false
	}

	logsys.Printf("userAccount.TenantId %v RoleList.Data[0].TenantId %v", userAccount.TenantId, RoleList.Data[0].TenantId)

	return userAccount.TenantId == RoleList.Data[0].TenantId
}

// 获取后台管理 [菜单] 列表
func (ut *DepartmentUtil) GetRoleInfoById(Page, PageSize, id int, TenantId string) (*erps.GcRoleList, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT  id, role_name, type, price_limit, description,enabled,sort,BIN_TO_UUID(tenant_id) as tenant_id,created_time,delete_flag FROM `fm_role` "

		total  int
		roles  = erps.GcRoleList{}
		DBConn = databases.DBMysql
	)
	//条件
	Filters := " WHERE delete_flag != '1' "
	if id != 0 {
		Filters += fmt.Sprintf(" AND id = '%v' ", id)
	}
	if TenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_role`  " + Filters

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)
		return &roles, err
	}

	Limits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)

	//admin 查询全部角色数据，具有编辑权限

	baseSql += Filters
	baseSql += Limits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&roles.Data, baseSql)
	logsys.Printf("from role select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	roles.Total = uint(total)
	roles.Size = uint(PageSize)
	if len(roles.Data) > 0 {
		roles.Columns = FiledOper.FmRolesMap()
	}
	return &roles, nil
}

// 获取后台管理 [菜单] 列表
func (ut *DepartmentUtil) GeMenuFullList(Page, PageSize int, login_name, menu_number, description string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT id, menu_type, menu_name, menu_number,level,description,leader_id FROM  `fm_menu`   "

		total    int
		useMenus = erps.GcMenuList{}
		DBConn   = databases.DBMysql
	)

	var Filter string = " WHERE delete_flag != '1' "
	if menu_number != "" {
		Filter += fmt.Sprintf(" AND menu_number like '%v' ", "%"+menu_number+"%")
	}
	if description != "" {
		Filter += fmt.Sprintf(" AND description like '%v' ", "%"+description+"%")
	}

	totalSql := "SELECT COUNT(*) AS Total FROM `fm_menu` " + Filter

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &useMenus, err
	}

	//admin 查询全部菜单数据，具有编辑权限
	if login_name == "admin" {
		usedep := erps.GcFullMenuList{}
		baseSql = fmt.Sprintf("SELECT id, menu_type, menu_name, menu_number,created_time,level,description,leader_id,delete_flag,path FROM  `fm_menu` %v GROUP BY menu_number ORDER BY menu_number ;", Filter)

		err = DBConn.Select(&usedep.Data, baseSql)
		logsys.Printf("from db select:%#v  err:%v \n", usedep, err)
		if err != nil {
			return nil, err
		}

		usedep.Total = uint(total)
		usedep.Size = uint(len(usedep.Data))
		return &usedep, nil
	} else {
		return ut.GeMenuLimitsList(login_name, "")
	}

	// Limits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	// // 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	// baseSql += Filter
	// baseSql += Limits
	// err = DBConn.Select(&useMenus.Data, baseSql)
	// logsys.Printf("sql:%#v from db select:%#v  err:%v \n", baseSql, useMenus, err)
	// if err != nil {
	// 	return nil, err
	// }

	// useMenus.Total = uint(total)
	// useMenus.Size = uint(len(useMenus.Data))
	// return &useMenus, nil
}

// 更新角色数据
func (ut *DepartmentUtil) PostRoleUpdate(roles *erps.GcRole) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	if roles.CreatedAt == "" {
		roles.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	if roles.Type == "" {
		roles.Type = "RoleFunction" // 角色功能模块
	}

	if roles.PriceLimit == "" {
		roles.PriceLimit = "0" //不屏蔽
	}

	if roles.DeleteFlag != "" {
		rst, err := DBConn.NamedExec(`UPDATE fm_role SET
			delete_flag=:delete_flag 
			WHERE id=:id`, roles)

		logsys.Printf("fm_role update:%#v, with db rst:%v err:%v\n", roles, rst, err)
		if err != nil {

			return 0, err
		}

		rst, err = DBConn.NamedExec(`UPDATE fm_role_user SET delete_flag=:delete_flag WHERE role_name=:role_name`, roles)
		if err != nil {
			return 0, fmt.Errorf("Update fm_role_limits failed:%#v ", err)
		}
		return rst.RowsAffected()
	}

	/*
		`UPDATE fm_role_user SET type=:type,delete_flag=:delete_flag, role_name=:role_name, login_name=:login_name
		WHERE id=:id`, role_user
	*/

	var OriginRole = erps.GcRole{}
	var upTotal = 0
	err := DBConn.Get(&OriginRole, `SELECT * FROM fm_role WHERE id=?`, roles.ID)
	if err != nil {
		return 0, fmt.Errorf("Origin Role fm_role failed:%#v ", err)
	}

	//优先更新 enabled
	if OriginRole.Enabled != roles.Enabled {
		enabledSql := fmt.Sprintf(`UPDATE fm_role SET enabled=%v WHERE id=%v`, roles.Enabled, roles.ID)
		rst, err := DBConn.Exec(enabledSql)
		logsys.Printf("gc role enabled update rst:%#v err:%#v", rst, err)
		if err != nil {
			return 0, fmt.Errorf("Update Role fm_role failed:%#v ", err)
		}
		return rst.RowsAffected()
	}

	if roles.Type != "" && OriginRole.Type != roles.Type {
		OriginRole.Type = roles.Type
		upTotal += 1
	}
	if roles.RoleName != "" && OriginRole.RoleName != roles.RoleName {
		OriginRole.RoleName = roles.RoleName
		upTotal += 1
	}
	if roles.PriceLimit != "" && OriginRole.PriceLimit != roles.PriceLimit {
		OriginRole.PriceLimit = roles.PriceLimit
		upTotal += 1
	}
	if roles.Description != "" && OriginRole.Description != roles.Description {
		OriginRole.Description = roles.Description
		upTotal += 1
	}
	if roles.Sort != "" && OriginRole.Sort != roles.Sort {
		OriginRole.Sort = roles.Sort
		upTotal += 1
	}

	OriginRole.CreatedAt = roles.CreatedAt
	rst, err := DBConn.NamedExec(`UPDATE fm_role SET
		type=:type,
		role_name=:role_name,
		price_limit=:price_limit,
		created_time=:created_time, 
		description=:description, 
		sort=:sort
		WHERE id=:id`, OriginRole)

	if rst == nil {
		return 0, fmt.Errorf("Update fm_role failed:%#v ", err)
	}

	total, err = rst.RowsAffected()
	logsys.Printf("new fm_role insert:%#v, with db rst:%v err:%v\n", roles, total, err)
	if err == nil {

		return total, err
	}

	return 0, fmt.Errorf("Update fm_role failed:%#v ", err)

}

// 新增后台管理 [角色] 列表
func (ut *DepartmentUtil) InsertRole(user string, roles *erps.GcRoleParam) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	if roles.CreateTime == "" {
		roles.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if roles.DeleteFlag == "" {
		roles.DeleteFlag = "0" //不删除
	}

	if roles.Type == "" {
		roles.Type = "RoleFunction" // 角色功能模块
	}

	if roles.PriceLimit == "" {
		roles.PriceLimit = "0" //不屏蔽
	}

	if roles.Sort == "" {
		roles.Sort = "0" //不排序
	}

	baseInsertSql := fmt.Sprintf(`INSERT INTO fm_role ( 
		type,
		role_name,
		price_limit, 
		description,
		enabled,
		sort,
		tenant_id) VALUES ('%v','%v','%v','%v','%v','%v',%v)`,
		roles.Type, roles.RoleName, roles.PriceLimit, roles.Description, 1, roles.Sort,
		fmt.Sprintf("(UUID_TO_BIN('%v'))", uuid.MustParse(roles.TenantId)))
	rst, err := DBConn.Exec(baseInsertSql)

	logsys.Printf("new fm_role insert:%#v, with db rst:%v err:%v\n", baseInsertSql, rst, err)
	if err != nil {
		return total, err
	}

	return rst.LastInsertId()
}

// 更新菜单权限
func (ut *DepartmentUtil) GeMenuUpdate(name string, qmenu []*erps.GcMenu) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	// admin 查询全部菜单数据，具有编辑权限
	var (
		DBConn = databases.DBMysql
	)
	if name == "admin" {
		for _, v := range qmenu {
			rst, err := DBConn.NamedExec(`UPDATE fm_menu SET delete_flag=:delete_flag, description=:description, level=:level
        WHERE menu_number=:menu_number`, v)
			logsys.Printf("from db select:%#v  err:%v ", rst, err)
			if err != nil {
				return nil, err
			}
		}
		return len(qmenu), nil
	}

	return len(qmenu), fmt.Errorf("nothing done")
}

// 获取后台管理 [权限组] 列表 doing
func (ut *DepartmentUtil) GeMenuLimitsList(login_name, tenant_id string) (*erps.GcMenuList, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = ""
		total   int
		menudep = erps.GcMenuList{}
		DBConn  = databases.DBMysql
	)

	if login_name == "" {
		return nil, fmt.Errorf("登陆认证信息 %v 错误", login_name)
	}

	totalSql := "SELECT COUNT(*) AS Total FROM `fm_menu`; "

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v", err, totalSql)
		return &menudep, err
	}

	// var tenFilter string
	// if tenant_id != "" {
	// 	tenFilter = fmt.Sprintf(" AND BIN_TO_UUID(mt.tenant_id)  = '%v' ", tenant_id)
	// }
	//根据用户名 筛选菜单名称列表
	// baseSql = fmt.Sprintf(`SELECT distinct b.menu_name AS menu_name,
	// b.menu_number AS menu_number,b.leader_id as leader_id, b.level as level,
	// b.description AS description,b.menu_type AS menu_type
	// FROM fm_menu b
	// left JOIN  fm_role_limits c ON c.btn_id = b.id
	// left JOIN  fm_limits t ON t.limits_id = c.limits_id
	// left JOIN  fm_role_user d ON d.role_name = c.role_name

	// left JOIN  fm_user u ON ( d.login_name = u.login_name)
	// left JOIN fm_tenant mt ON BIN_TO_UUID(mt.tenant_id) = BIN_TO_UUID(t.tenant_id)
	// WHERE  b.delete_flag != '1'  AND d.login_name = '%v'  %v
	//  ORDER BY b.menu_number;`, login_name, tenFilter)

	baseSql = fmt.Sprintf(`
	
		WITH matched_menus AS (
			-- 原始匹配相关功能点（主或子）
			SELECT distinct b.menu_name AS menu_name, b.id AS id, b.menu_number AS menu_number,b.leader_id as leader_id, b.level as level, b.description AS description,b.menu_type AS menu_type 
			FROM fm_menu b 
			left JOIN  fm_role_limits c ON c.btn_id = b.id 
			left JOIN  fm_limits t ON t.limits_id = c.limits_id 
		left JOIN  fm_role_user d ON d.role_name = c.role_name  
		left JOIN  fm_user u ON ( d.login_name = u.login_name)
		left JOIN fm_tenant mt ON BIN_TO_UUID(mt.tenant_id) = BIN_TO_UUID(t.tenant_id) 
		WHERE   b.delete_flag != '1'  AND d.login_name = '%v'    
			),
			missing_leaders AS (
		SELECT distinct b.menu_name AS menu_name,b.id AS id,  b.menu_number AS menu_number,b.leader_id as leader_id, b.level as level, b.description AS description,b.menu_type AS menu_type 
			FROM fm_menu sub
			JOIN fm_menu b ON b.menu_number = sub.leader_id
			WHERE sub.leader_id IS NOT NULL
			AND sub.menu_number IN (SELECT menu_number FROM matched_menus)
			AND b.menu_number NOT IN (SELECT menu_number FROM matched_menus)
			)
			-- 合并原始查询和缺失的主功能
		SELECT * FROM matched_menus
		UNION
		SELECT * FROM missing_leaders
		ORDER BY menu_number;`, login_name)

	err = DBConn.Select(&menudep.Data, baseSql)
	logsys.Printf("menu limit from db select:%#v  data:%#v err:%v ", baseSql, len(menudep.Data), err)
	if err != nil {
		return nil, err
	}

	menudep.Total = uint(total)
	menudep.Size = uint(len(menudep.Data))
	return &menudep, nil
}

// 根据部门 获取用户信息列表
func (ut *UserUtil) GetUsersByDepName(dep_name string) (*erps.UserList, error) {

	var (
		DBConn = databases.DBMysql
	)
	if dep_name == "" {
		return nil, fmt.Errorf("Get users list by department name fail:%v", dep_name)
	}
	// 查询数据 You can also get a single result, a la QueryRow

	usep := erps.UserList{}
	err := DBConn.Select(&usep.Data, BaseUserSql+"  WHERE department=?", dep_name)

	logsys.Printf("from db get user:%#v err:%v \n", usep, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil

}

// 根据 ctx 上下文获取用户信息
func (ut *UserUtil) GetUserByCtx(ctx *gin.Context) (*erps.FmUser, error) {

	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) <= len(configs.AuthHeader) {
		return nil, fmt.Errorf("Authorization MISSING")
	}
	tokenString := authHeader[len(configs.AuthHeader):]
	//根据token获取用户id
	tk, err := utils.ValidateToken(tokenString)
	if err != nil {
		err := fmt.Errorf("Token is invalid %#v", err)
		return nil, err
	}

	userId, _, err := utils.GetInfoFromToken(tk, "user_id")
	if err != nil {
		err := fmt.Errorf("user_id of Token error:  %#v", err)
		return nil, err
	}

	// 由id获取用户信息
	gu, err := ut.GetUserById(fmt.Sprintf("%v", userId))
	if err != nil {

		err := fmt.Errorf("user does not exit %#v", err)
		return nil, err
	}

	return gu, nil
}

// 根据id获取用户信息
func (ut *UserUtil) GetUserById(ids string) (*erps.FmUser, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		usep   = erps.FmUser{}
		DBConn = databases.DBMysql
		err    error
	)
	if ids == "" {
		return nil, fmt.Errorf("参数账号id为空")
	}

	filter := fmt.Sprintf("  WHERE id='%v' ", ids)
	baseSql := BaseUserSql + filter
	err = DBConn.Get(&usep, baseSql)
	logsys.Printf("from db get user:%#v err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil

}

// 根据登陆名获取用户信息
func (ut *UserUtil) GetUserByLoginName(login_name string) (*erps.FmUser, error) {

	var (
		usep   = erps.FmUser{}
		DBConn = databases.DBMysql
	)
	// 查询单个数据 You can also get a single result, a la QueryRow

	err := DBConn.Get(&usep, BaseUserSql+"  WHERE login_name=?", login_name)

	logsys.Printf("from db get:%v user:%#v err:%v \n", login_name, usep, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil

}

// 根据 电话 获取用户信息
func (ut *UserUtil) GetUserByPhone(phonenum string) (*erps.FmUser, error) {

	var (
		usep   = erps.FmUser{}
		DBConn = databases.DBMysql
	)
	// 查询单个数据 You can also get a single result, a la QueryRow

	err := DBConn.Get(&usep, BaseUserSql+"  WHERE phonenum=?", phonenum)

	logsys.Printf("from db get by phonenum:%v user:%#v err:%v \n", phonenum, usep, err)
	if err != nil {
		return nil, err
	}

	return &usep, nil
}

// 根据 用户名 或 电话 或 登录名 获取用户信息
func (ut *UserUtil) GetUserByInfo(LoginName, UserName, Phonenum string) (*erps.FmUser, error) {

	var (
		UserInfo   = []erps.FmUser{}
		UserBase   = BaseUserSql
		UserFilter = "WHERE "
		LIMITS     = "ORDER BY id DESC LIMIT 1"
		DBConn     = databases.DBMysql
	)

	if LoginName == "" && UserName == "" && Phonenum == "" {
		return nil, fmt.Errorf("参数错误，用户查询时,login_name,username,phonenum不能同时为空.")
	}

	// 查询条件
	if LoginName != "" {
		UserFilter += fmt.Sprintf("login_name = '%v'", LoginName)
	}

	if UserName != "" {
		if len(UserFilter) <= 6 {
			UserFilter += fmt.Sprintf("username = '%v'", UserName)
		} else {
			UserFilter += fmt.Sprintf("AND username = '%v'", UserName)
		}
	}

	if Phonenum != "" {
		if len(UserFilter) <= 6 {
			UserFilter += fmt.Sprintf("phonenum = '%v'", Phonenum)
		} else {
			UserFilter += fmt.Sprintf("AND phonenum = '%v'", Phonenum)
		}
	}

	if len(UserFilter) > 6 {
		UserBase += fmt.Sprintf(" %v %v", UserFilter, LIMITS)
		err := DBConn.Select(&UserInfo, UserBase)

		logsys.Printf("err:%#v userSql:%v  \n", err, UserBase)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("参数错误，用户查询时,login_name:%v,username:%v,phonenum:%v 无法同时匹配.", LoginName, UserName, Phonenum)
	}

	return &UserInfo[0], nil
}

// 根据 职位 获取批量用户信息
func (ut *UserUtil) GetUsersByPostion(Position string) (*string, error) {

	var (
		UserInfo   = []erps.FmUser{}
		UserBase   = BaseUserSql
		UserFilter = "WHERE "
		LIMITS     = "ORDER BY id DESC"
		DBConn     = databases.DBMysql
	)

	if Position == "" {
		return nil, fmt.Errorf("参数错误，根据职位进行用户查询时,Position不能同时为空")
	}

	// 查询条件
	if Position != "" {
		UserFilter += fmt.Sprintf("position = '%v'", Position)
	}

	if len(UserFilter) > 6 {
		UserBase += fmt.Sprintf(" %v %v", UserFilter, LIMITS)
		err := DBConn.Select(&UserInfo, UserBase)

		logsys.Printf("err:%#v userSql:%v  \n", err, UserBase)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("参数错误，用户查询时,postion:%v 无法同时匹配", Position)
	}

	var retStr string = ""
	if len(UserInfo) > 0 {
		for _, user := range UserInfo {
			if len(retStr) == 0 {
				retStr += fmt.Sprintf("'%v'", user.Id)
			}
			retStr += fmt.Sprintf(", '%v'", user.Id)
		}
	}
	return &retStr, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type BatchUtil struct{}

func NewPreBatchOrderUtil() *BatchUtil {
	return &BatchUtil{}
}
