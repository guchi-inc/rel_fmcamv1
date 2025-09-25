// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package code

//创建
type CreateRequest struct {
	Name                string `form:"name" binding:"required"`           // 任务名称
	Spec                string `form:"spec" binding:"required"`           // crontab 表达式
	Command             string `form:"command" binding:"required"`        // 执行命令
	Protocol            int32  `form:"protocol" binding:"required"`       // 执行方式 1:shell 2:http
	HttpMethod          int32  `form:"http_method"`                       // http 请求方式 1:get 2:post
	Timeout             int32  `form:"timeout" binding:"required"`        // 超时时间(单位:秒)
	RetryTimes          int32  `form:"retry_times" binding:"required"`    // 重试次数
	RetryInterval       int32  `form:"retry_interval" binding:"required"` // 重试间隔(单位:秒)
	NotifyStatus        int32  `form:"notify_status" binding:"required"`  // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyType          int32  `form:"notify_type"`                       // 通知类型 1:邮件 2:webhook
	NotifyReceiverEmail string `form:"notify_receiver_email"`             // 通知者邮箱地址(多个用,分割)
	NotifyKeyword       string `form:"notify_keyword"`                    // 通知匹配关键字(多个用,分割)
	Remark              string `form:"remark"`                            // 备注
	IsUsed              int32  `form:"is_used" binding:"required"`        // 是否启用 1:是  -1:否
}

type CreateResponse struct {
	Id int32 `json:"id"` // 主键ID
}

//查细节
type DetailRequest struct {
	Id string `uri:"id"` // HashID
}

type DetailResponse struct {
	Name                string `json:"name"`                  // 任务名称
	Spec                string `json:"spec"`                  // crontab 表达式
	Command             string `json:"command"`               // 执行命令
	Protocol            int32  `json:"protocol"`              // 执行方式 1:shell 2:http
	HttpMethod          int32  `json:"http_method"`           // http 请求方式 1:get 2:post
	Timeout             int32  `json:"timeout"`               // 超时时间(单位:秒)
	RetryTimes          int32  `json:"retry_times"`           // 重试次数
	RetryInterval       int32  `json:"retry_interval"`        // 重试间隔(单位:秒)
	NotifyStatus        int32  `json:"notify_status"`         // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyType          int32  `json:"notify_type"`           // 通知类型 1:邮件 2:webhook
	NotifyReceiverEmail string `json:"notify_receiver_email"` // 通知者邮箱地址(多个用,分割)
	NotifyKeyword       string `json:"notify_keyword"`        // 通知匹配关键字(多个用,分割)
	Remark              string `json:"remark"`                // 备注
	IsUsed              int32  `json:"is_used"`               // 是否启用 1:是  -1:否
}

//手动执行
type ExecuteRequest struct {
	Id string `uri:"id"` // HashID
}

type ExecuteResponse struct {
	Id int `json:"id"` // ID
}

//查cron 任务列表
type ListRequest struct {
	Page     int    `form:"page"`      // 第几页
	PageSize int    `form:"page_size"` // 每页显示条数
	Name     string `form:"name"`      // 任务名称
	Protocol int    `form:"protocol"`  // 执行方式 1:shell 2:http
	IsUsed   int    `form:"is_used"`   // 是否启用 1:是  -1:否
}

type ListData struct {
	Id               int    `json:"id"`                 // ID
	HashID           string `json:"hashid"`             // hashid
	Name             string `json:"name"`               // 任务名称
	Protocol         int    `json:"protocol"`           // 执行方式 1:shell 2:http
	ProtocolText     string `json:"protocol_text"`      // 执行方式
	Spec             string `json:"spec"`               // crontab 表达式
	Command          string `json:"command"`            // 执行命令
	HttpMethod       int    `json:"http_method"`        // http 请求方式 1:get 2:post
	HttpMethodText   string `json:"http_method_text"`   // http 请求方式
	Timeout          int    `json:"timeout"`            // 超时时间(单位:秒)
	RetryTimes       int    `json:"retry_times"`        // 重试次数
	RetryInterval    int    `json:"retry_interval"`     // 重试间隔(单位:秒)
	NotifyStatus     int    `json:"notify_status"`      // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyStatusText string `json:"notify_status_text"` // 执行结束是否通知
	IsUsed           int    `json:"is_used"`            // 是否启用 1=启用 2=禁用
	IsUsedText       string `json:"is_used_text"`       // 是否启用
	CreatedAt        string `json:"created_time"`       // 创建时间
	CreatedUser      string `json:"created_user"`       // 创建人
	UpdatedAt        string `json:"updated_time"`       // 更新时间
	UpdatedUser      string `json:"updated_user"`       // 更新人
}

type ListResponse struct {
	Total      int64      `json:"total"`
	Size       int64      `json:"size"`
	Data       []ListData `json:"list"`
	Pagination struct {
		Total        int `json:"total"`
		CurrentPage  int `json:"current_page"`
		PerPageCount int `json:"per_page_count"`
	} `json:"pagination"`
}

//修改编辑
type ModifyRequest struct {
	Id                  string `form:"id" binding:"required"`             // 任务ID
	Name                string `form:"name" binding:"required"`           // 任务名称
	Spec                string `form:"spec" binding:"required"`           // crontab 表达式
	Command             string `form:"command" binding:"required"`        // 执行命令
	Protocol            int32  `form:"protocol" binding:"required"`       // 执行方式 1:shell 2:http
	HttpMethod          int32  `form:"http_method"`                       // http 请求方式 1:get 2:post
	Timeout             int32  `form:"timeout" binding:"required"`        // 超时时间(单位:秒)
	RetryTimes          int32  `form:"retry_times" binding:"required"`    // 重试次数
	RetryInterval       int32  `form:"retry_interval" binding:"required"` // 重试间隔(单位:秒)
	NotifyStatus        int32  `form:"notify_status" binding:"required"`  // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyType          int32  `form:"notify_type"`                       // 通知类型 1:邮件 2:webhook
	NotifyReceiverEmail string `form:"notify_receiver_email"`             // 通知者邮箱地址(多个用,分割)
	NotifyKeyword       string `form:"notify_keyword"`                    // 通知匹配关键字(多个用,分割)
	Remark              string `form:"remark"`                            // 备注
	IsUsed              int32  `form:"is_used" binding:"required"`        // 是否启用 1:是  -1:否
}

type ModifyResponse struct {
	Id int32 `json:"id"` // 主键ID
}

//更新状态
type UpdateStateRequest struct {
	Id   string `form:"id"`   // 主键ID
	Used int32  `form:"used"` // 是否启用 1:是 -1:否
}

type UpdateStateResponse struct {
	Id int32 `json:"id"` // 主键ID
}
