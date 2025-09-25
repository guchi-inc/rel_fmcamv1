// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

//打印时调用的参数列表
type Params struct {
	Field string `json:"field" form:"field"` //字段
	Value string `json:"value" form:"value"` //字段的值
	Unit  string `json:"unit" form:"unit"`   //字段对应的单位
}

//打印参数默认值，通过配置文件可以配置
type Printer struct {
	Default string `db:"default" json:"default" form:"default"`
	Backup  string `db:"backup" json:"backup" form:"backup"`
}

//[{\"design_no\":\"130611230711179647\",
//\"date\": \"2023-07-30 17:13:21\",\"length\":\"189\",\"color\":\"red\"}]
type TmpDefault struct {
	DesignNo string `db:"design_no" json:"design_no" form:"design_no"`
	Date     string `db:"date" json:"date" form:"date"`
	Length   string `db:"length" json:"length" form:"length"`
	Color    string `db:"color" json:"color" form:"color"`
}

// `
// curl --location --request POST "https://dlabelplugga.ctaiot.com:9443/api/print/send_data_print.json" ^
// --header "accessToken: eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkdWRpYW4tYXBwLXNlY3VyaXR5Iiwic3ViIjoie1wiYXBwSWRcIjpcImRlbW8tYXBwLWtleVwiLFwidXNlcklkXCI6XCIxXCJ9IiwiYXVkIjoidXNlciIsImV4cCI6NDgzOTY0Mjk4MiwiaWF0IjoxNjg2MDQyOTgyfQ.iqJ46KlYAz7om6M_gD6BArNrZgvpZ_ezyVXgnIxn6fc" ^
// --header "User-Agent: Apix.cn)" ^
// --header "Content-Type: application/x-www-form-urlencoded " ^
// --header "Accept: */*" ^
// --header "Host: dlabelplugga.ctaiot.com:9443" ^
// --header "Connection: keep-alive" ^
// --data-urlencode "equip=1" ^
// --data-urlencode "doc_no=1686989129754" ^
// --data-urlencode "print_data=[{\"design_no\":\"130611230711179647\",
// \"date\": \"2023-07-30 17:13:21\",\"length\":\"189\",\"color\":\"red\",
//\"fact_name\":\"New yang tech.\",\"type\":\"JIng A0703\",\"grammage\":\"82de/298.2\",
//\"pre_batch\":\"230730\"}]" ^
// --data-urlencode "copies_number=1" ^
// --data-urlencode "brand_par_id=57"
// `
type TmpBackUp struct {
	DesignNo string `db:"design_no" json:"design_no" form:"design_no"`
	Date     string `db:"date" json:"date" form:"date"`
	Length   string `db:"length" json:"length" form:"length"`
	Color    string `db:"color" json:"color" form:"color"`
	FactName string `db:"fact_name" json:"fact_name" form:"fact_name"`
	Type     string `db:"type" json:"type" form:"type"`
	Grammage string `db:"grammage" json:"grammage" form:"grammage"`
	PreBatch string `db:"pre_batch" json:"pre_batch" form:"pre_batch"`
}

//打印日志
type PrintTaskLog struct {
	ID            int64  `db:"id" json:"id" form:"id"`
	PrintNumber   string `db:"print_number" json:"print_number" form:"print_number"`       // '票据号',
	Username      string `db:"username" json:"username" form:"username"`                   //'业主姓名',
	PrintDatas    string `db:"print_datas" json:"print_datas" form:"print_datas"`          // '数据号',
	TotalCounts   string `db:"total_counts" json:"total_counts" form:"total_counts"`       // '计算总额',
	TotalRececive string `db:"total_rececive" json:"total_rececive" form:"total_rececive"` // '实收总额',
	Remarks       string `db:"remarks" json:"remarks" form:"remarks"`                      //'票据备注',
	Tiaokuan      string `db:"tiaokuan" json:"tiaokuan" form:"tiaokuan"`                   // '条款',
	Creator       string `db:"creator" json:"creator" form:"creator"`                      //操作人
	CreatedAt     string `db:"created_time" json:"created_time" form:"created_time"`       //'新增时间',
}

// 打印日志记录列表
type PrintTaskLogList struct {
	Total   int64          `json:"total" db:"total" form:"total"`
	Size    int64          `json:"size" db:"size" form:"size"`
	Data    []PrintTaskLog `json:"data" db:"data" form:"data"`
	Columns []GcDesc       `json:"columns" db:"columns" form:"columns"`
}
