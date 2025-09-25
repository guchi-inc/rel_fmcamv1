package sysclients

import (
	"fmcam/common/databases"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/fmalertdefinition"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 人员类型 新增
func (oc *ORMCleint) ProfileTypeNew(c *gin.Context, PParm *seclients.ProfileType) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
		err    error
	)

	builder := client.ProfileType.Create()

	var changed = false
	if PParm.TypeName == "" && PParm.TypeCode == nil {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.TypeName != "" {
		builder.SetTypeName(PParm.TypeName)
		changed = true
	}
	if PParm.WarningLevel != nil {
		builder.SetWarningLevel(*PParm.WarningLevel)
		changed = true
	}
	if PParm.WarningEnabled != nil {
		builder.SetWarningEnabled(*PParm.WarningEnabled)
		changed = true
	}
	if PParm.TypeCode != nil {
		builder.SetTypeCode(*PParm.TypeCode)
		changed = true
	} else {
		builder.SetTypeCode(utils.GenerateRandomString(8))
	}

	if PParm.Description != nil {
		builder.SetDescription(*PParm.Description)
		changed = true
	}
	if PParm.FaceValidityHours != nil {
		builder.SetFaceValidityHours(*PParm.FaceValidityHours)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive user:%v save new profile type:%#v", user, PParm)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)
		Pthor, err := builder.Save(Ctx)
		sysdebug.Printf("receive:%v params:%v", Pthor, err)
		if err != nil {
			return 0, err
		}
		if PParm.TenantID != "" {
			upSql := fmt.Sprintf(" update ProfileType set  tenant_id = %v WHERE id = '%v' ", fmt.Sprintf("(UUID_TO_BIN('%v'))", uuid.MustParse(PParm.TenantID)), Pthor.ID)
			rst, err := DBConn.Exec(upSql)
			sysdebug.Printf("%v rst:%v params:%v", upSql, rst, err)
			if err != nil {
				return 0, err
			}
		}

		return Pthor.ID, nil
	}
	return 0, fmt.Errorf("nothing to do")

}

// 新增后台管理 [人员类型] 列表 doing
func (mt *ORMCleint) InsertProfileTypeList(c *gin.Context, MRRecord *seclients.ParamProfileTypeLists) ([]int64, error) {

	var (
		wg = sync.WaitGroup{}
	)
	if len(MRRecord.NewTypeList) == 0 {
		return nil, fmt.Errorf("人员类型为空")
	}

	//数据转换
	//批量写入 NamedExec 是基于 Prepare 和 Exec 的封装，不返回多条 ID
	newIds := []int64{}
	wg.Add(len(MRRecord.NewTypeList))
	for _, btnid := range MRRecord.NewTypeList {
		go func() {
			defer wg.Done()
			if btnid.TypeCode == nil {
				newCode := utils.GenerateRandomString(8)
				btnid.TypeCode = &newCode
			}

			gpId := int(MRRecord.ID) //人员分组id
			glevelInfo, err := mt.GroupAlertsInfo(c, gpId, "")
			if err != nil {
				sysdebug.Printf("group info:%#v err:%#v \n", glevelInfo, err)
				return
			}

			rms := seclients.ProfileType{
				TypeName:          btnid.TypeName,
				TypeCode:          btnid.TypeCode,
				TenantID:          MRRecord.TenantID,
				FaceValidityHours: btnid.FaceValidityHours}

			newId, err := mt.ProfileTypeNew(c, &rms)
			if err != nil {
				sysdebug.Printf("new prfile type:%#v err:%#v \n", rms, err)
				return
			}

			//map 组和profileType映射表
			_, err = mt.GroupProfileTypeMapNew(c, &seclients.ProfileTypeAlertInfoMap{GroupID: MRRecord.ID,
				ProfileTypeID: newId, TenantID: &MRRecord.TenantID})
			sysdebug.Printf("new prfile group map:%#v type id:%#v err:%#v \n", MRRecord.ID, newId, err)

			if err != nil {
				sysdebug.Printf("new map type:%#v err:%#v \n", MRRecord.ID, err)
				return
			}

			newIds = append(newIds, newId)

		}()
	}
	wg.Wait()

	if len(newIds) == 0 {
		return nil, fmt.Errorf("新增类型全部失败")
	}

	return newIds, nil
}

// 人员类型 修改
func (oc *ORMCleint) ProfileTypeUpdate(c *gin.Context, PParm *seclients.ProfileType) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
		err    error
	)

	sysdebug.Printf("pparm ProfileTypeUpdate data:%#v \n", PParm)
	if PParm.DeleteFlag != nil && *PParm.DeleteFlag == "1" {
		baseSql := fmt.Sprintf("DELETE FROM ProfileType WHERE id = '%v'; ", PParm.ID)
		rst, err := DBConn.Exec(baseSql)
		if err != nil {
			return 0, err
		}
		return rst.RowsAffected()
	}

	update := client.ProfileType.UpdateOneID(PParm.ID)

	var changed = false
	if PParm.TypeName == "" && PParm.TypeCode == nil {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.TypeName != "" {
		update.SetTypeName(PParm.TypeName)
		changed = true
	}
	if PParm.TypeCode != nil {
		update.SetTypeCode(*PParm.TypeCode)
		changed = true
	}
	if PParm.WarningEnabled != nil {
		update.SetWarningEnabled(*PParm.WarningEnabled)
		changed = true
	}
	if PParm.WarningLevel != nil {
		update.SetWarningLevel(*PParm.WarningLevel)
		changed = true
	}
	if PParm.Description != nil {
		update.SetDescription(*PParm.Description)
		changed = true
	}
	if PParm.FaceValidityHours != nil {
		update.SetFaceValidityHours(*PParm.FaceValidityHours)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive user:%v save new profile type:%#v", user, PParm)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)
		err = update.Exec(Ctx)
		sysdebug.Printf("receive:%v params:%v", PParm.ID, err)
		if err != nil {
			return 0, err
		}
		// if PParm.TenantID != "" {
		// 	upSql := fmt.Sprintf(" update ProfileType set  tenant_id = %v WHERE id = '%v' ", fmt.Sprintf("(UUID_TO_BIN('%v'))", uuid.MustParse(PParm.TenantID)), PParm.ID)
		// 	rst, err := DBConn.Exec(upSql)
		// 	sysdebug.Printf("%v rst:%v params:%v", upSql, rst, err)
		// 	if err != nil {
		// 		return 0, err
		// 	}
		// }

		return PParm.ID, nil
	}
	return 0, fmt.Errorf("nothing to do")

}

// 新增后台管理 [人员类型] 列表 doing
func (mt *ORMCleint) ProfileTypeListUpdate(c *gin.Context, MRRecord *seclients.ParamProfileTypeLists) ([]int64, error) {

	var (
		wg = sync.WaitGroup{}
	)
	if len(MRRecord.NewTypeList) == 0 {
		return nil, fmt.Errorf("人员类型为空")
	}

	//数据转换
	//批量写入 NamedExec 是基于 Prepare 和 Exec 的封装，不返回多条 ID
	newIds := []int64{}
	wg.Add(len(MRRecord.NewTypeList))
	for _, btnid := range MRRecord.NewTypeList {
		go func() {
			defer wg.Done()
			gpId := int(MRRecord.ID)

			glevelInfo, err := mt.GroupAlertsInfo(c, gpId, "")
			if err != nil {
				sysdebug.Printf("group info:%#v err:%#v \n", glevelInfo, err)
				return
			}
			rms := seclients.ProfileType{
				ID:          btnid.ID,
				DeleteFlag:  btnid.DeleteFlag,
				Enabled:     btnid.Enabled,
				Description: btnid.Description,

				TypeName:          btnid.TypeName,
				TypeCode:          btnid.TypeCode,
				TenantID:          MRRecord.TenantID,
				FaceValidityHours: btnid.FaceValidityHours}

			if btnid.ID == 0 {
				newId, err := mt.ProfileTypeNew(c, &rms)
				if err != nil {
					sysdebug.Printf("new prfile type:%#v err:%#v \n", rms, err)
					return
				}

				newIds = append(newIds, newId)
				//关联表map中间表信息
				_, err = mt.GroupProfileTypeMapNew(c, &seclients.ProfileTypeAlertInfoMap{GroupID: MRRecord.ID,
					ProfileTypeID: newId, TenantID: &MRRecord.TenantID})
				sysdebug.Printf("new prfile group map:%#v type id:%#v err:%#v \n", MRRecord.ID, newId, err)
				if err != nil {
					return
				}
			} else {

				//等级信息
				var conds []predicate.FmAlertDefinition
				conds = append(conds, fmalertdefinition.IDEQ(btnid.ID))
				levelInfo, err := mt.LevelAlertInfo(c, conds)
				if err == nil || levelInfo != nil {
					rms.WarningLevel = &levelInfo.Level
					rms.WarningEnabled = glevelInfo.Enabled
				}
				//执行更新
				newId, err := mt.ProfileTypeUpdate(c, &rms)
				if err != nil {
					sysdebug.Printf("update prfile type:%#v err:%#v \n", rms, err)
					return
				}
				newIds = append(newIds, newId)
			}

		}()
	}
	wg.Wait()

	if len(newIds) == 0 {
		return nil, fmt.Errorf("修改类型全部失败")
	}

	return newIds, nil
}
