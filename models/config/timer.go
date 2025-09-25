//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package config

type Timer struct {
	Spec        string   `mapstructure:"spec" json:"spec" yaml:"spec"` // CRON表达式
	Detail      []Detail `mapstructure:"detail" json:"detail" yaml:"detail"`
	WithSeconds bool     `mapstructure:"with_seconds" json:"with_seconds" yaml:"with_seconds"` // 是否精确到秒
	Start       bool     `mapstructure:"start" json:"start" yaml:"start"`                      // 是否启用

}
type Detail struct {
	TableName    string `mapstructure:"tableName" json:"tableName" yaml:"tableName"`          // 需要清理的表名
	CompareField string `mapstructure:"compareField" json:"compareField" yaml:"compareField"` // 需要比较时间的字段
	Interval     string `mapstructure:"interval" json:"interval" yaml:"interval"`             // 时间间隔
}
