//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package config

type Redis struct {
	DB           int  `mapstructure:"db" json:"db" yaml:"db"` // redis的哪个数据库
	MaxRetries   int  `mapstructure:"max_retries" json:"max_retries" yaml:"max_retries"`
	PoolSize     int  `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	MinIdleConns int  `mapstructure:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns"`
	Enable       bool `mapstructure:"enable" json:"enable" yaml:"enable"` //是否启用缓存

	IsLocal bool `mapstructure:"is_local" json:"is_local" yaml:"is_local"` //是否本地缓存

	Addr        string `mapstructure:"addr" json:"addr" yaml:"addr"`                         // 服务器地址:端口，如果启用本地缓存则会被程序执行时修改
	Password    string `mapstructure:"password" json:"password" yaml:"password"`             // 密码
	PathLocal   string `mapstructure:"path_local" json:"path_local" yaml:"path_local"`       //本地mqtt代理路径
	ConfigLocal string `mapstructure:"config_local" json:"config_local" yaml:"config_local"` //本地mqtt代理路径

}

type Clusters struct {
	Addrs    string `mapstructure:"addrs" json:"addrs" yaml:"addrs"`          // 服务器主节点地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库

}
