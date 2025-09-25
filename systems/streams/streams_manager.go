package streams

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/iotmodel"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// 组流设备修改
func (sp *StreamRep) StreamDeviceUpdate(fst *iotmodel.FmStreamTasks) (int64, error) {

	var (
		DBConn       = databases.DBMysql
		deviceOrigin = iotmodel.FmStreamTasks{}
		QuerySql     = "select * from fm_tasks "
	)

	if fst.ID < 0 {
		return 0, fmt.Errorf("仅支持按ID修改")
	}

	err := DBConn.Get(&deviceOrigin, QuerySql+fmt.Sprintf(" WHERE id = '%v' ", fst.ID))

	if err != nil {
		err = fmt.Errorf("数据不存在:%v", err)
		return 0, err
	}

	timeNow := time.Now().In(configs.Loc)
	tStr := timeutil.CSTLayoutString(&timeNow)

	fst.UpdatedAt = tStr
	//优先更新状态，如果enabled状态为1，只更新状态
	if deviceOrigin.Enable != fst.Enable && (fst.Enable == "0" || fst.Enable == "1") {

		upSql := fmt.Sprintf("UPDATE fm_user SET enabled = '%v', updated_time = '%v'  WHERE id = '%v';", fst.Enable, tStr, fst.ID)

		rst, err := DBConn.Exec(upSql)
		syslog.Printf("from db update enable:%#v  err:%v \n", rst, err)
		if err != nil {
			return 0, err
		}

		return rst.RowsAffected()
	}

	//更新状态删除标记，如果delete_flag状态为1，只更新状态
	if deviceOrigin.DeleteFlag != fst.DeleteFlag && (fst.DeleteFlag == "0" || fst.DeleteFlag == "1") {

		upSql := fmt.Sprintf("UPDATE fm_user SET delete_flag = '%v', updated_time = '%v'  WHERE id = '%v';", fst.DeleteFlag, tStr, fst.ID)

		rst, err := DBConn.Exec(upSql)
		syslog.Printf("from db update delete flag:%#v  err:%v \n", rst, err)
		if err != nil {
			return 0, err
		}

		return rst.RowsAffected()
	}

	//只允许 管理员 根据 id 更新
	upTotal := 0

	//更新其他信息
	if fst.TaskType != "" && (deviceOrigin.TaskType != fst.TaskType) {
		deviceOrigin.TaskType = fst.TaskType
		upTotal += 1
	}

	if fst.TaskName != "" && (deviceOrigin.TaskName != fst.TaskName) {
		deviceOrigin.TaskName = fst.TaskName
		upTotal += 1
	}

	if fst.PageId != "" && (deviceOrigin.PageId != fst.PageId) {
		deviceOrigin.PageId = fst.PageId
		upTotal += 1
	}

	if fst.Creator != "" && (deviceOrigin.Creator != fst.Creator) {
		deviceOrigin.Creator = fst.Creator
		upTotal += 1
	}

	if fst.BasicUnit != "" && (deviceOrigin.BasicUnit != fst.BasicUnit) {
		deviceOrigin.BasicUnit = fst.BasicUnit
		upTotal += 1
	}

	if fst.BasicInterval != "" && (deviceOrigin.BasicInterval != fst.BasicInterval) {
		deviceOrigin.BasicInterval = fst.BasicInterval
		upTotal += 1
	}
	if fst.FitHousing != "" && (deviceOrigin.FitHousing != fst.FitHousing) {
		deviceOrigin.FitHousing = fst.FitHousing
		upTotal += 1
	}
	if fst.PlayUrl != "" && (deviceOrigin.PlayUrl != fst.PlayUrl) {
		deviceOrigin.PlayUrl = fst.PlayUrl
		upTotal += 1
	}
	if fst.RtspUrl != "" && (deviceOrigin.RtspUrl != fst.RtspUrl) {
		deviceOrigin.RtspUrl = fst.RtspUrl
		upTotal += 1
	}

	if fst.RtspUrlOne != "" && (deviceOrigin.RtspUrlOne != fst.RtspUrlOne) {
		deviceOrigin.RtspUrlOne = fst.RtspUrlOne
		upTotal += 1
	}
	if fst.RtspUrlTwo != "" && (deviceOrigin.RtspUrlTwo != fst.RtspUrlTwo) {
		deviceOrigin.RtspUrlTwo = fst.RtspUrlTwo
		upTotal += 1
	}
	if fst.Description != "" && (deviceOrigin.Description != fst.Description) {
		deviceOrigin.Description = fst.Description
		upTotal += 1
	}

	if fst.CreatedAt != "" && (deviceOrigin.CreatedAt != fst.CreatedAt) {
		deviceOrigin.CreatedAt = fst.CreatedAt
		upTotal += 1
	}

	if fst.UpdatedAt != "" && (deviceOrigin.UpdatedAt != fst.UpdatedAt) {
		deviceOrigin.UpdatedAt = fst.UpdatedAt
		upTotal += 1
	}

	if upTotal == 0 {
		err = fmt.Errorf("数据无需修改。")
		return 0, err
	}

	deviceStructs := []iotmodel.FmStreamTasks{
		deviceOrigin,
	}

	/*
	 */
	rst, err := DBConn.NamedExec(`UPDATE fm_tasks SET 
	task_type=:task_type, 
	task_name=:task_name,
	page_id=:page_id, 
	enabled=:enabled,
	creator=:creator,
	description=:description,
	basic_unit=:basic_unit, 
	basic_interval=:basic_interval,
	fit_housing=:fit_housing,
	play_url=:play_url,
	url=:url,
	url_1=:url_1,
	url_2=:url_2,
	updated_time=:updated_time,
	created_time=:created_time 
	WHERE id=:id`, deviceStructs)
	syslog.Printf("from db select:%#v  err:%v \n", rst, err)
	if err != nil {
		return 0, err
	}

	syslog.Printf("from db get:%v user:%#v err:%v \n", fst, rst, err)
	if err != nil {
		return 0, err
	}

	return rst.RowsAffected()
}

// 组流设备新增
func (sp *StreamRep) StreamDeviceNew(fst *iotmodel.FmStreamTasks) (int64, error) {

	var (
		DBConn = databases.DBMysql
	)

	timeNow := time.Now().In(configs.Loc)
	tStr := timeutil.CSTLayoutString(&timeNow)

	if fst.CreatedAt == "" {
		fst.CreatedAt = tStr
		fst.UpdatedAt = tStr
	}

	if fst.Enable == "" {
		fst.Enable = "1"
	}

	deviceStructs := []iotmodel.FmStreamTasks{
		*fst,
	}

	rst, err := DBConn.NamedExec(`INSERT INTO fm_tasks (
		task_type,
		task_name,
		page_id,
		enabled,
		creator,
		basic_unit,
		basic_interval,
		fit_housing,
		play_url,
		url,  
		url_1,
		url_2,
		updated_time, 
		description,
		delete_flag,  
		created_time)
        VALUES (:task_type,
			:task_name,
			:page_id,
			:enabled,
			:creator,
			:basic_unit,
			:basic_interval,
			:fit_housing,
			:play_url,
			:url,  
			:url_1,
			:url_2,
			:updated_time, 
			:description,
			:delete_flag,  
			:created_time)`, deviceStructs)

	syslog.Printf("new device register:%#v, with db rst:%v err:%v\n", fst, rst, err)
	if err != nil {
		// make sure err is a mysql.MySQLError.
		if errMySQL, ok := err.(*mysql.MySQLError); ok {
			switch errMySQL.Number {
			case 1062:
				// TODO handle Error 1062: Duplicate entry '%s' for key %d
				return 0, fmt.Errorf("已添加该设备:%v 页面id:%v", fst.TaskName, fst.PageId)
			}
		}

		return 0, err
	}

	lastId, err := rst.LastInsertId()
	syslog.Printf("new device register:%#v, with db rst:%v \n", lastId, err)

	if err != nil {
		return 0, err
	}

	return lastId, err
}

// 查询
func (sp *StreamRep) StreamSelect(page, pageSize int, tenant_id, pageid, url string) (*iotmodel.StreamTaskList, error) {

	var (
		StreamsDevices        = &iotmodel.StreamTaskList{}
		totalSql       string = "select count(*) from fm_tasks "
		baseSql        string = "select * from fm_tasks "

		filter string = " WHERE delete_flag != '1' "
		DBConn        = databases.DBMysql
	)

	if pageid != "" {
		filter += fmt.Sprintf(" AND page_id = '%v' ", pageid)
	}
	if url != "" {
		filter += fmt.Sprintf(" AND url = '%v' ", url)
	}

	//总数
	var totals int
	totalSqls := totalSql + filter
	err := DBConn.Get(&totals, totalSqls)
	syslog.Printf("total is :%v sql:%v err:%#v", totals, totalSqls, err)

	if err != nil {
		return nil, err
	}

	//数据
	limits := fmt.Sprintf(" LIMIT %v OFFSET %v;", pageSize, pageSize*(page-1))
	baseSql += filter
	baseSql += limits

	err = DBConn.Select(&StreamsDevices.Data, baseSql)
	syslog.Printf("Data is :%v sql:%v err:%#v", StreamsDevices.Data, baseSql, err)

	if err != nil {
		return nil, err
	}

	StreamsDevices.Total = int64(totals)
	StreamsDevices.Size = int64(pageSize)
	StreamsDevices.Columns, _ = FeeStatusMap()

	return StreamsDevices, nil
}
