// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package code

import (
	"fmcam/common/configs"
)

var ByteCodeFile []byte
var (
	CronTables = map[string]bool{
		DefaultCronTable:   true,
		DefaultVideoTable:  true,
		DefaultScriptTable: true,
		DefaultLiveTable:   true,
	}
)

var (
	DefaultZsetPrefix  = "zset"
	DefaultCronTable   = "cron"
	DefaultVideoTable  = "videos"
	DefaultScriptTable = "python"
	DefaultLiveTable   = "live" //直播表
	AllowCronTables    = []string{DefaultCronTable, DefaultVideoTable, DefaultScriptTable}
	DefaultPriKey      = "id"

	//按顺序存入和解析cron 结构体
	OrderFields = []string{"name", "spec", "command"}
	//直播字段
	LiveFields = []string{"name", "url", "start_time"}
)

// Failure 错误时返回结构
type Failure struct {
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
}

var (
	Success = 100000 //100000：成功响应

	Failed = 100003 //：失败响应

	NullData = 100005 //：查询空值
)

const (
	ServerError        = 10101
	TooManyRequests    = 10102
	ParamBindError     = 10103
	AuthorizationError = 10104 //登录token错误
	UrlSignError       = 10105
	CacheSetError      = 10106
	CacheGetError      = 10107
	CacheDelError      = 10108
	CacheNotExist      = 10109
	ResubmitError      = 10110
	HashIdsEncodeError = 10111
	HashIdsDecodeError = 10112
	RBACError          = 10113
	RedisConnectError  = 10114
	MySQLConnectError  = 10115
	WriteConfigError   = 10116
	SendEmailError     = 10117
	MySQLExecError     = 10118
	GoVersionError     = 10119
	SocketConnectError = 10120
	SocketSendError    = 10121
	ParamError         = 10122 //参数校验错误

	AuthorizedCreateError    = 20101
	AuthorizedListError      = 20102
	AuthorizedDeleteError    = 20103
	AuthorizedUpdateError    = 20104
	AuthorizedDetailError    = 20105 //token 合法性详情校验失败
	AuthorizedCreateAPIError = 20106
	AuthorizedListAPIError   = 20107
	AuthorizedDeleteAPIError = 20108 //登出后 又试图 使用的token
	AuthorizedPowerError     = 20109 //权限不足
	AuthorizedPasswordError  = 20110 //密码错误

	AdminCreateError             = 20201
	AdminListError               = 20202
	AdminDeleteError             = 20203
	AdminUpdateError             = 20204
	AdminResetPasswordError      = 20205
	AdminLoginError              = 20206 //登陆失败 状态异常
	AdminLogOutError             = 20207
	AdminModifyPasswordError     = 20208
	AdminModifyPersonalInfoError = 20209 //管理个人信息失败
	AdminMenuListError           = 20210
	AdminMenuCreateError         = 20211
	AdminOfflineError            = 20212
	AdminDetailError             = 20213 //个人信息获取失败
	AdminModifyDataError         = 20214 //数据操作失败

	MenuCreateError       = 20301
	MenuUpdateError       = 20302
	MenuListError         = 20303
	MenuDeleteError       = 20304
	MenuDetailError       = 20305
	MenuCreateActionError = 20306
	MenuListActionError   = 20307
	MenuDeleteActionError = 20308

	TaskCreateError  = 20401
	TaskUpdateError  = 20402
	TaskListError    = 20403
	TaskDetailError  = 20404
	TaskExecuteError = 20405
)

func TextCode(code int) string {
	lang := configs.GCONFIG.Language.Local

	if lang == configs.ZhCN {
		return ZhCNText[code]
	}

	if lang == configs.EnUS {
		return EnUSText[code]
	}

	return ZhCNText[code]
}
