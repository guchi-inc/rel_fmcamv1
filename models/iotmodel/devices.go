package iotmodel

import "fmcam/models/erps"

//流组列表
type StreamTaskList struct {
	Total   int64           `db:"total" json:"total" form:"total"`
	Size    int64           `db:"size" json:"size" form:"size"`
	Data    []FmStreamTasks `json:"data"`
	Columns []erps.GcDesc   `json:"columns" form:"columns"`
}

//校验记录查询列表
type StreamLogList struct {
	Total   int64         `db:"total" json:"total" form:"total"`
	Size    int64         `db:"size" json:"size" form:"size"`
	Data    []FmStreamLog `json:"data"`
	Columns []erps.GcDesc `json:"columns" form:"columns"`
}

// 流组类型
type FmStreamTasks struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	TaskType string `db:"task_type" json:"task_type" form:"task_type"` //  流组类型
	TaskName string `db:"task_name" json:"task_name" form:"task_name"` //  流组任务类型

	PageId string `db:"page_id" json:"page_id" form:"page_id"` //  重点区域

	Enable string `db:"enabled" json:"enabled" form:"enabled"` // 缴纳状态 0 未交， 1 已缴

	Creator       string `db:"creator" json:"creator" form:"creator"`                      //  操作员
	BasicUnit     string `db:"basic_unit" json:"basic_unit" form:"basic_unit"`             //  单位,一层一个
	BasicInterval string `db:"basic_interval" json:"basic_interval" form:"basic_interval"` //  分,每分钟轮询一次
	FitHousing    string `db:"fit_housing" json:"fit_housing" form:"fit_housing"`          //  适用场景 酒店 园区
	RtspUrl       string `db:"url" json:"url" form:"url"`                                  //主帧流,(一般是stream_0)
	PlayUrl       string `db:"play_url" json:"play_url" form:"play_url"`                   //主帧流,(一般是stream_0)

	RtspUrlOne string `db:"url_1" json:"url_1" form:"url_1"` //低帧流,(一般是stream_0)
	RtspUrlTwo string `db:"url_2" json:"url_2" form:"url_2"` //性能流,(一般是stream_0)

	CreatedAt   string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
	UpdatedAt   string `db:"updated_time" json:"updated_time" form:"updated_time"` //  创建时间
	Description string `db:"description" json:"description" form:"description"`    //备注信息

}

// 校验记录
type FmStreamLog struct {
	ID     int64 `db:"id" json:"id" form:"id"`                //  主键
	TaskId int64 `db:"task_id" json:"task_id" form:"task_id"` //  流组任务id

	LoginName string `db:"login_name" json:"login_name" form:"login_name"` //登录用户名

	TaskType  string `db:"task_type" json:"task_type" form:"task_type"`          //  流组类型
	TaskName  string `db:"task_name" json:"task_name" form:"task_name"`          //  流组任务类型
	CreatedAt string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间

	ExpiryDate string `db:"expiry_date" json:"expiry_date" form:"expiry_date"` //  服务截止时间
	CashType   string `db:"cash_type" json:"cash_type" form:"cash_type"`       //缴费方式 现金，网银转账

	Description string `db:"description" json:"description" form:"description"` //备注信息

	Creator    string `db:"creator" json:"creator" form:"creator"`                //  操作员
	Enabled    string `db:"enabled" json:"enabled" form:"enabled"`                //  计算列，状态
	UpdateTime string `db:"updated_time" json:"updated_time" form:"updated_time"` //  更新时间

	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"` //  删除标记，0未删除，1删除

	Feed any `db:"feed" json:"feed" form:"feed"` //实际 缴纳金额

}
