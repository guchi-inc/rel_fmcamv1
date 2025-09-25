// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package sysaly

type QaInfo struct {
	Id    string `db:"id" json:"id" form:"id"`          // 主键ID
	Title string `db:"title" json:"title" form:"title"` //
	// setQas({ id,  title, created_time, introduction});
	Introduction string `db:"introduction" json:"introduction" form:"introduction"`
	CreatedAt    string `db:"created_time" json:"created_time" form:"created_time"`
}

type QaList struct {
	Total int      `db:"total" json:"total" form:"total"`
	Data  []QaInfo `db:"data" json:"data" form:"data"`
	Size  int      `db:"size" json:"size" form:"size"`
}

// 接口跟踪
type TraceInfo struct {
	Path   string `json:"path" form:"path"`     //服务路径
	Data   string `json:"data" form:"data"`     //错误发生时返回的数据
	Code   string `json:"code" form:"code"`     //错误代码
	Date   string `json:"date" form:"date"`     //日期
	User   string `json:"user" form:"user"`     //用户 token
	Brower string `json:"brower" form:"brower"` //Brower
	Timer  string `json:"timer" form:"timer"`   //耗时
}
