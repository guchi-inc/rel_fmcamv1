package seclients

import "fmcam/models/erps"

type FieldDiff struct {
	Field string      `json:"field" db:"field"`
	Old   interface{} `json:"old" db:"old"`
	New   interface{} `json:"new" db:"new"`
}

type SQLLogResponse struct {
	ID           int         `json:"id" db:"id"`
	PK           *int        `json:"pk_value" db:"pk_value"` //在账户信息里 就是 用户ID
	Action       string      `json:"action" db:"action"`
	Table        string      `json:"table" db:"table"`
	Args         *string     `json:"args" db:"args"`
	CreatedAt    *string     `json:"created_time" db:"created_time"` //操作发生时间
	Changes      []FieldDiff `json:"changes" db:"changes"`
	ChangesCname *string     `json:"changes_cname" db:"changes_cname"` //日志事件
	Creator      *string     `json:"creator" db:"creator"`             //操作人
	Position     *string     `json:"position" db:"position"`           //用户类型
	UserName     *string     `json:"username" db:"username"`           //用户名称
	Name         *string     `json:"name" db:"name"`                   //设备名称
}

//设备日志， 账户日志  用户新增日志 等  sql_logs 表
type SQLLogPages struct {
	Data     []any         `json:"data" db:"data"`
	Total    int           `json:"total" db:"total"`
	Page     int           `json:"page" db:"page"`
	PageSize int           `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

type DeviceLogResponse struct {
	ID           int         `json:"id" db:"id"`
	PK           *int        `json:"pk_value" db:"pk_value"` //在账户信息里 就是 用户ID
	Action       string      `json:"action" db:"action"`
	Table        string      `json:"table" db:"table"`
	Args         *string     `json:"args" db:"args"`
	CreatedAt    *string     `json:"created_time" db:"created_time"` //操作发生时间
	Changes      []FieldDiff `json:"changes" db:"changes"`
	ChangesCname *string     `json:"changes_cname" db:"changes_cname"`   //日志事件
	Creator      *string     `json:"creator" db:"creator"`               //操作人
	DeviceAttr   *string     `json:"device_attr" db:"device_attr"`       //设备属性
	UserName     *string     `json:"username" db:"username"`             //用户名称
	Name         *string     `json:"name" db:"name"`                     //设备名称
	DeviceHardId *string     `json:"device_hard_id" db:"device_hard_id"` //设备硬件ID
	DeviceModel  *string     `json:"device_model" db:"device_model"`     //设备型号
}

//设备日志， 账户日志  用户新增日志 等  sql_logs 表
type DeviceLogPages struct {
	Data     []DeviceLogResponse `json:"data" db:"data"`
	Total    int                 `json:"total" db:"total"`
	Page     int                 `json:"page" db:"page"`
	PageSize int                 `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc       `json:"columns" db:"columns"`
}
