package sysclients

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 省份 列表查询
func (oc *ORMCleint) ProvinceList(Page, PageSize int, Ctx *gin.Context, conds []predicate.Province) (*seclients.ProvincePage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.ProvincePage
	offset := (Page - 1) * PageSize

	total, err := client.Debug().Province.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().Province.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select("id", "code", "name", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//All(Ctx)
	list := result.Data
	sysdebug.Println("fm_gov_province retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	sysdebug.Printf("err:%#v \n", err)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("fm_gov_province", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 城市 列表查询
func (oc *ORMCleint) CityList(Page, PageSize int, Ctx *gin.Context, conds []predicate.GovCity) (*seclients.CityPage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.CityPage
	offset := (Page - 1) * PageSize

	total, err := client.Debug().GovCity.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().GovCity.Query().
		Where(conds...).
		Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).Select("id", "code", "name", "province_code", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//.All(Ctx)
	list := result.Data
	sysdebug.Println("fm_gov_city retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("fm_gov_city", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 区县 列表查询
func (oc *ORMCleint) AreasList(Page, PageSize int, Ctx *gin.Context, conds []predicate.GovArea) (*seclients.AreaPage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.AreaPage
	offset := (Page - 1) * PageSize

	total, err := client.Debug().GovArea.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().GovArea.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select("id", "code", "name", "province_code", "city_code", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//All(Ctx)
	list := result.Data
	sysdebug.Println("fm_gov_area retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("fm_gov_area", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 乡镇街道 列表查询
func (oc *ORMCleint) StreetsList(Page, PageSize int, Ctx *gin.Context, conds []predicate.GovStreet) (*seclients.StreetPage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.StreetPage
	offset := (Page - 1) * PageSize

	total, err := client.Debug().GovStreet.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().GovStreet.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select("id", "code", "name", "province_code", "city_code", "area_code", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//All(Ctx)
	list := result.Data

	sysdebug.Println("fm_gov_street retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("fm_gov_street", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 预警记录 列表查询 ids, device_id, alert_level, c, enabled, capture_log_id, tenantId, StartAt, EndAt
func (oc *ORMCleint) AlertsList(Page, PageSize, ids, device_id, alert_level, capture_log_id *int, Ctx *gin.Context, statusTr, tenantId, StartAt, EndAt string) (*seclients.AlertsList, error) {

	var (
		// client = databases.EntClient
		DBConn  = databases.DBMysql
		total   int
		baseSql = `
		SELECT a.id,
		BIN_TO_UUID(a.tenant_id) as tenant_id,   
		a.capture_log_id,
		a.device_id,
		a.alert_level,
		a.status,
		a.fm_user_id,
		a.handled_time,
		a.remarks,
		a.created_time,
		
		b.matched_profile_id,
		b.capture_image_url,
		d.type_id as profile_type_id,
		c.location
		FROM  Alerts  a
		LEFT JOIN CaptureLogs b  ON a.capture_log_id = b.id
		LEFT JOIN Profile d ON d.id = b.matched_profile_id 
		LEFT JOIN Device c ON c.id = b.device_id  `
	)

	// if tenantId != "" {
	// 	return oc.alertsListByTenantId(*Page, *PageSize, tenantId)
	// }

	var result seclients.AlertsList

	Filters := " WHERE a.capture_log_id != 0 "
	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id)  = '%v' ", uuid.MustParse(tenantId))
	}
	if alert_level != nil && *alert_level != 0 {
		Filters += fmt.Sprintf(" AND a.alert_level  = '%v' ", *alert_level)
	}

	if capture_log_id != nil && *capture_log_id != 0 {
		Filters += fmt.Sprintf(" AND a.capture_log_id  = '%v' ", *capture_log_id)
	}
	if ids != nil && *ids != 0 {
		Filters += fmt.Sprintf(" AND a.id  = '%v' ", *ids)
	}
	if device_id != nil && *device_id != 0 {
		Filters += fmt.Sprintf(" AND a.device_id  = '%v' ", *device_id)
	}
	if statusTr != "" {
		Filters += fmt.Sprintf(" AND a.status  = '%v' ", statusTr)
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		sysdebug.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		sysdebug.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)
		Filters += fmt.Sprintf(" AND (a.created_time between '%v' and '%v' )", StartAt, EndAt)

	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `Alerts` a " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	page := *Page
	page_size := *PageSize
	Limits := fmt.Sprintf(" ORDER BY a.id DESC LIMIT %v OFFSET %v ", page_size, (page-1)*page_size)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query Alerts by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("Alerts", &seclients.FieldTags{SortTagFilter: "asc"})
		result.Columns = cm
		result.Total = total
		result.Page = *Page
		result.PageSize = *PageSize
	}
	return &result, nil

}

// 预警 信息数据列表查询 租户信息
func (oc *ORMCleint) alertsListByTenantId(Page, PageSize int, tenantId string) (*seclients.AlertsList, error) {

	var (
		baseSql = `
		SELECT id,
		BIN_TO_UUID(tenant_id) as tenant_id,   
		capture_log_id,
		device_id,
		alert_level,
		status,
		fm_user_id,
		handled_time,
		remarks,
		created_time
		FROM  Alerts  `

		total  int
		result = seclients.AlertsList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE id > 0 "
	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(tenantId))

	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `Alerts`  " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query Alerts by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("Alerts", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil
}

// 抓拍记录 列表查询
func (oc *ORMCleint) CaptureLogInfoById(Ctx *gin.Context, conds []predicate.CaptureLogs) (*genclients.CaptureLogs, error) {

	var (
		client = databases.EntClient
	)

	dats, err := client.Debug().CaptureLogs.Query().
		Where(conds...).First(Ctx)

	sysdebug.Println("Device retrieved", dats, zap.Int("count", len(conds)))
	if err != nil {
		return nil, err
	}

	return dats, nil

}

// 抓拍日志信息数据列表查询 租户信息
func (oc *ORMCleint) CaptureLogsList(Page, PageSize int64, startAt, endAt string, Ctx *gin.Context, PParm *seclients.CaptureLogInfo, FilterTag *seclients.FieldTags) (*seclients.CaptureLogsList, error) {

	var (
		baseSql = `
		SELECT a.id,
		BIN_TO_UUID(a.tenant_id) as tenant_id,  
		a.device_id,
		a.capture_time,
		a.capture_image_url,
		a.has_alert,
		a.matched_profile_id, 
		a.device_name,
		a.device_location,
		a.content,
		a.match_score,
		a.func_type,
		t.type_name,
		b.name,
		b.type_id,
		b.room_id,
		IFNULL(c.id, '0') AS alert_id,
		c.status, 
		c.alert_level,
		IFNULL(d.id , '0') AS face_id 
		FROM  CaptureLogs a
		LEFT JOIN Profile b ON b.id = a.matched_profile_id 
		LEFT JOIN ProfileType t ON t.id = b.type_id 
		LEFT JOIN Alerts c ON c.capture_log_id = a.id 
		LEFT JOIN Face d ON d.profile_id = b.id 
		`
		total  int64
		result = seclients.CaptureLogsList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE a.id > 0 "
	if PParm.MatchedProfileID != nil {
		Filters += fmt.Sprintf(" AND a.matched_profile_id = '%v' ", *PParm.MatchedProfileID)
	}
	if PParm.TypeId != nil {
		Filters += fmt.Sprintf(" AND b.type_id = '%v' ", *PParm.TypeId)
	}
	if PParm.DeviceID != nil {
		Filters += fmt.Sprintf(" AND a.device_id = '%v' ", *PParm.DeviceID)
	}
	if PParm.TenantID != nil {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id)  = '%v' ", uuid.MustParse(*PParm.TenantID))
	}

	if PParm.Status != nil {
		Filters += fmt.Sprintf(" AND c.status  = '%v' ", *PParm.Status)
	}
	if PParm.FuncType != nil {
		Filters += fmt.Sprintf(" AND a.func_type  = '%v' ", *PParm.FuncType)
	}
	if PParm.FaceId != nil {
		Filters += fmt.Sprintf(" AND d.id  = '%v' ", *PParm.FaceId)
	}

	if startAt != "" && endAt != "" {
		Filters += fmt.Sprintf(" AND (a.capture_time between '%v' and '%v') ", startAt, endAt)
	}

	//总数
	totalSql := `SELECT COUNT(*) AS Total FROM CaptureLogs  a 
	LEFT JOIN Profile b ON b.id = a.matched_profile_id 
	LEFT JOIN Alerts c ON c.capture_log_id = a.id 
	LEFT JOIN Face d ON d.profile_id = b.id 
	` + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY a.id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query CaptureLogs ,  baseSql:%v err:%#v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cms := []erps.GcDesc{}
		if *PParm.FuncType == 1 {
			cms = append(cms, erps.GcDesc{Name: "face_id", CName: "面部ID", IsVisible: true, Sort: "caption92"})
		}
		cm, _ := databases.FmGlobalMap("CaptureLogs", FilterTag)
		cms = append(cms, cm...)
		if *PParm.FuncType == 0 {
			cms = append(cms, erps.GcDesc{Name: "capture_time", CName: "采集时间", IsVisible: true, IsSearchable: true})
		}
		result.Columns = cms
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
		sysdebug.Printf("data len:%v \n", len(result.Data))
	}
	return &result, nil
}

// 设备列表查询
func (oc *ORMCleint) DevicesList(Page, PageSize int, Ctx *gin.Context, conds []predicate.Devices, tenantId string) (*seclients.DevicesList, error) {

	var (
		client = databases.EntClient
	)

	if tenantId != "" {
		return oc.deviceListByTenantId(Page, PageSize, tenantId)
	}

	var result seclients.DevicesList
	offset := (Page - 1) * PageSize

	total, err := client.Debug().Devices.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	list, err := client.Debug().Devices.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).All(Ctx)

	sysdebug.Println("Device retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("Device", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize

	}
	return &result, err

}

// 设备信息数据列表查询 租户信息
func (oc *ORMCleint) deviceListByTenantId(Page, PageSize int, tenantId string) (*seclients.DevicesList, error) {

	var (
		baseSql = `
		SELECT id,
		BIN_TO_UUID(tenant_id) as tenant_id, 
		name,
		url,
		location,
		func_type,
		roi_x,
		roi_y,
		roi_width,
		roi_height,
		target_fps,
		enabled,
		created_time,
		updated_time
		FROM Device  `

		total  int
		result = seclients.DevicesList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE id > 0 "
	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(tenantId))

	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `Device`  " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query Device by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("Device", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil
}

// 面部信息列表查询
func (oc *ORMCleint) FacesList(Page, PageSize, profile_id, is_primary *int, Ctx *gin.Context, conds []predicate.Faces, tenantId string) (*seclients.FacesList, error) {

	var (
		client = databases.EntClient
	)

	if tenantId != "" {
		return oc.faceListByTenantId(Page, PageSize, profile_id, is_primary, tenantId)
	}

	var result seclients.FacesList
	offset := (*Page - 1) * (*PageSize)

	list, err := client.Debug().Faces.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(*PageSize).
		Offset(offset).All(Ctx)

	sysdebug.Printf("ent err:%v, list:%v  Ctx:%v conds:%#v \n", err, len(list), Ctx, conds)

	if err != nil {
		return nil, err
	}

	total, err := client.Debug().Faces.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	sysdebug.Println("profile retrieved", zap.Int("count", len(list)), zap.Int("total", total))

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("Face", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = *Page
		result.PageSize = *PageSize

	}
	return &result, err

}

// 面部信息数据列表查询
func (oc *ORMCleint) faceListByTenantId(Page, PageSize, profile_id, is_primary *int, tenantId string) (*seclients.FacesList, error) {

	var (
		baseSql = `
		SELECT id,
		BIN_TO_UUID(tenant_id) as tenant_id,
			profile_id,
			face_embedding,
			image_url,
			is_primary,
			created_time,
			updated_time
		FROM Face  `

		total  int
		result = seclients.FacesList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE id > 0 "
	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(tenantId))

	}
	if is_primary != nil {
		Filters += fmt.Sprintf(" AND is_primary = '%v' ", *is_primary)
	}
	if profile_id != nil {
		Filters += fmt.Sprintf(" AND profile_Id = '%v' ", profile_id)
	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `Face`  " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY id DESC LIMIT %v OFFSET %v ", *PageSize, (*Page-1)*(*PageSize))
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query Face by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("Face", nil)
		result.Columns = cm
		result.Total = total
		result.Page = *Page
		result.PageSize = *PageSize
	}
	return &result, nil
}

// 面部信息列表查询
func (oc *ORMCleint) FacesUpdate(c *gin.Context, PParm *seclients.FaceParam) (int, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	//上下文传递操作人信息
	if *PParm.DeleteFlag == "1" && PParm.ID != 0 {

		//记录操作
		Ctx := databases.WithOperatorFromGin(c)

		//查询人员信息  设置为
		profileInfo := seclients.FaceParam{}
		QuerySql := fmt.Sprintf(" SELECT id, profile_id FROM Face WHERE id = '%v'; ", PParm.ID)
		err := DBConn.Get(&profileInfo, QuerySql)
		if err != nil {
			return 0, err
		}

		//执行删除
		err = client.Debug().Faces.DeleteOneID(PParm.ID).Exec(Ctx)
		sysdebug.Printf("ent err:%v,  Ctx:%v conds:%#v \n", err, Ctx, PParm)
		if err != nil {
			return 0, err
		}

		// 更新人员信息的 状态为 0
		UpProfileSql := fmt.Sprintf(" update Profile Set enabled = '%v' WHERE id = '%v'; ", 0, *profileInfo.ProfileID)
		_, err = DBConn.Exec(UpProfileSql)
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	return 0, fmt.Errorf("ID或删除标记错误")
}

// 临时库查询
func (oc *ORMCleint) TemporaryFaceList(Page, PageSize int, ctx *gin.Context, conds []predicate.TemporaryFace, tenantId string) (*seclients.TemporaryFaceList, error) {

	var (
		client = databases.EntClient
	)

	if tenantId != "" {
		return oc.proTmpFaceListByTenantId(Page, PageSize, tenantId)
	}
	var result seclients.TemporaryFaceList

	offset := (Page - 1) * PageSize

	total, err := client.Debug().TemporaryFace.Query().Where(conds...).Count(ctx)
	if err != nil {
		return nil, err
	}

	list, err := client.Debug().TemporaryFace.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).All(ctx)

	sysdebug.Printf("ent err:%v,   conds:%#v \n", err, conds)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("TemporaryFace", nil)
		result.Data = list
		result.Total = total
		result.Columns = cm
		result.Page = 1
		result.PageSize = PageSize
	}

	return &result, nil

}

// 人员类型列表查询
func (oc *ORMCleint) proTmpFaceListByTenantId(Page, PageSize int, tenantId string) (*seclients.TemporaryFaceList, error) {

	var (
		baseSql = `SELECT id,
		BIN_TO_UUID(tenant_id) as tenant_id,
		face_embedding,
		img_url,
		first_seen_time,
		last_seen_time,
		last_seen_location,
		capture_count,
		expires_time
	FROM TemporaryFace  `

		total  int
		result = seclients.TemporaryFaceList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE id > 0 "
	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(tenantId))

	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `TemporaryFace`  " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query TemporaryFace by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("TemporaryFace", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil
}

// 人员列表查询
func (oc *ORMCleint) ProfilesList(Page, PageSize int, Ctx *gin.Context, conds []predicate.Profiles, StartAt, EndAt string, genProfile *genclients.Profiles) (*seclients.ProfilesList, error) {

	var (
		client = databases.EntClient
	)

	if genProfile.TenantID.String() != "" && genProfile.TenantID.String() != configs.DefTenantId {
		return oc.profileListByTenantId(Page, PageSize, StartAt, EndAt, genProfile)
	}

	var result seclients.ProfilesList
	offset := (Page - 1) * PageSize

	list, err := client.Debug().Profiles.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).All(Ctx)

	sysdebug.Printf("ent err:%v, list:%v  Ctx:%v conds:%#v \n", err, len(list), Ctx, conds)

	if err != nil {
		return nil, err
	}

	total, err := client.Debug().Profiles.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	sysdebug.Println("profile retrieved", zap.Int("count", len(list)), zap.Int("total", total))

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("Profile", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize

	}
	return &result, err

}

// 人员类型列表查询
func (oc *ORMCleint) profileListByTenantId(Page, PageSize int, StartAt, EndAt string, genProfile *genclients.Profiles) (*seclients.ProfilesList, error) {

	var (
		baseSql = `SELECT id, 
			name,
			BIN_TO_UUID(tenant_id) as tenant_id,
			type_id,
			id_card_number,
			phone_number,
			enabled,
			room_id,
			tmp_url,
			created_time,
			updated_time
			FROM  Profile `

		total  int
		result = seclients.ProfilesList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE id > 0 "
	if genProfile.TenantID.String() != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(genProfile.TenantID.String()))
	}
	if genProfile.TypeID != nil {
		Filters += fmt.Sprintf(" AND type_id = '%v' ", *genProfile.TypeID)
	}
	if genProfile.Name != nil {
		Filters += fmt.Sprintf(" AND name like '%v' ", *genProfile.Name+"%")
	}
	if genProfile.IDCardNumber != nil {
		Filters += fmt.Sprintf(" AND id_card_number = '%v' ", *genProfile.IDCardNumber)
	}
	if genProfile.PhoneNumber != nil {
		Filters += fmt.Sprintf(" AND phone_number = '%v' ", *genProfile.PhoneNumber)
	}
	if genProfile.RoomID != nil {
		Filters += fmt.Sprintf(" AND room_id = '%v' ", *genProfile.RoomID)
	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `Profile`  " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)
	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)
	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query Profile by tenantid:%v,  baseSql:%v total:%#v \n", genProfile.TenantID, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("Profile", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil
}

// 人员类型列表查询
func (oc *ORMCleint) ProfileTypesList(Page, PageSize, group_id int, enabled *int, ctx *gin.Context, conds []predicate.ProfileType, typeName, typeCode, tenantId string) (*seclients.ProfileSuppTypesList, error) {

	var (
		client  = databases.EntClient
		DBConn  = databases.DBMysql
		baseSql = `SELECT a.id,b.supplier,
		BIN_TO_UUID(a.tenant_id) as tenant_id, 
		a.type_name,a.warning_level,a.warning_enabled,a.type_code,a.description,
		a.enabled,a.face_validity_hours,
		a.created_time,a.updated_time
		FROM ProfileType a 
		LEFT JOIN group_profile_type_mapping d ON d.profile_type_id = a.id
		LEFT JOIN fm_alert_group c ON c.id = d.group_id 
		LEFT JOIN fm_tenant b ON b.tenant_id = a.tenant_id 
		`
	)

	var result seclients.ProfileSuppTypesList //ProfileTypesPointer

	total, err := client.Debug().ProfileType.Query().Where(conds...).Count(ctx)
	sysdebug.Printf("ent err:%v, list:%v conds:%#v \n", err, total, conds)
	if err != nil {
		return nil, err
	}

	offset := (Page - 1) * PageSize
	var Filter = " WHERE a.id > 0 "
	if typeName != "" {
		Filter += fmt.Sprintf(" AND a.type_name = '%v' ", typeName)
	}
	if typeCode != "" {
		Filter += fmt.Sprintf(" AND a.type_code = '%v' ", typeCode)
	}
	if group_id != 0 {
		Filter += fmt.Sprintf(" AND d.group_id  = '%v' ", group_id)
	}
	if enabled != nil {
		Filter += fmt.Sprintf(" AND c.enabled = '%v' ", *enabled)
	}
	if tenantId != "" {
		// return oc.ProfileTypesListByTenantId(Page, PageSize, tenantId)
		Filter += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id) = '%v' ", uuid.MustParse(tenantId))
	}

	// list, err := client.Debug().ProfileType.Query().
	// 	Where(conds...).
	// 	Limit(PageSize).
	// 	Offset(offset).All(ctx)
	limits := fmt.Sprintf(" ORDER BY a.id DESC LIMIT %v OFFSET %v ", PageSize, offset)

	baseSql += Filter
	baseSql += limits
	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("ent err:%v, list:%v sql:%v,   \n", err, len(result.Data), baseSql)
	if err != nil {
		return nil, err
	}

	sysdebug.Println("ProfileType retrieved", zap.Int("count", len(result.Data)), zap.Int("total", total))
	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("ProfileType", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil

}

// 人员类型列表查询
func (oc *ORMCleint) ProfileTypesListByTenantId(Page, PageSize int, tenantId string) (*seclients.ProfileSuppTypesList, error) {

	var (
		baseSql = `SELECT a.id,
		BIN_TO_UUID(a.tenant_id) as tenant_id,
		a.type_name,
		b.supplier,
		a.type_code,
		a.description,
		a.enabled,
		a.created_time,
		a.updated_time
		FROM ProfileType a 
		LEFT JOIN fm_tenant b ON a.tenant_id = b.tenant_id `

		total  int
		result = seclients.ProfileSuppTypesList{}
		DBConn = databases.DBMysql
	)

	Filters := " WHERE a.id > 0 "

	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id)  = '%v' ", uuid.MustParse(tenantId))

	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `ProfileType` a " + Filters
	err := DBConn.QueryRow(totalSql).Scan(&total)

	sysdebug.Printf("ent err:%v,  totalSql:%v total:%#v \n", err, totalSql, total)

	if err != nil {
		return nil, err
	}

	//数据列表
	Limits := fmt.Sprintf(" ORDER BY a.id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits

	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("query ProfileType by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("ProfileType", nil)
		cm = append(cm, erps.GcDesc{Name: "supplier", CName: "酒店名", IsVisible: true})
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil
}
