// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package code

import "encoding/json"

// SessionUserInfo 当前用户会话信息
type SessionUserInfo struct {
	UserID   int32  `json:"user_id"`   // 用户ID
	UserName string `json:"user_name"` // 用户名
}

// Marshal 序列化到JSON
func (user *SessionUserInfo) Marshal() (jsonRaw []byte) {
	jsonRaw, _ = json.Marshal(user)
	return
}

// AliasForRecordMetrics 对请求路径起个别名，用于记录指标。
// 如：Get /user/:username 这样的路径，因为 username 会有非常多的情况，这样记录指标非常不友好。
func AliasForRecordMetrics(path string) HandlerFunc {
	return func(ctx Context) {
		ctx.setAlias(path)
	}
}
