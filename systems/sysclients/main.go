package sysclients

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/storeobs"
	"fmcam/models/erps"
	"fmcam/models/iotmodel"
	"fmcam/systems/users"
	"log"
	"os"
	"strings"
)

var (
	sysdebug = log.New(os.Stdout, "sysclients INFO -", 13)
	logFile  = configs.NewLogger("system-clients:")

	UserService = users.DefUserRepo
	//图像文件
	OBSFile = databases.OBFileTask
	ObsApi  = storeobs.StorageObs
)

type SysCleint struct{}

// 将 fm_gov_province 结构体映射到 三元组 结构中   Name, CName
func FmGovProvinceMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
	)

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("fm_gov_province"))
	sysdebug.Printf("columns of fm_gov_province:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		sysdebug.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	sysdebug.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "code", CName: newCommentMap["code"]})
	descs = append(descs, erps.GcDesc{Name: "name", CName: newCommentMap["name"]})
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})
	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs, nil
}

// 将 fm_gov_city 结构体映射到 三元组 结构中   Name, CName
func FmGovCityMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
	)

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("fm_gov_city"))
	sysdebug.Printf("columns of fm_gov_city:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		sysdebug.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	sysdebug.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "code", CName: newCommentMap["code"]})
	descs = append(descs, erps.GcDesc{Name: "name", CName: newCommentMap["name"]})
	descs = append(descs, erps.GcDesc{Name: "province_code", CName: newCommentMap["province_code"]})
	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})
	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs, nil
}

// 将 fm_gov_city 结构体映射到 三元组 结构中   Name, CName
func FmGovAreaMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
	)

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("fm_gov_area"))
	sysdebug.Printf("columns of fm_gov_area:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		sysdebug.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	sysdebug.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "code", CName: newCommentMap["code"]})
	descs = append(descs, erps.GcDesc{Name: "name", CName: newCommentMap["name"]})
	descs = append(descs, erps.GcDesc{Name: "province_code", CName: newCommentMap["province_code"]})
	descs = append(descs, erps.GcDesc{Name: "city_code", CName: newCommentMap["city_code"]})

	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})
	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs, nil
}

// 将 fm_gov_city 结构体映射到 三元组 结构中   Name, CName
func FmGovStreetMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
	)

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("fm_gov_street"))
	sysdebug.Printf("columns of fm_gov_street:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		sysdebug.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	sysdebug.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "code", CName: newCommentMap["code"]})
	descs = append(descs, erps.GcDesc{Name: "name", CName: newCommentMap["name"]})
	descs = append(descs, erps.GcDesc{Name: "province_code", CName: newCommentMap["province_code"]})
	descs = append(descs, erps.GcDesc{Name: "city_code", CName: newCommentMap["city_code"]})
	descs = append(descs, erps.GcDesc{Name: "area_code", CName: newCommentMap["area_code"]})

	descs = append(descs, erps.GcDesc{Name: "creator", CName: newCommentMap["creator"]})
	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})

	return descs, nil
}

// 将 fm_field+matedata 结构体映射到 三元组 结构中   Name, CName
func FmFieldMetadataMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
		dbSql  = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required FROM field_metadata where table_name = 'field_metadata' ;"
	)

	var erpDescs = []erps.GcDesc{}
	err := DBConn.Select(&erpDescs, dbSql)
	sysdebug.Printf("all columns of field_metadata:%#v map erpDescs:%#v err:%#v\n", dbSql, erpDescs, err)

	if err == nil {
		return erpDescs, nil
	}

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("field_metadata"))
	sysdebug.Printf("columns of field_metadata:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		sysdebug.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	sysdebug.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "cname", CName: newCommentMap["cname"]})
	descs = append(descs, erps.GcDesc{Name: "name", CName: newCommentMap["name"]})
	descs = append(descs, erps.GcDesc{Name: "table_nanme", CName: newCommentMap["table_name"]})
	descs = append(descs, erps.GcDesc{Name: "data_type", CName: newCommentMap["data_type"]})
	descs = append(descs, erps.GcDesc{Name: "is_visible", CName: newCommentMap["is_visible"]})

	descs = append(descs, erps.GcDesc{Name: "is_searchable", CName: newCommentMap["is_searchable"]})
	descs = append(descs, erps.GcDesc{Name: "is_editable", CName: newCommentMap["is_editable"]})
	descs = append(descs, erps.GcDesc{Name: "is_required", CName: newCommentMap["is_required"]})
	descs = append(descs, erps.GcDesc{Name: "max_length", CName: newCommentMap["max_length"]})
	descs = append(descs, erps.GcDesc{Name: "default_value", CName: newCommentMap["default_value"]})

	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})
	descs = append(descs, erps.GcDesc{Name: "updated_time", CName: newCommentMap["updated_time"]})

	return descs, nil
}
