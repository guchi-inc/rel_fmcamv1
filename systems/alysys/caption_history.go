package alysys

import (
	"fmcam/common/databases"
	"fmcam/models/sysaly"
	"fmt"
)

// 抓拍 历史轨迹
func (mt *AlySysOperator) GetCaptionHostory(Page, PageSize, profile_id int, room_id, startAt, endAt, idCardNumber, TenantId string) (*sysaly.CaptionHistoryList, error) {

	var (
		captionHistory = sysaly.CaptionHistoryList{Size: PageSize}
		DBConn         = databases.DBMysql

		baseSql = `SELECT a.id, a.capture_time, a.capture_image_url,a.device_id,  
		d.location, a.matched_profile_id , pf.room_id, f.id AS face_id		
		FROM  CaptureLogs a
		LEFT JOIN Device d ON d.id = a.device_id
		LEFT JOIN Profile pf ON pf.id = a.matched_profile_id
		LEFT JOIN Face f ON f.profile_id = pf.id  
		`
		groupSql = " GROUP BY    a.id,face_id 	 "

		filter = " WHERE a.id > 0 "
	)

	if profile_id != 0 {
		filter += fmt.Sprintf(" AND a.matched_profile_id = %v  ", profile_id)
	}

	if idCardNumber != "" {
		filter += fmt.Sprintf(" AND pf.id_card_number = %v  ", idCardNumber)
	}

	//按房号或时间
	if room_id != "" {
		filter += fmt.Sprintf(" AND pf.room_id = '%v' ", room_id)
	}

	if startAt != "" && endAt != "" {
		filter += fmt.Sprintf(" AND (a.capture_time between '%v' AND '%v') ", startAt, endAt)
	}

	var total int
	totalSql := " SELECT COUNT(*) FROM CaptureLogs a LEFT JOIN Device d ON d.id = a.device_id LEFT JOIN Profile pf ON pf.id = a.matched_profile_id "
	totalSql += filter
	err := DBConn.Get(&total, totalSql)
	logsys.Printf("history totalSql:%v  err:%v \n", totalSql, err)
	if err != nil {
		return nil, err
	}

	Limits := fmt.Sprintf(" ORDER BY   a.id,face_id  ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += filter
	baseSql += groupSql
	baseSql += Limits
	err = DBConn.Select(&captionHistory.Data, baseSql)
	logsys.Printf("history baseSql:%v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(captionHistory.Data) > 0 {
		descColumns, _ := FmCaptionLogMap()
		captionHistory.Total = total
		captionHistory.Columns = descColumns
	}
	return &captionHistory, nil
}
