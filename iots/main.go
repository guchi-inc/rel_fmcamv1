// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package iots

import (
	"fmcam/ctrlib/clients/tasks"
	"log"
	"os"
)

/*
管理iot 层, 与 systems 数据管理层的下级。
监听 mqtt 设备服务的数据，并将相关数据发布到 redis，最后由redis转发到 相应接口。

1  转发 modbus mqtt 信号
2  （待扩展） 温度 湿度信号
3  打印机驱动和打印接口
*/

var (
	/*printer 相关参数初始化*/
	host      = "https://dlabelplugga.ctaiot.com:9443"
	path      = "/api/sdk/get/token.json"
	appKey    = "516051646435762"
	appSecret = "277fa7eb6461470d9a044732a088b031"
	//接入方用户id,接入客户端应用下用户id保持唯一
	userId    = "1"
	timestamp = "2023-06-06 11:40:58"
	allTask   = tasks.InitTasks()
	logger    = log.New(os.Stdout, "INFO -", 18)
)

type Iots struct{}

func NewIots() *Iots {
	return &Iots{}
}

var (
	IotsApp = NewIots()
)
