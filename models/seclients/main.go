package seclients

import (
	"fmcam/models/erps"
	"fmcam/systems/genclients"
	"time"
)

// 上下文操作人
type CtxOperator struct {
	UserId string `json:"user_id" db:"user_id"`
	IP     string `json:"ip" db:"ip"`
}

// Province is the model entity for the Province schema.
type Province struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 编号
	Code string `json:"code"  db:"code"`
	// 名称
	Name string `json:"name"  db:"name"`
	// 操作人
	Creator string `json:"creator"  db:"creator"`
	// 删除标记
	DeleteFlag string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间
	CreatedAt any `json:"created_time,omitempty" db:"created_time"`
}

// 省
type ProvincePage struct {
	Data     []*Province   `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

// GovCity is the model entity for the GovCity schema.
type GovCity struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 编号
	Code string `json:"code"  db:"code"`
	// 名称
	Name string `json:"name"  db:"name"`
	// 省编号
	ProvinceCode string `json:"province_code"  db:"province_code"`
	// 操作人
	Creator string `json:"creator"  db:"creator"`
	// 删除标记
	DeleteFlag string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间
	CreatedAt any `json:"created_time" db:"created_time"`
}

// 市
type CityPage struct {
	Data     []*GovCity    `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

// GovArea is the model entity for the GovArea schema.
type GovArea struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 编号
	Code string `json:"code"  db:"code"`
	// 名称
	Name string `json:"name"  db:"name"`
	// 省编号
	ProvinceCode string `json:"province_code"  db:"province_code"`
	// 市编号
	CityCode string `json:"city_code"  db:"city_code"`
	// 操作人
	Creator string `json:"creator"  db:"creator"`
	// 删除标记
	DeleteFlag string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间
	CreatedAt any `json:"created_time" db:"created_time"`
}

// 区县
type AreaPage struct {
	Data     []*GovArea    `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

// GovStreet is the model entity for the GovStreet schema.
type GovStreet struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 编号
	Code string `json:"code"  db:"code"`
	// 名称
	Name string `json:"name"  db:"name"`
	// 省编号
	ProvinceCode string `json:"province_code"  db:"province_code"`
	// 市编号
	CityCode string `json:"city_code"  db:"city_code"`
	// 区县编号
	AreaCode string `json:"area_code"  db:"area_code"`
	// 操作人
	Creator string `json:"creator"  db:"creator"`
	// 删除标记
	DeleteFlag string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间
	CreatedAt any `json:"created_time" db:"created_time"`
}

// 街道/乡镇
type StreetPage struct {
	Data     []*GovStreet  `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

// 字段管理
type FieldMetadataPage struct {
	Data     genclients.FieldMetadataSlice `json:"data"`
	Total    int                           `json:"total"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"page_size"`
	Columns  []erps.GcDesc                 `json:"columns" db:"columns"`
}

type TenantFull struct {
	ID           int64  `json:"id" db:"id" form:"id"`                                     // 主键
	Enabled      bool   `json:"enabled" db:"enabled" form:"enabled"`                      // 状态
	Isystem      bool   `json:"isystem" db:"isystem" form:"isystem"`                      // 系统自带, 0非系统 1系统
	AdvanceIn    string `json:"advance_in" db:"advance_in" form:"advance_in"`             // 预收款
	BeginNeedGet string `json:"begin_need_get" db:"begin_need_get" form:"begin_need_get"` // 期初应收
	BeginNeedPay string `json:"begin_need_pay" db:"begin_need_pay" form:"begin_need_pay"` // 期初应付
	AllNeedGet   string `json:"all_need_get" db:"all_need_get" form:"all_need_get"`       // 累计应收
	AllNeedPay   string `json:"all_need_pay" db:"all_need_pay" form:"all_need_pay"`       // 累计应付
	TaxRate      string `json:"tax_rate" db:"tax_rate" form:"tax_rate"`                   // 税率

	TenantId   string `json:"tenant_id" db:"tenant_id" form:"tenant_id"`       // 所属租户 ID
	DeleteFlag string `json:"delete_flag" db:"delete_flag" form:"delete_flag"` // 删除标记,0未删除,1删除

	Supplier      string `json:"supplier" db:"supplier" form:"supplier"`                   // 租户名
	Contacts      string `json:"contacts" db:"contacts" form:"contacts"`                   // 联系人
	Email         string `json:"email" db:"email" form:"email"`                            // 电子邮箱
	Type          string `json:"type" db:"type" form:"type"`                               // 类型
	Province      string `json:"province" db:"province" form:"province"`                   // 省份
	City          string `json:"city" db:"city" form:"city"`                               // 城市
	Area          string `json:"area" db:"area" form:"area"`                               // 区县
	Street        string `json:"street" db:"street" form:"street"`                         // 街道乡镇
	Address       string `json:"address" db:"address" form:"address"`                      // 门牌地址
	AddrCode      string `json:"addr_code" db:"addr_code" form:"addr_code"`                // 地址编码,比如邮政编码
	FullAddress   string `json:"full_address" db:"full_address" form:"full_address"`       // 详细地址
	Fax           string `json:"fax" db:"fax" form:"fax"`                                  // 传真
	PhoneNum      string `json:"phone_num" db:"phone_num" form:"phone_num"`                // 固定电话
	Telephone     string `json:"telephone" db:"telephone" form:"telephone"`                // 手机号
	TaxNum        string `json:"tax_num" db:"tax_num" form:"tax_num"`                      // 纳税人,识别号
	BankName      string `json:"bank_name" db:"bank_name" form:"bank_name"`                // 开户行
	AccountNumber string `json:"account_number" db:"account_number" form:"account_number"` // 账号
	Sort          string `json:"sort" db:"sort" form:"sort"`                               // 排序
	Description   string `json:"description" db:"description" form:"description"`          // 备注
	Creator       string `json:"creator" db:"creator" form:"creator"`                      // 操作人
	CreatedAt     string `json:"created_time" db:"created_time" form:"created_time"`       // 新增时间
	UpdateAt      string `json:"updated_time" db:"updated_time" form:"updated_time"`       // 更新时间

}

// 租户管理
type TenantsPage struct {
	Data     []*TenantFull `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

// FMPMSAPI 参数接收
type FMPMSApi struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// PMS商名
	PmsName string `json:"pms_name"  db:"pms_name"`
	// 对接地址
	PmsAPI string `json:"pms_api"  db:"pms_api"`
	// 启用
	Enabled *bool `json:"enabled" db:"enabled"`
	// 联系人
	Contact string `json:"contact"  db:"contact"`
	// 联系电话
	Phonenum *string `json:"phonenum"  db:"phonenum"`
	//功能描述
	Description *string `json:"description"  db:"description"`
	// 删除标记
	DeleteFlag *string `json:"delete_flag"  db:"delete_flag"`
	// 操作员
	Creator *string `json:"creator"  db:"creator"`
	// 创建时间
	CreatedAt *time.Time `json:"created_time" db:"created_time"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_time" db:"updated_time"`
}

// fm pms api 页面
type FMPMSApiPage struct {
	Data     []*FMPMSApi   `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}
