package seclients

import (
	"fmcam/models/erps"
	"fmcam/systems/genclients"
	"time"

	"github.com/google/uuid"
)

type ProfileType struct {

	// ID of the ent.
	ID int64 `json:"id" db:"id"`

	// 启用
	Enabled bool `json:"enabled" db:"enabled"`
	// 可删除
	DelEnabled *bool `json:"deleteable" db:"deleteable"`
	// 有效时长
	FaceValidityHours *int `json:"face_validity_hours"  db:"face_validity_hours"`
	//预警等级
	WarningLevel *int `json:"warning_level"  db:"warning_level"`
	//预警启用
	WarningEnabled *bool `json:"warning_enabled"  db:"warning_enabled"`
	// 类型名
	TypeName string `json:"type_name"  db:"type_name"`
	// 人员代码
	TypeCode *string `json:"type_code"  db:"type_code"`
	// 租户号
	TenantID string `json:"tenant_id"  db:"tenant_id"`
	// 描述
	Description *string `json:"description"  db:"description"`

	// 删除标记
	DeleteFlag *string `json:"delete_flag"  db:"delete_flag"`
	// 创建时间timeutil.JSONTime
	CreatedAt *time.Time `json:"created_time,omitempty" db:"created_time"`
	// 更新时间 sql.NullTime
	UpdatedAt *time.Time `json:"updated_time,omitempty" db:"updated_time"`
}

// 批量接收 人员类型
type ParamProfileTypeLists struct {
	ID int64 `json:"id" db:"id"`
	//人员类型列表
	NewTypeList []ProfileType `json:"profile_type_list" db:"profile_type_list"`
	// 租户号
	TenantID string `json:"tenant_id"  db:"tenant_id"`
}

// 客户端 人员类别
type ProfileTypeList struct {
	Data     []ProfileType `json:"data" db:"data"`
	Total    int           `json:"total" db:"total"`
	Page     int           `json:"page" db:"page"`
	PageSize int           `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

type ProfileSuppTypes struct {
	// ProfileType is the model entity
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 类型名
	TypeName *string `json:"type_name"  db:"type_name"`
	// 租户酒店名
	Supplier *string `json:"supplier"  db:"supplier"`
	// 租户号
	TenantID uuid.UUID `json:"tenant_id"  db:"tenant_id"`
	// 描述
	Description *string `json:"description"  db:"description"`
	// 启用
	Enabled bool `json:"enabled" db:"enabled"`
	//预警等级
	WarningLevel *int `json:"warning_level"  db:"warning_level"`
	//预警启用
	WarningEnabled *bool `json:"warning_enabled"  db:"warning_enabled"`
	// 人员类型
	TypeCode *string `json:"type_code"  db:"type_code"`
	// 有效时长
	FaceValidityHours *int `json:"face_validity_hours"  db:"face_validity_hours"`
	// 创建时间
	CreatedAt time.Time `json:"created_time,omitempty" db:"created_time"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_time,omitempty" db:"updated_time"`
}

type ProfileInfos struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 人员类型
	TypeID *int64 `json:"type_id" db:"type_id"`
	// 姓名
	Name *string `json:"name"  db:"name"`
	// 租户号
	TenantID *string `json:"tenant_id"  db:"tenant_id"`
	// 身份证号
	IDCardNumber *string `json:"id_card_number"  db:"id_card_number"`
	// 手机号
	PhoneNumber *string `json:"phone_number"  db:"phone_number"`
	// 启用
	Enabled *bool `json:"enabled" db:"enabled"`
	// 房间号
	RoomID *string `json:"room_id" db:"room_id"`
	// 面部临时地址
	TmpURL     *string `json:"tmp_url" db:"tmp_url"`
	DeleteFlag *string `json:"delete_flag" db:"delete_flag"`
	// 创建时间
	CreatedAt time.Time `json:"created_time" db:"created_time"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_time" db:"updated_time"`
}

// 客户端 人员类别
type ProfileSuppTypesList struct {
	Data     []*ProfileSuppTypes `json:"data" db:"data"`
	Total    int                 `json:"total" db:"total"`
	Page     int                 `json:"page" db:"page"`
	PageSize int                 `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc       `json:"columns" db:"columns"`
}

// 客户端 人员类别
type ProfileTypesPointer struct {
	Data     []*genclients.ProfileType `json:"data" db:"data"`
	Total    int                       `json:"total" db:"total"`
	Page     int                       `json:"page" db:"page"`
	PageSize int                       `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc             `json:"columns" db:"columns"`
}

// 客户端 人员
type ProfilesList struct {
	Data     []*genclients.Profiles `json:"data" db:"data"`
	Total    int                    `json:"total" db:"total"`
	Page     int                    `json:"page" db:"page"`
	PageSize int                    `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc          `json:"columns" db:"columns"`
}

// 客户端 面部数据
type FacesList struct {
	Data     []*genclients.Faces `json:"data" db:"data"`
	Total    int                 `json:"total" db:"total"`
	Page     int                 `json:"page" db:"page"`
	PageSize int                 `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc       `json:"columns" db:"columns"`
}

// 客户端 暂存面部数据
type TemporaryFaceList struct {
	Data     []*genclients.TemporaryFace `json:"data" db:"data"`
	Total    int                         `json:"total" db:"total"`
	Page     int                         `json:"page" db:"page"`
	PageSize int                         `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc               `json:"columns" db:"columns"`
}

// 客户端 暂存面部数据
type DevicesList struct {
	Data     []*genclients.Devices `json:"data" db:"data"`
	Total    int                   `json:"total" db:"total"`
	Page     int                   `json:"page" db:"page"`
	PageSize int                   `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc         `json:"columns" db:"columns"`
}

type CaptureLogInfo struct {
	ID int64 `json:"id" db:"id"`
	// 租户号
	TenantID *string `json:"tenant_id"  db:"tenant_id"`
	// 设备 ID
	DeviceID *int64 `json:"device_id" db:"device_id"`
	// 比中人员ID
	MatchedProfileID *int64 `json:"matched_profile_id" db:"matched_profile_id"`
	// 人员房间ID
	RoomId *string `json:"room_id" db:"room_id"`

	//人员信息 类型，id，姓名
	Name *string `json:"name"  db:"name"`

	TypeId   *int64  `json:"type_id" db:"type_id"`
	TypeName *string `json:"type_name"  db:"type_name"`

	//人员面部id
	FaceId *int64 `json:"face_id" db:"face_id"`
	//已预警ID
	AlertId    *int64 `json:"alert_id" db:"alert_id"`
	AlertLevel *int   `json:"alert_level" db:"alert_level"`
	// 是否抓拍
	FuncType *int8 `json:"func_type" db:"func_type"`
	// 比对得分
	MatchScore *float32 `json:"match_score" db:"match_score"`
	// 触发预警
	HasAlert *bool `json:"has_alert" db:"has_alert"`
	// 设备名称
	DeviceName *string `json:"device_name"  db:"device_name"`
	// 设备位置
	DeviceLocation *string `json:"device_location"  db:"device_location"`
	// 抓拍记录内容
	Content *string `json:"content"  db:"content"`
	// 抓拍的图片地址
	CaptureImageURL *string `json:"capture_image_url"  db:"capture_image_url"`
	// 抓拍时间
	CaptureTime *time.Time `json:"capture_time" db:"capture_time"`
	// 处理状态: 0=未处理, 1=已处理, 2=已忽略
	Status *int `json:"status" db:"status"`
}

// 客户端 抓拍数据
type CaptureLogsList struct {
	Data     []*CaptureLogInfo `json:"data" db:"data"`
	Total    int64             `json:"total" db:"total"`
	Page     int64             `json:"page" db:"page"`
	PageSize int64             `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc     `json:"columns" db:"columns"`
}

// Alerts is the model entity for the Alerts schema.
type Alerts struct {
	// ID of the ent.
	ID int64 `json:"id" db:"id"`
	// 处理意见或备注
	Remarks any `json:"remarks,omitempty" db:"remarks"`
	// 租户号
	TenantID string `json:"tenant_id"  db:"tenant_id"`
	// 处理人ID
	FmUserId any `json:"fm_user_id" db:"fm_user_id"`
	// 设备 ID
	DeviceID int64 `json:"device_id" db:"device_id"`

	// 抓拍记录 ID
	CaptureLogID int64 `json:"capture_log_id" db:"capture_log_id"`
	// 预警等级: 1=一级预警, 2=二级预警, 3=三级预警, 4=四级预警
	AlertLevel any `json:"alert_level" db:"alert_level"`
	// 处理状态: 0=未处理, 1=已处理, 2=已忽略
	Status int8 `json:"status" db:"status"`
	// 开始处理的时间
	HandledAt any `json:"handled_time" db:"handled_time"`
	// 预警产生时间
	CreatedAt time.Time `json:"created_time" db:"created_time"`

	ProfileTypeId   *int64  `json:"profile_type_id" db:"profile_type_id"`       //人员类型id
	ProfileId       *int64  `json:"matched_profile_id" db:"matched_profile_id"` //人员id
	CaptureImageUrl *string `json:"capture_image_url" db:"capture_image_url"`   //抓拍图像地址
	Location        *string `json:"location" db:"location"`                     //抓拍设备位置
}

// 客户端 告警数据
type AlertsList struct {
	Data     []*Alerts     `json:"data" db:"data"`
	Total    int           `json:"total" db:"total"`
	Page     int           `json:"page" db:"page"`
	PageSize int           `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc `json:"columns" db:"columns"`
}

// 定义后端发送的数据结构
type AlertNews struct {
	ID int64 `json:"id" db:"id"`

	// 处理人ID
	FmUserId *int64 `json:"fm_user_id" db:"fm_user_id"`
	// 设备 ID
	DeviceID int64 `json:"device_id" db:"device_id"`

	Level *int `json:"level" db:"level"` // 预警等级
	// 抓拍记录 ID
	CaptureLogID int64 `json:"capture_log_id" db:"capture_log_id"`
	// 预警等级: 1=一级预警, 2=二级预警, 3=三级预警, 4=四级预警
	AlertGroupId     *int64  `json:"alert_group_id" db:"alert_group_id"` // 预警组信息
	AlertLevel       *int64  `json:"alert_level" db:"alert_level"`
	MatchedProfileId *int64  `json:"matched_profile_id" db:"matched_profile_id"`
	ProfileTypeId    *int64  `json:"profile_type_id" db:"profile_type_id"` //人员类型id
	Action           *string `json:"action" db:"action"`                   // 客户端页面动作
	AlarmSound       *string `json:"alarm_sound" db:"alarm_sound"`         //客户端页面提示音
	Location         *string `json:"location" db:"location"`               //设备位置
	// 处理意见或备注
	Remarks *string `json:"remarks,omitempty" db:"remarks"`
	// 租户号
	TenantID        string `json:"tenant_id"  db:"tenant_id"`
	CaptureImageUrl string `json:"capture_image_url"  db:"capture_image_url"`

	// 预警产生时间
	CreatedAt *time.Time `json:"created_time" db:"created_time"`
}

// checkout profile
type CheckoutProfile struct {
	Id int64 `db:"id" json:"id" form:"id"` //主键

	LoginPhonenum *string `db:"login_phonenum" json:"login_phonenum,omitempty" form:"login_phonenum"` //鉴权电话
	LoginName     *string `db:"login_name" json:"login_name" form:"login_name"`                       //登录用户名
	LoginPassword *string `db:"login_password" json:"login_password" form:"login_password"`           //登陆密码

	PhoneGuest  *string `db:"phone_guest" json:"phone_guest,omitempty" form:"phone_guest"`       // 住客电话
	IdCardGuest *string `db:"id_card_guest" json:"id_card_guest,omitempty" form:"id_card_guest"` //证件号
	RoomId      *string `db:"room_id" json:"room_id,omitempty" form:"room_id"`                   // 入住房间号
	TmpUrl      *string `db:"tmp_url" json:"tmp_url,omitempty" form:"tmp_url"`                   // 证件照URL

}

// 查询账户 uuid checkout profile
type UserProfile struct {
	Id int64 `db:"id" json:"id" form:"id"` //主键

	//电子邮箱
	Phonenum  *string    `db:"phonenum" json:"phonenum,omitempty" form:"phonenum"` //用户电话
	LoginName *string    `db:"login_name" json:"login_name" form:"login_name"`     //登录用户名
	Password  *string    `db:"password" json:"password" form:"password"`           //登陆密码
	Email     *string    `db:"email" json:"email,omitempty" form:"email"`
	TenantId  *uuid.UUID `db:"tenant_id" json:"tenant_id,omitempty" form:"tenant_id"` //证件号
}

// 客户端  账号列表
type FmAccountsList struct {
	Data     []*genclients.FmUserAccount `json:"data" db:"data"`
	Total    int                         `json:"total" db:"total"`
	Page     int                         `json:"page" db:"page"`
	PageSize int                         `json:"page_size" db:"page_size"`
	Columns  []erps.GcDesc               `json:"columns" db:"columns"`
}

// 客户端  账号数据
type FmAccountsInfo struct {
	Account    *genclients.FmUserAccount `json:"account" db:"account"`
	TenantInfo *TenantFull               `json:"tenant_info" db:"tenant_info"`
}

// 面部表可修改参数 接收
type FaceParam struct {
	ID int64 `json:"id" db:"id"`
	// 租户号
	TenantID *string `json:"tenant_id"  db:"tenant_id"`
	// ProfileID holds the value of the "profile_id" field.
	ProfileID *int64 `json:"profile_id" db:"profile_id"`
	// 照片地址
	ImageURL *string `json:"image_url"  db:"image_url"`
	//删除标记
	DeleteFlag *string `json:"delete_flag" db:"delete_flag"`
}
