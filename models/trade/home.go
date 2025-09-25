//Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package trade

//首页相关参数
type HomePage struct {
	HomeText    string `json:"home_text,omitempty"`    //主页宣传语
	LogoUrl     string `json:"logo_url,omitempty"`     //logo地址
	DownloadUrl string `json:"download_url,omitempty"` //下载企业ppt简介的链接
	Contact     string `json:"contact,omitempty"`      //联系方式
	NewsList    any    `json:"news_list,omitempty"`    //新闻列表
}
