package seclients

import (
	"fmcam/models/erps"
	"fmcam/systems/genclients"
	"time"

	"github.com/google/uuid"
)

// apikeys管理
type GenApiKeyPages struct {
	Data     []*genclients.Apikeys `json:"data"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Columns  []erps.GcDesc         `json:"columns" db:"columns"`
}

// GenApikeys is the model entity for the GenApikeys schema.
type GenApikeys struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// UserID holds the value of the "user_id" field.
	UserID int64 `json:"user_id" db:"user_id"`

	// 使用计数
	UsageCount int64 `json:"usage_count" db:"usage_count"`
	// 启用
	Enabled *bool `json:"enabled" db:"enabled"`
	// 酒店号
	TenantID uuid.UUID `json:"tenant_id"  db:"tenant_id"`
	// 到期日期
	ExpiresTime time.Time `json:"expires_time" db:"expires_time"`
	// 创建时间
	CreatedTime time.Time `json:"created_time" db:"created_time"`
	// 创建时间
	UpdatedTime time.Time `json:"updated_time" db:"updated_time"`
	// 最近使用时间
	LastUsedTime time.Time `json:"last_used_time" db:"last_used_time"`
	//删除标记
	DeleteFlag string `json:"delete_flag" db:"delete_flag"`
	// APIKey字符串
	APIKey string `json:"api_key" db:"api_key"`
	// Key名称
	KeyName string `json:"key_name" db:"key_name"`

	//使用场景
	Type *int `json:"type,omitempty" db:"type"` //0 web，1 内部
	//扩展字段 参数接收
	Days      *int    `json:"days,omitempty" db:"days"`             //有效期
	LoginName *string `json:"login_name,omitempty" db:"login_name"` //key关联登录名
	Phonenum  *string `json:"phonenum,omitempty" db:"phonenum"`     //key管理账户手机号
}

// ApiKey 表
type ApiKey struct {
	ID int64 `json:"id" db:"id"`
	// UserID holds the value of the "user_id" field.
	UserID *int64 `json:"user_id" db:"user_id"`
	// 使用计数
	UsageCount *int64 `json:"usage_count" db:"usage_count"`
	// 启用
	Enabled *bool `json:"enabled" db:"enabled"`
	// 酒店号
	TenantID *string `json:"tenant_id"  db:"tenant_id"`
	// APIKey字符串
	APIKey string `json:"api_key" db:"api_key"`
	// Key名称
	KeyName *string `json:"key_name" db:"key_name"`
	// 到期日期
	ExpiresTime *time.Time `json:"expires_time" db:"expires_time"`
	// 创建时间
	CreatedTime *time.Time `json:"created_time" db:"created_time"`
	// 创建时间
	UpdatedTime *time.Time `json:"updated_time" db:"updated_time"`
	// 最近使用时间
	LastUsedTime *time.Time `json:"last_used_time" db:"last_used_time"`
	//使用场景
	Type *int `json:"type,omitempty" db:"type"` //0 web，1 内部
	//扩展字段 参数接收
	Days      *int    `json:"days,omitempty" db:"days"`             //有效期
	LoginName *string `json:"login_name,omitempty" db:"login_name"` //key关联登录名
	Phonenum  *string `json:"phonenum,omitempty" db:"phonenum"`     //key管理账户手机号

}

// 查询 ApiKey表
type ApiKeyList struct {
	Total   int           `json:"total" db:"total"`
	Size    int           `json:"page_size" db:"page_size"` //几天内的，size =1 表示今天，如果 size 为 -1，同时data不为空，则表示全部总的数据
	Data    []ApiKey      `json:"data" db:"data"`           //查询的数据
	Columns []erps.GcDesc `json:"columns" db:"columns"`     //字段说明
}
