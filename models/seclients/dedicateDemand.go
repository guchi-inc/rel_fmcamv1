package seclients

import (
	"fmcam/models/erps"
	"fmcam/systems/genclients"
)

type DedicateServiceInfo struct {
	ID int64 `json:"id" db:"id"`
	// 酒店id
	WorkID *int64 `json:"work_id" db:"work_id"`
	// 已方客服
	Contacts *string `json:"contacts"  db:"contacts"`
	// 店名
	Supply *string `json:"supplier"  db:"supplier"`
	// 联系电话
	Phonenum *string `json:"phonenum"  db:"phonenum"`
	// 联系电邮
	Email *string `json:"email"  db:"email"`
	// 传真号
	Fax     *string `json:"fax"  db:"fax"`
	Creator *string `json:"creator"  db:"creator"`
	// 描述
	Description *string `json:"description"  db:"description"`
	DeleteFlag  *string `json:"delete_flag" db:"delete_flag"`
}

// 服务人员列表
type DedicatedServicePage struct {
	Data     []genclients.FmDedicatedServices `json:"data" db:"data"`
	Total    int                              `json:"total"`
	Page     int                              `json:"page"`
	PageSize int                              `json:"page_size"`
	Columns  []erps.GcDesc                    `json:"columns" db:"columns"`
}

// 留言建议 列表
type DemandsPage struct {
	Data     []genclients.FmDemands `json:"data" db:"data"`
	Total    int                    `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Columns  []erps.GcDesc          `json:"columns" db:"columns"`
}
