package alysys

import (
	"fmcam/common/databases"
	"fmcam/models/sysaly"
	"fmt"
)

// 按一天24小时聚合 抓拍记录
func (mt *AlySysOperator) CaptionLogHoursData(startTime, endTime, TenantId string) (*sysaly.HourSisList[any], error) {

	var datas = make(map[string]any)
	var TenantVar string
	var DBConn = databases.DBMysql
	var UseSysaly = sysaly.HourSisList[any]{TotalData: make(map[string]any), Data: datas}

	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}

	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
	//限制筛选条件
	if TenantId != "" {
		TenantVar = fmt.Sprintf("  AND BIN_TO_UUID(a.tenant_id) = '%v' ", TenantId)
	}

	var baseSql = fmt.Sprintf(`
		WITH RECURSIVE hours AS (
			SELECT 0 AS h
			UNION ALL
			SELECT h + 1 FROM hours WHERE h < 23
		),
		agg AS (
			SELECT 
				HOUR(a.capture_time) AS hr,
				COUNT(*) AS cnt
			FROM CaptureLogs a
			WHERE a.capture_time BETWEEN 
				COALESCE(@start_time, (SELECT MIN(capture_time) FROM CaptureLogs)) 
				AND @end_time
				%v
			GROUP BY HOUR(a.capture_time)
		)
		SELECT 
			h.h AS stats,
			COALESCE(e.cnt, 0) AS count,
			SUM(COALESCE(e.cnt, 0)) OVER (
				ORDER BY h.h
			) AS all_count
		FROM hours h
		LEFT JOIN agg e
			ON e.hr = h.h
		ORDER BY h.h;
	`, TenantVar)

	//累计总数 到当前时间, 不传日期区间  则表示页面首次加载，返回累计总数
	if startTime == "" && endTime == "" {
		tx.Exec("set @end_time = CURDATE();")
		hourAccuredRange := []sysaly.PositionStats{}
		err = tx.Select(&hourAccuredRange, baseSql)
		logsys.Printf("dateDatas select:%#v  err:%v \n", baseSql, err)
		if err != nil {
			tx.Rollback()
			logsys.Printf("查询失败:%#v", err)
			return nil, err
		}
		UseSysaly.TotalData["accrued_hour"] = hourAccuredRange
	}

	//日期区间数 默认 30天内
	if startTime != "" {
		tx.Exec(fmt.Sprintf("set @start_time = '%v';", startTime))
	} else {
		tx.Exec("set @start_time = CURDATE() - 30;")
	}
	if endTime != "" {
		tx.Exec(fmt.Sprintf("set @end_time = '%v';", endTime))
	} else {
		tx.Exec("set @end_time = CURDATE();")
	}

	hourDateRange := []sysaly.PositionStats{}
	err = tx.Select(&hourDateRange, baseSql)
	logsys.Printf("dateDatas select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		tx.Rollback()
		logsys.Printf("查询失败:%#v", err)
		return nil, err
	}
	UseSysaly.Data["hour"] = hourDateRange

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &UseSysaly, nil
}

// 按天/周/月 聚合
func (mt *AlySysOperator) CaptionLogsWeekMonth(DateRange, TenantId, endAt, sub_type string, devIds *string, twoModel *bool) (*sysaly.SysClientList[any], error) {

	// var ModeVar string
	var datas = make(map[string]any)
	var TenantVar string
	var DBConn = databases.DBMysql
	var UseSysaly = sysaly.SysClientList[any]{Data: datas}

	if DateRange == "" {
		DateRange = "month"
	}
	//dateRange 按天 date，周 week，月 month
	if !DefDateRange[DateRange] {
		return nil, fmt.Errorf("不支持的时间模式查询方式")
	}

	if sub_type == "" {
		sub_type = "date" //该函数查 时间类型的统计
	}

	/////////////////////////////////////////////////环境变量设置

	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	///////////存储函数 生成 id 统计维度，只在twoModel有效时执行
	tx.MustExec("DROP PROCEDURE IF EXISTS split_to_table;")
	_, err = tx.Exec(` 
		CREATE PROCEDURE split_to_table(IN idList TEXT)
		BEGIN
			DECLARE pos INT DEFAULT 1;
			DECLARE commaPos INT;
			DECLARE curId TEXT;
			DROP TEMPORARY TABLE IF EXISTS tmp_ids;
			CREATE TEMPORARY TABLE tmp_ids (device_id INT);
			WHILE pos > 0 DO
				SET commaPos = LOCATE(',', idList, pos);

				IF commaPos > 0 THEN
					SET curId = SUBSTRING(idList, pos, commaPos - pos);
					INSERT INTO tmp_ids VALUES (CAST(curId AS UNSIGNED));
					SET pos = commaPos + 1;
				ELSE
					SET curId = SUBSTRING(idList, pos);
					IF curId <> '' THEN
						INSERT INTO tmp_ids VALUES (CAST(curId AS UNSIGNED));
					END IF;
					SET pos = 0;
				END IF;
			END WHILE;
		END; `)

	if err != nil {
		logsys.Printf("存储过程创建失败:%#v", err)
		return nil, err
	}
	if endAt == "" {
		endAt = "CURDATE()"
	} else {
		endAt = "'" + endAt + "'"
	}
	enabTwo := 0
	if *twoModel {
		enabTwo = 1
	}

	// 设置会话变量
	envSetSql := fmt.Sprintf(`SET @range_type = '%v';`, DateRange)

	tx.MustExec(envSetSql)
	envSetSql1 := fmt.Sprintf(`SET @agg = '%v'; `, sub_type)
	tx.MustExec(envSetSql1)
	envSetSql2 := fmt.Sprintf(`SET @endDate = %v; `, endAt)
	tx.MustExec(envSetSql2)
	envSetSql3 := fmt.Sprintf(`	SET @idList = '%v'; ; `, *devIds)
	tx.MustExec(envSetSql3)
	envSetSql4 := fmt.Sprintf(`SET @twoModel = '%v';`, enabTwo)
	_, err = tx.Exec(envSetSql4)

	logsys.Printf("环境设置%verr:%#v", envSetSql, err)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//调用存储过程
	_, err = tx.Exec("CALL split_to_table(@idList);")
	if err != nil {
		tx.Rollback()
		logsys.Printf("存储过程调用%v失败:%#v", devIds, err)
		return nil, err
	}

	//////////筛选条件
	if TenantId != "" {
		TenantVar = fmt.Sprintf("  AND BIN_TO_UUID(a.tenant_id) = '%v' ", TenantId)
	}

	// DateRange week  month
	var (
		baseSql = fmt.Sprintf(`  
			WITH RECURSIVE periods AS (
				SELECT 1 AS step
				UNION ALL
				SELECT step + 1
				FROM periods
				WHERE step < 7
			),
			time_dim AS (
				SELECT 
					CASE 
						WHEN @range_type = 'date'  
							THEN DATE_FORMAT(@endDate - INTERVAL (7 - p.step) DAY, '%%Y-%%m-%%d')
						WHEN @range_type = 'week'  
							THEN DATE_FORMAT(
									STR_TO_DATE(CONCAT(YEARWEEK(@endDate - INTERVAL (7 - p.step) WEEK, 1), ' Monday'), '%%X%%V %%W'),
									'%%x-%%v'
								)
						WHEN @range_type = 'month' 
							THEN DATE_FORMAT(@endDate - INTERVAL (p.step - 1) MONTH, '%%Y-%%m')
					END AS period_val,
					CASE 
						WHEN @range_type = 'date'  
							THEN @endDate - INTERVAL (7 - p.step) DAY
						WHEN @range_type = 'week'  
							THEN @endDate - INTERVAL (7 - p.step) WEEK
						WHEN @range_type = 'month' 
							THEN @endDate - INTERVAL (p.step - 1) MONTH
					END AS real_date,
					p.step
				FROM periods p
			),
			id_dim AS (
				-- 如果传入 idList 非空，取 split_to_table 的结果；否则返回一行 NULL
				SELECT device_id FROM tmp_ids
				UNION ALL
				SELECT NULL WHERE @idList = ''
			)
			SELECT 
			CASE
				WHEN @twoModel = 0 THEN 
					CASE 
						WHEN   @idList = '' THEN t.period_val
						ELSE i.device_id
					END
				ELSE t.period_val
			END AS stats,
			(CASE WHEN @twoModel = 1 AND @idList != '' THEN i.device_id ELSE NULL END) AS device_id,
			COUNT(DISTINCT a.id) AS count,
				@range_type AS time_range,
				@agg  AS data_type
			FROM time_dim t
			JOIN id_dim i
			LEFT JOIN CaptureLogs a  
				ON (i.device_id IS NULL OR a.device_id = i.device_id)
			AND (
					-- 日维度
					(@range_type = 'date' AND 
						(
						(@agg = 'date'  AND DATE(a.capture_time) = t.real_date)
					OR (@agg = 'accrued' AND DATE(a.capture_time) <= t.real_date)
						)
					)
				OR (@range_type = 'week' AND 
						(
						(@agg = 'date'  AND YEARWEEK(a.capture_time, 1) = YEARWEEK(t.real_date, 1))
					OR (@agg = 'accrued' AND YEARWEEK(a.capture_time, 1) <= YEARWEEK(t.real_date, 1))
						)
					)
				OR (@range_type = 'month' AND 
						(
						(@agg = 'date'  AND DATE_FORMAT(a.capture_time, '%%Y-%%m') = DATE_FORMAT(t.real_date, '%%Y-%%m'))
					OR (@agg = 'accrued' AND DATE_FORMAT(a.capture_time, '%%Y-%%m') <= DATE_FORMAT(t.real_date, '%%Y-%%m'))
						)
					)
				)
			%v
			GROUP BY 
			CASE 
				WHEN @twoModel = 0 THEN 
				CASE 
					WHEN   @idList = '' THEN t.period_val
					ELSE i.device_id
				END
				ELSE t.period_val
			END ,
			 CASE WHEN @twoModel = 1 AND @idList != '' THEN i.device_id ELSE NULL END 
			ORDER BY stats, device_id;
		`, TenantVar)
	)

	type7Range := []sysaly.PositionStats{}
	err = tx.Select(&type7Range, baseSql)
	logsys.Printf("dateDatas:%#v  select:%#v  err:%v \n", type7Range, baseSql, err)

	if err != nil {
		tx.Rollback()
		logsys.Printf("查询失败:%#v", err)
		return nil, err
	}

	tx.MustExec(`DROP PROCEDURE IF EXISTS split_to_table;`)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	datas[sub_type] = type7Range

	UseSysaly.Data = datas
	return &UseSysaly, nil
}

// 按一天24小时聚合
func (mt *AlySysOperator) AlertHoursDay(startTime, endTime, TenantId string) (*sysaly.HourSisList[any], error) {

	var datas = make(map[string]any)
	var TenantVar string
	var DBConn = databases.DBMysql
	var UseSysaly = sysaly.HourSisList[any]{TotalData: make(map[string]any), Data: datas}

	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}

	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
	//限制筛选条件
	if TenantId != "" {
		TenantVar = fmt.Sprintf("  AND BIN_TO_UUID(a.tenant_id) = '%v' ", TenantId)
	}

	var baseSql = fmt.Sprintf(`
		WITH RECURSIVE hours AS (
			SELECT 0 AS h
			UNION ALL
			SELECT h + 1 FROM hours WHERE h < 23
		),
		agg AS (
			SELECT 
				HOUR(a.created_time) AS hr,
				COUNT(*) AS cnt
			FROM Alerts a
			WHERE a.created_time BETWEEN 
				COALESCE(@start_time, (SELECT MIN(created_time) FROM Alerts)) 
				AND @end_time
				%v
			GROUP BY HOUR(a.created_time)
		)
		SELECT 
			h.h AS stats,
			COALESCE(e.cnt, 0) AS count,
			SUM(COALESCE(e.cnt, 0)) OVER (
				ORDER BY h.h
			) AS all_count
		FROM hours h
		LEFT JOIN agg e
			ON e.hr = h.h
		ORDER BY h.h;
	`, TenantVar)

	//累计总数 到当前时间, 不传日期区间  则表示页面首次加载，返回累计总数
	if startTime == "" && endTime == "" {
		tx.Exec("set @end_time = CURDATE();")
		hourAccuredRange := []sysaly.PositionStats{}
		err = tx.Select(&hourAccuredRange, baseSql)
		logsys.Printf("dateDatas select:%#v  err:%v \n", baseSql, err)
		if err != nil {
			tx.Rollback()
			logsys.Printf("查询失败:%#v", err)
			return nil, err
		}
		UseSysaly.TotalData["accrued_hour"] = hourAccuredRange
	}

	//日期区间数 默认 30天内
	if startTime != "" {
		tx.Exec(fmt.Sprintf("set @start_time = '%v';", startTime))
	} else {
		tx.Exec("set @start_time = CURDATE() - 30;")
	}
	if endTime != "" {
		tx.Exec(fmt.Sprintf("set @end_time = '%v';", endTime))
	} else {
		tx.Exec("set @end_time = CURDATE();")
	}

	hourDateRange := []sysaly.PositionStats{}
	err = tx.Select(&hourDateRange, baseSql)
	logsys.Printf("dateDatas select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		tx.Rollback()
		logsys.Printf("查询失败:%#v", err)
		return nil, err
	}
	UseSysaly.Data["hour"] = hourDateRange

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &UseSysaly, nil
}

// 按天/周/月 聚合
func (mt *AlySysOperator) AlertsWeekMonth(DateRange, TenantId, endAt, sub_type string, devIds *string, twoModel *bool) (*sysaly.SysClientList[any], error) {

	// var ModeVar string
	var datas = make(map[string]any)
	var TenantVar string
	var DBConn = databases.DBMysql
	var UseSysaly = sysaly.SysClientList[any]{Data: datas}

	//dateRange 按天 date，周 week，月 month
	if DateRange == "" {
		DateRange = "month" //默认月度
	}
	if !DefDateRange[DateRange] {
		return nil, fmt.Errorf("不支持的时间模式查询方式")
	}

	if sub_type == "" {
		sub_type = "date" //该函数查 时间类型的统计
	}

	/////////////////////////////////////////////////环境变量设置

	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	///////////存储函数 生成 id 统计维度，只在twoModel有效时执行
	tx.MustExec("DROP PROCEDURE IF EXISTS split_to_table;")
	_, err = tx.Exec(` 
		CREATE PROCEDURE split_to_table(IN idList TEXT)
		BEGIN
			DECLARE pos INT DEFAULT 1;
			DECLARE commaPos INT;
			DECLARE curId TEXT;
			DROP TEMPORARY TABLE IF EXISTS tmp_ids;
			CREATE TEMPORARY TABLE tmp_ids (device_id INT);
			WHILE pos > 0 DO
				SET commaPos = LOCATE(',', idList, pos);

				IF commaPos > 0 THEN
					SET curId = SUBSTRING(idList, pos, commaPos - pos);
					INSERT INTO tmp_ids VALUES (CAST(curId AS UNSIGNED));
					SET pos = commaPos + 1;
				ELSE
					SET curId = SUBSTRING(idList, pos);
					IF curId <> '' THEN
						INSERT INTO tmp_ids VALUES (CAST(curId AS UNSIGNED));
					END IF;
					SET pos = 0;
				END IF;
			END WHILE;
		END; `)

	if err != nil {
		logsys.Printf("存储过程创建失败:%#v", err)
		return nil, err
	}
	if endAt == "" {
		endAt = "CURDATE()"
	} else {
		endAt = "'" + endAt + "'"
	}
	enabTwo := 0
	if *twoModel {
		enabTwo = 1
	}

	// 设置会话变量
	envSetSql := fmt.Sprintf(`SET @range_type = '%v';`, DateRange)

	tx.MustExec(envSetSql)
	envSetSql1 := fmt.Sprintf(`SET @agg = '%v'; `, sub_type)
	tx.MustExec(envSetSql1)
	envSetSql2 := fmt.Sprintf(`SET @endDate = %v; `, endAt)
	tx.MustExec(envSetSql2)
	envSetSql3 := fmt.Sprintf(`	SET @idList = '%v'; ; `, *devIds)
	tx.MustExec(envSetSql3)
	envSetSql4 := fmt.Sprintf(`SET @twoModel = '%v';`, enabTwo)
	_, err = tx.Exec(envSetSql4)

	logsys.Printf("环境设置%verr:%#v", envSetSql, err)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//调用存储过程
	_, err = tx.Exec("CALL split_to_table(@idList);")
	if err != nil {
		tx.Rollback()
		logsys.Printf("存储过程调用%v失败:%#v", devIds, err)
		return nil, err
	}

	//////////筛选条件
	if TenantId != "" {
		TenantVar = fmt.Sprintf("  AND BIN_TO_UUID(a.tenant_id) = '%v' ", TenantId)
	}

	// DateRange week  month
	var (
		baseSql = fmt.Sprintf(`  
			WITH RECURSIVE periods AS (
				SELECT 1 AS step
				UNION ALL
				SELECT step + 1
				FROM periods
				WHERE step < 7
			),
			time_dim AS (
				SELECT 
					CASE 
						WHEN @range_type = 'date'  
							THEN DATE_FORMAT(@endDate - INTERVAL (7 - p.step) DAY, '%%Y-%%m-%%d')
						WHEN @range_type = 'week'  
							THEN DATE_FORMAT(
									STR_TO_DATE(CONCAT(YEARWEEK(@endDate - INTERVAL (7 - p.step) WEEK, 1), ' Monday'), '%%X%%V %%W'),
									'%%x-%%v'
								)
						WHEN @range_type = 'month' 
							THEN DATE_FORMAT(@endDate - INTERVAL (p.step - 1) MONTH, '%%Y-%%m')
					END AS period_val,
					CASE 
						WHEN @range_type = 'date'  
							THEN @endDate - INTERVAL (7 - p.step) DAY
						WHEN @range_type = 'week'  
							THEN @endDate - INTERVAL (7 - p.step) WEEK
						WHEN @range_type = 'month' 
							THEN @endDate - INTERVAL (p.step - 1) MONTH
					END AS real_date,
					p.step
				FROM periods p
			),
			id_dim AS (
				-- 如果传入 idList 非空，取 split_to_table 的结果；否则返回一行 NULL
				SELECT device_id FROM tmp_ids
				UNION ALL
				SELECT NULL WHERE @idList = ''
			)
			SELECT 
			CASE
				WHEN @twoModel = 0 THEN 
				CASE 
				WHEN   @idList = '' THEN t.period_val
            	ELSE i.device_id
				END
				ELSE t.period_val
			END AS stats,
			(CASE WHEN @twoModel = 1 AND @idList != '' THEN i.device_id ELSE NULL END) AS device_id,
			COUNT(DISTINCT a.id) AS count,
				@range_type AS time_range,
				@agg  AS data_type
			FROM time_dim t
			JOIN id_dim i
			LEFT JOIN Alerts a  
				ON (i.device_id IS NULL OR a.device_id = i.device_id)
			AND (
					-- 日维度
					(@range_type = 'date' AND 
						(
						(@agg = 'date'  AND DATE(a.created_time) = t.real_date)
					OR (@agg = 'accrued' AND DATE(a.created_time) <= t.real_date)
						)
					)
				OR (@range_type = 'week' AND 
						(
						(@agg = 'date'  AND YEARWEEK(a.created_time, 1) = YEARWEEK(t.real_date, 1))
					OR (@agg = 'accrued' AND YEARWEEK(a.created_time, 1) <= YEARWEEK(t.real_date, 1))
						)
					)
				OR (@range_type = 'month' AND 
						(
						(@agg = 'date'  AND DATE_FORMAT(a.created_time, '%%Y-%%m') = DATE_FORMAT(t.real_date, '%%Y-%%m'))
					OR (@agg = 'accrued' AND DATE_FORMAT(a.created_time, '%%Y-%%m') <= DATE_FORMAT(t.real_date, '%%Y-%%m'))
						)
					)
				)
			%v
			GROUP BY 
			CASE 
				WHEN @twoModel = 0 THEN 
				CASE 
				WHEN   @idList = '' THEN t.period_val
				ELSE i.device_id
				END
				ELSE t.period_val
				END ,
			CASE WHEN @twoModel = 1 AND @idList != '' THEN i.device_id ELSE NULL END 
			ORDER BY stats, device_id;
		`, TenantVar)
	)

	//stats, (CASE WHEN @twoModel = 1 THEN i.device_id ELSE NULL END)

	type7Range := []sysaly.PositionStats{}
	err = tx.Select(&type7Range, baseSql)
	logsys.Printf("dateDatas:%#v  select:%#v  err:%v \n", type7Range, baseSql, err)

	if err != nil {
		tx.Rollback()
		logsys.Printf("查询失败:%#v", err)
		return nil, err
	}

	tx.MustExec(`DROP PROCEDURE IF EXISTS split_to_table;`)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	datas[sub_type] = type7Range

	UseSysaly.Data = datas
	return &UseSysaly, nil
}

// 按天/周/月 聚合 采集数据
func (mt *AlySysOperator) GathersWeekMonth(DateRange, TenantId, endAt, sub_type string, ProfileTypeIds *string) (*sysaly.SysClientList[any], error) {

	// var ModeVar string
	var datas = make(map[string]any)
	var TenantVar string
	var DBConn = databases.DBMysql
	var UseSysaly = sysaly.SysClientList[any]{Data: datas}

	//dateRange 按天 date，周 week，月 month
	if !DefDateRange[DateRange] {
		return nil, fmt.Errorf("不支持的时间模式查询方式")
	}

	if sub_type == "" { // 按时间增量 date 或 累计 accrued
		sub_type = "date" //该函数查 时间类型的统计
	}

	/////////////////////////////////////////////////环境变量设置

	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	///////////存储函数 生成 id 统计维度，只在twoModel有效时执行
	tx.MustExec("DROP PROCEDURE IF EXISTS split_to_profile;")
	_, err = tx.Exec(` 
		CREATE PROCEDURE split_to_profile(IN idList TEXT)
		BEGIN
			DECLARE pos INT DEFAULT 1;
			DECLARE commaPos INT;
			DECLARE curId TEXT;
			DROP TEMPORARY TABLE IF EXISTS tmp_ids;
			CREATE TEMPORARY TABLE tmp_ids (type_id INT);
			WHILE pos > 0 DO
				SET commaPos = LOCATE(',', idList, pos);

				IF commaPos > 0 THEN
					SET curId = SUBSTRING(idList, pos, commaPos - pos);
					INSERT INTO tmp_ids VALUES (CAST(curId AS UNSIGNED));
					SET pos = commaPos + 1;
				ELSE
					SET curId = SUBSTRING(idList, pos);
					IF curId <> '' THEN
						INSERT INTO tmp_ids VALUES (CAST(curId AS UNSIGNED));
					END IF;
					SET pos = 0;
				END IF;
			END WHILE;
		END; `)

	if err != nil {
		logsys.Printf("存储过程创建失败:%#v", err)
		return nil, err
	}
	if endAt == "" {
		endAt = "CURDATE()"
	} else {
		endAt = "'" + endAt + "'"
	}
	// 设置会话变量
	envSetSql := fmt.Sprintf(`SET @range_type = '%v';`, DateRange)

	tx.MustExec(envSetSql)
	envSetSql1 := fmt.Sprintf(`SET @agg = '%v'; `, sub_type)
	tx.MustExec(envSetSql1)
	envSetSql2 := fmt.Sprintf(`SET @endDate = %v; `, endAt)
	tx.MustExec(envSetSql2)
	envSetSql3 := fmt.Sprintf(`	SET @idList = '%v'; ; `, *ProfileTypeIds)
	tx.MustExec(envSetSql3)
	envSetSql4 := fmt.Sprintf(`SET @twoModel = '%v';`, 1)
	_, err = tx.Exec(envSetSql4)

	logsys.Printf("环境设置%verr:%#v", envSetSql, err)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//调用存储过程
	_, err = tx.Exec("CALL split_to_profile(@idList);")
	if err != nil {
		tx.Rollback()
		logsys.Printf("存储过程调用%v失败:%#v", ProfileTypeIds, err)
		return nil, err
	}

	//////////筛选条件
	if TenantId != "" {
		TenantVar = fmt.Sprintf("  AND BIN_TO_UUID(a.tenant_id) = '%v' ", TenantId)
	}

	// DateRange week  month
	var (
		baseSql = fmt.Sprintf(`  
			WITH RECURSIVE periods AS (
				SELECT 1 AS step
				UNION ALL
				SELECT step + 1
				FROM periods
				WHERE step < 7
			),
			time_dim AS (
				SELECT 
					CASE 
						WHEN @range_type = 'date'  
							THEN DATE_FORMAT(@endDate - INTERVAL (7 - p.step) DAY, '%%Y-%%m-%%d')
						WHEN @range_type = 'week'  
							THEN DATE_FORMAT(
									STR_TO_DATE(CONCAT(YEARWEEK(@endDate - INTERVAL (7 - p.step) WEEK, 1), ' Monday'), '%%X%%V %%W'),
									'%%x-%%v'
								)
						WHEN @range_type = 'month' 
							THEN DATE_FORMAT(@endDate - INTERVAL (p.step - 1) MONTH, '%%Y-%%m')
					END AS period_val,
					CASE 
						WHEN @range_type = 'date'  
							THEN @endDate - INTERVAL (7 - p.step) DAY
						WHEN @range_type = 'week'  
							THEN @endDate - INTERVAL (7 - p.step) WEEK
						WHEN @range_type = 'month' 
							THEN @endDate - INTERVAL (p.step - 1) MONTH
					END AS real_date,
					p.step
				FROM periods p
			),
			id_dim AS (
				-- 如果传入 idList 非空，取 split_to_table 的结果；否则返回一行 NULL
				SELECT type_id FROM tmp_ids
				UNION ALL
				SELECT NULL WHERE @idList = ''
			)
			SELECT 
			CASE
				WHEN @twoModel = 0 THEN 
				CASE 
				WHEN   @idList = '' THEN t.period_val
            	ELSE i.type_id
				END
				ELSE t.period_val
			END AS stats,
			(CASE WHEN @twoModel = 1 AND @idList != '' THEN i.type_id ELSE NULL END) AS type_id,
			COUNT(DISTINCT a.id) AS count,
				@range_type AS time_range,
				@agg  AS data_type
			FROM time_dim t
			JOIN id_dim i
			JOIN Profile pf ON (i.type_id IS NULL OR pf.type_id = i.type_id)
			LEFT JOIN CaptureLogs a  
				ON a.matched_profile_id = pf.id
			AND (
					-- 日维度
					(@range_type = 'date' AND 
						(
						(@agg = 'date'  AND DATE(a.capture_time) = t.real_date)
					OR (@agg = 'accrued' AND DATE(a.capture_time) <= t.real_date)
						)
					)
				OR (@range_type = 'week' AND 
						(
						(@agg = 'date'  AND YEARWEEK(a.capture_time, 1) = YEARWEEK(t.real_date, 1))
					OR (@agg = 'accrued' AND YEARWEEK(a.capture_time, 1) <= YEARWEEK(t.real_date, 1))
						)
					)
				OR (@range_type = 'month' AND 
						(
						(@agg = 'date'  AND DATE_FORMAT(a.capture_time, '%%Y-%%m') = DATE_FORMAT(t.real_date, '%%Y-%%m'))
					OR (@agg = 'accrued' AND DATE_FORMAT(a.capture_time, '%%Y-%%m') <= DATE_FORMAT(t.real_date, '%%Y-%%m'))
						)
					)
				)
			%v
			GROUP BY 
			CASE 
				WHEN @twoModel = 0 THEN 
				CASE 
				WHEN   @idList = '' THEN t.period_val
				ELSE i.type_id
				END
				ELSE t.period_val
				END ,
			CASE WHEN @twoModel = 1 AND @idList != '' THEN i.type_id ELSE NULL END 
			ORDER BY stats, type_id;
		`, TenantVar)
	)

	//stats, (CASE WHEN @twoModel = 1 THEN i.type_id ELSE NULL END)

	type7Range := []sysaly.PositionStats{}
	err = tx.Select(&type7Range, baseSql)
	logsys.Printf("dateDatas:%#v  week monthSql:%#v  err:%v \n", type7Range, baseSql, err)

	if err != nil {
		tx.Rollback()
		logsys.Printf("查询失败:%#v", err)
		return nil, err
	}

	tx.MustExec(`DROP PROCEDURE IF EXISTS split_to_profile;`)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	datas[sub_type] = type7Range

	UseSysaly.Data = datas
	return &UseSysaly, nil
}
