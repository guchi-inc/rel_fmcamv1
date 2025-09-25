package alysys

import (
	"fmcam/common/configs"
	"fmcam/models/erps"
	"log"
	"os"
)

var (
	logsys         = log.New(os.Stdout, "INFO -", 13)
	AlysysOperator = AlySysOperator{}
	DefDateRange   = configs.DefDateRange
)

// 客户端操作绑定 结构体
type AlySysOperator struct{}

//组合 轨距查询 字段说明

// 将 captionlog 轨迹查询 结构体映射到 三元组 结构中   Name, CName
func FmCaptionLogMap() ([]erps.GcDesc, error) {

	var (
	// DBConn  = databases.DBMysql
	// baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required,data_type,max_length FROM field_metadata  "
	)

	var erpDescs = []erps.GcDesc{}
	// baseSql += " WHERE table_name = 'CaptureLogs'   "

	// limits := " ORDER BY id asc "
	// baseSql += limits
	// err := DBConn.Select(&erpDescs, baseSql)
	// logsys.Printf("field select of table:%v sql: %v ,err :%v \n", "CaptureLogs", baseSql, err)
	// if err != nil {
	// 	return nil, err
	// }
	erpDescs = append(erpDescs, erps.GcDesc{Name: "face_id", CName: "面部ID", IsVisible: true, IsSearchable: false})
	erpDescs = append(erpDescs, erps.GcDesc{Name: "device_id", CName: "设备ID", IsVisible: true, IsSearchable: false})
	erpDescs = append(erpDescs, erps.GcDesc{Name: "location", CName: "设备位置", IsVisible: true, IsSearchable: false})
	erpDescs = append(erpDescs, erps.GcDesc{Name: "capture_time", CName: "抓拍时间", IsVisible: true, IsSearchable: true})
	erpDescs = append(erpDescs, erps.GcDesc{Name: "capture_image_url", CName: "抓拍图像", IsVisible: true, IsSearchable: false})

	erpDescs = append(erpDescs, erps.GcDesc{Name: "id", CName: "抓拍号", IsVisible: false, IsSearchable: false})
	erpDescs = append(erpDescs, erps.GcDesc{Name: "room_id", CName: "房间号", IsVisible: false, IsSearchable: true})

	erpDescs = append(erpDescs, erps.GcDesc{Name: "id_card_number", CName: "证件号", IsVisible: false, IsSearchable: false})
	erpDescs = append(erpDescs, erps.GcDesc{Name: "sort", CName: "排序", IsVisible: false, IsSearchable: false})

	erpDescs = append(erpDescs, erps.GcDesc{Name: "matched_profile_id", CName: "人员号", IsVisible: false, IsSearchable: true})

	return erpDescs, nil
}
