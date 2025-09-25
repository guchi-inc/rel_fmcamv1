// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package models

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

// 在此处注册查询字段,
var (
	//表名 和 前缀，fm_supplier
	TP_fm_supplier        = map[string][]string{"fm_supplier": {"phone_num", "supplier", "email", "contacts"}}
	TP_fm_user            = []string{`id`, `username`, "login_name", `leader_flag`, `position`, `department`, "email", `phonenum`, `ismanager`, `isystem`, `enabled`, `description`, `ethnicity`, `gender`, `local`, `localhost`, `m2_localhost`, `member_id`, `leader_id`, `device_time`, "BIN_TO_UUID(tenant_id) as tenant_id", `created_time`, `updated_time`, `deleted_time`}
	TP_fm_material_extend = map[string][]string{"fm_material_extend": {"bar_code", "material_id", "delete_flag", "tenant_id"}}
	TP_fm_depot_head      = map[string][]string{"fm_depot_head": {"source", "creator", "type"}}
	TP_printer            = map[string][]string{"tmp_backup": {"design_no", "date", "length", "color", "fact_name", "type", "grammage", "pre_batch"},
		"tmp_default": {"design_no", "date", "length", "color"}}
)

// 在此处注册全部表名，作为管理普通用户消息订阅的依据
var (
	Topics = map[string][]string{"account": {"fm_account", "fm_account_head", "fm_account_item"},
		"depot":       {"fm_depot", "fm_depot_head", "fm_depot_item"},
		"device_func": {"fm_devices", "fm_function", "fm_in_out_item"},
		"material": {"fm_material", "fm_material_timetribute", "fm_material_category", "fm_material_current_stock",
			"fm_material_extend", "fm_material_initial_stock", "fm_material_property"},
		"system_user": {"fm_supplier", "fm_tenant", "fm_user", "fm_user_business"},
		// "msg_log":     {"fm_msg", "fm_log"},
		// "platform_role_supplier": {"fm_orga_user_rel", "fm_organization", "fm_person", "fm_platform_config", "fm_role",
		// "fm_sequence", "fm_serial_number"},
		// "system_config": {"fm_system_config", "fm_unit"},
	}
)

var (
	logger = log.New(os.Stdout, "INFO -", 18)
	Ctx    = context.Background()
)

// 所有redis 客户端，包括单点和集群
type Caches struct {
	Client       *redis.Client
	UniverClient redis.UniversalClient
}

func (db *Caches) GetConn() (*redis.Conn, context.Context) {
	return db.Client.Conn(Ctx), Ctx
}

type Storage struct {
	*sqlx.DB
}

// dev测试
func (s *Storage) Insert(t any) bool {

	return true
}

func (s *Storage) IsConnected() bool {
	return true
}

func NewStorages(db *sqlx.DB) *Storage {
	return &Storage{db}
}

// 全部关系型数据库
type DBMangers struct {
	DBMysql   *Storage
	OBSClient *minio.Client
	DBRedis   *Caches
}

// any 一定要为 指针, 给table 赋值名称
func CallFuncs(tys any) any {

	refVal := []reflect.Value{}

	ntys := reflect.ValueOf(tys).MethodByName("TableName").Call(refVal)

	// logger.Printf("CallFuncs Success:%#v   call ref valus:%#v tval:%#v, after isTack:%#v \n", m, refVals, tval, tys.IsTack)
	logger.Printf("new struct:%#v nrst now:%#v \n", tys, ntys)
	return tys

}
