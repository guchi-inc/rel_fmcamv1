package sysclients

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 乡镇街道 更新
func (oc *ORMCleint) StreetUpdate(c *gin.Context, PParm *genclients.GovStreet) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	//鉴权
	user, err := UserService.GetUserByCtxToken(c)

	sysdebug.Printf("receive params:%v", user)
	if err != nil || user.LoginName == "" {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
		return 0, err
	}

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)
	originPm, err := client.GovStreet.Get(Ctx, PParm.ID)
	if err != nil {
		err = fmt.Errorf("%v:GovStreet:%v ", code.ZhCNText[code.AdminDetailError], PParm)
		return 0, err
	}

	update := client.GovStreet.UpdateOneID(PParm.ID)
	//	// Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != "" && PParm.DeleteFlag != originPm.DeleteFlag {
		update = update.SetDeleteFlag(PParm.DeleteFlag)
		err := update.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}

	var changed = false
	if PParm.Code != "" && PParm.Code != originPm.Code {
		update = update.SetCode(PParm.Code)
		changed = true
	}

	if PParm.Name != "" && PParm.Name != originPm.Name {
		update = update.SetName(PParm.Name)
		changed = true
	}

	if PParm.ProvinceCode != "" && PParm.ProvinceCode != originPm.ProvinceCode {
		update = update.SetProvinceCode(PParm.ProvinceCode)
		changed = true
	}
	if PParm.CityCode != "" && PParm.CityCode != originPm.CityCode {
		update = update.SetCityCode(PParm.CityCode)
		changed = true
	}
	if PParm.AreaCode != "" && PParm.AreaCode != originPm.AreaCode {
		update = update.SetAreaCode(PParm.AreaCode)
		changed = true
	}

	if PParm.Creator != "" && PParm.Creator != originPm.Creator {
		update = update.SetCreator(PParm.Creator)
		changed = true
	}

	PParm.Creator = user.LoginName
	if changed {
		//更新执行
		err = update.Exec(Ctx)
		if err != nil {
			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminModifyDataError]+err.Error())
		}
		return PParm.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}

// 乡镇街道  新增
func (oc *ORMCleint) StreetNew(c *gin.Context, PParm *genclients.GovStreet) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.GovStreet.Create()
	var changed = false
	if PParm.Name == "" && PParm.Code == "" {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.Name != "" {
		builder.SetName(PParm.Name)
		changed = true
	}
	if PParm.Code != "" {
		builder.SetCode(PParm.Code)
		changed = true
	}
	if PParm.ProvinceCode != "" {
		builder.SetProvinceCode(PParm.ProvinceCode)
		changed = true
	}
	if PParm.AreaCode != "" {
		builder.SetAreaCode(PParm.AreaCode)
		changed = true
	}
	if PParm.CityCode != "" {
		builder.SetCityCode(PParm.CityCode)
		changed = true
	}
	if PParm.Creator != "" {
		builder.SetCreator(PParm.Creator)
		changed = true
	}
	if PParm.DeleteFlag != "" {
		builder.SetDeleteFlag(PParm.DeleteFlag)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		t0 := time.Now().Local().In(configs.Loc)
		PParm.Creator = user.LoginName
		PParm.CreatedTime = &t0

		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)

		builder.SetCreator(PParm.Creator)
		builder.SetCreatedTime(*PParm.CreatedTime)
		Pthor, err := builder.Save(Ctx)
		// 执行
		if err != nil {
			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminCreateError]+err.Error())
		}
		return Pthor.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}

// 省 更新
func (oc *ORMCleint) ProvenceUpdate(c *gin.Context, PParm *genclients.Province) (int64, error) {

	var (
		client = databases.EntClient
	)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)
	originPm, err := client.Province.Get(Ctx, PParm.ID)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AdminModifyDataError], PParm)
		return 0, err
	}
	update := client.Province.UpdateOneID(PParm.ID)

	//	// Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != "" && PParm.DeleteFlag != originPm.DeleteFlag {
		update = update.SetDeleteFlag(PParm.DeleteFlag)
		err := update.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}

	var changed = false
	if PParm.Code != "" && PParm.Code != originPm.Code {
		update = update.SetCode(PParm.Code)
		changed = true
	}

	if PParm.Name != "" && PParm.Name != originPm.Name {
		update = update.SetName(PParm.Name)
		changed = true
	}
	if PParm.Creator != "" && PParm.Creator != originPm.Creator {
		update = update.SetCreator(PParm.Creator)
		changed = true
	}

	if changed {

		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}

		t0 := time.Now().Local().In(configs.Loc)
		PParm.Creator = user.LoginName
		PParm.CreatedTime = &t0

		//更新执行
		err = update.Exec(Ctx)
		if err != nil {
			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminModifyDataError]+err.Error())
		}
		return PParm.ID, nil

	}
	return 0, fmt.Errorf("nothing to do")
}

// 省  新增
func (oc *ORMCleint) ProvenceNew(c *gin.Context, PParm *genclients.Province) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.Province.Create()

	var changed = false
	if PParm.Name == "" && PParm.Code == "" {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.Name != "" {
		builder.SetName(PParm.Name)
		changed = true
	}
	if PParm.Code != "" {
		builder.SetCode(PParm.Code)
		changed = true
	}
	if PParm.Creator != "" {
		builder.SetCreator(PParm.Creator)
		changed = true
	}
	if PParm.DeleteFlag != "" {
		builder.SetDeleteFlag(PParm.DeleteFlag)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		t0 := time.Now().Local().In(configs.Loc)
		PParm.Creator = user.LoginName
		PParm.CreatedTime = &t0

		// 执行
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)

		builder.SetCreator(PParm.Creator)
		builder.SetCreatedTime(*PParm.CreatedTime)
		Pthor, err := builder.Save(Ctx)

		if err != nil {
			err = fmt.Errorf("%v  ", code.ZhCNText[code.AdminModifyDataError]+err.Error())
			return 0, err
		}
		return Pthor.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}

// 区县 更新
func (oc *ORMCleint) AreaUpdate(c *gin.Context, PParm *genclients.GovArea) (int64, error) {

	var (
		client = databases.EntClient
	)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)
	originPm, err := client.GovArea.Get(Ctx, PParm.ID)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AdminModifyDataError], PParm)
		return 0, err
	}
	update := client.GovArea.UpdateOneID(PParm.ID)

	//	// Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != "" && PParm.DeleteFlag != originPm.DeleteFlag {
		update = update.SetDeleteFlag(PParm.DeleteFlag)
		err := update.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}

	var changed = false
	if PParm.Code != "" && PParm.Code != originPm.Code {
		update = update.SetCode(PParm.Code)
		changed = true
	}

	if PParm.Name != "" && PParm.Name != originPm.Name {
		update = update.SetName(PParm.Name)
		changed = true
	}

	if PParm.ProvinceCode != "" && PParm.ProvinceCode != originPm.ProvinceCode {
		update = update.SetProvinceCode(PParm.ProvinceCode)
		changed = true
	}
	if PParm.CityCode != "" && PParm.CityCode != originPm.CityCode {
		update = update.SetCityCode(PParm.CityCode)
		changed = true
	}

	if PParm.Creator != "" && PParm.Creator != originPm.Creator {
		update = update.SetCreator(PParm.Creator)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		PParm.Creator = user.LoginName

		//更新执行
		//上下文传递操作人信息
		err = update.Exec(Ctx)
		if err != nil {

			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminModifyDataError]+err.Error())
		}
		return PParm.ID, nil
	}
	return 0, fmt.Errorf("nothing to do")
}

// 区县 新增
func (oc *ORMCleint) AreaNew(c *gin.Context, PParm *genclients.GovArea) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.GovArea.Create()
	var changed = false
	if PParm.Name == "" && PParm.Code == "" {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}
	if PParm.Name != "" {
		builder.SetName(PParm.Name)
		changed = true
	}
	if PParm.Code != "" {
		builder.SetCode(PParm.Code)
		changed = true
	}
	if PParm.ProvinceCode != "" {
		builder.SetProvinceCode(PParm.ProvinceCode)
		changed = true
	}
	if PParm.CityCode != "" {
		builder.SetCityCode(PParm.CityCode)
		changed = true
	}
	if PParm.Creator != "" {
		builder.SetCreator(PParm.Creator)
		changed = true
	}
	if PParm.DeleteFlag != "" {
		builder.SetDeleteFlag(PParm.DeleteFlag)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		t0 := time.Now().Local().In(configs.Loc)
		PParm.Creator = user.LoginName
		PParm.CreatedTime = &t0

		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)

		builder.SetCreator(PParm.Creator)
		builder.SetCreatedTime(*PParm.CreatedTime)
		Pthor, err := builder.Save(Ctx)
		// 执行
		if err != nil {
			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminCreateError]+err.Error())
		}
		return Pthor.ID, nil
	}
	return 0, fmt.Errorf("nothing to do")
}

// 城市 列表更新
func (oc *ORMCleint) CityUpdate(c *gin.Context, PParm *genclients.GovCity) (int64, error) {

	var (
		client = databases.EntClient
	)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)
	originPm, err := client.GovCity.Get(Ctx, PParm.ID)

	sysdebug.Printf("originPm:%#v receive params:%#v", originPm, PParm)
	if err != nil {
		return 0, err
	}

	update := client.GovCity.UpdateOneID(PParm.ID)
	//	// Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != "" && PParm.DeleteFlag != originPm.DeleteFlag {
		update = update.SetDeleteFlag(PParm.DeleteFlag)
		err := update.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}

	var changed = false
	if PParm.Code != "" && PParm.Code != originPm.Code {
		update = update.SetCode(PParm.Code)
		changed = true
	}

	if PParm.Name != "" && PParm.Name != originPm.Name {
		update = update.SetName(PParm.Name)
		changed = true
	}

	if PParm.ProvinceCode != "" && PParm.ProvinceCode != originPm.ProvinceCode {
		update = update.SetProvinceCode(PParm.ProvinceCode)
		changed = true
	}
	if PParm.Creator != "" && PParm.Creator != originPm.Creator {
		update = update.SetCreator(PParm.Creator)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive auth:%v", user)
		if err != nil || user.LoginName == "" {
			return 0, err
		}

		t0 := time.Now().Local().In(configs.Loc)
		PParm.Creator = user.LoginName
		PParm.CreatedTime = &t0

		//更新执行
		err = update.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}
	return 0, fmt.Errorf("nothing to do")

}

// 城市 新增
func (oc *ORMCleint) CityNew(c *gin.Context, PParm *genclients.GovCity) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.GovCity.Create()

	var changed = false
	if PParm.Name == "" && PParm.Code == "" {
		err = fmt.Errorf("%v:name code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.Name != "" {
		builder.SetName(PParm.Name)
		changed = true
	}
	if PParm.Code != "" {
		builder.SetCode(PParm.Code)
		changed = true
	}
	if PParm.ProvinceCode != "" {
		builder.SetProvinceCode(PParm.ProvinceCode)
		changed = true
	}
	if PParm.Creator != "" {
		builder.SetCreator(PParm.Creator)
		changed = true
	}
	if PParm.DeleteFlag != "" {
		builder.SetDeleteFlag(PParm.DeleteFlag)
		changed = true
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		t0 := time.Now().Local().In(configs.Loc)
		PParm.Creator = user.LoginName
		PParm.CreatedTime = &t0

		// 执行
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)

		builder.SetCreator(PParm.Creator)
		builder.SetCreatedTime(*PParm.CreatedTime)
		Pthor, err := builder.Save(Ctx)

		if err != nil {
			return 0, err
		}
		return Pthor.ID, nil
	}
	return 0, fmt.Errorf("nothing to do")

}

// 预警信息 更新
func (oc *ORMCleint) AlertUpdate(c *gin.Context, PParm *seclients.Alerts) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)
	originPm, err := client.Alerts.Get(Ctx, PParm.ID)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AdminModifyDataError], PParm)
		return 0, err
	}

	update := client.Alerts.UpdateOneID(PParm.ID)
	var changed = false

	//根据人员类型 查预警等级
	// var condDefs = []predicate.FmAlertDefinition{}
	// condDefs = append(condDefs, fmalertdefinition.ProfileTypeIDEQ(*PParm.ProfileTypeId))
	// alertDef, err := oc.LevelAlertInfo(c, condDefs)
	// sysdebug.Printf("receive alertDef:%v \n err %#v \n", alertDef, err)
	// if err == nil {
	// 	if alertDef != nil && alertDef.Level != 0 && alertDef.Level != int(originPm.AlertLevel) {
	// 		update = update.SetAlertLevel(int8(alertDef.Level))
	// 		changed = true
	// 	}
	// }

	//更新人员表的类型  为用户指定类型
	bsSql := fmt.Sprintf(`
		update Profile set type_id = '%v' WHERE id = (select matched_profile_id FROM CaptureLogs where id = '%v') 
	`, *PParm.ProfileTypeId, originPm.CaptureLogID)

	_, terr := DBConn.Exec(bsSql)
	sysdebug.Printf("update profile info bsSql:%v \n err %#v \n", bsSql, terr)
	if terr != nil {
		return 0, fmt.Errorf("人员信息修改失败" + terr.Error())
	} else {
		changed = true
	}

	update = update.SetStatus(PParm.Status)
	changed = true
	//最后更新预警记录 处理状态
	if changed {

		//鉴权
		user, err := UserService.GetUserByCtxToken(c)

		sysdebug.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		t0 := time.Now().Local().In(configs.Loc)

		//已处理
		update = update.SetFmUserID(user.Id)
		update = update.SetHandledTime(t0)

		//更新执行
		err = update.Exec(Ctx)
		sysdebug.Printf("update alert :%v ProfileTypeId:%v \n err %#v \n", PParm.ID, PParm.ProfileTypeId, err)
		if err != nil {
			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminModifyDataError]+err.Error())
		}

		return PParm.ID, nil

	}
	return 0, fmt.Errorf("nothing to do")
}
