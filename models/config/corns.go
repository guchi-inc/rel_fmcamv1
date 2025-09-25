//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package config

//本地语言
type Language struct {
	Local string `mapstructure:"local" json:"local" yaml:"local"`
}

type HashIds struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Length int    `mapstructure:"length" json:"length" yaml:"length"`
}

// 解析 进程
type Anysis struct {
	Host       string `mapstructure:"host" json:"host" yaml:"host"`                      // "192.168.30.131"  #报告文件服务
	CondaPath  string `mapstructure:"conda_path" json:"conda_path" yaml:"conda_path"`    // "C:\\Users\\kukeg\\Miniconda3\\scripts\\activate.bat"  # 执行者环境
	AnysisProc string `mapstructure:"anysis_proc" json:"anysis_proc" yaml:"anysis_proc"` // "D:\\faceswaps\\faceswap.py"  # 提取工具路径
	AnysisOut  string `mapstructure:"anysis_out" json:"anysis_out" yaml:"anysis_out"`    // "D:\\vedio\\deep_face_out\\"  # 输出路径，不是最终路径，需要对每次提取创建所属路径
}
