// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package trade

//Trade Print
type TradePrint struct {
	Id           int    `db:"id" json:"id" form:"id"`
	DesignNo     string `db:"design_no" json:"design_no" form:"design_no"`
	Color        string `db:"color" json:"color"  form:"color"`
	Length       string `db:"length" json:"length"  form:"length"`
	LotNo        string `db:"lot_no" json:"lot_no"  form:"lot_no"`
	RollNo       string `db:"roll_no" json:"roll_no"  form:"roll_no"`
	CreateSerial string `db:"create_serial" json:"create_serial"  form:"create_serial"`
	CreatedAt    string `db:"created_time" json:"created_time" from:"created_time"`
	TenantId     string `db:"tenant_id" json:"tenant_id" form:"tenant_id"`
	DeleteFlag   string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`
}

//多条同时打印
type TradePrintList struct {
	Total uint         `db:"total" json:"total"`
	Size  uint         `db:"size" json:"size"`
	Data  []TradePrint `db:"data" json:"data" form:"data"`
}

// 打印机调用接口
type DCloudPrint struct {
	Code       uint   `json:"code"`
	Data       string `json:"data,omitempty"`
	Message    string `json:"message,omitempty"`
	Attachment string `json:"attachment,omitempty"`
}
