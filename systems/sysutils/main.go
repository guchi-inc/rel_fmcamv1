// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package sysutils

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/models/erps"
	"fmcam/models/iotmodel"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
快捷工具结构体
*/
var (
	logsys  = log.New(os.Stdout, "system-util: ", 13)
	logFile = configs.NewLogger("system-util:")

	UserUtils       = &UserUtil{}
	FiledOper       = &SysOperator{}
	UserDefaultPath = "./default" //用户默认文件路径

	BaseTenantSql = `SELECT id,
	supplier,
	contacts,
	email,
	BIN_TO_UUID(tenant_id) as tenant_id,
	type,
	province,
	city,
	area,
	street,
	address,
	addr_code,
	full_address,
	fax,
	phone_num,
	telephone,
	tax_num,
	bank_name,
	account_number,
	sort,
	description,
	enabled,
	delete_flag,
	isystem,
	advance_in,
	begin_need_get,
	begin_need_pay,
	all_need_get,
	all_need_pay,
	tax_rate,
	creator,
	created_time,
	updated_time FROM fm_tenant `
)

type DepartmentUtil struct{}

func NewDepartmentUtil() *DepartmentUtil {
	return &DepartmentUtil{}
}

type UserUtil struct{}

func NewUserUtil() *UserUtil {
	return &UserUtil{}
}

// 客户端操作绑定 结构体
type SysOperator struct{}

func (u *SysOperator) DepartmentMap(full bool) []erps.GcDesc {

	var (
		DBConn  = databases.DBMysql
		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	baseSql += " WHERE table_name = " + "'" + "fm_department" + "';"

	err := DBConn.Select(&erpDescs, baseSql)
	logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "fm_department", baseSql, err)
	if err == nil {
		return erpDescs
	}

	var matOrderMap []iotmodel.SimpleStruct
	rst2 := DBConn.Select(&matOrderMap, configs.GetSchemaSql("fm_department"))
	fmt.Printf("columns of fm_production_order:%#v map struct:%#v\n", rst2, matOrderMap)

	var newCommentMap = map[string]string{}
	//赋值表注释
	for _, v := range matOrderMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "id", CName: newCommentMap["id"]})
	descs = append(descs, erps.GcDesc{Name: "dep_type", CName: newCommentMap["dep_type"]})
	descs = append(descs, erps.GcDesc{Name: "dep_name", CName: newCommentMap["dep_name"]})
	descs = append(descs, erps.GcDesc{Name: "manager_name", CName: newCommentMap["manager_name"]})
	if !full {
		return descs
	}
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})
	descs = append(descs, erps.GcDesc{Name: "delete_flag", CName: newCommentMap["delete_flag"]})
	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs
}

func (u *SysOperator) TeantCustomerMap(full bool) []erps.GcDesc {

	var (
		DBConn  = databases.DBMysql
		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	baseSql += " WHERE table_name = " + "'" + "fm_tenant" + "';"

	err := DBConn.Select(&erpDescs, baseSql)
	logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "fm_tenant", baseSql, err)
	if err == nil {
		return erpDescs
	}

	var matTenantMap []iotmodel.SimpleStruct
	rst2 := DBConn.Select(&matTenantMap, configs.GetSchemaSql("fm_tenant"))
	fmt.Printf("columns of fm_tenant:%#v map struct:%#v\n", rst2, matTenantMap)

	var newCommentMap = map[string]string{}
	//赋值表注释
	for _, v := range matTenantMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	var descs = []erps.GcDesc{}

	//descs = append(descs, erps.GcDesc{Name: "id", CName: newCommentMap["id"]})
	descs = append(descs, erps.GcDesc{Name: "supplier", CName: newCommentMap["supplier"]})
	descs = append(descs, erps.GcDesc{Name: "contacts", CName: newCommentMap["contacts"]})
	descs = append(descs, erps.GcDesc{Name: "email", CName: newCommentMap["email"]})
	descs = append(descs, erps.GcDesc{Name: "tenant_id", CName: newCommentMap["tenant_id"]})
	descs = append(descs, erps.GcDesc{Name: "type", CName: newCommentMap["type"]})
	//descs = append(descs, erps.GcDesc{Name: "enabled", CName: newCommentMap["enabled"]})
	descs = append(descs, erps.GcDesc{Name: "province", CName: newCommentMap["province"]})
	descs = append(descs, erps.GcDesc{Name: "city", CName: newCommentMap["city"]})
	descs = append(descs, erps.GcDesc{Name: "area", CName: newCommentMap["area"]})
	descs = append(descs, erps.GcDesc{Name: "street", CName: newCommentMap["street"]})
	descs = append(descs, erps.GcDesc{Name: "address", CName: newCommentMap["address"]})
	descs = append(descs, erps.GcDesc{Name: "addr_code", CName: newCommentMap["addr_code"]})
	descs = append(descs, erps.GcDesc{Name: "fax", CName: newCommentMap["fax"]})
	descs = append(descs, erps.GcDesc{Name: "phone_num", CName: newCommentMap["phone_num"]})
	descs = append(descs, erps.GcDesc{Name: "telephone", CName: newCommentMap["telephone"]})
	descs = append(descs, erps.GcDesc{Name: "tax_num", CName: newCommentMap["tax_num"]})
	descs = append(descs, erps.GcDesc{Name: "bank_name", CName: newCommentMap["bank_name"]})
	descs = append(descs, erps.GcDesc{Name: "account_number", CName: newCommentMap["account_number"]})
	descs = append(descs, erps.GcDesc{Name: "sort", CName: newCommentMap["sort"]})
	descs = append(descs, erps.GcDesc{Name: "isystem", CName: newCommentMap["isystem"]})
	descs = append(descs, erps.GcDesc{Name: "advance_in", CName: newCommentMap["advance_in"]})
	descs = append(descs, erps.GcDesc{Name: "begin_need_get", CName: newCommentMap["begin_need_get"]})
	descs = append(descs, erps.GcDesc{Name: "begin_need_pay", CName: newCommentMap["begin_need_pay"]})
	descs = append(descs, erps.GcDesc{Name: "all_need_get", CName: newCommentMap["all_need_get"]})
	descs = append(descs, erps.GcDesc{Name: "all_need_pay", CName: newCommentMap["all_need_pay"]})
	descs = append(descs, erps.GcDesc{Name: "tax_rate", CName: newCommentMap["tax_rate"]})
	descs = append(descs, erps.GcDesc{Name: "full_address", CName: newCommentMap["full_address"]})
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})

	descs = append(descs, erps.GcDesc{Name: "description", CName: newCommentMap["description"]})
	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})
	descs = append(descs, erps.GcDesc{Name: "updated_time", CName: newCommentMap["updated_time"]})

	return descs
}

// fm_limits
// 权限组
func (u *SysOperator) FmLimitsMap() []erps.GcDesc {

	var (
		DBConn  = databases.DBMysql
		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	baseSql += " WHERE table_name = " + "'" + "fm_limits" + "';"

	err := DBConn.Select(&erpDescs, baseSql)
	logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "fm_limits", baseSql, err)
	if err == nil {
		return erpDescs
	}

	var matTenantMap []iotmodel.SimpleStruct
	rst2 := DBConn.Select(&matTenantMap, configs.GetSchemaSql("fm_limits"))
	fmt.Printf("columns of fm_limits:%#v map struct:%#v\n", rst2, matTenantMap)

	var newCommentMap = map[string]string{}
	//赋值表注释
	for _, v := range matTenantMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "limits_id", CName: newCommentMap["limits_id"]})
	descs = append(descs, erps.GcDesc{Name: "tenant_id", CName: newCommentMap["tenant_id"]})
	descs = append(descs, erps.GcDesc{Name: "limits_type", CName: newCommentMap["limits_type"]})
	descs = append(descs, erps.GcDesc{Name: "limits_name", CName: newCommentMap["limits_name"]})
	descs = append(descs, erps.GcDesc{Name: "level", CName: newCommentMap["level"]})
	descs = append(descs, erps.GcDesc{Name: "description", CName: newCommentMap["description"]})

	descs = append(descs, erps.GcDesc{Name: "perm_mask", CName: newCommentMap["perm_mask"]})

	return descs
}

// 权限组和角色
func (u *SysOperator) FmRoleLimitsMap() []erps.GcDesc {

	var (
		DBConn  = databases.DBMysql
		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	baseSql += " WHERE table_name = " + "'" + "fm_role_limits" + "';"

	err := DBConn.Select(&erpDescs, baseSql)
	logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "fm_role_limits", baseSql, err)
	if err == nil {
		return erpDescs
	}

	var matTenantMap []iotmodel.SimpleStruct
	rst2 := DBConn.Select(&matTenantMap, configs.GetSchemaSql("fm_role_limits"))
	fmt.Printf("columns of fm_role_limits:%#v map struct:%#v\n", rst2, matTenantMap)

	var newCommentMap = map[string]string{}
	//赋值表注释
	for _, v := range matTenantMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "btn_id", CName: newCommentMap["btn_id"]})
	descs = append(descs, erps.GcDesc{Name: "role_name", CName: newCommentMap["role_name"]})
	descs = append(descs, erps.GcDesc{Name: "limits_id", CName: newCommentMap["limits_id"]})
	descs = append(descs, erps.GcDesc{Name: "limits_name", CName: newCommentMap["limits_name"]})
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})

	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs
}

// 角色账号
func (u *SysOperator) FmRoleUsersMap() []erps.GcDesc {

	var (
		DBConn  = databases.DBMysql
		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	baseSql += " WHERE table_name = " + "'" + "fm_role_user" + "';"

	err := DBConn.Select(&erpDescs, baseSql)
	logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "fm_role_user", baseSql, err)
	if err == nil {
		return erpDescs
	}

	var matRoleUserMap []iotmodel.SimpleStruct
	rst2 := DBConn.Select(&matRoleUserMap, configs.GetSchemaSql("fm_role_user"))
	fmt.Printf("columns of fm_role_user:%#v map struct:%#v\n", rst2, matRoleUserMap)

	var newCommentMap = map[string]string{}
	//赋值表注释
	for _, v := range matRoleUserMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "type", CName: newCommentMap["type"]})
	descs = append(descs, erps.GcDesc{Name: "role_name", CName: newCommentMap["role_name"]})
	descs = append(descs, erps.GcDesc{Name: "login_name", CName: newCommentMap["login_name"]})
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})

	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs
}

// 角色列表
func (u *SysOperator) FmRolesMap() []erps.GcDesc {

	var (
		DBConn = databases.DBMysql

		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	baseSql += " WHERE table_name = " + "'" + "fm_role" + "';"

	err := DBConn.Select(&erpDescs, baseSql)
	logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "fm_role", baseSql, err)
	if err == nil {
		return erpDescs
	}

	var matRoleUserMap []iotmodel.SimpleStruct
	rst2 := DBConn.Select(&matRoleUserMap, configs.GetSchemaSql("fm_role"))
	fmt.Printf("columns of fm_role:%#v map struct:%#v\n", rst2, matRoleUserMap)

	var newCommentMap = map[string]string{}
	//赋值表注释
	for _, v := range matRoleUserMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "type", CName: newCommentMap["type"]})
	descs = append(descs, erps.GcDesc{Name: "role_name", CName: newCommentMap["role_name"]})
	descs = append(descs, erps.GcDesc{Name: "price_limit", CName: newCommentMap["price_limit"]})
	descs = append(descs, erps.GcDesc{Name: "description", CName: newCommentMap["description"]})
	descs = append(descs, erps.GcDesc{Name: "enabled", CName: newCommentMap["enabled"]})
	descs = append(descs, erps.GcDesc{Name: "tenant_id", CName: newCommentMap["tenant_id"]})

	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs
}
