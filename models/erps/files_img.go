// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

//base64格式
type ImageUploadRequest struct {
	ImgParent  string `json:"img_parent"`
	Base64Data string `json:"image_base64"`
}

// 图像图形
type ImgInfo struct {
	ID      int64 `db:"id" json:"id" form:"id"`                //  主键
	Enabled bool  `db:"enabled" json:"enabled" form:"enabled"` //状态，1：正常，0 封禁

	Picture []byte `db:"picture" json:"picture" form:"picture"`

	ImgParent string `db:"img_parent" json:"img_parent" form:"img_parent"`
	LocalPath string `db:"local_path" json:"local_path" form:"local_path"`
	Creator   string `db:"creator" json:"creator" form:"creator"`

	CreatedAt string `db:"created_time" json:"created_time" form:"created_time"`
	UpdatedAt string `db:"updated_time" json:"updated_time" form:"updated_time"`
}

type ImgBase64 struct {
	ImgParent string `db:"img_parent" json:"img_parent" form:"img_parent"`
	Picture   string `db:"picture" json:"picture" form:"picture"`
	Url       string `db:"url" json:"url,omitempty" form:"url"`
}

//二进制数据
type ImgBinary struct {
	ImgParent string `db:"img_parent" json:"img_parent" form:"img_parent"`
	Picture   []byte `db:"picture" json:"picture" form:"picture"`
	Url       string `db:"url" json:"url,omitempty" form:"url"`
}

// 促销商品信息列表
type CutDownInfoList struct {
	Total   uint          `json:"total" db:"total"`
	Size    uint          `json:"size" db:"size"`
	Data    []CutDownInfo `json:"data" db:"data"`
	Columns []GcDesc      `json:"columns" db:"columns"`
}

// 促销商品信息
type CutDownInfo struct {
	MaterialTypeId string `db:"material_type_id" json:"material_type_id" form:"material_type_id"` //类型ID
	ImgParent      string `db:"img_parent" json:"img_parent" form:"img_parent"`                   //图片名
	Picture        []byte `db:"picture" json:"picture" form:"picture"`                            //类图信息
	Name           string `db:"name" json:"name" form:"name"`                                     //商品名称
	Remark         string `db:"remark" json:"remark" form:"remark"`                               //折扣率  -20 百分之二十折扣
	UpdateTime     string `db:"updated_time" json:"updated_time" form:"updated_time"`             //更新时间
}

// 促销商品信息
type CutDownBase64List struct {
	Code    int                 `json:"code" `
	Total   uint                `json:"total" db:"total"`
	Size    uint                `json:"size" db:"size"`
	Data    []CutDownBase64Info `json:"data" db:"data"`
	Columns []GcDesc            `json:"columns" db:"columns"`
}

type CutDownBase64Info struct {
	MaterialTypeId string `db:"material_type_id" json:"material_type_id" form:"material_type_id"` //类型ID
	ImgParent      string `db:"img_parent" json:"img_parent" form:"img_parent"`
	Picture        string `db:"picture" json:"picture" form:"picture"`
	Name           string `db:"name" json:"name" form:"name"`                         //商品名称
	Remark         string `db:"remark" json:"remark" form:"remark"`                   //折扣率  -20 百分之二十折扣
	UpdateTime     string `db:"updated_time" json:"updated_time" form:"updated_time"` //更新时间
}
