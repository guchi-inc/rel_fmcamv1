package seclients

import (
	"fmcam/models/erps"
	"time"
)

// 需要排除的查询字段条件
type FieldTags struct {
	SortTagFilter string `json:"tag_sort" db:"tag_sort"`
}

type GroupTypeMapping struct {
	ID int64 `json:"id" db:"id"`
	// 分组ID
	GroupID int64 `json:"group_id" db:"group_id"`
	// 设备 ID
	ProfileTypeID int64 `json:"profile_type_id" db:"profile_type_id"`
	// 操作人
	Creator string `json:"creator" db:"creator"`
	// 操作时间
	CreatedTime *time.Time `json:"created_time" db:"created_time"`
}

// 人员类型组 映射列表
type FmGroupTypeMappingPages struct {
	Data     []GroupTypeMapping `json:"data" db:"data"`
	Total    int                `json:"total" db:"total"`
	Page     int                `json:"page" db:"page"`
	PageSize int                `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc      `json:"columns" db:"columns"`
}

// 正在使用的 从某个人员类型  查询的 相关人员分组  预警分组，等级信息
type ProfileTypeAlertInfoMap struct {
	// 关联的分组id
	GroupID int64 `json:"group_id"  db:"group_id"`
	// 人员类型id
	ProfileTypeID int64 `json:"profile_type_id"  db:"profile_type_id"`
	// 租户号
	TenantID *string `json:"tenant_id"  db:"tenant_id"`
	//人员类别名称
	TypeName *string `json:"type_name" db:"type_name"`
	//人员类别代码
	TypeCode *string `json:"type_code" db:"type_code"`
	//人员类别描述
	Description *string `json:"description" db:"description"`
	//人员类别图像有效时长
	FaceValidityHours *int `json:"face_validity_hours" db:"face_validity_hours"`
	//人员类别预警等级
	WarningLevel *int `json:"warning_level" db:"warning_level"`
	//人员类别预警开关
	WarningEnabled *bool `json:"warning_enabled" db:"warning_enabled"`
	//人员类别图 可删除
	Deleteable *bool `json:"deleteable" db:"deleteable"`
	//人员类别状态
	Enabled *int `json:"enabled" db:"enabled"`

	//组类别
	GroupType *int8 `json:"group_type" db:"group_type"`
	// 组名
	GroupName *string `json:"group_name"  db:"group_name"`
	// 定制化
	Customization *string `json:"customization"  db:"customization"`
	// 预警定义等级
	Level *int `json:"level"  db:"level"`
	// 预警处理操作
	Action *string `json:"action"  db:"action"`
	// 预警声音
	AlarmSound *string `json:"alarm_sound"  db:"alarm_sound"`
	// 创建时间
	CreatedAt *time.Time `json:"created_time,omitempty" db:"created_time"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_time,omitempty" db:"updated_time"`
}

// 人员类型组 映射列表
type FmProfileTypeAlertInfoPages struct {
	Data     []ProfileTypeAlertInfoMap `json:"data" db:"data"`
	Total    int                       `json:"total" db:"total"`
	Page     int                       `json:"page" db:"page"`
	PageSize int                       `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc             `json:"columns" db:"columns"`
}

// FmAlertGroup is the model entity for the FmAlertGroup schema.
type FmAlertGroup struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	//组类别
	GroupType int8 `json:"group_type" db:"group_type"`
	// 启用
	Enabled *bool `json:"enabled"  db:"enabled"`
	//启用类型 null 未启用， 租户号
	EnabledGroup *string `json:"uniq_enabled_group"  db:"uniq_enabled_group"`
	// 组名
	GroupName string `json:"group_name"  db:"group_name"`
	//合作方酒名
	Supplier string `json:"supplier"  db:"supplier"`
	// 租户号
	TenantID string `json:"tenant_id"  db:"tenant_id"`
	// 定制化
	Customization string `json:"customization"  db:"customization"`
	// 描述
	Description string `json:"description"  db:"description"`
	//操作员
	Creator string `json:"creator"  db:"creator"`
	// 删除标记
	DeleteFlag string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间
	CreatedAt any `json:"created_time,omitempty" db:"created_time"`
	// 更新时间
	UpdatedAt any `json:"updated_time,omitempty" db:"updated_time"`

	//扩展字段 方便接收相关等级 level 数据
	NewLevelList []FmAlertLevel `json:"level_list" db:"level_list"`
	//扩展字段 方便接收相关人员类型 数据
	NewProfileTypeList []ProfileType `json:"profile_type_list" db:"profile_type_list"`
}

// 预警组列表
type FmAlertGroupPages struct {
	Data     []FmAlertGroup `json:"data" db:"data"`
	Total    int            `json:"total" db:"total"`
	Page     int            `json:"page" db:"page"`
	PageSize int            `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc  `json:"columns" db:"columns"`
}

// 预警等级 FmAlertGroup is the model entity for the FmAlertGroup schema.
type FmAlertLevel struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 关联的分组id
	AlertGroupID int64 `json:"alert_group_id"  db:"alert_group_id"`
	// 人员类型id
	ProfileTypeID int64 `json:"profile_type_id"  db:"profile_type_id"`
	// 预警等级
	Level int `json:"level"  db:"level"`
	//预警组是否启用中
	AlertGroupEnabled *bool `json:"enabled"  db:"enabled"`
	// 处理操作
	Action string `json:"action"  db:"action"`
	// 预警声音
	AlarmSound string `json:"alarm_sound"  db:"alarm_sound"`
	// 删除标记
	DeleteFlag string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间
	CreatedAt *time.Time `json:"created_time,omitempty" db:"created_time"`
	// 更新时间
	UpdatedAt *time.Time `json:"updated_time,omitempty" db:"updated_time"`
}

// 预警等级列表
type FmAlertLevelPages struct {
	Data     []FmAlertLevel `json:"data" db:"data"`
	Total    int            `json:"total" db:"total"`
	Page     int            `json:"page" db:"page"`
	PageSize int            `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc  `json:"columns" db:"columns"`
}

// 批量接收 预警等级
type ParamAlertList struct {
	//预警等级列表
	ID           int64          `json:"id" db:"id"` //预警组id
	NewLevelList []FmAlertLevel `json:"level_list" db:"level_list"`
}
