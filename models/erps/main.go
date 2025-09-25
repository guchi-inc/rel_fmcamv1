// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "INFO -", 18)
)

// 表名记录
type TableName struct {
	Name string
}

// 通用返回结构。按字段封装值
type GcDesc struct {
	// Value any    `json:"value"`
	MaxLength    int    `json:"max_length" db:"max_length"`       //数据长度
	IsVisible    bool   `json:"is_visible" db:"is_visible"`       //可访问
	IsSearchable bool   `json:"is_searchable" db:"is_searchable"` //可搜索
	IsEditable   bool   `json:"is_editable" db:"is_editable"`     //可编辑
	IsRequired   bool   `json:"is_required" db:"is_required"`     //必填
	Name         string `json:"name" db:"name"`                   //字段
	CName        string `json:"cname" db:"cname"`                 //字段中文
	DataType     string `json:"data_type" db:"data_type"`         //数据类型
	Sort         string `json:"sort" db:"sort"`                   //排序规则
}

/*
Then type CTRL+SHIFT+p

And run the DDLTOGS: Convert the DDL (SQL) instructions to Go Struct to convert the sql text
selected to go struct

Or just select it and press CTRL+SHIFT+i
*/

// 部门表 原始表结构, 使用指针类型接受可能为空的字段，并且将隐藏,omitempty 的json字段。
type GcDepartmentStruct struct {
	ID          int64  `db:"id" json:"id,omitempty" form:"id"`                               //  主键
	DepType     string `db:"dep_type" json:"dep_type,omitempty" form:"dep_type"`             //  部门类型
	CreatedAt   string `db:"created_time" json:"created_time,omitempty" form:"created_time"` //  创建时间
	Creator     string `db:"creator" json:"creator,omitempty" form:"creator"`                //  操作员
	ManagerName string `db:"manager_name" json:"manager_name,omitempty" form:"manager_name"` //  负责人
	DepName     string `db:"dep_name" json:"dep_name" form:"dep_name"`                       //  部门名称
	DeleteFlag  string `db:"delete_flag" json:"delete_flag,omitempty" form:"delete_flag"`    //  删除标记，0未删除，1删除
}

// 用户表 原始表结构, 使用指针类型接受可能为空的字段，并且将隐藏,omitempty 的json字段。
type FmUser struct { //不返回空
	Id int64 `db:"id" json:"id" form:"id"` //主键

	Username    string `db:"username" json:"username" form:"username"`                    //用户姓名--例如张三
	LoginName   string `db:"login_name" json:"login_name" form:"login_name"`              //登录用户名
	Password    string `db:"password" json:"password" form:"password"`                    //登陆密码
	LeaderFlag  string `db:"leader_flag" json:"leader_flag,omitempty" form:"leader_flag"` //是否经理，0否，1是 /**string*/
	Position    string `db:"position" json:"position,omitempty" form:"position"`          //职位
	Department  string `db:"department" json:"department,omitempty" form:"department"`    //所属部门
	Email       string `db:"email" json:"email,omitempty" form:"email"`                   //电子邮箱
	Phonenum    string `db:"phonenum" json:"phonenum,omitempty" form:"phonenum"`          //电话号码
	Description string `db:"description" json:"description,omitempty" form:"description"` //用户描述信息
	Ethnicity   string `db:"ethnicity" json:"ethnicity" form:"ethnicity"`                 //国籍民族
	Gender      string `db:"gender" json:"gender" form:"gender"`                          //性别
	Local       string `db:"local" json:"local" form:"local"`                             //小区名
	LocalHost   string `db:"localhost" json:"localhost" form:"localhost"`                 //楼栋和房间
	M2LocalHost string `db:"m2_localhost" json:"m2_localhost" form:"m2_localhost"`        //面积 平方米

	MemberId   string `db:"member_id" json:"member_id"`                  //member_id 社区id
	LeaderId   string `db:"leader_id" json:"leader_id"`                  //父级账号id
	DeviceTime string `db:"device_time" json:"device_time"`              //最近登陆时间或 设备id
	TenantId   string `db:"tenant_id" json:"tenant_id" form:"tenant_id"` //  租商号

	CreatedAt string `db:"created_time" json:"created_time"` //创建时间
	UpdatedAt string `db:"updated_time" json:"updated_time"` //最近登录时间
	DeletedAt string `db:"deleted_time" json:"deleted_time"` ////禁用时间

	Ismanager string `db:"ismanager" json:"ismanager" form:"ismanager"` //是否为管理者 0==管理者 1==员工
	Isystem   string `db:"isystem" json:"isystem" form:"isystem"`       //是否系统自带数据

	DeleteFlag string `db:"delete_flag" json:"delete_flag,omitempty" form:"delete_flag"` //  删除标记，0未删除，1删除
	Enabled    *bool  `db:"enabled" json:"enabled" form:"enabled"`                       //状态，1：正常，0 封禁

	IsSms bool `db:"is_sms" json:"is_sms" form:"is_sms"` //短信验证状态，0 false 不开启。1 开启

}

// 用户注册 Register Param
type GcParamUser struct {
	ID int64 `db:"id" json:"id" form:"id"` //主键

	Supplier string `json:"supplier" db:"supplier" form:"supplier"` // 用户公司名租户名

	Username    string `db:"username" json:"username" form:"username"`                    //用户姓名--例如张三
	LoginName   string `db:"login_name" json:"login_name" form:"login_name"`              //登录用户名
	Password    string `db:"password" json:"password" form:"password"`                    //登陆密码
	LeaderFlag  string `db:"leader_flag" json:"leader_flag,omitempty" form:"leader_flag"` //是否经理，0否，1是 /**string*/
	Position    string `db:"position" json:"position,omitempty" form:"position"`          //职位
	Department  string `db:"department" json:"department,omitempty" form:"department"`    //所属部门
	Email       string `db:"email" json:"email,omitempty" form:"email"`                   //电子邮箱
	Phonenum    string `db:"phonenum" json:"phonenum,omitempty" form:"phonenum"`          //电话号码
	Description string `db:"description" json:"description,omitempty" form:"description"` //用户描述信息
	Ethnicity   string `db:"ethnicity" json:"ethnicity" form:"ethnicity"`                 //国籍民族
	Gender      string `db:"gender" json:"gender" form:"gender"`                          //性别

	Local       string `db:"local" json:"local" form:"local"`                      //小区名
	LocalHost   string `db:"localhost" json:"localhost" form:"localhost"`          //楼栋和房间
	M2LocalHost string `db:"m2_localhost" json:"m2_localhost" form:"m2_localhost"` //面积 平方米

	//member_id 证件号id
	MemberId   string `db:"member_id" json:"member_id"`                  //int_id 向量id
	DeviceTime string `db:"device_time" json:"device_time"`              //最近登陆时间和 设备id
	TenantId   string `db:"tenant_id" json:"tenant_id" form:"tenant_id"` //  租商号

	CreatedAt string `db:"created_time" json:"created_time"` //创建时间
	UpdatedAt string `db:"updated_time" json:"updated_time"` //最近登录时间
	DeletedAt string `db:"deleted_time" json:"deleted_time"` ////禁用时间

	LeaderId any `db:"leader_id" json:"leader_id"`

	Ismanager  string `db:"ismanager" json:"ismanager" form:"ismanager"`                 //是否为管理者 0==管理者 1==员工
	Isystem    string `db:"isystem" json:"isystem" form:"isystem"`                       //是否系统自带数据
	DeleteFlag string `db:"delete_flag" json:"delete_flag,omitempty" form:"delete_flag"` //  删除标记，0未删除，1删除
	Enabled    bool   `db:"enabled" json:"enabled" form:"enabled"`                       //状态，1：正常，0 封禁

	IsSms   bool   `db:"is_sms" json:"is_sms" form:"is_sms"`    //短信验证状态，0 false 不开启。1 开启
	Creator string `db:"creator" json:"creator" form:"creator"` //  操作员

}

// 批次的颜色和长度信息查询
type ColorResponse struct {
	BatchId  string `db:"batch_id" json:"batch_id"`
	Creator  string `db:"creator" json:"creator"`
	Length   string `db:"length" json:"length"`
	TypeName string `db:"type_name" json:"type_name"`
	Color    string `db:"color" json:"color"`
}

type ColorResponses struct {
	Total int64           `json:"total"`
	Size  int64           `json:"size"`
	Data  []ColorResponse `json:"data"`
}
