// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package iotmodel

import "fmcam/models/erps"

// 订单数据表 返回结构体
type OrderResp struct {
	OrderNumber string `json:"order_number" db:"order_number"` //订单号
	BatchId     string `json:"batch_id" db:"batch_id"`         //批号
	UpdatedAt   string `json:"updated_time" db:"updated_time"` //上机时间
	CardNumber  string `json:"card_number" db:"card_number"`   //卡号
	TypeName    string `json:"type_name" db:"type_name"`       //品名
	MachineId   string `json:"machine_id" db:"machine_id"`     //机器id
	CraftType   string `json:"craft_type" db:"craft_type"`     //工艺类型  281.32
	Level       string `json:"level" db:"level"`               //品质等级
	Length      string `json:"Length" db:"length"`             //长度位置
	Category    string `json:"category" db:"category"`         //品种
}

// 工单任务统计报表， 返回结构体,  以瑕疵表 flaw 为主
type BatchFlawResp struct {
	OrderNumber string `json:"order_number" db:"order_number"`   //订单号
	LinkBatchId string `json:"link_batch_id" db:"link_batch_id"` //批号
	TypeName    string `json:"type_name" db:"type_name"`         //品名
	TypeLevel   string `json:"type_level" db:"type_level"`       //种类等级

	LinkMachineId string `json:"link_machine_id" db:"link_machine_id"` //机器id
	// CraftType  string `json:"craft_type" db:"craft_type"`                      //工艺类型  281.32
	// Level  string `json:"level" db:"level"`   //品质等级
	Length     string `json:"length" db:"length"`                                 //批次总长度
	Category   string `json:"category" db:"category"`                             //品种特征，工艺克重
	FlawId     string `json:"flaw_id" db:"flaw_id"`                               //标记id
	FlawScore  string `db:"flaw_score" json:"flaw_score" form:"flaw_score"`       //  减分多少
	FlawTypes  string `db:"flaw_types" json:"flaw_types" form:"flaw_types"`       //  瑕疵类型
	FlawName   string `db:"flaw_name" json:"flaw_name" form:"flaw_name"`          //  瑕疵名词 合格 不合格
	StaffId    string `db:"staff_id" json:"staff_id" form:"staff_id"`             //  瑕疵处生产负责的职员编号
	CreatedAt  string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
	Creator    string `db:"creator" json:"creator" form:"creator"`                //  操作员
	Stage      string `db:"stage" json:"stage" form:"stage"`                      //  状态，0初始，1 生产 2 检验 3 完成
	Days       string `db:"days" json:"days" form:"days"`                         //  天数
}

//返回结构体
type TotalBatchFlawResp struct {
	Total int64           `json:"total" `
	Size  int64           `json:"size"`
	Data  []BatchFlawResp `json:"data"`

	Columns []erps.GcDesc `json:"columns" db:"columns"` //展示字段说明
}

//分部门统计 数据
/*
gu.department AS department,
		gu.username AS username,
        COUNT(distinct b.link_order_number) AS order_total,
        COUNT(distinct a.link_batch_id) AS batch_total,
        COUNT(distinct a.creator) AS creator_total,
        COUNT(distinct a.staff_id) AS staff_total,
        SUM(a.flaw_score) AS score_total,
        SUM(b.length) AS total_length
*/
type DepartmentResp struct {
	DepName  string `json:"department" db:"department"`             //部门名
	UserName string `db:"username" json:"username" form:"username"` //user 表 username

	OrderTotal   string `json:"order_total" db:"order_total"`     //订单数
	BatchTotal   string `json:"batch_total" db:"batch_total"`     //料号数
	CreatorTotal string `json:"creator_total" db:"creator_total"` //检验操作者 人 数
	StaffTotal   string `json:"staff_total" db:"staff_total"`     //生产责任人 数
	ScoreTotal   string `json:"score_total" db:"score_total"`     //分数
	TotalLength  string `json:"total_length" db:"total_length"`   //料号总长
}

//分部门统计 返回结构体
type TotalDepartmentResp struct {
	Total int64            `json:"total" `
	Size  int64            `json:"size"`
	Data  []DepartmentResp `json:"data"`

	Columns []erps.GcDesc `json:"columns" db:"columns"` //展示字段说明
}
