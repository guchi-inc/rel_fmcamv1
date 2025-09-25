package erps

//省份查询列表
type FmGovProvinceList struct {
	Total uint            `json:"total" db:"total"`
	Size  uint            `json:"size" db:"size"`
	Data  []FmGovProvince `json:"data" db:"data"`
}

// 省份
type FmGovProvince struct {
	ID   int64  `db:"id" json:"id" form:"id"`       //  主键
	Code string `db:"code" json:"code" form:"code"` //  编号
	Name string `db:"name" json:"name" form:"name"` //  名称
}

//区县查询列表
type FmGovCityList struct {
	Total uint        `json:"total" db:"total"`
	Size  uint        `json:"size" db:"size"`
	Data  []FmGovCity `json:"data" db:"data"`
}

// 城市
type FmGovCity struct {
	ID           int64  `db:"id" json:"id" form:"id"`                                  //  主键
	Code         string `db:"code" json:"code" form:"code"`                            //  编号
	Name         string `db:"name" json:"name" form:"name"`                            //  名称
	ProvinceCode string `db:"province_code" json:"province_code" form:"province_code"` //  省份编号

}

//区县查询列表
type FmGovAreaList struct {
	Total uint        `json:"total" db:"total"`
	Size  uint        `json:"size" db:"size"`
	Data  []FmGovArea `json:"data" db:"data"`
}

// 区县
type FmGovArea struct {
	ID           int64  `db:"id" json:"id" form:"id"`                                  //  主键
	Code         string `db:"code" json:"code" form:"code"`                            //  编号
	Name         string `db:"name" json:"name" form:"name"`                            //  名称
	ProvinceCode string `db:"province_code" json:"province_code" form:"province_code"` //  省份编号
	CityCode     string `db:"city_code" json:"city_code" form:"city_code"`             //  城市编号

}

//街道查询列表
type FmGovStreetList struct {
	Total uint          `json:"total" db:"total"`
	Size  uint          `json:"size" db:"size"`
	Data  []FmGovStreet `json:"data" db:"data"`
}

// 街道
type FmGovStreet struct {
	ID           int64  `db:"id" json:"id" form:"id"`                                  //  主键
	Code         string `db:"code" json:"code" form:"code"`                            //  编号
	Name         string `db:"name" json:"name" form:"name"`                            //  名称
	ProvinceCode string `db:"province_code" json:"province_code" form:"province_code"` //  省份编号
	CityCode     string `db:"city_code" json:"city_code" form:"city_code"`             //  城市编号
	AreaCode     string `db:"area_code" json:"area_code" form:"area_code"`             //  城市编号

}
