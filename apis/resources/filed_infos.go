package resources

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/fieldmetadata"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 字段状态查询
func (that *LimitsApi) GetFieldsState(c *gin.Context) {

	var (
		client = databases.EntClient
		Ctx    = configs.Ctx
	)

	Name := c.DefaultQuery("name", "")
	Cname := c.DefaultQuery("cname", "")
	TableName := c.DefaultQuery("table_name", "")
	DataType := c.DefaultQuery("data_type", "")
	Description := c.DefaultQuery("description", "")
	DefaultValue := c.DefaultQuery("default_value", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if TableName == "" {
		helpers.JSONs(c, code.ParamError, fmt.Errorf("table_name required"))
		return
	}
	fmt.Printf("reveice Name:%v Codes:%v TableName:%v, DataType:%v Description:%v, DefaultValue:%v, PageSize:%v, Page:%v StartAt:%v EndAt:%v \n",
		Name, Cname, TableName, DataType, Description, DefaultValue, PageSize, Page, StartAt, EndAt)

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 10 {
		PageSize = 10
	}

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	offset := (Page - 1) * PageSize

	var conds = []predicate.FieldMetadata{}

	if Name != "" {
		conds = append(conds, fieldmetadata.NameContains(Name))
	}
	if Cname != "" {
		conds = append(conds, fieldmetadata.CnameContains(Cname))
	}
	if TableName != "" {
		conds = append(conds, fieldmetadata.TableName(TableName))
	}
	if DataType != "" {
		conds = append(conds, fieldmetadata.DataTypeContains(DataType))
	}
	if Description != "" {
		conds = append(conds, fieldmetadata.DefaultValueContains(Description))
	}
	if DefaultValue != "" {
		conds = append(conds, fieldmetadata.DefaultValueContains(DefaultValue))
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apilog.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, fieldmetadata.CreatedTimeGT(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apilog.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, fieldmetadata.CreatedTimeLT(EndAtDate))
	}

	list, err := client.Debug().FieldMetadata.Query().
		Where(conds...).
		Limit(PageSize).
		Offset(offset).
		All(Ctx)

	apilog.Printf("ent conds:%#v to list:%#v err:%#v \n", conds, list, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := client.Debug().FieldMetadata.Query().Where(conds...).Count(Ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	apilog.Println("FieldMetadata retrieved", zap.Int("count", len(list)), zap.Int("total", total))

	if len(list) > 0 {
		cm, _ := databases.FmGlobalMap(TableName, nil)
		helpers.JSONs(c, code.Success, seclients.FieldMetadataPage{
			Data:     list,
			Total:    total,
			Page:     Page,
			PageSize: PageSize,
			Columns:  cm,
		})
	} else {
		helpers.JSONs(c, code.Success, seclients.AreaPage{})
	}
}

// 字段状态更新
func (that *LimitsApi) UpdateFieldMetadata(c *gin.Context) {

	var (
		client = databases.EntClient
		PParm  = genclients.FieldMetadata{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v::%v::%#v ", code.ZhCNText[code.ParamBindError], PParm, err)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	//上下文传递操作人信息
	Ctx := databases.WithOperatorFromGin(c)

	apilog.Printf("gin ctx creator:%v \n", Ctx)
	originPm, err := client.FieldMetadata.Get(Ctx, PParm.ID)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AdminModifyDataError], PParm)
		helpers.JSONs(c, code.AdminModifyDataError, err)
		return
	}
	update := client.FieldMetadata.UpdateOneID(PParm.ID)

	var changed = false
	if PParm.Cname != "" && PParm.Cname != originPm.Cname {
		update = update.SetCname(PParm.Cname)
		changed = true
	}

	if PParm.DataType != "" && PParm.DataType != originPm.DataType {
		update = update.SetDataType(PParm.DataType)
		changed = true
	}

	if PParm.IsVisible != originPm.IsVisible {
		update = update.SetIsVisible(PParm.IsVisible)
		changed = true
	}
	if PParm.IsEditable != originPm.IsEditable {
		update = update.SetIsEditable(PParm.IsEditable)
		changed = true
	}
	if PParm.IsSearchable != originPm.IsSearchable {
		update = update.SetIsSearchable(PParm.IsSearchable)
		changed = true
	}
	if PParm.IsRequired != originPm.IsRequired {
		update = update.SetIsRequired(PParm.IsRequired)
		changed = true
	}
	if PParm.DefaultValue != originPm.DefaultValue {
		update = update.SetDefaultValue(PParm.DefaultValue)
		changed = true
	}
	if PParm.Description != originPm.Description {
		update = update.SetDescription(PParm.Description)
		changed = true
	}

	if changed {
		//鉴权 TODO
		user, err := UserService.GetUserByCtxToken(c)

		apilog.Printf("receive params:%v", user)
		if err != nil || user.LoginName == "" {
			err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.AuthorizationError], user)
			helpers.JSONs(c, code.AuthorizationError, err)
			return
		}
		t0 := time.Now().Local().In(configs.Loc)
		// PParm.CreatedAt = user.LoginName
		PParm.UpdatedTime = &t0

		//更新执行
		err = update.Exec(Ctx)
		if err != nil {
			helpers.JSONs(c, code.AdminModifyDataError, code.ZhCNText[code.AdminModifyDataError]+err.Error())
			return
		}
		helpers.JSONs(c, code.Success, gin.H{"data": PParm.ID, "message": "success"})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": PParm.ID, "message": "nothing todo"})

}
