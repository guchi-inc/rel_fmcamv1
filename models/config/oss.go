//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package config

type Local struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`                   // 本地文件访问路径
	StorePath string `mapstructure:"store-path" json:"store-path" yaml:"store-path"` // 本地文件存储路径
}

type Qiniu struct {
	Zone          string `mapstructure:"zone" json:"zone" yaml:"zone"`                                  // 存储区域
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`                            // 空间名称
	ImgPath       string `mapstructure:"img-path" json:"img-path" yaml:"img-path"`                      // CDN加速域名
	AccessKey     string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`                // 秘钥AK
	SecretKey     string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`                // 秘钥SK
	UseCdnDomains bool   `mapstructure:"use-cdn-domains" json:"use-cdn-domains" yaml:"use-cdn-domains"` // 上传是否使用CDN上传加速
	UseHTTPS      bool   `mapstructure:"use-https" json:"use-https" yaml:"use-https"`                   // 是否使用https

}
type AliyunOSS struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" json:"access-key-secret" yaml:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name" json:"bucket-name" yaml:"bucket-name"`
	BucketUrl       string `mapstructure:"bucket-url" json:"bucket-url" yaml:"bucket-url"`
	BasePath        string `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}

type HuaWeiObs struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Endpoint  string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKey string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
}
type TencentCOS struct {
	Bucket     string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region     string `mapstructure:"region" json:"region" yaml:"region"`
	SecretID   string `mapstructure:"secret-id" json:"secret-id" yaml:"secret-id"`
	SecretKey  string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	BaseURL    string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	PathPrefix string `mapstructure:"path-prefix" json:"path-prefix" yaml:"path-prefix"`
}
type AwsS3 struct {
	Bucket           string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region           string `mapstructure:"region" json:"region" yaml:"region"`
	Endpoint         string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	SecretID         string `mapstructure:"secret-id" json:"secret-id" yaml:"secret-id"`
	SecretKey        string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	BaseURL          string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	PathPrefix       string `mapstructure:"path-prefix" json:"path-prefix" yaml:"path-prefix"`
	S3ForcePathStyle bool   `mapstructure:"s3-force-path-style" json:"s3-force-path-style" yaml:"s3-force-path-style"`
	DisableSSL       bool   `mapstructure:"disable-ssl" json:"disable-ssl" yaml:"disable-ssl"`
}
