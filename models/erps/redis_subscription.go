//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

type KeysOfRedisPub struct {
	Action   string `json:"action,omitempty" form:"action,omitempty"`
	Topic    string `json:"topic,omitempty" form:"topic,omitempty"`
	Meter    string `json:"meter" form:"meter"`
	Operator string `json:"operator" form:"operator"`
}

// ws实时消息的 redis 格式
func NewKeysOfRedisPub(name, detail string) *KeysOfRedisPub {
	var kr = KeysOfRedisPub{}
	kr.Meter = detail
	return &kr
}

// 标记结构体
type FlawOfRedisPub struct {
	Action   string `json:"action,omitempty" form:"action,omitempty"`
	Topic    string `json:"topic,omitempty" form:"topic,omitempty"`
	Flaw     string `json:"flaw_name" form:"flaw_name"`
	Operator string `json:"operator" form:"operator"`
}
