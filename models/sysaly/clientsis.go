package sysaly

import "fmcam/models/erps"

//客户端 查询数据返回聚合到 24小时
type HourSisList[T any] struct {
	TotalData map[string]T  `json:"total_data" db:"total_data"`         //返回的总数
	Size      uint          `json:"page_size" db:"page_size"`           //几天内的，size =1 表示今天，如果 size 为 -1，同时data不为空，则表示全部总的数据
	Monitor   uint          `json:"monitor" db:"monitor"`               //0 不展示，1 展示
	DataType  string        `json:"data_type" db:"data_type"`           // 数据类型， 累计 accrued 或 区间 date
	Data      map[string]T  `json:"data" db:"data"`                     //数据
	Columns   []erps.GcDesc `json:"columns" db:"columns"`               //字段说明
	TenantId  *string       `json:"tenant_id,omitempty" db:"tenant_id"` //租户号
}

//客户端 查询数据返回
type SysClientList[T any] struct {
	TotalData map[string]T `json:"total_data" db:"total_data"`     //返回的总数
	Size      uint         `json:"page_size" db:"page_size"`       //几天内的，size =1 表示今天，如果 size 为 -1，同时data不为空，则表示全部总的数据
	Monitor   uint         `json:"monitor,omitempty" db:"monitor"` //0 不展示，1 展示
	DataType  string       `json:"data_type" db:"data_type"`       // 数据类型， user 或 date
	Data      map[string]T `json:"data" db:"data"`                 //数据
	Columns   any          `json:"columns" db:"columns"`           //字段说明
}

//按人员或日期 类型的统计值
// 筛选抓拍面部数据
type DailyCaptureCount struct {
	Total        int     `json:"count" db:"count"`
	TimeRange    string  `json:"time_range" db:"time_range"`
	DataType     string  `json:"data_type" db:"data_type"`
	Stats        *string `json:"stats,omitempty"  db:"stats"`
	UserCategory *string `json:"user_category,omitempty" db:"user_category"`
	Url          *string `json:"url,omitempty" db:"url"`
	Position     *string `json:"position,omitempty" db:"position"`
	TenantId     string  `json:"tenant_id,omitempty" db:"tenant_id"`
}

// 抓拍数据统计
type CaptureStats struct {
	Count    int    `json:"count"`                    // 抓拍数量
	DataType string `json:"data_type" db:"data_type"` // 数据类型， user 或 date

	Date      string `json:"date"`       // yyyy-mm-dd
	TimeRange string `json:"time_range"` // 时间范围标识，例如 "7d"
}

// 人员分类数据统计
type PersonStats struct {
	Count    int    `json:"count"`                    // 抓拍数量
	DataType string `json:"data_type" db:"data_type"` // 数据类型， user 或 date

	Category  string `json:"user_category"` // 比如 "visitor", "employee", "blacklist"
	TimeRange string `json:"time_range"`    // 时间范围标识，例如 "7d"

	Url string `json:"url"` // 图形地址
}

// 位置分类数据统计
type PositionStats struct {
	ID       int    `json:"id" db:"id"`               // 序号
	Count    *int   `json:"count" db:"count"`         //  数量
	AllCount *int64 `json:"all_count" db:"all_count"` //  该列累计数量

	DataType  *string `json:"data_type" db:"data_type"`          // 数据类型， position 或 date
	Category  string  `json:"user_category"  db:"user_category"` // 比如 "visitor", "employee", "blacklist"
	Position  *string `json:"position" db:"position"`            // 比如 "visitor", "employee", "blacklist"
	TimeRange string  `json:"time_range" db:"time_range"`        // 时间范围标识，例如 "D7"
	TypeId    *string `json:"type_id" db:"type_id"`              //人员类型 id 统计维度

	DeviceId *string `json:"device_id" db:"device_id"` //二维设备id 统计维度
	Stats    *string `json:"stats" db:"stats"`         // 统计维度  可能是设备id 或 天/周/月时间rang_type
	Url      string  `json:"url" db:"url"`             // 图形地址
}

// 预警分类数据统计
type WarnStats struct {
	Count    int    `json:"count" db:"count"`         // 抓拍数量
	DataType string `json:"data_type" db:"data_type"` // 数据类型， position 或 date

	LevelType string `json:"level_type" db:"level_type"` // 数据类型，  level 1 ~ 4
	Position  string `json:"position" db:"position"`     // 比如 "visitor", "employee", "blacklist"

	DateRange string `json:"date_range" db:"date_range"` // 比如 今日，累计，待处理
	TimeRange string `json:"time_range" db:"time_range"` // 时间范围标识，例如 "7d"
}

// 抓拍轨迹 封装 查询结果
type CaptionHistoryList struct {
	Total   int                  `json:"total" db:"total"`
	Size    int                  `json:"page_size" db:"page_size"` //几天内的，size =1 表示今天，如果 size 为 -1，同时data不为空，则表示全部总的数据
	Data    []CaptionHistoryData `json:"data" db:"data"`           //查询的数据
	Columns []erps.GcDesc        `json:"columns" db:"columns"`     //字段说明
}

type CaptionHistoryData struct {
	ID               int     `json:"id" db:"id"`                                 // 序号
	DeviceId         *int    `json:"device_id" db:"device_id"`                   // 设备id
	FaceId           *int    `json:"face_id" db:"face_id"`                       // 面部id
	Location         *string `json:"location" db:"location"`                     // 抓拍设备地址
	CaptureTime      *string `json:"capture_time" db:"capture_time"`             // 抓拍时间
	MatchedProfileId *string `json:"matched_profile_id" db:"matched_profile_id"` //matched_profile_id
	RoomId           *string `json:"room_id" db:"room_id"`                       //room_id
	CaptureImageUrl  *string `json:"capture_image_url" db:"capture_image_url"`   // 抓拍图片地址

}

/*
{
  "capture_by_time": {
    "7d": [
      {"date": "2025-07-01", "count": 120, "time_range": "7d"},
      {"date": "2025-07-02", "count": 98, "time_range": "7d"}
    ],
    "30d": [
      {"date": "2025-06-10", "count": 50, "time_range": "30d"}
    ]
  },
  "person_by_category": {
    "7d": [
      {"date": "2025-07-01", "category": "visitor", "count": 10, "time_range": "7d"},
      {"date": "2025-07-01", "category": "employee", "count": 30, "time_range": "7d"}
    ],
    "30d": [
      {"date": "2025-06-15", "category": "blacklist", "count": 5, "time_range": "30d"}
    ]
  }
}

*/
