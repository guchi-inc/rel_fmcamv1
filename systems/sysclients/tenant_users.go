package sysclients

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/models/code"
	"fmcam/models/schema"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ORMCleint struct{}

func (oc *ORMCleint) TenantList(Page, PageSize int, ctx *gin.Context, conds []predicate.Tenants) (*seclients.TenantsPage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.TenantsPage

	offset := (Page - 1) * PageSize
	total, err := client.Debug().Tenants.Query().Where(conds...).Count(ctx)
	if err != nil {
		return nil, err
	}

	sysdebug.Println("tenants retrieved", zap.Int("count", len(result.Data)), zap.Int("total", total))
	err = client.Debug().Tenants.
		Query().
		Where(conds...).
		Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select(`id,
		supplier,
		contacts,
		email,
		BIN_TO_UUID(tenant_id) as tenant_id,
		type,
		province,
		city,
		area,
		street,
		address,
		addr_code,
		full_address,
		fax,
		phone_num,
		telephone,
		tax_num,
		bank_name,
		account_number,
		sort,
		description,
		enabled,
		delete_flag,
		isystem,
		advance_in,
		begin_need_get,
		begin_need_pay,
		all_need_get,
		all_need_pay,
		tax_rate,
		creator,
		created_time,
		updated_time`).
		Scan(ctx, &result.Data)

	sysdebug.Println("tenants data", err, zap.Int("count", len(result.Data)), zap.Int("total", total))
	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_tenant", nil)
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
		result.Columns = cm
	}
	return &result, err

}

// 更新租户信息
func (oc *ORMCleint) TenantInfoUpdate(ctx *gin.Context, PParm *seclients.TenantFull) (int64, error) {
	var (
		client = databases.EntClient
	)

	update := client.Tenants.UpdateOneID(PParm.ID)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(ctx)
	originPm, err := client.Tenants.Get(Ctx, PParm.ID)
	sysdebug.Printf("get tenant info:%#v err:%#v \n", originPm, err)
	if err != nil {
		logFile.Infof("获取租户信息失败:%#v \n", PParm)
		return 0, err
	}

	//  Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != "" && PParm.DeleteFlag != originPm.DeleteFlag {
		update = update.SetDeleteFlag(PParm.DeleteFlag)
		err := update.Exec(Ctx)
		sysdebug.Printf("update tenant info:%#v err:%#v \n", update, err)
		if err != nil {
			logFile.Infof("删除标记更新失败:%#v \n", PParm)

			return 0, err
		}
		return PParm.ID, nil
	}

	var changed = false
	//基本信息
	if PParm.Supplier != "" && PParm.Supplier != originPm.Supplier {
		update = update.SetSupplier(PParm.Supplier)
		changed = true
	}
	if PParm.Contacts != "" && PParm.Contacts != originPm.Contacts {
		update = update.SetContacts(PParm.Contacts)
		changed = true
	}
	if PParm.Email != "" && PParm.Email != originPm.Email {
		update = update.SetEmail(PParm.Email)
		changed = true
	}
	if PParm.Telephone != "" && PParm.Telephone != originPm.Telephone {
		update = update.SetTelephone(PParm.Telephone)
		changed = true
	}
	if PParm.PhoneNum != "" && PParm.PhoneNum != originPm.PhoneNum {
		update = update.SetPhoneNum(PParm.PhoneNum)
		changed = true
	}

	if PParm.Province != "" && PParm.Province != originPm.Province {
		update = update.SetProvince(PParm.Province)
		changed = true
	}
	if PParm.City != "" && PParm.City != originPm.City {
		update = update.SetCity(PParm.City)
		changed = true
	}
	if PParm.Area != "" && PParm.Area != originPm.Area {
		update = update.SetArea(PParm.Area)
		changed = true
	}
	if PParm.Street != "" && PParm.Street != originPm.Street {
		update = update.SetStreet(PParm.Street)
		changed = true
	}
	if PParm.Address != "" && PParm.Address != originPm.Address {
		update = update.SetAddress(PParm.Address)
		changed = true
	}
	if PParm.AddrCode != "" && PParm.AddrCode != originPm.AddrCode {
		update = update.SetAddrCode(PParm.AddrCode)
		changed = true
	}
	if PParm.Creator != "" && PParm.Creator != originPm.Creator {
		update = update.SetCreator(PParm.Creator)
		changed = true
	}
	if PParm.Description != "" && PParm.Description != originPm.Description {
		update = update.SetDescription(PParm.Description)
		changed = true
	}
	//财务信息
	if PParm.Fax != "" && PParm.Fax != originPm.Fax {
		update = update.SetFax(PParm.Fax)
		changed = true
	}
	if PParm.TaxNum != "" && PParm.TaxNum != originPm.TaxNum {
		update = update.SetTaxNum(PParm.TaxNum)
		changed = true
	}
	if PParm.BankName != "" && PParm.BankName != originPm.BankName {
		update = update.SetBankName(PParm.BankName)
		changed = true
	}
	if PParm.AccountNumber != "" && PParm.AccountNumber != originPm.AccountNumber {
		update = update.SetAccountNumber(PParm.AccountNumber)
		changed = true
	}
	if PParm.Sort != "" && PParm.Sort != originPm.Sort {
		update = update.SetSort(PParm.Sort)
		changed = true
	}

	if PParm.TaxRate != "" {
		adFDec, err := schema.NewDecimalFromString(PParm.TaxRate)
		if err == nil {
			if originPm.TaxRate != adFDec {
				update = update.SetTaxRate(adFDec)
				changed = true
			}
		}
	}

	if PParm.AdvanceIn != "" {
		adFDec, err := schema.NewDecimalFromString(PParm.AdvanceIn)
		if err == nil {
			if originPm.AdvanceIn != adFDec {
				update = update.SetAdvanceIn(adFDec)
				changed = true
			}
		}
	}

	if PParm.BeginNeedGet != "" {
		adFDec, err := schema.NewDecimalFromString(PParm.BeginNeedGet)
		if err == nil {
			if originPm.BeginNeedGet != adFDec {
				update = update.SetBeginNeedGet(adFDec)
				changed = true
			}
		}
	}
	if PParm.BeginNeedPay != "" {
		adFDec, err := schema.NewDecimalFromString(PParm.BeginNeedPay)
		if err == nil {
			if originPm.BeginNeedPay != adFDec {
				update = update.SetBeginNeedPay(adFDec)
				changed = true
			}
		}
	}

	if PParm.AllNeedGet != "" {
		adFDec, err := schema.NewDecimalFromString(PParm.AllNeedGet)
		if err == nil {
			if originPm.AllNeedGet != adFDec {
				update = update.SetAllNeedGet(adFDec)
				changed = true
			}
		}
	}

	if PParm.AllNeedPay != "" {
		adFDec, err := schema.NewDecimalFromString(PParm.AllNeedPay)
		if err == nil {
			if originPm.AllNeedPay != adFDec {
				update = update.SetAllNeedPay(adFDec)
				changed = true
			}
		}
	}

	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(ctx)
		sysdebug.Printf("receive user:%v save new profile type:%#v", user, PParm)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(ctx)

		//更新执行
		err = update.Exec(Ctx)
		sysdebug.Printf("update:%#v  err:%#v ", update, err)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}
	return PParm.ID, nil

}

// pms api 列表查询
func (oc *ORMCleint) FMPMSAPIList(Page, PageSize int, Ctx *gin.Context, conds []predicate.FMPMSApi) (*seclients.FMPMSApiPage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.FMPMSApiPage
	offset := (Page - 1) * PageSize

	total, err := client.Debug().FMPMSApi.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().FMPMSApi.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select("id", "pms_name", "contact", "phonenum", "enabled", "pms_api", "creator", "created_time", "updated_time"). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//All(Ctx)
	list := result.Data
	sysdebug.Println("fm_pms_apis retrieved", zap.Int("count", len(list)), zap.Int("total", total))
	sysdebug.Printf("err:%#v \n", err)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap("fm_pms_apis", nil)
		result.Data = list
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 更新FMPMSAPI 信息
func (oc *ORMCleint) FMPMSApiUpdate(ctx *gin.Context, PParm *seclients.FMPMSApi) (int64, error) {
	var (
		client = databases.EntClient
	)

	update := client.FMPMSApi.UpdateOneID(PParm.ID)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(ctx)
	originPm, err := client.FMPMSApi.Get(Ctx, PParm.ID)
	sysdebug.Printf("get FMPMSApi info:%#v err:%#v \n", originPm, err)
	if err != nil {
		logFile.Infof("获取租户信息失败:%#v \n", PParm)
		return 0, err
	}

	//状态 enabled 更新直接返回
	if PParm.Enabled != nil && *PParm.Enabled != *originPm.Enabled {
		update = update.SetEnabled(*PParm.Enabled)
		err := update.Exec(Ctx)
		sysdebug.Printf("update FMPMSApi info:%#v err:%#v \n", update, err)
		if err != nil {
			logFile.Infof("更新状态失败:%#v \n", PParm.ID)

			return 0, err
		}
		return PParm.ID, nil
	}

	//  Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != nil && *PParm.DeleteFlag != *originPm.DeleteFlag {
		update = update.SetDeleteFlag(*PParm.DeleteFlag)
		err := update.Exec(Ctx)
		sysdebug.Printf("update FMPMSApi info:%#v err:%#v \n", update, err)
		if err != nil {
			logFile.Infof("删除标记更新失败:%#v \n", PParm.ID)
			return 0, err
		}
		return PParm.ID, nil
	}

	var changed = false
	//基本信息
	if PParm.PmsName != "" && PParm.PmsName != originPm.PmsName {
		update = update.SetPmsName(PParm.PmsName)
		changed = true
	}
	if PParm.Contact != "" && PParm.Contact != originPm.Contact {
		update = update.SetContact(PParm.Contact)
		changed = true
	}
	if PParm.Phonenum != nil && *PParm.Phonenum != *originPm.Phonenum {
		update = update.SetPhonenum(*PParm.Phonenum)
		changed = true
	}
	if PParm.PmsAPI != "" && PParm.PmsAPI != originPm.PmsAPI {
		update = update.SetPmsAPI(PParm.PmsAPI)
		changed = true
	}
	if changed {
		//鉴权
		user, err := UserService.GetUserByCtxToken(ctx)
		sysdebug.Printf("receive user:%v save new profile type:%#v", user, PParm)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			return 0, err
		}
		update = update.SetCreator(user.LoginName)
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(ctx)
		//更新执行
		err = update.Exec(Ctx)
		sysdebug.Printf("update:%#v  err:%#v ", update, err)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}
	return PParm.ID, nil

}

// 新增 FMPMSAPI 信息
func (oc *ORMCleint) FMPMSApiNew(c *gin.Context, PParm *seclients.FMPMSApi) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.FMPMSApi.Create()

	var changed = false
	if PParm.PmsName == "" && PParm.PmsAPI == "" {
		err = fmt.Errorf("%v:name api:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.PmsName != "" {
		builder.SetPmsName(PParm.PmsName)
		changed = true
	}
	if PParm.PmsAPI != "" {
		builder.SetPmsAPI(PParm.PmsAPI)
		changed = true
	}

	if PParm.Contact != "" {
		builder.SetContact(PParm.Contact)
		changed = true
	}
	if PParm.Phonenum != nil {
		builder.SetPhonenum(*PParm.Phonenum)
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
		PParm.Creator = &user.LoginName
		if PParm.Creator != nil {
			builder.SetCreator(*PParm.Creator)
			changed = true
		}

		to := time.Now().Local().In(configs.Loc)
		builder.SetCreatedTime(to)
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)
		Pthor, err := builder.Save(Ctx)
		if err != nil {
			err = fmt.Errorf("%v  ", code.ZhCNText[code.AdminModifyDataError]+err.Error())
			return 0, err
		}
		return Pthor.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}
