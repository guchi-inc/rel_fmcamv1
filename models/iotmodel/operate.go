// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package iotmodel

import "fmcam/models/erps"

/*

操作界面 交互数据
*/
//主表
type PmOrders struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	OrderNumber string `db:"order_number" json:"order_number" form:"order_number"` //  订单编号
	PreBatch    string `db:"pre_batch" json:"pre_batch" form:"pre_batch"`          //  预分批次
	Name        string `db:"name" json:"name" form:"name"`                         //  商品名(商品表的name)

	TypeName   string `db:"type_name" json:"type_name" form:"type_name"`       //  商品类型()
	TypeLevel  string `db:"type_level" json:"type_level" form:"type_level"`    //  类型级别()
	Category   string `db:"category" json:"category" form:"category"`          //工艺特征值
	Total      string `db:"total" json:"total" form:"total"`                   //总数
	PriceTotal string `db:"price_total" json:"price_total" form:"price_total"` //订单总价

	OrderType string `db:"order_type" json:"order_type" form:"order_type"` //  订单类型()

	Creator string `db:"creator" json:"creator" form:"creator"` //  操作员
	UnitId  string `db:"unit_id" json:"unit_id" form:"unit_id"` //  计量单位id

	TenantId  string `db:"tenant_id" json:"tenant_id" form:"tenant_id"`          //  租商id
	Stage     string `db:"stage" json:"stage" form:"stage"`                      //  阶段 4个
	Days      string `db:"days" json:"days" form:"days"`                         //  预期时长
	CreatedAt string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	UpdateAt  string `db:"updated_time" json:"updated_time" form:"updated_time"` //  创建时间

	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"` //  删除标记，0未删除，1删除

}

//订单信息 包括图片信息 描述
type PmOrderInfoDesc struct {
	Total   int64         `db:"total" json:"total" form:"total"`
	Size    int64         `db:"size" json:"size" form:"size"`
	Data    PmOrderInfo   `json:"data"`
	Columns []erps.GcDesc `json:"columns" form:"columns"`
}

//订单信息 包括图片信息
type PmOrderInfo struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	OrderNumber string `json:"order_number" db:"order_number"` //订单号
	Name        string `json:"name" db:"name"`                 //批次号
	ImgParent   string `json:"img_parent" db:"img_parent"`     //图名
	UserName    string `json:"username"  db:"username"`        //user 表 username
	LoginName   string `json:"login_name"  db:"login_name"`    //登陆账户
	PriceTotal  string `json:"price_total" db:"price_total"`   //金额
	Total       string `json:"total" db:"total"`               //总数
	Stage       string `json:"stage" db:"stage"  form:"stage"` //  阶段 4个

	Phonenum    string `json:"phonenum,omitempty" db:"phonenum" form:"phonenum"`    //电话号码
	Local       string `json:"local,omitempty" db:"local" form:"local"`             //楼栋
	Localhost   string `json:"localhost,omitempty" db:"localhost" form:"localhost"` //楼栋
	CreatedTime string `json:"created_time" db:"created_time" form:"created_time"`  //  创建时间

	Image *erps.ImgBase64 `json:"image" db:"image"` //图信息
}

//订单类型
//查询已分派记录列表结果封装
type PmOrderTypesList struct {
	Total int64          `db:"total" json:"total" form:"total"`
	Size  int64          `db:"size" json:"size" form:"size"`
	Data  []PmOrderTypes `db:"data" json:"data" form:"data"`

	Columns []erps.GcDesc `json:"columns" form:"columns"`
}

type PmOrderTypes struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	OrderType    string `db:"order_type" json:"order_type" form:"order_type"`          //  商品类型()
	TypeLevel    string `db:"type_level" json:"type_level" form:"type_level"`          //  类型级别()
	Category     string `db:"category" json:"category" form:"category"`                //特征值
	PriceDefault string `db:"price_default" json:"price_default" form:"price_default"` //服务费
	Creator      string `db:"creator" json:"creator" form:"creator"`                   //  操作员
	Enable       string `db:"enabled" json:"enabled" form:"enabled"`                   //  启用
	CreatedTime  string `db:"created_time" json:"created_time" form:"created_time"`    //  创建时间
	DeleteFlag   string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`       //  删除标记，0未删除，1删除

}

//订单分配 记录表
type GcAssignOrder struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	Creator     string `db:"creator" json:"creator" form:"creator"`                //  操作员
	OrderNumber string `db:"order_number" json:"order_number" form:"order_number"` //  订单编号
	PreBatch    string `db:"pre_batch" json:"pre_batch" form:"pre_batch"`          //  预分批次
	BatchId     string `db:"batch_id" json:"batch_id" form:"batch_id"`             //  批次号
	LoginName   string `db:"login_name" json:"login_name" form:"login_name"`       //  负责人登陆名

	AssignNumber string `db:"assign_number" json:"assign_number" form:"assign_number"` //  类型级别()
	Days         string `db:"days" json:"days" form:"days"`                            //  时长
	CreatedAt    string `db:"created_time" json:"created_time" form:"created_time"`    //  创建时间
	Stage        string `db:"stage" json:"stage" form:"stage"`                         //  阶段 4个，0初始 1生产中 2检验中 3完成
	Status       string `db:"status" json:"status" form:"status"`                      // 任务状态 0 正常 1 取消 2 禁用
	DeleteFlag   string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`       //  删除标记，0未删除，1删除
}

//查询已分派记录列表结果封装
type GcAssignList struct {
	Total int64           `db:"total" json:"total" form:"total"`
	Size  int64           `db:"size" json:"size" form:"size"`
	Data  []GcAssignOrder `db:"data" json:"data" form:"data"`

	Columns []erps.GcDesc `json:"columns" form:"columns"`
}

//查询订单列表结果封装
type GcOrderList struct {
	Total int64      `db:"total" json:"total" form:"total"`
	Size  int64      `db:"size" json:"size" form:"size"`
	Data  []PmOrders `db:"data" json:"data" form:"data"`

	Columns []erps.GcDesc `json:"columns" form:"columns"`
}

//查询参数，按部门统计查询 结构体， 聚合多个表的数据返回。
type GcDepartmentQuery struct {
	Page      int    `json:"page"`                               //当前分页号
	PageSize  int    `json:"page_size"`                          //分页大小
	DepName   string `json:"dep_name"`                           //部门名称
	MachineId string `json:"machine_id"`                         //机器号
	StartAt   string `json:"start_time"`                         //开始时间
	EndAt     string `json:"end_time"`                           //结束时间
	UserName  string `json:"username"`                           //user 表 username
	Phonenum  string `json:"phonenum,omitempty" form:"phonenum"` //电话号码，联系方式
}

//返回结构体: 部门 工单任务统计报表，  以瑕疵表 flaw 为主
type DepartmentBatchFlawResp struct {
	OrderNumber string `json:"order_number" db:"order_number" form:"order_number"` //订单号
	TypeName    string `json:"type_name" db:"type_name" form:"type_name"`          //品名
	TypeLevel   string `json:"type_level" db:"type_level" form:"type_level"`       //种类等级
	DepName     string `json:"dep_name" db:"dep_name" form:"dep_name"`             //部门名称

	LinkMachineId    string `json:"link_machine_id" db:"link_machine_id" form:"link_machine_id"`          //机器id
	Length           string `json:"length" db:"length" form:"length"`                                     //批次总长度
	Category         string `json:"category" db:"category" form:"category"`                               //品种特征，工艺克重
	StaffId          string `db:"staff_id" json:"staff_id" form:"staff_id"`                               //  瑕疵处生产负责的职员编号
	UserName         string `db:"username" json:"username" form:"username"`                               //user 表 username
	TotalBatchLength string `json:"total_batch_length" db:"total_batch_length" form:"total_batch_length"` //全部料号 批次 总长度
	TotalScore       string `json:"total_length" db:"total_length" form:"total_length"`                   //全部料号 总扣分
	TotalFlawLength  string `json:"total_flaw_length" db:"total_flaw_length" form:"total_flaw_length"`    //全部料号 不合格 标记 总长度
}

//部门统计：返回结构体
type TotalDepartmentFlawResp struct {
	Total int64                     `json:"total" `
	Size  int64                     `json:"size"`
	Data  []DepartmentBatchFlawResp `json:"data"`
}

//查询参数，细则查询结构体， 聚合多个表的数据返回。
type GcOrderQuery struct {
	Page        int    `json:"page"`                               //当前分页号
	PageSize    int    `json:"page_size"`                          //分页大小
	OrderNumber string `json:"order_number"`                       //订单号
	PreBatch    string `json:"pre_batch"`                          //批次号
	StartAt     string `json:"start_time"`                         //开始时间
	EndAt       string `json:"end_time"`                           //结束时间
	UserName    string `json:"username"`                           //user 表 username
	LoginName   string `json:"login_name"`                         //登陆账户
	Position    string `json:"position"`                           //职位
	Phonenum    string `json:"phonenum,omitempty" form:"phonenum"` //电话号码

}

//瑕疵信息记录
//返回参数：生产批次主表
type PmBatch struct {
	ID              int64  `db:"id" json:"id" form:"id"`                                              //  主键
	LinkOrderNumber string `db:"link_order_number" json:"link_order_number" form:"link_order_number"` //  订单编号-与订单表关联
	PreBatch        string `db:"pre_batch" json:"pre_batch" form:"pre_batch"`                         //  预分批次
	BatchId         string `db:"batch_id" json:"batch_id" form:"batch_id"`                            //  批次id
	TypeName        string `db:"type_name" json:"type_name" form:"type_name"`                         //  类型(生产批次)
	Name            string `db:"name" json:"name" form:"name"`                                        //  商品名
	CreatedTime     string `db:"created_time" json:"created_time" form:"created_time"`                //  创建时间
	ExpressName     string `db:"express_name" json:"express_name" form:"express_name"`                //  快递
	Stage           string `db:"stage" json:"stage" form:"stage"`                                     //  所处阶段
	Creator         string `db:"creator" json:"creator" form:"creator"`                               //  操作员
	Attributes      string `db:"attributes" json:"attributes" form:"attributes"`                      //  5
	DeleteFlag      string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`                   //  删除标记，0未删除，1删除
}

//批次信息查询返回结构体
type BatchList struct {
	Total int64     `db:"total" json:"total" form:"total"`
	Size  int64     `db:"size" json:"size" form:"size"`
	Data  []PmBatch `db:"data" json:"data" form:"data"`

	Columns []erps.GcDesc `db:"columns" json:"columns" form:"columns"` //字段标识
}

//料号 瑕疵主表
type GcBatchFlaw struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	OrderNumber string `db:"order_number" json:"order_number" form:"order_number"`    // 订单号
	LinkBatchId string `db:"link_batch_id" json:"link_batch_id" form:"link_batch_id"` //  订单编号-与订单表关联

	FlawName        string `db:"flaw_name" json:"flaw_name" form:"flaw_name"`                         //  批次id
	FlawId          string `db:"flaw_id" json:"flaw_id" form:"flaw_id"`                               //  瑕疵编号
	FlawNumber      string `db:"flaw_number" json:"flaw_number" form:"flaw_number"`                   //  瑕疵个数
	FlawLevel       string `db:"flaw_level" json:"flaw_level" form:"flaw_level"`                      //  瑕疵程度
	FlawScore       string `db:"flaw_score" json:"flaw_score" form:"flaw_score"`                      //  减分多少
	FlawTypes       string `db:"flaw_types" json:"flaw_types" form:"flaw_types"`                      //  瑕疵类型
	FlawStart       string `db:"flaw_start" json:"flaw_start" form:"flaw_start"`                      //  瑕疵开始位置
	FlawEnd         string `db:"flaw_end" json:"flaw_end" form:"flaw_end"`                            //  瑕疵结束位置
	FlawLeftPostion string `db:"flaw_left_postion" json:"flaw_left_postion" form:"flaw_left_postion"` //  瑕疵左边位置

	LinkFlawId     string `db:"link_flaw_id" json:"link_flaw_id" form:"link_flaw_id"`          //  关联标记 表示此标记来自谁的复制
	LinkMachineId  string `db:"link_machine_id" json:"link_machine_id" form:"link_machine_id"` //  关联标记 表示此标记来自谁的复制
	StaffId        string `db:"staff_id" json:"staff_id" form:"staff_id"`                      //  瑕疵处生产负责的职员编号
	CreateTime     string `db:"created_time" json:"created_time" form:"created_time"`          //  创建时间
	Creator        string `db:"creator" json:"creator" form:"creator"`                         //  操作员
	DeleteFlag     string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`             //  删除标记，0未删除，1删除
	FlawContinuous string `db:"flaw_continuous" json:"flaw_continuous" form:"flaw_continuous"` //  0 不连续  1 连续
	FlawFront      string `db:"flaw_front" json:"flaw_front" form:"flaw_front"`                //  瑕疵是否在正面 0 不在正面  1 在正面
	FlawHole       string `db:"flaw_hole" json:"flaw_hole" form:"flaw_hole"`                   //  瑕疵是否漏洞 0 不是漏洞  1 是漏洞

}

//批次信息查询返回结构体
type GcFlawTypesList struct {
	Total   int64         `db:"total" json:"total" form:"total"`
	Size    int64         `db:"size" json:"size" form:"size"`
	Data    []GcFlawTypes `db:"data" json:"data" form:"data"`
	Columns []erps.GcDesc `db:"columns" json:"columns" form:"columns"`
}

type GcFlawTypes struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	FlawName  string `db:"flaw_name" json:"flaw_name" form:"flaw_name"`    //  标记名
	FlawTypes string `db:"flaw_types" json:"flaw_types" form:"flaw_types"` //  标记类型
	FlawScore string `db:"flaw_score" json:"flaw_score" form:"flaw_score"` // 标记扣分
	FlawTimes string `db:"flaw_times" json:"flaw_times" form:"flaw_times"` //  标记使用次数

	CreateTime string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Creator    string `db:"creator" json:"creator" form:"creator"`                //  操作员
	Status     string `db:"status" json:"status" form:"status"`                   // 任务状态 0 正常 1 取消 2 禁用

	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"` //  删除标记，0未删除，1删除

}

// 查询返回结构体 批次对应的瑕疵列表
type GcBatchFlawList struct {
	Total int64         `db:"total" json:"total" form:"total"`
	Size  int64         `db:"size" json:"size" form:"size"`
	Data  []GcBatchFlaw `db:"data" json:"data" form:"data"`

	Columns []erps.GcDesc `db:"columns" json:"columns" form:"columns"`
}

// 删除参数结构体
type DeleteFlaw struct {
	ID     int64  `db:"id" json:"id" form:"id"`                //  主键
	FlawId string `db:"flaw_id" json:"flaw_id" form:"flaw_id"` //  瑕疵编号
}

type GcBatchFlawDelete struct {
	Data []DeleteFlaw `db:"data" json:"data" form:"data"`
}

//GcMachine 机器管理列表
type GcProductionMachine struct {
	ID          int64  `db:"id" json:"id" form:"id"`                                  //  主键
	MachineId   string `db:"machine_id" json:"machine_id" form:"machine_id"`          //  机器id
	Type        string `db:"type" json:"type" form:"type"`                            //  类型(机器类型)
	CreateTime  string `db:"created_time" json:"created_time" form:"created_time"`    //  创建时间
	Creator     string `db:"creator" json:"creator" form:"creator"`                   //  操作员
	LinkBatchId string `db:"link_batch_id" json:"link_batch_id" form:"link_batch_id"` //  批次编号-与批次表关联
	OrderNumber string `db:"order_number" json:"order_number" form:"order_number"`    // 订单号
	Owner       string `db:"owner" json:"owner" form:"owner"`                         // 负责人
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`       //  删除标记，0未删除，1删除
}

//查询返回列表
type GcProductionMachineList struct {
	Total int64                 `db:"total" json:"total" form:"total"`
	Size  int64                 `db:"size" json:"size" form:"size"`
	Data  []GcProductionMachine `db:"data" json:"data" form:"data"`
}

//GcCacheInfo 当前批检验信息
type GcCacheInfo struct {
	Color       string `db:"color" json:"color" form:"color"`                      //  颜色
	Category    string `db:"category" json:"category" form:"category"`             //  工艺值
	TypeName    string `db:"type_name" json:"type_name" form:"type_name"`          //  类型(订单产品类型)
	CreatedTime string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Creator     string `db:"creator" json:"creator" form:"creator"`                //  操作员
	BatchId     string `db:"batch_id" json:"batch_id" form:"batch_id"`             //  批次id
	OrderNumber string `db:"order_number" json:"order_number" form:"order_number"` // 订单号
	MachineId   string `db:"machine_id" json:"machine_id" form:"machine_id"`       //  机器id
	Quantity    string `db:"quantity" json:"quantity" form:"quantity"`             //  品质
	Length      string `db:"length" json:"length" form:"length"`                   //  总长

	Machines []GcProductionMachine `db:"machines" json:"machines" form:"machines"` //最近更新的机器列表 5个
	Columns  []erps.GcDesc         `json:"columns" db:"columns"`                   //展示字段说明
}

// 查询字段信息 返回结构体
type SimpleStruct struct {
	Field       string  `db:"field" json:"field"`
	Description string  `db:"description" json:"description"`
	Type        *string `db:"type" json:"type"`
	Keys        *string `db:"keys" json:"keys"`
}

// 前端操作 菜单列表
type GcOperateMenuList struct {
	Total   int64           `db:"total" json:"total" form:"total"`
	Size    int64           `db:"size" json:"size" form:"size"`
	Data    []GcOperateMenu `db:"data" json:"data" form:"data"`
	Columns []erps.GcDesc   `json:"columns" db:"columns"` //展示字段说明

}

//前端 单个操作菜单表 查询返回
type GcOperateMenu struct {
	ID          int64  `db:"id" json:"id" form:"id"`                               //  主键
	MenuType    string `db:"menu_type" json:"menu_type" form:"menu_type"`          //  菜单类型 业务 b，管理 m
	MenuName    string `db:"menu_name" json:"menu_name" form:"menu_name"`          //  菜单名
	MenuNumber  string `db:"menu_number" json:"menu_number" form:"menu_number"`    //  菜单编号
	CreateTime  string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Level       string `db:"level" json:"level" form:"level"`                      //  默认权限级别 low,middle,high
	Description string `db:"description" json:"description" form:"description"`    //  描述
	LeaderId    string `db:"leader_id" json:"leader_id" form:"leader_id"`          // 上级菜单id
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
	// Path        string `db:"path" json:"path" form:"path"`                      //  接口路径默认空

}

//瑕疵信息记录
//返回参数：生产批次主表
type PmAssess struct {
	ID              int64  `db:"id" json:"id" form:"id"`                                              //  主键
	LinkOrderNumber string `db:"link_order_number" json:"link_order_number" form:"link_order_number"` //  订单编号-与订单表关联
	ImgParent       string `db:"img_parent" json:"img_parent" form:"img_parent"`                      //  订单相关图片名
	Description     string `db:"description" json:"description" form:"description"`                   //  描述
	CreatedTime     string `db:"created_time" json:"created_time" form:"created_time"`                //  创建时间
	Creator         string `db:"creator" json:"creator" form:"creator"`                               //  操作员
	Attributes      string `db:"attributes" json:"attributes" form:"attributes"`                      //  5
	DeleteFlag      string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`                   //  删除标记，0未删除，1删除
}

//批次信息查询返回结构体
type PmAssessList struct {
	Total int64      `db:"total" json:"total" form:"total"`
	Size  int64      `db:"size" json:"size" form:"size"`
	Data  []PmAssess `db:"data" json:"data" form:"data"`

	Columns []erps.GcDesc `db:"columns" json:"columns" form:"columns"` //字段标识
}
