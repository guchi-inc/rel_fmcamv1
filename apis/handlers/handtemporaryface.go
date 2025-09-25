package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/systems/genclients/predicate"
	"fmcam/systems/genclients/temporaryface"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询临时库人脸
func (that *FmClientApi) GetTemporaryFace(c *gin.Context) {

	imageUrl := c.DefaultQuery("img_url", "")
	Ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	tenantId := c.DefaultQuery("tenant_id", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apilog.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	var conds = []predicate.TemporaryFace{}

	if imageUrl != "" {
		conds = append(conds, temporaryface.ImgURLEQ(imageUrl))
	}
	if Ids != 0 {
		conds = append(conds, temporaryface.IDEQ(int64(Ids)))
	}

	//过期时间区间
	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apilog.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, temporaryface.ExpiresTimeGT(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apilog.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, temporaryface.ExpiresTimeLT(EndAtDate))
	}

	datas, err := SerEntORMApp.TemporaryFaceList(Page, PageSize, c, conds, tenantId)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
	}

	helpers.JSONs(c, code.Success, datas)
}
