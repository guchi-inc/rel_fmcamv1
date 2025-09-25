// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

// 供应商
type FmTenant struct {
	Id      int64 `json:"id" db:"id" form:"id"`                // 主键
	Enabled bool  `db:"enabled" json:"enabled" form:"enabled"` //状态，1：正常，0 封禁

	Supplier      string `json:"supplier" db:"supplier" form:"supplier"`                   // 租户名
	Contacts      string `json:"contacts" db:"contacts" form:"contacts"`                   // 联系人
	Email         string `json:"email" db:"email" form:"email"`                            // 电子邮箱
	Description   string `json:"description" db:"description" form:"description"`          // 备注
	TenantId      string `json:"tenant_id" db:"tenant_id" form:"tenant_id"`                // 租户号
	Type          string `json:"type" db:"type" form:"type"`                               // 类型
	Province      string `json:"province" db:"province" form:"province"`                   // 省份
	City          string `json:"city" db:"city" form:"city"`                               // 城市
	Area          string `json:"area" db:"area" form:"area"`                               // 区县
	Street        string `json:"street" db:"street" form:"street"`                         // 街道乡镇
	Address       string `json:"address" db:"address" form:"address"`                      // 街道地址
	AddrCode      string `json:"addr_code" db:"addr_code" form:"addr_code"`                // 地址编码,比如邮政编码
	Fax           string `json:"fax" db:"fax" form:"fax"`                                  // 传真
	PhoneNum      string `json:"phone_num" db:"phone_num" form:"phone_num"`                // 固定电话
	Telephone     string `json:"telephone" db:"telephone" form:"telephone"`                // 手机号
	TaxNum        string `json:"tax_num" db:"tax_num" form:"tax_num"`                      // 纳税人,识别号
	BankName      string `json:"bank_name" db:"bank_name" form:"bank_name"`                // 开户行
	AccountNumber string `json:"account_number" db:"account_number" form:"account_number"` // 账号
	Sort          string `json:"sort" db:"sort" form:"sort"`                               // 排序
	Isystem       string `json:"isystem" db:"isystem" form:"isystem"`                      // 系统自带, 0非系统 1系统
	AdvanceIn     string `json:"advance_in" db:"advance_in" form:"advance_in"`             // 预收款
	BeginNeedGet  string `json:"begin_need_get" db:"begin_need_get" form:"begin_need_get"` // 期初应收
	BeginNeedPay  string `json:"begin_need_pay" db:"begin_need_pay" form:"begin_need_pay"` // 期初应付
	AllNeedGet    string `json:"all_need_get" db:"all_need_get" form:"all_need_get"`       // 累计应收
	AllNeedPay    string `json:"all_need_pay" db:"all_need_pay" form:"all_need_pay"`       // 累计应付
	TaxRate       string `json:"tax_rate" db:"tax_rate" form:"tax_rate"`                   // 税率
	FullAddress   string `json:"full_address" db:"full_address" form:"full_address"`       // 详细地址
	CreatedAt     string `db:"created_time" json:"created_time" form:"created_time"`       //  创建时间
	Creator       string `db:"creator" json:"creator" form:"creator"`                      //  操作员
	UpdatedAt     string `db:"updated_time" json:"updated_time" form:"updated_time"`       //  更新时间
	DeleteFlag    string `json:"delete_flag" db:"delete_flag" form:"delete_flag"`          // 删除标记,0未删除,1删除

}

// 供应商列表
type FmSupTenantList struct {
	Total   uint       `json:"total" db:"total"`
	Size    uint       `json:"size" db:"size"`
	Data    []FmTenant `json:"data" db:"data"`
	Columns []GcDesc   `json:"Columns" db:"Columns"`
}
