// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package config

import (
	"fmcam/models/erps"

	"go.uber.org/zap"
)

type Server struct {
	JWT JWT `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	// Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Clusters Clusters `mapstructure:"clusters" json:"clusters" yaml:"clusters"`

	//mqtt
	Mqtts Mqtts `mapstructure:"mqtts" json:"mqtts" yaml:"mqtts"`

	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	AutoCode Autocode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
	// sqlx
	//MysqlDsn string  `mapstructure:"mysql_dsn" json:"mysql_dsn"  `
	Mysql MysqlDB `mapstructure:"mysql" json:"mysql" yaml:"mysql"`

	// Pgsql  Pgsql           `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	Local Local `mapstructure:"local" json:"local" yaml:"local"`
	AwsS3 AwsS3 `mapstructure:"aws-s3" json:"aws-s3" yaml:"aws-s3"`

	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	//本地化
	Language Language `mapstructure:"language" json:"language" yaml:"language"`

	HashIds HashIds `mapstructure:"hashIds" json:"hashIds" yaml:"hashIds"`
	//默认日志程序
	ServerLog ServerLog `mapstructure:"server_log" json:"server_log" yaml:"server_log"`

	//定时任务 解析配置
	Anysis Anysis `mapstructure:"anysis" json:"anysis" yaml:"anysis"`

	//打印模板配置
	PrinterTmp erps.Printer `mapstructure:"printer" json:"printer" yaml:"printer"`
}

type ServerLog struct {
	*zap.SugaredLogger
}

func NewServerLog(log *zap.SugaredLogger) *ServerLog {
	return &ServerLog{log}
}
