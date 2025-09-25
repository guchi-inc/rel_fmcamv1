package streams

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

var (
	syslog = log.New(os.Stdout, "system streams info -", 13)
)

type StreamRep struct{}

func NewStreamRepo() *StreamRep {
	return &StreamRep{}
}

// 将 fm_tasks 结构体映射到 三元组 结构中   Name, CName
func FeeStatusMap() ([]erps.GcDesc, error) {

	var (
		DBConn = databases.DBMysql
	)

	//查标记表 名称注释
	var matMap []iotmodel.SimpleStruct //[]structMap
	rstErr := DBConn.Select(&matMap, configs.GetSchemaSql("fm_tasks"))
	fmt.Printf("columns of fm_tasks:%#v map struct:%#v\n", rstErr, matMap)
	if rstErr != nil {
		return nil, rstErr
	}

	//将map更换为更安全的类型?
	var newCommentMap = map[string]string{}
	//赋值
	for _, v := range matMap {
		slis := strings.Split(v.Description, ",")
		fmt.Printf("field:%v slis:%v len:%v\n", v.Field, slis, len(slis))
		if len(slis) >= 1 {
			newCommentMap[v.Field] = slis[0]
		}
	}

	fmt.Printf("new comment map:%#v \n", newCommentMap)

	var descs = []erps.GcDesc{}

	descs = append(descs, erps.GcDesc{Name: "task_type", CName: newCommentMap["task_type"]})
	descs = append(descs, erps.GcDesc{Name: "task_name", CName: newCommentMap["task_name"]})
	descs = append(descs, erps.GcDesc{Name: "enabled", CName: "状态"})

	descs = append(descs, erps.GcDesc{Name: "page_id", CName: newCommentMap["page_id"]})
	descs = append(descs, erps.GcDesc{Name: "basic_unit", CName: newCommentMap["basic_unit"]})
	descs = append(descs, erps.GcDesc{Name: "basic_interval", CName: newCommentMap["basic_interval"]})
	descs = append(descs, erps.GcDesc{Name: "fit_housing", CName: newCommentMap["fit_housing"]})

	descs = append(descs, erps.GcDesc{Name: "created_time", CName: newCommentMap["created_time"]})
	descs = append(descs, erps.GcDesc{Name: "expiry_date", CName: newCommentMap["expiry_date"]})

	descs = append(descs, erps.GcDesc{Name: "updated_time", CName: newCommentMap["updated_time"]})
	descs = append(descs, erps.GcDesc{Name: "url", CName: newCommentMap["url"]})

	descs = append(descs, erps.GcDesc{Name: "play_url", CName: newCommentMap["play_url"]})

	descs = append(descs, erps.GcDesc{Name: "url_1", CName: newCommentMap["url_1"]})
	descs = append(descs, erps.GcDesc{Name: "url_2", CName: newCommentMap["url_2"]})

	descs = append(descs, erps.GcDesc{Name: "description", CName: newCommentMap["description"]})
	return descs, nil
}
