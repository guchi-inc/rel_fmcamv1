// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package iotmodel

import "fmcam/models/erps"

//查询参数，进度信息等。
type GcCacheProcing struct {
	Count int64 `json:"count"` //当前账户班次

	PreBatch    string `json:"pre_batch"`                        //预批次
	OrderNumber string `json:"order_number"`                     //订单号
	Category    string `json:"category" form:"category"`         //工艺克重
	TypeLevel   string `json:"type_level"`                       //等级
	Today       string `json:"today"`                            //进度
	TodayTotal  string `json:"today_total"`                      //今天总数
	Doing       string `json:"doing"`                            //已完成
	Total       string `json:"total"`                            //总数
	Company     string `json:"company"`                          //公司
	Address     string `json:"address,omitempty" form:"address"` // 联系地址
	MachineId   string `json:"machine_id"`                       //机器号
	SubId       string `json:"sub_id"`                           //子批次料号
	TypeName    string `json:"type_name"`                        //子批次料号
	MonthTotal  string `json:"month_total"`                      //月检总数

	Columns []erps.GcDesc `json:"columns" db:"columns"` //展示字段说明

}
