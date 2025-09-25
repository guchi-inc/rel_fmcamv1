package sysclients

import (
	"database/sql"
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 人员类型 和分组 列表关联查询
func (oc *ORMCleint) GroupProfileTypeMapList(Page, PageSize int, Ctx *gin.Context, conds []predicate.GrouProfileTypeMapping) (*seclients.FmGroupTypeMappingPages, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.FmGroupTypeMappingPages
	offset := (Page - 1) * PageSize

	total, err := client.Debug().GrouProfileTypeMapping.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().GrouProfileTypeMapping.Query().
		Where(conds...).
		Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).Select("id", "group_id", "profile_type_id", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//.All(Ctx)
	list := result.Data
	sysdebug.Println("GrouProfileTypeMapping retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("group_profile_type_mapping", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 新增 人员类型和组的映射关系
func (oc *ORMCleint) GroupProfileTypeMapNew(ctx *gin.Context, PParm *seclients.ProfileTypeAlertInfoMap) (int64, error) {

	var (
		client = databases.EntClient
	)
	builder := client.GrouProfileTypeMapping.Create()

	if PParm.GroupID == 0 || PParm.ProfileTypeID == 0 {
		return 0, fmt.Errorf("参数错误")
	}

	if PParm.GroupID != 0 {
		builder.SetGroupID(PParm.GroupID)
	}
	if PParm.ProfileTypeID != 0 {
		builder.SetProfileTypeID(PParm.ProfileTypeID)
	}

	//上下文传递操作人信息
	user, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		return 0, fmt.Errorf("用户鉴权错误:%v", err.Error())
	}
	Ctx := databases.WithOperatorFromGin(ctx)
	builder.SetCreator(user.LoginName)
	newMap, err := builder.Save(Ctx)
	if err != nil {
		return 0, err
	}

	return newMap.ID, nil
}

// 人员类型 和分组 信息关联查询，包括预警等级，预警组信息
func (oc *ORMCleint) GroupProfileTypeInfoList(Page, PageSize, group_id, profile_type_id int, Ctx *gin.Context) (*seclients.FmProfileTypeAlertInfoPages, error) {

	var (
		DBConn  = databases.DBMysql
		total   int
		baseSql = `
		SELECT
		a.group_id,
		a.profile_type_id,
		b.type_name,
		b.type_code,
		b.description,
		b.face_validity_hours,
		b.warning_level,
		b.warning_enabled,
		b.deleteable,

		b.enabled,
		c.group_name,
		c.customization, 
		c.group_type,
		d.level,
		d.action,
		d.alarm_sound,
		BIN_TO_UUID(b.tenant_id) as tenant_id,    
		a.created_time
		FROM  group_profile_type_mapping  a
		LEFT JOIN ProfileType b ON b.id = a.profile_type_id
		LEFT JOIN fm_alert_group c ON c.id = a.group_id
		LEFT JOIN fm_alert_definition d ON  d.profile_type_id = a.profile_type_id`
	)

	var result = seclients.FmProfileTypeAlertInfoPages{}

	Filters := " WHERE a.id > 0 "

	if group_id != 0 {
		Filters += fmt.Sprintf(" AND a.group_id = '%v' ", group_id)
	}
	if profile_type_id != 0 {
		Filters += fmt.Sprintf(" AND a.profile_type_id = '%v' ", profile_type_id)
	}

	//总数
	totalSql := fmt.Sprintf(`
	SELECT COUNT(*) 	FROM  group_profile_type_mapping  a
	LEFT JOIN ProfileType b ON b.id = a.profile_type_id
	LEFT JOIN fm_alert_group c ON c.id = a.group_id
	LEFT JOIN fm_alert_definition d ON  d.profile_type_id = a.profile_type_id
	` + Filters)
	err := DBConn.Get(&total, totalSql)
	if err != nil {
		return nil, err
	}

	//查数据
	Limits := fmt.Sprintf(" ORDER BY a.id DESC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql += Filters
	baseSql += Limits
	err = DBConn.Select(&result.Data, baseSql)
	sysdebug.Printf("get fm_alert_group mapping:%v group_id:%#v,  profile_type_id:%#v err:%v\n", baseSql, group_id, profile_type_id, err)

	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("group_profile_type_mapping", nil)

		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 预警组 新增定义信息
func (oc *ORMCleint) NewGroupAlertsInfo(Ctx *gin.Context, alertGpInfo *seclients.FmAlertGroup) (int64, error) {

	var (
		// client = databases.EntClient
		DBConn = databases.DBMysql
	)

	//  SQL
	newuuid := uuid.New().String()
	query := fmt.Sprintf(`
		INSERT INTO fm_alert_group ( group_name,group_type,tenant_id, customization, description, creator) VALUES ('%v','%v', %v, '%v', '%v', '%v')
		`, alertGpInfo.GroupName, alertGpInfo.GroupType, fmt.Sprintf("(UUID_TO_BIN('%v'))", alertGpInfo.TenantID), alertGpInfo.Customization, alertGpInfo.Description, alertGpInfo.Creator)

	//写入 执行
	rst, err := DBConn.Exec(query)
	sysdebug.Printf("new fm_alert_group:%v insert:%#v, with db rst:%#v err:%v\n", newuuid, query, rst, err)

	if err != nil {
		sysdebug.Println("新增失败: ", err)
		return 0, err
	}

	return rst.LastInsertId()

}

// 预警组信息 修改
func (oc *ORMCleint) GroupAlertsInfoUpdate(ctx *gin.Context, PParm *seclients.FmAlertGroup) ([]int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	var ids []int64

	orginGroup, err := client.FmAlertGroup.Get(ctx, PParm.ID)
	if err != nil {
		return nil, fmt.Errorf("不存在该分组:%v \n", PParm.ID)
	}

	PParm.GroupType = orginGroup.GroupType
	//删除 预警分组
	if PParm.DeleteFlag == "1" {
		if PParm.TenantID == configs.DefTenantId {
			return nil, fmt.Errorf("内置组和人员类型不可删除")
		}
		baseSql := fmt.Sprintf("DELETE FROM fm_alert_group WHERE id = %v; ", PParm.ID)
		_, err := DBConn.Exec(baseSql)
		if err != nil {
			return nil, err
		}

		var rst sql.Result
		if PParm.GroupType == 1 {
			//同时也删除分组下的预警等级 level定义
			baseLevelSql := fmt.Sprintf("DELETE FROM fm_alert_definition WHERE alert_group_id = %v; ", PParm.ID)
			rst, err = DBConn.Exec(baseLevelSql)
			if err != nil {
				return nil, err
			}

		} else if PParm.GroupType == 2 {

			//同时也删除分组下的人员类型 定义
			baseTypeSql := fmt.Sprintf("DELETE FROM ProfileType WHERE deleteable = '1' AND id in ( SELECT profile_type_id FROM group_profile_type_mapping WHERE group_id = '%v') ; ", PParm.ID)
			rst, err = DBConn.Exec(baseTypeSql)
			if err != nil {
				return nil, err
			}

			//删除关联表信息
			delGroupMapSql := fmt.Sprintf("DELETE FROM group_profile_type_mapping WHERE  group_id = '%v'  ; ", PParm.ID)
			rst, err = DBConn.Exec(delGroupMapSql)
			if err != nil {
				return nil, err
			}
		} else {
			rst = nil
		}
		if rst != nil {
			delids, _ := rst.RowsAffected()
			ids = append(ids, delids)
		}
		return ids, nil
	}

	update := client.FmAlertGroup.UpdateOneID(PParm.ID)

	var changed = false
	if PParm.GroupName != "" {
		update.SetGroupName(PParm.GroupName)
		changed = true
	}
	if PParm.Customization != "" {
		update.SetCustomization(PParm.Customization)
		changed = true
	}

	if PParm.Description != "" {
		update.SetDescription(PParm.Description)
		changed = true
	}

	if PParm.Enabled != nil {
		update.SetEnabled(*PParm.Enabled)
		changed = true
	}
	//等级更新
	if PParm.GroupType == 1 && len(PParm.NewLevelList) > 0 {
		datas, err := oc.LevelAlertListUpdate(ctx, &seclients.ParamAlertList{ID: PParm.ID,
			NewLevelList: PParm.NewLevelList})
		sysdebug.Printf("update LevelAlertListUpdate:%v err:%v", datas, err)

		if err != nil {
			return nil, err
		}
		ids = append(ids, datas...)
	}
	if PParm.GroupType == 2 && len(PParm.NewProfileTypeList) > 0 {
		//人员类型组 更新
		datas, err := oc.ProfileTypeListUpdate(ctx,
			&seclients.ParamProfileTypeLists{ID: PParm.ID,
				NewTypeList: PParm.NewProfileTypeList, TenantID: PParm.TenantID})
		if err != nil {
			return nil, err
		}
		ids = append(ids, datas...)
	}
	//组更新
	if changed {
		//鉴权
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(ctx)
		err := update.Exec(Ctx)
		sysdebug.Printf("update group:%v err:%v", PParm.ID, err)
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				return nil, fmt.Errorf("%v%v", code.ZhCNText[code.AdminModifyDataError], "模板冲突!已有使用中同类型模板1个")
			}
			return nil, err
		}
		ids = append(ids, PParm.ID)
	}

	//如果是启用 [预警组]，需要更新人员类型表相关数据
	if PParm.Enabled != nil && PParm.GroupType == 1 {
		upTypeSql := fmt.Sprintf(`
			UPDATE ProfileType SET warning_enabled = %v WHERE id in ( SELECT profile_type_id FROM fm_alert_definition WHERE alert_group_id = '%v') 
		`, *PParm.Enabled, PParm.ID)
		rst, err := DBConn.Exec(upTypeSql)
		sysdebug.Printf("update profile type:%v err:%v", upTypeSql, err)
		if err != nil {
			return nil, err
		}
		idss, _ := rst.RowsAffected()
		sysdebug.Printf("update profile type success:%v", idss)
	}
	return ids, nil

}

// 查询某个预警组定义信息
func (oc *ORMCleint) GroupAlertsInfo(Ctx *gin.Context, ids int, groupName string) (*seclients.FmAlertGroup, error) {

	var (
		// client = databases.EntClient
		DBConn = databases.DBMysql
		// total   int
		baseSql = `
		SELECT a.id,
		BIN_TO_UUID(a.tenant_id) as tenant_id,   
		a.group_name,
		a.group_type,
		a.customization,
		a.description,
		a.enabled,
		a.creator,
		a.updated_time,
		a.created_time
		FROM  fm_alert_group a `
	)

	var result seclients.FmAlertGroup

	var Filter = " WHERE id > 0 "
	if ids != 0 {
		Filter += fmt.Sprintf(" AND a.id = '%v' ", ids)
	}
	if groupName != "" {
		Filter += fmt.Sprintf(" AND a.group_name = '%v' ", groupName)
	}

	//查询
	Limits := " ORDER BY ID DESC LIMIT 1"
	baseSql = baseSql + Filter + Limits
	err := DBConn.Get(&result, baseSql)
	sysdebug.Printf("query fm_alert_group:%v,  baseSql:%v total:%#v \n", result, baseSql, err)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

// 预警记录 列表查询 ids, device_id, alert_level, c, enabled, capture_log_id, tenantId, StartAt, EndAt
func (oc *ORMCleint) GroupAlertsList(Page, PageSize, ids, groupType int, Ctx *gin.Context, groupName, customization, tenantId, StartAt, EndAt string) (*seclients.FmAlertGroupPages, error) {

	var (
		// client = databases.EntClient
		DBConn  = databases.DBMysql
		total   int
		baseSql = `
		SELECT a.id,

		a.group_name,
		a.group_type,
		b.supplier,
		a.enabled,
		BIN_TO_UUID(a.tenant_id) as tenant_id,   
		a.customization,
		a.description,
		a.creator,
		a.updated_time,
		a.created_time
		FROM  fm_alert_group  a
		LEFT JOIN fm_tenant b ON b.tenant_id = a.tenant_id `
	)

	var result seclients.FmAlertGroupPages

	Filters := " WHERE a.id > 0 "
	if tenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(a.tenant_id)  = '%v' ", uuid.MustParse(tenantId))
	}
	if groupType != 0 {
		Filters += fmt.Sprintf(" AND a.group_type = '%v' ", groupType)
	}

	if groupName != "" {
		Filters += fmt.Sprintf(" AND a.group_name like '%v' ", "%"+groupName+"%")
	}

	if customization != "" {
		Filters += fmt.Sprintf(" AND a.customization = '%v' ", customization)
	}

	if ids != 0 {
		Filters += fmt.Sprintf(" AND a.id  = '%v' ", ids)
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		sysdebug.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		sysdebug.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)
		Filters += fmt.Sprintf(" AND (created_time between '%v' and '%v' )", StartAt, EndAt)

	}

	//总数
	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_alert_group`  a " + Filters
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
	sysdebug.Printf("query Alerts by tenantid:%v,  baseSql:%v total:%#v \n", tenantId, baseSql, err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_alert_group", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, nil

}

// 预警 信息数据列表查询 租户信息
func (oc *ORMCleint) GroupalertsListByTenantId(Page, PageSize int, tenantId string) (*seclients.AlertsList, error) {

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

// 查询某个预警等级定义信息
func (oc *ORMCleint) LevelAlertInfo(Ctx *gin.Context, conds []predicate.FmAlertDefinition) (*genclients.FmAlertDefinition, error) {

	var (
		client = databases.EntClient

		// DBConn = databases.DBMysql
		// total   int
		// baseSql = `
		// SELECT a.id,
		// a.alert_group_id,
		// a.customization,
		// a.description,
		// a.enabled,
		// a.creator,
		// a.updated_time,
		// a.created_time
		// FROM  fm_alert_definition a `  fm_alert_definition fm_alert_definition

	)

	results, err := client.Debug().FmAlertDefinition.Query().Where(conds...).Only(Ctx) // 这里必须加指定字段才能使用Scan

	sysdebug.Println("FmAlertDefinition retrieved", results, "err:", err)
	if err != nil {
		return nil, err
	}

	return results, nil

}

// 查询某个预警等级定义信息
func (oc *ORMCleint) LevelAlertList(Page, PageSize int, Ctx *gin.Context, conds []predicate.FmAlertDefinition) (*seclients.FmAlertLevelPages, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.FmAlertLevelPages

	total, err := client.Debug().FmAlertDefinition.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().FmAlertDefinition.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Select("id", "alert_group_id", "level", "profile_type_id", "action", "alarm_sound", "created_time", "updated_time"). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//.All(Ctx)
	sysdebug.Println("fm_alert_definition retrieved", zap.Int("count", len(result.Data)), zap.Int("total", total))
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_alert_definition", nil)

		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
		return &result, nil
	}
	return nil, nil
}

// 新增预警等级信息
func (oc *ORMCleint) LevelAlertNew(ctx *gin.Context, PParm *seclients.FmAlertLevel) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.FmAlertDefinition.Create()
	//预警等级 定义必与 预警组关联
	if PParm.Level == 0 || PParm.AlertGroupID == 0 {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}
	if PParm.AlertGroupID != 0 {
		builder.SetAlertGroupID(PParm.AlertGroupID)
	}
	if PParm.Level != 0 {
		builder.SetLevel(PParm.Level)
	}

	if PParm.ProfileTypeID != 0 {
		builder.SetProfileTypeID(PParm.ProfileTypeID)
	}
	if PParm.Action != "" {
		builder.SetAction(PParm.Action)
	}
	if PParm.AlarmSound != "" {
		builder.SetAlarmSound(PParm.AlarmSound)
	}

	// 执行
	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(ctx)
	Pthor, err := builder.Save(Ctx)
	sysdebug.Printf("builder new level:%v err:%#v \n", Pthor.ID, err)
	if err != nil {
		return 0, err
	}

	//更新人员类型关联信息
	PTypeUp := client.ProfileType.UpdateOneID(PParm.ProfileTypeID)
	PTypeUp.SetWarningLevel(PParm.Level)
	//默认不启用
	PTypeUp.SetWarningEnabled(false)
	err = PTypeUp.Exec(ctx)
	sysdebug.Printf("profile type 更新:%#v \v", err)
	if err != nil {
		return 0, err
	}
	return Pthor.ID, nil
}

// 新增后台管理 [预警等级] 列表新增 doing
func (mt *ORMCleint) LevelAlertNewList(c *gin.Context, MRRecord *seclients.ParamAlertList) ([]int64, error) {

	var (
		wg = sync.WaitGroup{}
	)
	if len(MRRecord.NewLevelList) == 0 {
		return nil, fmt.Errorf("预警等级为空")
	}

	//数据转换

	//批量写入 NamedExec 是基于 Prepare 和 Exec 的封装，不返回多条 ID
	newIds := []int64{}
	wg.Add(len(MRRecord.NewLevelList))
	for _, btnid := range MRRecord.NewLevelList {
		go func() {
			defer wg.Done()
			rms := seclients.FmAlertLevel{
				AlertGroupID:  MRRecord.ID,
				ProfileTypeID: btnid.ProfileTypeID,
				Level:         btnid.Level,
				Action:        btnid.Action,
				AlarmSound:    btnid.AlarmSound}

			newId, err := mt.LevelAlertNew(c, &rms)
			if err != nil {
				sysdebug.Printf("new alert level  :%#v err:%#v \n", rms, err)
				return
			}
			newIds = append(newIds, newId)

		}()
	}
	wg.Wait()

	if len(newIds) == 0 {
		return nil, fmt.Errorf("新增预警等级全部失败")
	}

	return newIds, nil
}

//批量修改

// 新增后台管理 [预警等级] 列表新增 doing
func (mt *ORMCleint) LevelAlertListUpdate(c *gin.Context, MRRecord *seclients.ParamAlertList) ([]int64, error) {

	var (
		wg = sync.WaitGroup{}
	)
	if len(MRRecord.NewLevelList) == 0 {
		return nil, fmt.Errorf("预警等级为空")
	}

	//数据转换
	//批量写入 NamedExec 是基于 Prepare 和 Exec 的封装，不返回多条 ID
	newIds := []int64{}
	wg.Add(len(MRRecord.NewLevelList))
	for _, btnid := range MRRecord.NewLevelList {
		go func() {
			defer wg.Done()
			if btnid.Level == 0 || (MRRecord.ID == 0 && btnid.AlertGroupID == 0) {
				sysdebug.Printf("new alert level  :%#v AlertGroupID:%#v \n", btnid.Level, btnid.AlertGroupID)
				return
			}
			if MRRecord.ID == 0 {
				MRRecord.ID = btnid.AlertGroupID
			}
			rms := seclients.FmAlertLevel{
				ID:            btnid.ID,
				AlertGroupID:  MRRecord.ID,
				ProfileTypeID: btnid.ProfileTypeID,
				DeleteFlag:    btnid.DeleteFlag,
				Level:         btnid.Level,
				Action:        btnid.Action,
				AlarmSound:    btnid.AlarmSound}

			newId, err := mt.LevelAlertUpdate(c, &rms)
			if err != nil {
				sysdebug.Printf("new alert level  :%#v err:%#v \n", rms, err)
				return
			}
			newIds = append(newIds, newId)

		}()
	}
	wg.Wait()

	if len(newIds) == 0 {
		return nil, fmt.Errorf("修改预警等级全部失败")
	}

	return newIds, nil
}

// 修改
func (oc *ORMCleint) LevelAlertUpdate(ctx *gin.Context, PParm *seclients.FmAlertLevel) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
		err    error
	)

	if PParm.ID == 0 {
		err = fmt.Errorf("%v:id:%v ", code.ZhCNText[code.ParamError], PParm.ID)
		return 0, err
	}

	sysdebug.Printf("pparm update data:%#v \n", PParm)
	//执行删除 优先
	if PParm.DeleteFlag == "1" {
		delSql := fmt.Sprintf("DELETE FROM fm_alert_definition WHERE id = '%v' ", PParm.ID)
		rst, err := DBConn.Exec(delSql)
		sysdebug.Printf("已删除:%v rst:%v \n", PParm.ID, err)
		if err != nil {
			return 0, err
		}
		return rst.RowsAffected()
	}

	updates := client.FmAlertDefinition.UpdateOneID(PParm.ID)
	// 执行
	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(ctx)

	//预警等级 定义必与 预警组关联
	if PParm.Level == 0 || PParm.AlertGroupID == 0 {
		err = fmt.Errorf("%v:level:%v ", code.ZhCNText[code.ParamError], PParm.Level)
		return 0, err
	}

	var changed bool
	if PParm.AlertGroupID != 0 {
		updates.SetAlertGroupID(PParm.AlertGroupID)
		changed = true
	}
	if PParm.Level != 0 {
		updates.SetLevel(PParm.Level)
		changed = true
	}

	if PParm.ProfileTypeID != 0 {
		updates.SetProfileTypeID(PParm.ProfileTypeID)
		changed = true
	}

	if PParm.AlarmSound != "" {
		updates.SetAlarmSound(PParm.AlarmSound)
		changed = true
	}
	if PParm.Action != "" {
		updates.SetAction(PParm.Action)
		changed = true
	}

	if changed {

		//执行更新
		err := updates.Exec(Ctx)
		sysdebug.Printf("updates:%v rst:%v \n", PParm.ID, err)

		if err != nil {
			return 0, err
		}

		//更新对应 人员类型 预警等级
		TypeUpdates := client.Debug().ProfileType.UpdateOneID(PParm.ProfileTypeID)
		TypeUpdates.SetWarningLevel(PParm.Level)
		err = TypeUpdates.Exec(ctx)
		sysdebug.Printf("update types:%v level:%v err:%v \n", PParm.ID, PParm.Level, err)
		if err != nil {
			err = fmt.Errorf("update profiletype:%v err:%v level:%v ", PParm.ProfileTypeID, err, PParm.Level)
			return 0, err
		}
		return PParm.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}
