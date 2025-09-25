package alysis

import (
	"fmcam/common/configs"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/sysaly"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DBs1Datas(dataType, tenant_id, StartAt, EndAt, SubType, RangeType string, Device *string, twoModel *bool) (*sysaly.SysClientList[any], error) {

	var (
		userCategory = []string{} //"陌生人", "住客", "上门服务", "访客", "工作人员", "其他白名单人员", "黑名单人员"}
	)

	//查询全部
	if tenant_id == "" {
		return nil, fmt.Errorf("不支持的租户号:%v", tenant_id)
	}

	// 人员类型
	CategoryList, _ := UserUtil.UserTypeList(1, 50, tenant_id)
	if CategoryList == nil {
		CategoryList, _ := UserUtil.UserTypeList(1, 50, configs.DefTenantId)
		logapi.Println("default tenant_id: ", tenant_id, "\n default CategoryList:", CategoryList != nil)
		for _, cate := range CategoryList.Data {
			userCategory = append(userCategory, cate.TypeName)
		}
	} else {
		for _, cate := range CategoryList.Data {
			userCategory = append(userCategory, cate.TypeName)
		}
	}

	if !dataTypes[dataType] {
		return nil, fmt.Errorf("不支持的查询类型:%v", dataType)
	}

	switch dataType {
	case "capture":

		clientSys, err := AlysysOperator.GetCaptureLogsWeek(1, 7, tenant_id, StartAt, EndAt, SubType, RangeType, Device, twoModel)
		logapi.Printf("dates capture nil?:%#v err:%v \n", clientSys != nil, err)
		if err != nil {
			return nil, err
		}
		//用户统计
		logapi.Println("\n CategoryList:", CategoryList)
		return clientSys, nil
	case "warn":

		//今日统计  累计  待处理
		clientSys, err := AlysysOperator.TotalAlertsPostion(1, 7, tenant_id, StartAt, EndAt, SubType, RangeType, Device, twoModel)
		logapi.Println("Alerts weeks tenant_id: ", tenant_id, "\n Alerts:")
		logapi.Printf("dates warn nil?:%#v err:%v \n", clientSys != nil, err)
		if err != nil {
			return nil, err
		}
		//用户统计
		logapi.Println("\n Alerts:", clientSys)
		return clientSys, nil

	case "gather":

		//最近采集 信息统计
		clientSys, err := AlysysOperator.GetGatherLogWeek(1, 7, tenant_id, StartAt, EndAt)
		logapi.Println("weeks tenant_id: ", tenant_id, "\n CategoryList:", CategoryList != nil)
		logapi.Printf("dates gather nil?:%#v err:%v \n", clientSys != nil, err)
		if err != nil {
			return nil, err
		}
		//用户统计
		logapi.Println("\n CategoryList:", CategoryList)
		return clientSys, nil
	case "accrued_gather":

		clientSys, err := AlysysOperator.GetGatherLogAccrued(1, 7, tenant_id, StartAt, EndAt, SubType, RangeType, Device, twoModel)
		logapi.Printf("dates accrued_gather nil?:%#v err:%v \n", clientSys != nil, err)
		if err != nil {
			return nil, err
		}
		//用户统计
		logapi.Println("\n accrued_gather CategoryList:", CategoryList)
		return clientSys, nil
	case "warning": //近期预警

		// 筛选租户号
		clientSys, err := AlysysOperator.GetAlertsWeek(1, 7, tenant_id)
		logapi.Println("weeks GetAlertsWeek: ", tenant_id, "\n CategoryList:")
		logapi.Printf("dates warning nil?:%#v GetAlertsWeek err:%v \n", clientSys != nil, err)
		if err != nil {
			return nil, err
		}
		return clientSys, nil
	default:
		logapi.Printf("nothing to do with %v", dataType)
		return nil, fmt.Errorf("DOING")
	}
}

// 轨迹查询
func (that *AlysisApi) CaptionHistory(c *gin.Context) {

	roomId := c.DefaultQuery("room_id", "")
	idCardNumber := c.DefaultQuery("id_card_number", "")
	profileId, _ := strconv.Atoi(c.DefaultQuery("matched_profile_id", ""))

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		logapi.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	user, err := UserUtil.GetUserByCtx(c)
	if err != nil {
		helpers.JSONs(c, code.AuthorizationError, code.ZhCNText[code.AuthorizationError])
		return
	}

	datas, err := AlysysOperator.GetCaptionHostory(Page, PageSize, profileId, roomId, StartAt, EndAt, idCardNumber, user.TenantId)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}
