package sysclients

import (
	"fmcam/common/databases"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 酒店 需求留言记录 查询
func (oc *ORMCleint) Demandsist(Page, PageSize int, Ctx *gin.Context, conds []predicate.FmDemands) (*seclients.DemandsPage, error) {

	var (
		client = databases.EntClient
	)

	var result seclients.DemandsPage // []genclients.FmDedicatedServices
	offset := (Page - 1) * PageSize

	total, err := client.Debug().FmDemands.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().FmDemands.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select("id", "supplier", "username", "phonenum", "email", "province", "city", "area", "street", "message", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	sysdebug.Println("FmDemands retrieved", zap.Int("data", total))
	sysdebug.Printf("err:%#v \n", err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_demands", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err
}

// 需求留言 新增
func (oc *ORMCleint) DemandsNew(c *gin.Context, PParm *genclients.FmDemands) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.FmDemands.Create()

	var changed = false
	if PParm.Supplier == nil && PParm.Username == nil {
		err = fmt.Errorf("%v:店名 姓名不能为空:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.Province != nil {
		builder.SetProvince(*PParm.Province)
		changed = true
	}
	if PParm.City != nil {
		builder.SetCity(*PParm.City)
		changed = true
	}
	if PParm.Area != nil {
		builder.SetArea(*PParm.Area)
		changed = true
	}
	if PParm.Street != nil {
		builder.SetStreet(*PParm.Street)
		changed = true
	}
	if PParm.Phonenum != nil {
		builder.SetPhonenum(*PParm.Phonenum)
		changed = true
	}
	if PParm.Supplier != nil {
		builder.SetSupplier(*PParm.Supplier)
		changed = true
	}
	if PParm.Email != nil {
		builder.SetEmail(*PParm.Email)
		changed = true
	}
	if PParm.Message != nil {
		builder.SetMessage(*PParm.Message)
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

		// 执行
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)
		builder.SetCreator(*PParm.Creator)
		Pthor, err := builder.Save(Ctx)
		if err != nil {
			err = fmt.Errorf("%v  ", code.ZhCNText[code.AdminModifyDataError]+err.Error())
			return 0, err
		}
		return Pthor.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}

// 客服 人员列表查询
func (oc *ORMCleint) DedicateServiceList(Page, PageSize int, Ctx *gin.Context, conds []predicate.FmDedicatedServices, TenantId string) (*seclients.DedicatedServicePage, error) {

	var (
		client = databases.EntClient
		DBconn = databases.DBMysql
	)

	var result seclients.DedicatedServicePage // []genclients.FmDedicatedServices
	offset := (Page - 1) * PageSize

	if TenantId != "" {
		var OwnTenatInfo = genclients.FmDedicatedServices{}
		baseSql := fmt.Sprintf(`SELECT a.id, a.contacts , a.work_id , 
		a.supplier, a.email, a.fax, a.phonenum, a.description, a.creator, a.created_time FROM fm_dedicated_services a
		LEFT JOIN fm_tenant b ON b.id = a.work_id AND BIN_TO_UUID(tenant_id) = '%v' ORDER BY a.id DESC limit 1 `, uuid.MustParse(TenantId))
		err := DBconn.Get(&OwnTenatInfo, baseSql)
		if err != nil {
			return nil, err
		}
		result.Data = []genclients.FmDedicatedServices{OwnTenatInfo}
		return &result, nil
	}

	total, err := client.Debug().FmDedicatedServices.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	err = client.Debug().FmDedicatedServices.Query().
		Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		Select("id", "contacts", "work_id", "supplier", "email", "fax", "phonenum", "description", "creator", `created_time`). // 这里必须加指定字段才能使用Scan
		Scan(Ctx, &result.Data)

	//All(Ctx)

	sysdebug.Println("FmDedicatedServices retrieved", zap.Int("data", total))
	sysdebug.Printf("err:%#v \n", err)
	if err != nil {
		return nil, err
	}

	if len(result.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_dedicated_services", nil)
		result.Columns = cm
		result.Total = total
		result.Page = Page
		result.PageSize = PageSize
	}
	return &result, err

}

// 客服人员 更新
func (oc *ORMCleint) DedicateUpdate(c *gin.Context, PParm *seclients.DedicateServiceInfo) (int64, error) {

	var (
		client = databases.EntClient
	)

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)

	//	// Delete 示例 删除标签更新 直接返回
	if PParm.DeleteFlag != nil && *PParm.DeleteFlag == "1" {
		delOne := client.FmDedicatedServices.DeleteOneID(PParm.ID)
		err := delOne.Exec(Ctx)
		if err != nil {
			return 0, err
		}
		return PParm.ID, nil
	}

	originPm, err := client.FmDedicatedServices.Get(Ctx, PParm.ID)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AdminModifyDataError], PParm)
		return 0, err
	}
	update := client.FmDedicatedServices.UpdateOneID(PParm.ID)

	var changed = false
	if PParm.Contacts != nil && *PParm.Contacts != *originPm.Contacts {
		update = update.SetContacts(*PParm.Contacts)
		changed = true
	}

	if PParm.Description != nil && *PParm.Description != *originPm.Description {
		update = update.SetDescription(*PParm.Description)
		changed = true
	}
	if PParm.Email != nil && *PParm.Email != *originPm.Email {
		update = update.SetEmail(*PParm.Email)
		changed = true
	}
	if PParm.Phonenum != nil && *PParm.Phonenum != *originPm.Phonenum {
		update = update.SetPhonenum(*PParm.Phonenum)
		changed = true
	}
	if PParm.WorkID != nil && *PParm.WorkID != *originPm.WorkID {
		update = update.SetWorkID(*PParm.WorkID)
		changed = true
	}
	if PParm.Supply != nil && *PParm.Supply != *originPm.Supplier {
		update = update.SetSupplier(*PParm.Supply)
		changed = true
	}
	if PParm.Fax != nil && *PParm.Fax != *originPm.Fax {
		update = update.SetFax(*PParm.Fax)
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

		//更新执行
		err = update.Exec(Ctx)
		if err != nil {
			return 0, fmt.Errorf("%v", code.ZhCNText[code.AdminModifyDataError]+err.Error())
		}
		return PParm.ID, nil

	}
	return 0, fmt.Errorf("nothing to do")
}

// 客服人员 新增
func (oc *ORMCleint) DedicateNew(c *gin.Context, PParm *genclients.FmDedicatedServices) (int64, error) {

	var (
		client = databases.EntClient
		err    error
	)

	builder := client.FmDedicatedServices.Create()

	var changed = false
	if PParm.WorkID == nil && PParm.Contacts == nil {
		err = fmt.Errorf("%v:work id contacts code:%v ", code.ZhCNText[code.ParamError], PParm)
		return 0, err
	}

	if PParm.WorkID != nil {
		builder.SetWorkID(*PParm.WorkID)
		changed = true
	}
	if PParm.Contacts != nil {
		builder.SetContacts(*PParm.Contacts)
		changed = true
	}
	if PParm.Description != nil {
		builder.SetDescription(*PParm.Description)
		changed = true
	}
	if PParm.Email != nil {
		builder.SetEmail(*PParm.Email)
		changed = true
	}
	if PParm.Fax != nil {
		builder.SetFax(*PParm.Fax)
		changed = true
	}
	if PParm.Phonenum != nil {
		builder.SetPhonenum(*PParm.Phonenum)
		changed = true
	}
	if PParm.Supplier != nil {
		builder.SetSupplier(*PParm.Supplier)
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

		// 执行
		//上下文传递操作人信息
		Ctx := databases.WithOperatorFromGin(c)
		builder.SetCreator(*PParm.Creator)
		Pthor, err := builder.Save(Ctx)
		if err != nil {
			err = fmt.Errorf("%v  ", code.ZhCNText[code.AdminModifyDataError]+err.Error())
			return 0, err
		}
		return Pthor.ID, nil
	}

	return 0, fmt.Errorf("nothing to do")
}
