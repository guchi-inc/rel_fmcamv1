package alysys

import (
	"fmcam/common/databases"
	"fmcam/models/sysaly"
	"fmt"

	"github.com/google/uuid"
)

// 统计查询执行 获取抓拍记录统计
func (mt *AlySysOperator) GetCaptureLogsWeek(Page, PageSize int, TenantId, StartAt, EndAt, sub_type, RangeType string, devIds *string, twoModel *bool) (*sysaly.SysClientList[any], error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		datas     = make(map[string]any)
		useSysaly = sysaly.SysClientList[any]{Data: datas}
		DBConn    = databases.DBMysql
		N         = 7
		filters   string
	)

	//单独查hour
	if sub_type == "hour" {
		hourData, err := mt.CaptionLogHoursData(StartAt, EndAt, TenantId)
		if err != nil {
			return nil, err
		}
		useSysaly.Data["hour"] = hourData.Data["hour"]
		return &useSysaly, nil
	}
	//已选择 时间范围格式  周 和 月
	if RangeType != "" || (devIds != nil && sub_type == "accrued") {
		return mt.CaptionLogsWeekMonth(RangeType, TenantId, EndAt, sub_type, devIds, twoModel)
	}
	//日期
	if TenantId != "" {
		filters += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	if StartAt != "" && EndAt != "" {
		filters += fmt.Sprintf(" AND (a.capture_time between '%v' AND '%v' )", StartAt, EndAt)
	}

	DatePoint := fmt.Sprintf("%v", "CURDATE()")
	if EndAt != "" {
		// DATE('2025-08-30 00:00:00')
		DatePoint = fmt.Sprintf("DATE('%v')", EndAt)
	}
	opSql := fmt.Sprintf(`
		 
			WITH RECURSIVE dates AS (
				SELECT %v - INTERVAL %d DAY AS days
				UNION ALL
				SELECT days + INTERVAL 1 DAY
				FROM dates
				WHERE days + INTERVAL 1 DAY <= %v
			)
			SELECT 
			    'date' AS data_type,
				'd7' AS time_range, 
				DATE_FORMAT(d.days, '%%Y-%%m-%%d') as stats,
				COALESCE(COUNT(a.id), 0) AS count
			FROM dates d
			 
	`, DatePoint, N-1, DatePoint)

	joinLeft := " LEFT JOIN CaptureLogs a  ON DATE(a.capture_time) "
	joinRight := " d.days AND a.func_type = '1' "
	filters += " GROUP BY stats ORDER BY stats "

	//事务开启 减少网络传输
	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	// 筛选抓拍面部数据
	if sub_type == "" || sub_type == "date" || sub_type == "accrued" {
		if sub_type == "" {
			//查两次
			//日期值
			dateDatas := []sysaly.DailyCaptureCount{}
			opSql2 := opSql
			opSql2 += joinLeft + " = " + joinRight
			opSql2 += filters
			err := tx.Select(&dateDatas, opSql2)
			logsys.Printf("dateDatas:%#v  opSql2:%#v  err:%v \n", dateDatas, opSql2, err)
			if err != nil {
				return nil, err
			}
			datas["date"] = dateDatas

			//累计值
			dateDataAccrued := []sysaly.DailyCaptureCount{}
			opSql3 := opSql
			opSql3 += joinLeft + " <= " + joinRight
			opSql3 += filters
			err = tx.Select(&dateDataAccrued, opSql3)
			logsys.Printf("dateDatas:%#v  opSql3:%#v  err:%v \n", dateDatas, opSql3, err)
			if err != nil {
				return nil, err
			}
			datas["accrued"] = dateDataAccrued
		} else {
			//查指定 图
			if sub_type == "date" {
				dateDatas := []sysaly.DailyCaptureCount{}
				opSql4 := opSql
				opSql4 += joinLeft + " = " + joinRight
				opSql4 += filters
				err := tx.Select(&dateDatas, opSql4)
				logsys.Printf("sub_type date Datas:%#v  opSql4:%#v  err:%v \n", dateDatas, opSql4, err)
				if err != nil {
					return nil, err
				}
				datas["date"] = dateDatas
			}
			if sub_type == "accrued" {
				//累计值
				dateDataAccrued := []sysaly.DailyCaptureCount{}
				opSql5 := opSql
				opSql5 += joinLeft + "<= " + joinRight
				opSql5 += filters
				err := tx.Select(&dateDataAccrued, opSql5)
				logsys.Printf("dateDatas:%#v  opSql5:%#v  err:%v \n", dateDataAccrued, opSql5, err)
				if err != nil {
					return nil, err
				}
				datas["accrued"] = dateDataAccrued
			}
		}
	}

	// 按人员类型 的采集统计  全部
	var FilterTypes string = " WHERE  c.group_type = 2  "
	var JoinCaptionFilter string = ""
	if TenantId != "" {
		FilterTypes += fmt.Sprintf(" AND BIN_TO_UUID(pt.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	if StartAt != "" && EndAt != "" {
		JoinCaptionFilter += fmt.Sprintf(" AND (a.capture_time between '%v' AND '%v' )", StartAt, EndAt)
	}

	// 按人员类型统计 7 day
	baseTypeSql := `SELECT 
	? AS time_range,
	'user' AS data_type,
	pt.type_name AS user_category,
	COALESCE(COUNT(a.id), 0) AS count
	FROM   fm_alert_group c 
	LEFT JOIN group_profile_type_mapping ptmp ON ptmp.group_id = c.id
	LEFT JOIN ProfileType pt ON pt.id = ptmp.profile_type_id
	LEFT JOIN Profile b ON b.type_id = pt.id
	LEFT JOIN CaptureLogs a ON a.matched_profile_id = b.id AND a.func_type = '1'
				`
	alyRange := fmt.Sprintf("AND a.capture_time >= DATE_SUB(CURDATE(), INTERVAL %v DAY) ", N)
	groups := "GROUP BY user_category ORDER BY user_category "
	p7typeSql := fmt.Sprintf(`
			%v	%v
			%v
			%v
			%v
	`, baseTypeSql, JoinCaptionFilter, alyRange, FilterTypes, groups)

	if sub_type == "" || sub_type == "user" {
		type7Datas := []sysaly.DailyCaptureCount{}
		err := tx.Select(&type7Datas, p7typeSql, "d7")
		logsys.Printf("type7Datas:%#v  select:%#v  err:%v \n", type7Datas, p7typeSql, err)
		if err != nil {
			return nil, err
		}
		datas["user"] = type7Datas
	}

	//默认返回小时聚合
	hourData, err := mt.CaptionLogHoursData(StartAt, EndAt, TenantId)
	if err != nil {
		return nil, err
	}
	datas["hour"] = hourData.Data["hour"]

	useSysaly.Data = datas

	///查询总数
	dataTotals := make(map[string]any)

	//1天内按类型
	alyRange = fmt.Sprintf("AND a.capture_time >= DATE_SUB(CURDATE(), INTERVAL %v DAY) ", 1)
	p1typeSql := fmt.Sprintf(`
		%v	 %v 
		%v
		%v
		%v
	`, baseTypeSql, JoinCaptionFilter, alyRange, FilterTypes, groups)
	type1Datas := []sysaly.DailyCaptureCount{}
	err = tx.Select(&type1Datas, p1typeSql, "d1")
	logsys.Printf("type1Datas:%#v  select:%#v  err:%v \n", type1Datas, p1typeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["d1"] = type1Datas

	//全部时间的 按类型
	pAlltypeSql := fmt.Sprintf(`
		%v	 
		%v
		%v
	`, baseTypeSql, FilterTypes, groups)
	typeAllDatas := []sysaly.DailyCaptureCount{}
	err = tx.Select(&typeAllDatas, pAlltypeSql, "all")
	logsys.Printf("typeAllDatas:%#v  select:%#v  err:%v \n", typeAllDatas, pAlltypeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["all"] = typeAllDatas

	//返回小时 累计 聚合
	dataTotals["accrued_hour"] = hourData.TotalData["accrued_hour"]
	useSysaly.TotalData = dataTotals

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &useSysaly, nil
}

// 获取近期七个预警 按位置
func (mt *AlySysOperator) GetAlertsWeek(Page, PageSize int, TenantId string) (*sysaly.SysClientList[any], error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		datas     = make(map[string]any)
		useSysaly = sysaly.SysClientList[any]{Data: datas}
		DBConn    = databases.DBMysql
		N         = 7
	)

	//日期
	dateFormat := "'%Y-%m-%d %H:%i:%s'"
	filterTenant := ""
	if TenantId != "" {
		filterTenant += fmt.Sprintf(" WHERE BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	opSql := fmt.Sprintf(`
			SELECT 
			a.id AS id,   
			'position' AS data_type,
			de.location AS position,  
			DATE_FORMAT(a.created_time, %v) as stats,
			c.type_name AS user_category,
			d.capture_image_url AS url
			FROM Alerts a 
			LEFT JOIN CaptureLogs d ON d.id = a.capture_log_id
			LEFT JOIN Profile p ON p.id = d.matched_profile_id
			LEFT JOIN ProfileType c ON c.id = p.type_id
			LEFT JOIN Device de ON de.id = a.device_id
			%v
			ORDER BY a.id DESC 
			LIMIT %v
	`, dateFormat, filterTenant, N)

	// 筛选抓拍面部数据
	positionDatas := []sysaly.PositionStats{}
	err := DBConn.Select(&positionDatas, opSql)
	logsys.Printf("dateDatas:%#v  select:%#v  err:%v \n", positionDatas, opSql, err)
	if err != nil {
		return nil, err
	}
	datas["position"] = positionDatas

	useSysaly.Data = datas
	return &useSysaly, nil
}

// 获取预警 按设备和时段透视
func (mt *AlySysOperator) TotalDeviceTimeRange(Page, PageSize int, TenantId, startAt, endAt, sub_type string) (*sysaly.SysClientList[any], error) {
	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		// datas = make(map[string]any)
		// useSysaly        = sysaly.SysClientList[any]{Data: datas}
		// DBConn           = databases.DBMysql
		// N                = PageSize
		filters string = " WHERE "
	)

	//租户条件
	if TenantId != "" {
		filters += fmt.Sprintf("   BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	return nil, nil
}

// 获取预警 按位置统计透视
func (mt *AlySysOperator) TotalAlertsPostion(Page, PageSize int, TenantId, startAt, endAt, sub_type, RangeType string, devIds *string, twoModel *bool) (*sysaly.SysClientList[any], error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		datas            = make(map[string]any)
		useSysaly        = sysaly.SysClientList[any]{Data: datas}
		DBConn           = databases.DBMysql
		N                = PageSize
		filters   string = " WHERE "
	)

	//单独查hour
	if sub_type == "hour" {
		hourData, err := mt.AlertHoursDay(startAt, endAt, TenantId)
		if err != nil {
			return nil, err
		}
		useSysaly.Data["hour"] = hourData.Data["hour"]
		return &useSysaly, nil
	}
	//已选择 时间范围格式  周 和 月
	if RangeType != "" || (devIds != nil && sub_type == "accrued") {
		return mt.AlertsWeekMonth(RangeType, TenantId, endAt, sub_type, devIds, twoModel)
	}
	//租户条件
	if TenantId != "" {
		filters += fmt.Sprintf(" BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	if startAt != "" && endAt != "" {
		if len(filters) > 5 {
			filters += fmt.Sprintf("  AND (a.created_time between '%v' and '%v') ", startAt, endAt)
		} else {
			filters += fmt.Sprintf(" (a.created_time between '%v' and '%v') ", startAt, endAt)
		}
	}

	//位置
	groups := " GROUP BY position ORDER BY position DESC"
	opSql := fmt.Sprintf(`
			SELECT
			'position' AS data_type,
			COUNT(DISTINCT a.id) AS count,
			de.location AS position
			FROM Alerts a
			LEFT JOIN Device de ON de.id = a.device_id 
			%v
			%v 
	`, filters, groups)

	///事务开启
	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	// 筛选预警数据 按位置
	positionDatas := []sysaly.WarnStats{}
	if sub_type == "" || sub_type == "position" {
		err := tx.Select(&positionDatas, opSql)
		logsys.Printf("dateDatas:%#v  select:%#v  err:%v \n", positionDatas, opSql, err)
		if err != nil {
			return nil, err
		}
		datas["position"] = positionDatas
	}

	// 按结束日期 筛选预警数据 按近7天
	DatePoint := fmt.Sprintf("%v", "CURDATE()")
	if endAt != "" {
		// DATE('2025-08-30 00:00:00')
		DatePoint = fmt.Sprintf("DATE('%v')", endAt)
	}

	if TenantId != "" {
		filters = fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	//isIncre 增量或累计 统计符
	// var isIncre string = "="
	// if sub_type == "accrued" {
	// 	isIncre = "<="
	// }
	dateSql := fmt.Sprintf(`
		WITH RECURSIVE dates AS (
			SELECT  %v - INTERVAL %v DAY AS days
			UNION ALL
			SELECT days + INTERVAL 1 DAY
			FROM dates
			WHERE days + INTERVAL 1 DAY <= %v
		)
		SELECT 
			'date' AS data_type,
			'd7' AS time_range, 
			DATE_FORMAT(d.days, '%%Y-%%m-%%d') as stats,
			COUNT(a.id) AS count
		FROM dates d
		
	`, DatePoint, N-1, DatePoint)

	joinLeft := " LEFT JOIN Alerts a ON  DATE(a.created_time)   "
	joinRight := "  d.days "
	filters += " GROUP BY stats ORDER BY stats;"
	// 筛选抓拍面部数据
	if sub_type == "" || sub_type == "date" || sub_type == "accrued" {
		if sub_type == "" {
			//查两次
			//日期值
			dateDatas := []sysaly.DailyCaptureCount{}
			dateSql1 := dateSql
			dateSql1 += joinLeft + " = " + joinRight
			dateSql1 += filters
			err := tx.Select(&dateDatas, dateSql1)
			logsys.Printf("dateDatas:%#v  dateSql1:%#v  err:%v \n", dateDatas, dateSql1, err)
			if err != nil {
				return nil, err
			}
			datas["date"] = dateDatas

			//累计值
			dateDataAccrued := []sysaly.DailyCaptureCount{}
			dateSql2 := dateSql
			dateSql2 += joinLeft + "<= " + joinRight
			dateSql2 += filters
			err = tx.Select(&dateDataAccrued, dateSql2)
			logsys.Printf("dateDatas:%#v  dateSql2:%#v  err:%v \n", dateDatas, dateSql2, err)
			if err != nil {
				return nil, err
			}
			datas["accrued"] = dateDataAccrued
		} else {
			//查指定 图
			if sub_type == "date" {
				dateDatas := []sysaly.DailyCaptureCount{}
				dateSql3 := dateSql
				dateSql3 += joinLeft + " = " + joinRight
				dateSql3 += filters
				err := tx.Select(&dateDatas, dateSql3)
				logsys.Printf("dateDatas:%#v  dateSql3:%#v  err:%v \n", dateDatas, dateSql3, err)
				if err != nil {
					return nil, err
				}
				datas["date"] = dateDatas
			}
			if sub_type == "accrued" {
				//累计值
				dateDataAccrued := []sysaly.DailyCaptureCount{}
				dateSql4 := dateSql
				dateSql4 += joinLeft + "<= " + joinRight
				dateSql4 += filters
				err := tx.Select(&dateDataAccrued, dateSql4)
				logsys.Printf("dateDatas:%#v  dateSql4:%#v  err:%v \n", dateDataAccrued, dateSql4, err)
				if err != nil {
					return nil, err
				}
				datas["accrued"] = dateDataAccrued
			}
		}
	}

	levelTypeSql := ` SELECT 
	? AS date_range, 
	a.alert_level as level_type,  
	COUNT(a.id) as count 
	FROM Alerts a
	`
	groups = " GROUP BY level_type  ORDER BY level_type "

	//查询日期的 1天内
	var filtersLevel string
	if TenantId != "" {
		filtersLevel = fmt.Sprintf(" WHERE BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}
	if filtersLevel == "" {
		filtersLevel = fmt.Sprintf(" WHERE a.created_time >= DATE_SUB(%v, INTERVAL %v DAY) ", DatePoint, 7)
	} else {
		filtersLevel += fmt.Sprintf(" AND a.created_time >= DATE_SUB(%v, INTERVAL %v DAY) ", DatePoint, 7)
	}
	if startAt != "" && endAt != "" {
		if len(filters) > 5 {
			filtersLevel += fmt.Sprintf("  AND (a.created_time between '%v' and '%v') ", startAt, endAt)
		} else {
			filtersLevel = fmt.Sprintf(" WHERE (a.created_time between '%v' and '%v') ", startAt, endAt)
		}
	}

	pLeveltypeSql := fmt.Sprintf(`
		%v	
		%v 
		%v
	`, levelTypeSql, filtersLevel, groups)

	level7Datas := []sysaly.WarnStats{}
	err = tx.Select(&level7Datas, pLeveltypeSql, "d7")
	logsys.Printf("level7Datas:%#v  pLeveltypeSql:%#v  err:%v \n", level7Datas, pLeveltypeSql, err)
	if err != nil {
		return nil, err
	}
	// 7天内按预警登记的统计 预警数量
	datas["user"] = level7Datas

	//查 小时累计 和 区间
	hourData, err := mt.AlertHoursDay(startAt, endAt, TenantId)
	if err != nil {
		return nil, err
	}
	datas["hour"] = hourData.Data["hour"]
	useSysaly.Data = datas

	/////////////////////////////////////////////////////////////////////////
	// ///查询总数
	baseTypeSql := ` SELECT 
	? AS date_range, 
	a.alert_level as level_type,  
	COUNT(a.id) as count 
	FROM Alerts a
	`
	groups = " GROUP BY level_type  ORDER BY level_type "
	dataTotals := make(map[string]any)

	//查询日期的 1天内
	var alyRange string
	if TenantId != "" {
		filters = fmt.Sprintf(" WHERE BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}
	if filters == "" {
		alyRange = fmt.Sprintf(" WHERE a.created_time >= DATE_SUB(%v, INTERVAL %v DAY) ", DatePoint, 1)

	} else {
		alyRange = fmt.Sprintf(" AND a.created_time >= DATE_SUB(%v, INTERVAL %v DAY) ", DatePoint, 1)
	}
	p1typeSql := fmt.Sprintf(`
		%v	
		%v
		%v
		%v
	`, baseTypeSql, filters, alyRange, groups)

	type1Datas := []sysaly.WarnStats{}
	err = tx.Select(&type1Datas, p1typeSql, "d1")
	logsys.Printf("type1Datas:%#v  select:%#v  err:%v \n", type1Datas, p1typeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["d1"] = type1Datas

	//全部
	if startAt != "" && endAt != "" {
		if filters == "" {
			filters = fmt.Sprintf(" WHERE (a.created_time between  '%v' and '%v' ) ", startAt, endAt)
		} else {
			filters += fmt.Sprintf(" AND (a.created_time between '%v' and '%v' ) ", startAt, endAt)
		}
	}
	pAlltypeSql := fmt.Sprintf(`
		%v	
		%v 
		%v
	`, baseTypeSql, filters, groups)
	typeAllDatas := []sysaly.WarnStats{}
	err = tx.Select(&typeAllDatas, pAlltypeSql, "all")
	logsys.Printf("typeAllDatas:%#v  select:%#v  err:%v \n", typeAllDatas, pAlltypeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["all"] = typeAllDatas

	//待处理类型 todo
	if filters == "" {
		alyRange = fmt.Sprintf(" WHERE a.status = '%v' ", 0)
	} else {
		alyRange = fmt.Sprintf(" AND a.status = '%v' ", 0)
	}

	pTotoAlterSql := fmt.Sprintf(`
			%v	
			%v 
			%v
			%v
	`, baseTypeSql, filters, alyRange, groups)
	pAlterTodo := []sysaly.WarnStats{}
	err = tx.Select(&pAlterTodo, pTotoAlterSql, "todo")
	logsys.Printf("pAlterTodo:%#v  select:%#v  err:%v \n", typeAllDatas, pAlltypeSql, err)
	if err != nil {
		return nil, err
	}

	dataTotals["todo"] = pAlterTodo

	dataTotals["accrued_hour"] = hourData.TotalData["accrued_hour"]
	useSysaly.TotalData = dataTotals

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &useSysaly, nil
}

// 获取 采集日志统计
func (mt *AlySysOperator) GetGatherLogWeek(Page, PageSize int, TenantId, StartAt, EndAt string) (*sysaly.SysClientList[any], error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		datas       = make(map[string]any)
		useSysaly   = sysaly.SysClientList[any]{Data: datas}
		DBConn      = databases.DBMysql
		N           = 7
		filterDates string
	)

	dateFormat := "'%Y-%m-%d'"
	//日期
	if TenantId != "" {
		filterDates += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}
	if StartAt != "" && EndAt != "" {
		filterDates += fmt.Sprintf(" AND (a.capture_time between '%v' AND '%v' )", StartAt, EndAt)
	}

	DatePoint := fmt.Sprintf("%v", "CURDATE()")
	if EndAt != "" {
		// DATE('2025-08-30 00:00:00')
		DatePoint = fmt.Sprintf("DATE('%v')", EndAt)
	}

	opSql := fmt.Sprintf(`
			WITH RECURSIVE dates AS (
				SELECT %v - INTERVAL %d DAY AS days
				UNION ALL
				SELECT days + INTERVAL 1 DAY
				FROM dates
				WHERE days + INTERVAL 1 DAY <= %v
			)
			SELECT 
			    'date' AS data_type,
				'd7' AS time_range, 
				DATE_FORMAT(d.days, %v) as stats,
				COALESCE(COUNT(a.id), 0) AS count
			FROM dates d
			LEFT JOIN CaptureLogs a 
				ON DATE(a.capture_time) = d.days AND a.func_type = '0'
			%v
			GROUP BY stats
			ORDER BY stats;
	`, DatePoint, N-1, DatePoint, dateFormat, filterDates)

	// 筛选 采集 面部数据
	dateDatas := []sysaly.DailyCaptureCount{}
	err := DBConn.Select(&dateDatas, opSql)
	logsys.Printf("dateDatas:%#v  select:%#v  err:%v \n", dateDatas, opSql, err)
	if err != nil {
		return nil, err
	}
	datas["date"] = dateDatas
	useSysaly.Data = datas

	///查询总数
	dataTotals := make(map[string]any)
	// 按人员类型 的采集统计  全部
	var FilterTypes string = " WHERE  c.group_type = 2  "
	var JoinCaptionFilter string = ""
	if TenantId != "" {
		FilterTypes += fmt.Sprintf(" AND BIN_TO_UUID(pt.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	if StartAt != "" && EndAt != "" {
		JoinCaptionFilter += fmt.Sprintf(" AND (a.capture_time between '%v' AND '%v' )", StartAt, EndAt)
	}

	baseTypeSql := `SELECT 
	? AS time_range,
	'user' AS data_type,
	pt.type_name AS user_category,
	COALESCE(COUNT(a.id), 0) AS count
	FROM   fm_alert_group c 
	LEFT JOIN group_profile_type_mapping ptmp ON ptmp.group_id = c.id
	LEFT JOIN ProfileType pt ON pt.id = ptmp.profile_type_id 
	LEFT JOIN Profile b ON b.type_id = pt.id 
	LEFT JOIN CaptureLogs a ON a.matched_profile_id = b.id AND a.func_type = '0'
	`
	groups := " GROUP BY user_category ORDER BY user_category "
	p7typeSql := fmt.Sprintf(`
		%v	%v
		%v
		%v
		`, baseTypeSql, JoinCaptionFilter, FilterTypes, groups)

	type7Datas := []sysaly.DailyCaptureCount{}
	err = DBConn.Select(&type7Datas, p7typeSql, "all")
	logsys.Printf("type7Datas:%#v  select:%#v  err:%v \n", type7Datas, p7typeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["user"] = type7Datas

	//1天内按类型
	alyRange := fmt.Sprintf("AND a.capture_time >= DATE_SUB(CURDATE(), INTERVAL %v DAY) ", 1)
	totalSql := ` SELECT 
	? AS time_range,
	'user' AS data_type,
	pt.type_name AS user_category,
	COALESCE(COUNT(a.id), 0) AS count
	FROM  fm_alert_group c 
	LEFT JOIN group_profile_type_mapping ptmp ON ptmp.group_id = c.id
	LEFT JOIN ProfileType pt ON pt.id = ptmp.profile_type_id  
	LEFT JOIN Profile b ON b.type_id = pt.id 
	LEFT JOIN CaptureLogs a ON a.matched_profile_id = b.id AND a.func_type = '0'
	`
	p1typeSql := fmt.Sprintf(`
	%v
	%v %v 	
	%v
	%v 
	`, totalSql, filterDates, alyRange,
		FilterTypes, groups)
	gatherDatas := []sysaly.DailyCaptureCount{}
	err = DBConn.Select(&gatherDatas, p1typeSql, "d1")
	logsys.Printf("type1Datas:%#v \n p1typeSql:%#v  err:%v \n", gatherDatas, p1typeSql, err)
	if err != nil {
		return nil, err
	}

	//全部时间的 按类型
	pAlltypeSql := fmt.Sprintf(`
	%v
		%v %v
	`, totalSql, FilterTypes, groups)
	typeAllDatas := []sysaly.DailyCaptureCount{}
	err = DBConn.Select(&typeAllDatas, pAlltypeSql, "all")
	logsys.Printf("typeAllDatas:%#v  \n pAlltypeSql:%#v  err:%v \n", typeAllDatas, pAlltypeSql, err)
	if err != nil {
		return nil, err
	}

	gatherDatas = append(gatherDatas, typeAllDatas...)
	dataTotals["gather"] = gatherDatas

	useSysaly.TotalData = dataTotals
	return &useSysaly, nil
}

// 获取 采集日志统计
func (mt *AlySysOperator) GetGatherLogAccrued(Page, PageSize int, TenantId, StartAt, EndAt, sub_type, RangeType string, profileIds *string, twoModel *bool) (*sysaly.SysClientList[any], error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		datas     = make(map[string]any)
		useSysaly = sysaly.SysClientList[any]{Data: datas}
		DBConn    = databases.DBMysql
	)

	type TypeInfo struct {
		ID       int    `json:"id" db:"id"`
		TypeName string `json:"type_name" db:"type_name"`
	}

	///事务开启
	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	newTypes := []TypeInfo{}
	if profileIds == nil || *profileIds == "" {
		profileIdListSql := fmt.Sprintf(`SELECT id,type_name FROM ProfileType WHERE BIN_TO_UUID(tenant_id) = '%v';`, TenantId)

		err := tx.Select(&newTypes, profileIdListSql)
		if err != nil || len(newTypes) == 0 {
			return nil, nil //避免首次加载失败 报错
		}

		var pIdStr string
		for _, pid := range newTypes {
			if len(pIdStr) == 0 {
				pIdStr += fmt.Sprintf("%v", pid.ID)
			} else {
				pIdStr += fmt.Sprintf(",%v", pid.ID)

			}
		}
		profileIds = &pIdStr
		useSysaly.Columns = newTypes
	}

	logsys.Printf("RangeType:%#v    sub_type:%v for profile type ids:%v \n", RangeType, sub_type, *profileIds)

	//已选择 时间范围格式  周 和 月
	if RangeType != "" {
		useSysaly, err := mt.GathersWeekMonth(RangeType, TenantId, EndAt, sub_type, profileIds)
		useSysaly.Columns = newTypes
		return useSysaly, err
	}

	// 筛选抓拍面部数据
	if sub_type == "" || sub_type == "date" || sub_type == "accrued" {
		//初始加载页面 单图 按天统计， 新增
		dateTwoFile, err := mt.GathersWeekMonth("date", TenantId, "", "date", profileIds)
		if err != nil {
			return nil, err
		}
		datas["date"] = dateTwoFile.Data["date"]

		//初始加载页面 单图 按天统计， 累计
		dateTwoFileAcc, err := mt.GathersWeekMonth("date", TenantId, "", "accrued", profileIds)
		if err != nil {
			return nil, err
		}
		datas["accrued"] = dateTwoFileAcc.Data["accrued"]
	}

	// 按人员类型 的采集统计  全部  2 分组类型为 人员分组
	var FilterTypes string = " WHERE  c.group_type = 2  AND pt.type_name is not NULL"
	var JoinGatherTimeFilter string = ""
	if TenantId != "" {
		FilterTypes += fmt.Sprintf(" AND BIN_TO_UUID(c.tenant_id) = '%v' ", uuid.MustParse(TenantId))
	}

	if StartAt != "" && EndAt != "" {
		JoinGatherTimeFilter += fmt.Sprintf(" AND (a.capture_time between '%v' AND '%v' )", StartAt, EndAt)
	}

	// 按人员类型统计 7 day 累计
	baseTypeSql := `SELECT 
	? AS time_range,
	'user' AS data_type,
	pt.type_name AS user_category,
	COUNT(a.id) AS count
	FROM   fm_alert_group c 
	LEFT JOIN group_profile_type_mapping ptmp ON ptmp.group_id = c.id
	LEFT JOIN ProfileType pt ON pt.id = ptmp.profile_type_id
	LEFT JOIN Profile b ON b.type_id = pt.id
	LEFT JOIN CaptureLogs a ON a.matched_profile_id = b.id AND a.func_type = '0'
				`
	alyRange7 := fmt.Sprintf("AND a.capture_time >= DATE_SUB(CURDATE(), INTERVAL %v DAY) ", 7)

	groups := "GROUP BY user_category ORDER BY user_category "
	pTimetypeSql := fmt.Sprintf(`
			%v	%v %v
			%v 
			%v
	`, baseTypeSql, JoinGatherTimeFilter, alyRange7, FilterTypes, groups)

	if sub_type == "" || sub_type == "user" {
		type7Datas := []sysaly.DailyCaptureCount{}
		err := tx.Select(&type7Datas, pTimetypeSql, "date")
		logsys.Printf("type7Datas:%#v  pTimetypeSql:%#v  err:%v \n", type7Datas, pTimetypeSql, err)
		if err != nil {
			return nil, err
		}
		datas["user"] = type7Datas
	}

	useSysaly.Data = datas
	///查询总数
	dataTotals := make(map[string]any)
	//1天内按类型
	alyRange1 := fmt.Sprintf("AND a.capture_time >= DATE_SUB(CURDATE(), INTERVAL %v DAY) ", 1)
	p1typeSql := fmt.Sprintf(`
		%v	 %v 
		%v
		%v
		%v
	`, baseTypeSql, JoinGatherTimeFilter, alyRange1, FilterTypes, groups)
	type1Datas := []sysaly.DailyCaptureCount{}
	err = tx.Select(&type1Datas, p1typeSql, "d1")
	logsys.Printf("type1Datas:%#v  p1typeSql:%#v  err:%v \n", type1Datas, p1typeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["d1"] = type1Datas

	//全部时间的 按类型
	pAlltypeSql := fmt.Sprintf(`
		%v	 
		%v
		%v
	`, baseTypeSql, FilterTypes, groups)
	typeAllDatas := []sysaly.DailyCaptureCount{}
	err = tx.Select(&typeAllDatas, pAlltypeSql, "all")
	logsys.Printf("typeAllDatas:%#v  pAlltypeSql:%#v  err:%v \n", typeAllDatas, pAlltypeSql, err)
	if err != nil {
		return nil, err
	}
	dataTotals["all"] = typeAllDatas

	//当前生效面部数据 Face 库
	baseFaceSql := `SELECT 
	? AS time_range,
	'face' AS data_type, 
	COALESCE(COUNT(a.id), 0) AS count
	FROM   Face a  
	LEFT JOIN Profile b ON b.id = a.profile_id AND b.enabled = 1
				`
	FilterFaceTypes := fmt.Sprintf(" WHERE BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))

	allFaceSql := fmt.Sprintf(`
			%v	 
			%v 
	`, baseFaceSql, FilterFaceTypes)

	if sub_type == "" {
		typeFaceDatas := []sysaly.DailyCaptureCount{}
		err := tx.Select(&typeFaceDatas, allFaceSql, "all")
		logsys.Printf("typeFaceDatas:%#v  select:%#v  err:%v \n", typeFaceDatas, allFaceSql, err)
		if err != nil {
			return nil, err
		}
		dataTotals["face"] = typeFaceDatas
	}

	//当前临时面部数据 TmpFace 库
	//当前生效面部数据 Face 库
	baseTmpFaceSql := `SELECT 
		? AS time_range,
		'tmp_face' AS data_type, 
		COALESCE(COUNT(a.id), 0) AS count
		FROM TemporaryFace a   
					`
	FilterTmpFaceTypes := fmt.Sprintf(" WHERE BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(TenantId))

	allTmpFaceSql := fmt.Sprintf(`
			%v	 
			%v 
	`, baseTmpFaceSql, FilterTmpFaceTypes)

	if sub_type == "" {
		tmpFaceDatas := []sysaly.DailyCaptureCount{}
		err := tx.Select(&tmpFaceDatas, allTmpFaceSql, "all")
		logsys.Printf("typeFaceDatas:%#v  select:%#v  err:%v \n", tmpFaceDatas, allTmpFaceSql, err)
		if err != nil {
			return nil, err
		}
		dataTotals["tmp_face"] = tmpFaceDatas
	}

	useSysaly.TotalData = dataTotals

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &useSysaly, nil
}
