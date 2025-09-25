//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package mqtts

const (
	HealhStatusUp   = "UP" //如果数据库连接和mqtt服务连接都完好，则返回UP
	HealhStatusDown = "DOWN"
)

// 健康检查
type Health struct {
	Status string `json:"status"`
}

// gcproxyer 的结构体 接受
type ChipEvent struct {
	Action  string `json:"action,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Message string `json:"message"`
}

// 赋值新消息
func InitChipEvent(topic string, ce *ChipEvent) *ChipEvent {
	ce.Topic = topic
	return ce
}
