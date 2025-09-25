// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package sysutils

import (
	"database/sql"
	"fmcam/common/databases"
	"fmcam/models/erps"
	"fmt"
	"time"
)

// 根据 用户登陆名 获取后台管理 [角色组] 列表
// 三个条件 相互之间为且的关系
func (ut *DepartmentUtil) GcRoleUserList(Page, PageSize int, role_type string, login_name, role_name string) (*erps.GcRoleUserList, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT id, role_id, user_id,type, role_name, login_name,created_time,creator,delete_flag FROM  `fm_role_user` WHERE delete_flag != '1' "

		total   int
		useRole = erps.GcRoleUserList{}
		DBConn  = databases.DBMysql
		filter  = ""
	)

	//根据用户名 筛选菜单名称列表
	if login_name != "" {
		filter += fmt.Sprintf(` AND login_name = '%v'  `, login_name)
	}
	if role_name != "" {
		filter += fmt.Sprintf(` AND role_name like '%v'  `, role_name+"%")
	}
	if role_type != "" {
		filter += fmt.Sprintf(`  AND type like '%v'  `, role_type+"%")
	}

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_role_user` WHERE delete_flag != '1' " + filter

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)
		return &useRole, err
	}

	//数据
	Limits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)

	baseSql += filter
	baseSql += Limits
	err = DBConn.Select(&useRole.Data, baseSql)
	logsys.Printf("from db select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useRole.Total = uint(total)
	useRole.Size = uint(PageSize)
	if len(useRole.Data) > 0 {
		useRole.Columns = FiledOper.FmRoleUsersMap()
	}
	return &useRole, nil
}

// 新增后台管理 [角色/用户信息] 列表 doing
func (ut *DepartmentUtil) InsertRoleUsers(user string, role_user *erps.GcRoleUser) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	if role_user.CreateTime == "" {
		role_user.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if role_user.DeleteFlag == "" {
		role_user.DeleteFlag = "0"
	}

	if role_user.Creator == "" {
		role_user.Creator = user
	}

	recordStructs := []erps.GcRoleUser{
		*role_user,
	}

	//角色和账号必 属于同一个租户
	if !ut.SamceTenantUserAndRole(int(role_user.UserID), int(role_user.RoleID)) {
		return 0, fmt.Errorf("用户 %v 与 分配角色 %v 不同区域", role_user.UserID, role_user.RoleID)
	}

	//同步 名称 和 id
	var GcRoleInfo = erps.GcRole{}
	roleInfo, _ := ut.GetRoleInfoById(1, 10, int(role_user.RoleID), "")
	if roleInfo != nil {
		if len(roleInfo.Data) > 0 {
			GcRoleInfo = roleInfo.Data[0]
			if GcRoleInfo.RoleName != role_user.RoleName {
				role_user.RoleName = GcRoleInfo.RoleName
			}
		} else {
			return 0, fmt.Errorf("Role角色不存在:%v", role_user.RoleID)
		}
	}

	userInfo, _ := UserUtils.GetUserById(fmt.Sprintf("%v", role_user.UserID))
	if userInfo != nil {
		role_user.LoginName = userInfo.LoginName
	} else {
		return 0, fmt.Errorf("User账号不存在:%v", role_user.UserID)
	}

	//查询用户是否已经分配角色
	roleUserOrigin, err := ut.GcRoleUserByLoginName(int(role_user.UserID), "")
	fmt.Printf("用户是否存在角色？%#v\n", roleUserOrigin)
	if roleUserOrigin != nil && err == nil {
		if roleUserOrigin.DeleteFlag != "1" {
			return 0, fmt.Errorf("用户%v 已分配角色:%#v", role_user.LoginName, roleUserOrigin.RoleName)
		} else {
			delSql := fmt.Sprintf("DELETE FROM  `fm_role_user` WHERE (`login_name` = '%v' OR `user_id` = '%v' );", role_user.LoginName, role_user.UserID)
			DBConn.Exec(delSql)
		}
	}

	rst, err := DBConn.NamedExec(`INSERT INTO fm_role_user ( 
		role_id,
		user_id,
		type,
		role_name,
		login_name,
		created_time,
		creator)
        VALUES (:role_id, :user_id, :type,:role_name, :login_name,:created_time, :creator)`, recordStructs)

	logsys.Printf("new fm_role_user insert:%#v, with db rst:%v err:%v\n", recordStructs, rst, err)
	if err == nil {
		total, err = rst.LastInsertId()
		return total, err
	}

	return 0, fmt.Errorf("Insert fm_role_user failed:%v", err)
}

// 根据 用户登陆名 获取后台管理 [角色组] 列表
// 三个条件 相互之间为且的关系
func (ut *DepartmentUtil) GcRoleUserById(ids int) (*erps.GcRoleUser, error) {

	// 查询单个数据
	var (
		DBConn  = databases.DBMysql
		baseSql = fmt.Sprintf("SELECT id, type, role_name, login_name,created_time,creator,delete_flag FROM  `fm_role_user` WHERE id = '%v'", ids)

		usedep = []erps.GcRoleUser{}
	)

	err := DBConn.Select(&usedep, baseSql)
	logsys.Printf("from db select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(usedep) == 0 {
		return nil, fmt.Errorf("No data of user ids:%v", ids)
	}
	return &usedep[0], nil
}

// 根据登陆用户名查询。 一个用户账号只能在系统只有一个分配角色
func (ut *DepartmentUtil) GcRoleUserByLoginName(user_id int, login_name string) (*erps.GcRoleUser, error) {

	// 查询单个数据
	var (
		DBConn  = databases.DBMysql
		baseSql = fmt.Sprintf("SELECT id, role_id, user_id, type, role_name, login_name,created_time,creator,delete_flag FROM  `fm_role_user` WHERE login_name = '%v' OR user_id = '%v' ", login_name, user_id)

		usedep = []erps.GcRoleUser{}
	)

	err := DBConn.Select(&usedep, baseSql)
	logsys.Printf("from db select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}
	if len(usedep) >= 1 {
		return &usedep[0], nil
	}
	return nil, fmt.Errorf("No role data of user:%v", login_name)
}

// 新增后台管理 [角色组] 列表 doing
func (ut *DepartmentUtil) UpdateOneRoleUsers(role_users *erps.GcRoleUser, UpdateDeleteFlag, TenantId string) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		err    error
		rst    sql.Result
		DBConn = databases.DBMysql
	)

	if role_users == nil || role_users.ID == 0 {
		return 0, fmt.Errorf("数据不存在:%#v", role_users)
	}

	//更新删除标记
	if UpdateDeleteFlag != "" {

		delSql := fmt.Sprintf("UPDATE fm_role_user SET  delete_flag='%v'  WHERE id='%v'", UpdateDeleteFlag, role_users.ID)
		rst, err = DBConn.Exec(delSql)

		logsys.Printf("new fm_role_user insert:%#v, with db rst:%v err:%v\n", role_users, rst, err)

		return rst.RowsAffected()
	}

	roleUserOrigin, err := ut.GcRoleUserById(int(role_users.ID))
	if err != nil {
		return 0, fmt.Errorf("数据不存在:%v", role_users)
	}

	var Uptotal int
	if role_users.Type != "" && role_users.Type != roleUserOrigin.Type {
		roleUserOrigin.Type = role_users.Type
		Uptotal += 1
	}

	//角色和账号必 属于同一个租户
	if !ut.SamceTenantUserAndRole(int(role_users.UserID), int(role_users.RoleID)) {
		return 0, fmt.Errorf("用户 %v 与 分配角色 %v 不同区域", role_users.UserID, role_users.RoleID)
	}

	if role_users.RoleID != 0 && role_users.RoleID != roleUserOrigin.RoleID {
		roleUserOrigin.RoleID = role_users.RoleID
		Uptotal += 1

		var GcRoleInfo = erps.GcRole{}
		roleInfo, _ := ut.GetRoleInfoById(1, 10, int(role_users.RoleID), "")
		if roleInfo != nil {
			if len(roleInfo.Data) > 0 {
				GcRoleInfo = roleInfo.Data[0]
				roleUserOrigin.RoleName = GcRoleInfo.RoleName
				Uptotal += 1
			} else {
				return 0, fmt.Errorf("Role角色不存在:%v", role_users.RoleID)
			}
		}
	}

	if role_users.UserID != 0 && role_users.UserID != roleUserOrigin.UserID {

		roleUserOrigin.UserID = role_users.UserID

		userInfo, _ := UserUtils.GetUserById(fmt.Sprintf("%v", role_users.UserID))
		if userInfo != nil {
			roleUserOrigin.LoginName = userInfo.LoginName
			Uptotal += 1
		} else {
			return 0, fmt.Errorf("User账号不存在:%v", role_users.UserID)
		}

	}

	if Uptotal == 0 {
		return 0, fmt.Errorf("数据未改动:%v", role_users)
	}

	recordStructs := []erps.GcRoleUser{
		*roleUserOrigin,
	}

	rst, err = DBConn.NamedExec(`UPDATE fm_role_user SET 
	role_id=:role_id,
	user_id=:user_id,
	type=:type, 
	role_name=:role_name, 
	login_name=:login_name
	WHERE id=:id`, recordStructs)

	logsys.Printf("fm_role_user UPDATE:%#v, with db rst:%v err:%v\n", roleUserOrigin, rst, err)

	//出错，提前返回
	if err != nil {
		return 0, fmt.Errorf("Update  role user failed:%v", err)
	}

	return rst.RowsAffected()
}

// 批量更新后台管理 [角色组] 列表 doing
func (ut *DepartmentUtil) UpdateRoleUsers(role_users []erps.GcRoleUser) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		err    error
		rst    sql.Result
		DBConn = databases.DBMysql
	)

	if role_users == nil {
		return 0, fmt.Errorf("UPDATE role user len:%v", role_users)
	}

	rst, err = DBConn.NamedExec(`UPDATE fm_role_user SET 
	type=:type, 
	role_name=:role_name, 
	delete_flage=:delete_flag,
	login_name=:login_name
	WHERE id=:id`, role_users)

	logsys.Printf("new fm_role_user insert:%#v, with db rst:%v err:%v\n", role_users, rst, err)

	//出错，提前返回
	if err != nil {
		return 0, fmt.Errorf("Update  role user failed:%v", err)
	}

	if err == nil {
		total, err := rst.RowsAffected()
		return total, err
	}
	return 0, fmt.Errorf("Update  role user failed:%v", err)
}
