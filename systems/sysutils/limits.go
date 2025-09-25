package sysutils

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/erps"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LimitsUtil struct{}

// 获取后台管理 [菜单] 列表
func (mt *LimitsUtil) GetLimitsList(Page, PageSize int, name, LimitsId, TenantId string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT id, limits_id,BIN_TO_UUID(tenant_id) as tenant_id, limits_type,limits_name,level,description,perm_mask FROM `fm_limits` "

		total     int
		useLimits = erps.GcLimitsList{}
		DBConn    = databases.DBMysql
	)

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_limits`  "

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)
		return &useLimits, err
	}

	//条件
	Filters := " WHERE delete_flag != '1' "
	if LimitsId != "" {
		Filters += " AND limits_id LIKE '%" + LimitsId + "%' "
	}
	if TenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}

	Limits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	//admin 查询全部角色数据，具有编辑权限
	if name == "admin" {
		baseSql = "SELECT id, limits_id,BIN_TO_UUID(tenant_id) as tenant_id, limits_type,limits_name,level,description,perm_mask FROM  `fm_limits`  "
		baseSql += Filters
		baseSql += Limits
		err = DBConn.Select(&useLimits.Data, baseSql)
		logsys.Printf("from db admin baseSql:%#v  err:%v \n", baseSql, err)
		if err != nil {
			return nil, err
		}

		useLimits.Total = uint(total)
		useLimits.Size = uint(len(useLimits.Data))
		if len(useLimits.Data) > 0 {
			cm, _ := databases.FmGlobalMap("fm_limits", nil)
			useLimits.Columns = cm
		}
		return &useLimits, nil
	}

	baseSql += Filters
	baseSql += Limits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&useLimits.Data, baseSql)
	logsys.Printf("from limits select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useLimits.Total = uint(total)
	useLimits.Size = uint(PageSize)
	if len(useLimits.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_limits", nil)
		useLimits.Columns = cm
	}
	return &useLimits, nil
}

// 更新角色数据
func (mt *LimitsUtil) PostLimitsUpdate(limits *erps.GcLimits) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	if limits.DeleteFlag == "1" {
		delSql := fmt.Sprintf(" DELETE FROM fm_limits WHERE id = '%v' ", limits.ID)
		rst, err := DBConn.Exec(delSql)

		total, err = rst.RowsAffected()
		logsys.Printf("fm_role update:%#v, with db rst:%v err:%v\n", limits, total, err)
		if err != nil {
			return 0, err
		}

	}

	/*
		`UPDATE fm_role_user SET type=:type,delete_flag=:delete_flag, role_name=:role_name, login_name=:login_name
		WHERE id=:id`, role_user
	*/

	var OriginLimits = erps.GcLimits{}
	var upTotal = 0
	err := DBConn.Get(&OriginLimits, `SELECT * FROM fm_limits WHERE id=?`, limits.ID)
	if err != nil {
		return 0, err
	}

	if limits.LimitsId != "" && OriginLimits.LimitsId != limits.LimitsId {
		OriginLimits.LimitsId = limits.LimitsId
		upTotal += 1
	}
	if limits.TenantId != "" && OriginLimits.TenantId != limits.TenantId {
		OriginLimits.TenantId = limits.TenantId
		upTotal += 1
	}
	if limits.LimitsType != "" && OriginLimits.LimitsType != limits.LimitsType {
		OriginLimits.LimitsType = limits.LimitsType
		upTotal += 1
	}
	if limits.Description != "" && OriginLimits.Description != limits.Description {
		OriginLimits.Description = limits.Description
		upTotal += 1
	}
	if limits.LimitsName != "" && OriginLimits.LimitsName != limits.LimitsName {
		OriginLimits.LimitsName = limits.LimitsName
		upTotal += 1
	}
	if limits.Level != "" && OriginLimits.Level != limits.Level {
		OriginLimits.Level = limits.Level
		upTotal += 1
	}
	if limits.Description != "" && OriginLimits.Description != limits.Description {
		OriginLimits.Description = limits.Description
		upTotal += 1
	}
	if limits.PermMask != "" && OriginLimits.PermMask != limits.PermMask {
		OriginLimits.PermMask = limits.PermMask
		upTotal += 1
	}
	rst, err := DBConn.NamedExec(`UPDATE fm_limits SET
		limits_id=:limits_id,
		tenant_id=:tenant_id,
		limits_type=:limits_type,
		limits_name=:limits_name, 
		description=:description, 
		level=:level,
		perm_mask=:perm_mask
		WHERE id=:id`, OriginLimits)

	if rst == nil {
		return 0, fmt.Errorf("Update fm_limits failed:%v", err)
	}

	total, err = rst.RowsAffected()
	logsys.Printf("fm_limits update:%#v, with db rst:%v err:%v\n", limits, total, err)
	if err == nil {

		return total, err
	}

	return 0, fmt.Errorf("Update fm_role_limits failed:%#v\n", err)

}

// 新增后台管理 [资源组名] 信息
func (mt *LimitsUtil) LimitsNew(username string, limits *erps.GcLimits) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	if limits.DeleteFlag == "" {
		limits.DeleteFlag = "0" //不删除
	}

	LStructs := []erps.GcLimits{
		*limits,
	}

	rst, err := DBConn.NamedExec(`INSERT INTO fm_limits ( 
		limits_id,
		tenant_id,
		limits_type,
		limits_name,
		level,
		description,  
	    delete_flag)
        VALUES (:limits_id,:tenant_id, :limits_type,:limits_name, :level,
			:description, :delete_flag)`, LStructs)

	logsys.Printf("new fm_limits insert:%#v, with db rst:%v err:%v\n", LStructs, rst, err)
	if err != nil {
		return total, err
	}

	return rst.LastInsertId()
}

// 根据 角色名 获取后台管理 [角色权限关联组] 列表 doing
func (mt *LimitsUtil) GetRoleLimitsList(Page, PageSize, btn_id int, role_name, username, StartAt, EndAt, tenant_id string) (*erps.GcRoleLimitsList, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int
		usedep = erps.GcRoleLimitsList{}
		DBConn = databases.DBMysql
	)

	//条件
	Filter := " WHERE a.delete_flag != '1' "
	if btn_id != 0 {
		Filter += fmt.Sprintf(" AND a.btn_id = '%v' ", btn_id)
	}
	if role_name != "" {
		Filter += fmt.Sprintf(" AND a.role_name = '%v' ", role_name)
	}

	// var JoinFilter string
	if username != "" {
		Filter += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM fm_role_user d
			JOIN fm_user b ON b.login_name = d.login_name
			WHERE d.role_name = a.role_name
			  AND (b.username LIKE '%v' OR b.login_name = '%v')
		)`, username+"%", username)
	}
	if tenant_id != "" {
		Filter += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM fm_role r
			WHERE r.role_name = a.role_name
			  AND BIN_TO_UUID(r.tenant_id) = '%v'
		) `, uuid.MustParse(tenant_id))
	}
	if StartAt != "" && EndAt != "" {
		Filter += fmt.Sprintf(" AND (a.created_time between '%v' AND '%v') ", StartAt, EndAt)
	}

	totalSql := fmt.Sprintf(`SELECT COUNT(*) AS Total FROM  fm_role_limits  a    %v `, Filter)

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)
		return &usedep, err
	}

	//根据用户名 筛选菜单名称列表
	Limits := fmt.Sprintf("  ORDER BY a.id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)

	baseSql := fmt.Sprintf(`SELECT a.id,a.btn_id, a.role_name, a.limits_id, 
		a.limits_name,a.creator,a.created_time
		FROM   fm_role_limits a 
		%v
		`, Filter)
	baseSql += Limits
	err = DBConn.Select(&usedep.Data, baseSql)
	logsys.Printf("from db select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	usedep.Total = uint(total)
	usedep.Size = uint(PageSize)
	if len(usedep.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_role_limits", nil)
		usedep.Columns = cm
	}

	return &usedep, nil
}

// 根据 角色名 获取后台管理 [角色权限关联组] 列表 doing
func (mt *LimitsUtil) GetRoleBtnMenus(role_name string) (*erps.GcMenuList, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = `SELECT  
      	  		a.id AS id,
				a.menu_name,
				a.leader_id,
				a.description,
				a.menu_number
			FROM 
				fm_role_limits b
			JOIN 
				fm_menu a ON b.btn_id = a.id
			`

		usedep = erps.GcMenuList{}
		DBConn = databases.DBMysql
	)

	//条件
	Filter := " WHERE a.id > 0 "
	if role_name != "" {
		Filter += fmt.Sprintf("  AND  b.role_name = '%v'  ", role_name)
	}

	//根据用户名 筛选菜单名称列表
	Limits := " ORDER BY a.menu_number   "
	baseSql += Filter
	baseSql += Limits
	err := DBConn.Select(&usedep.Data, baseSql)
	logsys.Printf("from db select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	usedep.Total = uint(len(usedep.Data))
	usedep.Size = uint(len(usedep.Data))
	if len(usedep.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_menu", nil)
		usedep.Columns = cm
	}

	return &usedep, nil
}

// 新增后台管理 [权限组] 列表 doing
func (mt *LimitsUtil) InsertRoleLimits(login_name string, MRRecord *erps.FmRoleLimitNews) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		DBConn = databases.DBMysql
	)

	if len(MRRecord.BtnIdList) == 0 {
		return 0, fmt.Errorf("功能无法为空")
	}

	t0 := time.Now().Local().In(configs.Loc)
	ts := timeutil.CSTLayoutString(&t0)

	//数据转换
	recordStructs := []erps.GcRoleLimits{}
	for _, btnid := range MRRecord.BtnIdList {
		rms := erps.GcRoleLimits{BtnId: btnid,
			RoleName:   MRRecord.RoleName,
			LimitsId:   MRRecord.LimitsId,
			LimitsName: MRRecord.LimitsName,
			Creator:    login_name,
			CreateTime: ts}
		recordStructs = append(recordStructs, rms)
	}

	//批量写入 NamedExec 是基于 Prepare 和 Exec 的封装，不返回多条 ID
	rst, err := DBConn.NamedExec(`INSERT INTO fm_role_limits ( 
		btn_id,
		role_name,
		limits_id,
		limits_name,
		created_time,
		creator)
        VALUES (:btn_id,:role_name, :limits_id,:limits_name,:created_time, 
			:creator)`, recordStructs)

	logsys.Printf("new fm_role_limits insert:%#v, with db rst:%v err:%v\n", recordStructs, rst, err)
	if err != nil {
		return 0, err

	}

	return rst.LastInsertId()
}

// 删除资源组角色关联
func (mt *LimitsUtil) DelRoleLimits(role_name string, btn_id, id int64) (int64, error) {
	var (
		delSql = "DELETE FROM fm_role_limits "
		DBConn = databases.DBMysql
	)

	if id != 0 {
		delSql += fmt.Sprintf(" WHERE id = '%v' ", id)
	} else {
		if role_name == "" || btn_id == 0 {
			return 0, fmt.Errorf("角色资源信息不存在%v:%v", role_name, btn_id)
		} else {
			delSql += fmt.Sprintf(" WHERE role_name = '%v' AND btn_id = '%v' ", role_name, btn_id)
		}
	}
	rst := DBConn.MustExec(delSql)
	return rst.RowsAffected()
}

// 删除角色功能数据
func (mt *LimitsUtil) PostLimitRolesUpdate(roleLimits *erps.GcRoleLimits) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow

	if roleLimits.DeleteFlag == "1" {
		return mt.DelRoleLimits("", 0, roleLimits.ID)
	}
	return 0, fmt.Errorf("nothing to do")
}
