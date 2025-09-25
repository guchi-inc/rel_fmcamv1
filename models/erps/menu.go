// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

//后台角色列表
type GcRoleList struct {
	Total   uint     `json:"total" db:"total"`
	Size    uint     `json:"size" db:"size"`
	Data    []GcRole `json:"data" db:"data"`
	Columns []GcDesc `json:"Columns" db:"Columns"`
}

//后台角色管理
type GcRole struct {
	ID      int64 `db:"id" json:"id" form:"id"`                //  主键
	Enabled bool  `db:"enabled" json:"enabled" form:"enabled"` //状态，1：正常，0 封禁

	RoleName    string `db:"role_name" json:"role_name" form:"role_name"`          //  名称
	Type        string `db:"type" json:"type" form:"type"`                         //  类型
	PriceLimit  string `db:"price_limit" json:"price_limit" form:"price_limit"`    //  价格屏蔽 1-屏蔽采购价 2-屏蔽零售价 3-屏蔽销售价
	Description string `db:"description" json:"description" form:"description"`    //  描述
	Sort        string `db:"sort" json:"sort" form:"sort"`                         //  排序
	TenantId    string `db:"tenant_id" json:"tenant_id" form:"tenant_id"`          //  用户id
	CreatedAt   string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
}

//后台角色新增参数
type GcRoleParam struct {
	ID      int64 `db:"id" json:"id" form:"id"`                //  主键
	Enabled bool  `db:"enabled" json:"enabled" form:"enabled"` //状态，1：正常，0 封禁

	RoleName    string `db:"role_name" json:"role_name" form:"role_name"`          //  名称
	Type        string `db:"type" json:"type" form:"type"`                         //  类型
	PriceLimit  string `db:"price_limit" json:"price_limit" form:"price_limit"`    //  价格屏蔽 1-屏蔽采购价 2-屏蔽零售价 3-屏蔽销售价
	Description string `db:"description" json:"description" form:"description"`    //  描述
	Sort        string `db:"sort" json:"sort" form:"sort"`                         //  排序
	TenantId    string `db:"tenant_id" json:"tenant_id" form:"tenant_id"`          //  租户号
	CreateTime  string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
	PermMask    string `db:"perm_mask" json:"perm_mask" form:"perm_mask"`          //  权限值
}

//后台 [角色/用户]  信息列表
type GcRoleUserList struct {
	Total   uint         `json:"total" db:"total"`
	Size    uint         `json:"size" db:"size"`
	Data    []GcRoleUser `json:"data" db:"data"`
	Columns []GcDesc     `json:"Columns" db:"Columns"`
}

//后台 [角色/用户] 管理
type GcRoleUser struct {
	ID     int64 `db:"id" json:"id" form:"id"`                //  主键
	RoleID int64 `db:"role_id" json:"role_id" form:"role_id"` //  角色主键
	UserID int64 `db:"user_id" json:"user_id" form:"user_id"` //  账号主键

	Type       string `db:"type" json:"type" form:"type"`                         //  类别
	RoleName   string `db:"role_name" json:"role_name" form:"role_name"`          //  角色名
	LoginName  string `db:"login_name" json:"login_name" form:"login_name"`       //  用户登陆名值
	CreateTime string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Creator    string `db:"creator" json:"creator" form:"creator"`                //  操作员
	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
}

//后台菜单类型 查询列表
type GcMenuList struct {
	Total   uint          `json:"total" db:"total"`
	Size    uint          `json:"size" db:"size"`
	Data    []GcQueryMenu `json:"data" db:"data"`
	Columns []GcDesc      `json:"Columns" db:"Columns"`
}

//后台菜单类型 查询列表
type GcFullMenuList struct {
	Total uint     `json:"total" db:"total"`
	Size  uint     `json:"size" db:"size"`
	Data  []GcMenu `json:"data" db:"data"`
}

//后端管理菜单表 查询返回
type GcMenu struct {
	ID          int64  `db:"id" json:"id" form:"id"`                               //  主键
	MenuType    string `db:"menu_type" json:"menu_type" form:"menu_type"`          //  菜单类型 业务 b，管理 m
	MenuName    string `db:"menu_name" json:"menu_name" form:"menu_name"`          //  菜单名
	MenuNumber  string `db:"menu_number" json:"menu_number" form:"menu_number"`    //  菜单编号
	CreateTime  string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Level       string `db:"level" json:"level" form:"level"`                      //  默认权限级别 low,middle,high
	Description string `db:"description" json:"description" form:"description"`    //  描述
	LeaderId    string `db:"leader_id" json:"leader_id" form:"leader_id"`          // 上级菜单id
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
	Path        string `db:"path" json:"path" form:"path"`                         //  接口路径默认空
}

//后端管理菜单表 查询返回
type GcQueryMenu struct {
	ID int64 `db:"id" json:"id" form:"id"` //  主键

	MenuType    string `db:"menu_type" json:"menu_type" form:"menu_type"`       //  菜单类型 业务 b，管理 m
	MenuName    string `db:"menu_name" json:"menu_name" form:"menu_name"`       //  菜单名
	MenuNumber  string `db:"menu_number" json:"menu_number" form:"menu_number"` //  菜单编号
	Level       string `db:"level" json:"level" form:"level"`                   //  默认权限级别 low,middle,high
	Description string `db:"description" json:"description" form:"description"` //  描述
	LeaderId    string `db:"leader_id" json:"leader_id" form:"leader_id"`       // 上级菜单id
	Path        string `db:"path" json:"path" form:"path"`                      //  接口路径默认空
}

// 权限组定义
type GcLimits struct {
	ID          int64  `db:"id" json:"id" form:"id"`                            //  主键
	LimitsId    string `db:"limits_id" json:"limits_id" form:"limits_id"`       //  权限编号
	TenantId    string `db:"tenant_id" json:"tenant_id" form:"tenant_id"`       //  租户号
	LimitsType  string `db:"limits_type" json:"limits_type" form:"limits_type"` //  权限类型 业务，管理
	LimitsName  string `db:"limits_name" json:"limits_name" form:"limits_name"` //  权限组名
	Level       string `db:"level" json:"level" form:"level"`                   //  默认权限级别 读/写
	DeleteFlag  string `db:"delete_flag" json:"delete_flag" form:"delete_flag"` //  删除标记，0未删除，1删除
	Description string `db:"description" json:"description" form:"description"` //  描述
	PermMask    string `db:"perm_mask" json:"perm_mask" form:"perm_mask"`       //  权限掩码
}

//后台菜单 权限组关联记录 查询列表
type GcLimitsList struct {
	Total   uint       `json:"total" db:"total"`
	Size    uint       `json:"size" db:"size"`
	Data    []GcLimits `json:"data" db:"data"`
	Columns []GcDesc   `json:"Columns" db:"Columns"`
}

//后台 用户 权限组关联记录 查询列表
type GcRoleLimitsList struct {
	Total   uint           `json:"total" db:"total"`
	Size    uint           `json:"size" db:"size"`
	Data    []GcRoleLimits `json:"data" db:"data"`
	Columns []GcDesc       `json:"Columns" db:"Columns"`
}

// 权限组 用户关联记录
type GcRoleLimits struct {
	ID    int64 `db:"id" json:"id" form:"id"`             //  主键
	BtnId int64 `db:"btn_id" json:"btn_id" form:"btn_id"` //  功能权限号

	RoleName   string `db:"role_name" json:"role_name" form:"role_name"`          //  角色名
	LimitsId   string `db:"limits_id" json:"limits_id" form:"limits_id"`          //  权限编号
	LimitsName string `db:"limits_name" json:"limits_name" form:"limits_name"`    //  权限组名
	CreateTime string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Creator    string `db:"creator" json:"creator" form:"creator"`                //  操作员
	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
}

// 权限组 用户关联记录
type FmRoleLimitNews struct {
	ID        int64   `db:"id" json:"id" form:"id"`                            //  主键
	BtnIdList []int64 `db:"btn_id_list" json:"btn_id_list" form:"btn_id_list"` //  要新增的 功能权限号列表

	RoleName   string `db:"role_name" json:"role_name" form:"role_name"`          //  角色名
	LimitsId   string `db:"limits_id" json:"limits_id" form:"limits_id"`          //  权限编号
	LimitsName string `db:"limits_name" json:"limits_name" form:"limits_name"`    //  权限组名
	CreateTime string `db:"created_time" json:"created_time" form:"created_time"` //  创建时间
	Creator    string `db:"creator" json:"creator" form:"creator"`                //  操作员
	DeleteFlag string `db:"delete_flag" json:"delete_flag" form:"delete_flag"`    //  删除标记，0未删除，1删除
}

// 角色已关联功能号记录
type FmRoleLimitHistoryLists struct {
	ID        int64   `db:"id" json:"id" form:"id"`                            //  主键
	BtnIdList []int64 `db:"btn_id_list" json:"btn_id_list" form:"btn_id_list"` //  要新增的 功能权限号列表

	RoleName string `db:"role_name" json:"role_name" form:"role_name"` //  角色名
	LimitsId string `db:"limits_id" json:"limits_id" form:"limits_id"` //  权限编号
}
