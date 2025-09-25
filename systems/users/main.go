// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"database/sql"
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/erps"
	"fmcam/models/iotmodel"
	"fmcam/models/seclients"
	"fmcam/systems/sysutils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	rows      *sql.Rows
	UsersData = &erps.FmUser{}
	logger    = log.New(os.Stderr, "INFO -", 13)

	//用户部门 工具
	UserUtil = sysutils.NewDepartmentUtil() //用户工具

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

var (
	DefUserRepo = NewUserRepo()
)

// 账号管理体系
type AccountsRes struct{}

// 系统内部管理结构体
type UserRepository struct{}

func NewUserRepo() *UserRepository {
	return &UserRepository{}
}

// 将 fm_user 结构体映射到 三元组 结构中   Name, CName
func FmUserMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
	)

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("fm_user"))
	logger.Printf("columns of fm_user:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		logger.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	logger.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "username", CName: newCommentMap["username"]})
	descs = append(descs, erps.GcDesc{Name: "login_name", CName: newCommentMap["login_name"]})
	descs = append(descs, erps.GcDesc{Name: "leader_flag", CName: newCommentMap["leader_flag"]})

	descs = append(descs, erps.GcDesc{Name: "position", CName: newCommentMap["position"]})
	descs = append(descs, erps.GcDesc{Name: "department", CName: newCommentMap["department"]})
	descs = append(descs, erps.GcDesc{Name: "email", CName: newCommentMap["email"]})
	descs = append(descs, erps.GcDesc{Name: "phonenum", CName: newCommentMap["phonenum"]})

	descs = append(descs, erps.GcDesc{Name: "ismanager", CName: newCommentMap["ismanager"]})
	descs = append(descs, erps.GcDesc{Name: "status", CName: newCommentMap["status"]})

	descs = append(descs, erps.GcDesc{Name: "description", CName: newCommentMap["description"]})
	descs = append(descs, erps.GcDesc{Name: "local", CName: newCommentMap["local"]})

	descs = append(descs, erps.GcDesc{Name: "localhost", CName: newCommentMap["localhost"]})

	descs = append(descs, erps.GcDesc{Name: "member_id", CName: newCommentMap["member_id"]})
	descs = append(descs, erps.GcDesc{Name: "int_id", CName: newCommentMap["int_id"]})

	descs = append(descs, erps.GcDesc{Name: "device_id", CName: newCommentMap["device_id"]})
	return descs, nil
}

// 注册简易 [租户] 信息
func NewSimpleCustomer(login_name string, suppliers *seclients.TenantFull) (string, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		DBConn = databases.DBMysql
	)

	tn := time.Now().Local().In(configs.Loc)
	ts := timeutil.CSTLayoutString(&tn)

	suppliers.Creator = login_name
	suppliers.CreatedAt = ts

	// 生成一个 UUID，并转为 binary[16]
	// u := uuid.New()
	// tenantIDBytes, err := u.MarshalBinary()
	// if err != nil {
	// 	logger.Fatal("UUID 转换失败: ", err)
	// }
	// fmt.Println("Tenant UUID:", u)
	// fmt.Println("Tenant Binary:", hex.EncodeToString(tenantIDBytes))

	// 插入 SQL
	newuuid := uuid.New().String()
	fmt.Println("Tenant Binary:", newuuid)

	query := fmt.Sprintf(`
		INSERT INTO fm_tenant (supplier,contacts,phone_num,tenant_id,creator) VALUES ('%v', '%v','%v', %v, '%v')
		`, suppliers.Supplier, suppliers.Contacts, suppliers.PhoneNum, fmt.Sprintf("(UUID_TO_BIN('%v'))", newuuid), suppliers.Creator)

	//写入 执行
	rst, err := DBConn.Exec(query)
	logger.Printf("new fm_tenant:%v insert:%#v, with db rst:%#v err:%v\n", newuuid, query, rst, err)

	if err != nil {
		logger.Println("插入失败: ", err)
		return "", err
	}

	return newuuid, err
}
