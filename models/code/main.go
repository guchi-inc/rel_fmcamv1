// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package code

import (
	"log"
	"os"
	"time"
)

var (
	Logger                       = log.New(os.Stderr, "DEBUG -", 13)
	FieldCreateCronTaskWithTable = []string{"name", "spec", "command"}
)

type CreateCronTaskWithTable struct {
	Id                  int64  `json:"id"`
	TableName           string `json:"table_name"`
	PriKey              string `json:"pri_key"`
	Name                string `json:"name"`    // 任务名称
	Spec                string `json:"spec"`    // crontab 表达式
	Command             string `json:"command"` // 执行命令
	Protocol            int64  // 执行方式 1:shell 2:http
	HttpMethod          int64  // http 请求方式 1:get 2:post
	Timeout             int64  // 超时时间(单位:秒)
	RetryTimes          int64  // 重试次数
	RetryInterval       int64  // 重试间隔(单位:秒)
	NotifyStatus        int64  // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyType          int64  // 通知类型 1:邮件 2:webhook
	NotifyReceiverEmail string // 通知者邮箱地址(多个用,分割)
	NotifyKeyword       string // 通知匹配关键字(多个用,分割)
	Remark              string // 备注
	IsUsed              int64  // 是否启用 1:是  -1:否
}

type RedisCronTask struct {
	Id                  int64     // 主键
	Name                string    // 任务名称
	Spec                string    // crontab 表达式
	Command             string    // 执行命令
	Protocol            int64     // 执行方式 1:shell 2:http
	HttpMethod          int64     // http 请求方式 1:get 2:post
	Timeout             int64     // 超时时间(单位:秒)
	RetryTimes          int64     // 重试次数
	RetryInterval       int64     // 重试间隔(单位:秒)
	NotifyStatus        int64     // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyType          int64     // 通知类型 1:邮件 2:webhook
	NotifyReceiverEmail string    // 通知者邮箱地址(多个用,分割)
	NotifyKeyword       string    // 通知匹配关键字(多个用,分割)
	Remark              string    // 备注
	IsUsed              int64     // 是否启用 1:是  -1:否
	CreatedAt           time.Time `json:"time"` // 创建时间
	CreatedUser         string    // 创建人
	UpdatedAt           time.Time `json:"update_time"` // 更新时间
	UpdatedUser         string    // 更新人
}

type RedisCronTaskData struct {
	Total int64
	Data  []CreateCronTaskWithTable
	Size  int64
}

var EnUSText = map[int]string{
	ServerError:        "Internal server error",
	TooManyRequests:    "Too many requests",
	ParamBindError:     "Parameter error",
	AuthorizationError: "Authorization error",
	UrlSignError:       "URL signature error",
	CacheSetError:      "Failed to set cache",
	CacheGetError:      "Failed to get cache",
	CacheDelError:      "Failed to del cache",
	CacheNotExist:      "Cache does not exist",
	ResubmitError:      "Please do not submit repeatedly",
	HashIdsEncodeError: "HashID encryption failed",
	HashIdsDecodeError: "HashID decryption failed",
	RBACError:          "No access",
	RedisConnectError:  "Failed to connection Redis",
	MySQLConnectError:  "Failed to connection MySQL",
	WriteConfigError:   "Failed to write configuration file",
	SendEmailError:     "Failed to send mail",
	MySQLExecError:     "SQL execution failed",
	GoVersionError:     "Go Version mismatch",
	SocketConnectError: "Socket not connected",
	SocketSendError:    "Socket message sending failed",

	AuthorizedCreateError:    "Failed to create caller",
	AuthorizedListError:      "Failed to get caller list",
	AuthorizedDeleteError:    "Failed to delete caller",
	AuthorizedUpdateError:    "Failed to update caller",
	AuthorizedDetailError:    "Failed to get caller details",
	AuthorizedCreateAPIError: "Failed to create caller API address",
	AuthorizedListAPIError:   "Failed to get caller API address list",
	AuthorizedDeleteAPIError: "Failed to delete caller API address",

	AdminCreateError:             "Failed to create administrator",
	AdminListError:               "Failed to get administrator list",
	AdminDeleteError:             "Failed to delete administrator",
	AdminUpdateError:             "Failed to update administrator",
	AdminResetPasswordError:      "Reset password failed",
	AdminLoginError:              "Login failed",
	AdminLogOutError:             "Exit failed",
	AdminModifyPasswordError:     "Failed to modify password",
	AdminModifyPersonalInfoError: "Failed to modify personal information",
	AdminMenuListError:           "Failed to get administrator menu authorization list",
	AdminMenuCreateError:         "Administrator menu authorization failed",
	AdminOfflineError:            "Offline administrator failed",
	AdminDetailError:             "Failed to get personal information",

	MenuCreateError:       "Failed to create menu",
	MenuUpdateError:       "Failed to update menu",
	MenuDeleteError:       "Failed to delete menu",
	MenuListError:         "Failed to get menu list",
	MenuDetailError:       "Failed to get menu details",
	MenuCreateActionError: "Failed to create menu action",
	MenuListActionError:   "Failed to get menu action list",
	MenuDeleteActionError: "Failed to delete menu action",

	TaskCreateError:  "Failed to create cron",
	TaskUpdateError:  "Failed to update menu",
	TaskListError:    "Failed to get cron list",
	TaskDetailError:  "Failed to get cron detail",
	TaskExecuteError: "Failed to execute cron",
}

var ZhCNText = map[int]string{
	ServerError:        "内部服务器错误",
	TooManyRequests:    "请求过多",
	ParamBindError:     "参数信息错误",
	AuthorizationError: "认证签名信息错误",
	UrlSignError:       "参数签名错误",
	CacheSetError:      "设置缓存失败",
	CacheGetError:      "获取缓存失败",
	CacheDelError:      "删除缓存失败",
	CacheNotExist:      "缓存不存在",
	ResubmitError:      "请勿重复提交",
	HashIdsEncodeError: "HashID 加密失败",
	HashIdsDecodeError: "HashID 解密失败",
	RBACError:          "暂无访问权限",
	RedisConnectError:  "Redis 连接失败",
	MySQLConnectError:  "MySQL 连接失败",
	WriteConfigError:   "写入配置文件失败",
	SendEmailError:     "发送邮件失败",
	MySQLExecError:     "SQL 执行失败",
	GoVersionError:     "Go 版本不满足要求",
	SocketConnectError: "Socket 未连接",
	SocketSendError:    "Socket 消息发送失败",
	ParamError:         "参数校验错误",

	AuthorizedCreateError:    "创建调用方失败",
	AuthorizedListError:      "获取调用方信息失败",
	AuthorizedDeleteError:    "删除调用方失败",
	AuthorizedUpdateError:    "更新调用方失败",
	AuthorizedDetailError:    "调用方认证失败,请重新认证。",
	AuthorizedCreateAPIError: "创建调用方 API 地址失败",
	AuthorizedListAPIError:   "获取调用方 API 地址列表失败",
	AuthorizedDeleteAPIError: "删除调用方 API 地址失败",
	AuthorizedPowerError:     "调用方权限不足,请联系管理员。",
	AuthorizedPasswordError:  "密码错误,请重试。",

	AdminCreateError:             "创建管理员失败",
	AdminListError:               "获取列表失败",
	AdminDeleteError:             "删除管理员失败",
	AdminUpdateError:             "更新管理员失败",
	AdminResetPasswordError:      "重置密码失败",
	AdminLoginError:              "状态异常，登录失败",
	AdminLogOutError:             "退出失败",
	AdminModifyPasswordError:     "修改密码失败",
	AdminModifyPersonalInfoError: "管理个人信息失败",
	AdminMenuListError:           "获取管理员菜单授权列表失败",
	AdminMenuCreateError:         "管理员菜单授权失败",
	AdminOfflineError:            "下线管理员失败",
	AdminDetailError:             "个人信息获取失败",
	AdminModifyDataError:         "数据修改失败",

	MenuCreateError:       "创建菜单失败",
	MenuUpdateError:       "更新菜单失败",
	MenuDeleteError:       "删除菜单失败",
	MenuListError:         "获取菜单列表失败",
	MenuDetailError:       "获取菜单详情失败",
	MenuCreateActionError: "创建菜单栏功能权限失败",
	MenuListActionError:   "获取菜单栏功能权限列表失败",
	MenuDeleteActionError: "删除菜单栏功能权限失败",

	TaskCreateError:  "创建后台任务失败",
	TaskUpdateError:  "更新后台任务失败",
	TaskListError:    "获取任务列表失败",
	TaskDetailError:  "获取任务详情失败",
	TaskExecuteError: "手动执行任务失败",
}
