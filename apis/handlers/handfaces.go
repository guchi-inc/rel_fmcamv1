package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/models/seclients"
	"fmcam/systems/genclients/faces"
	"fmcam/systems/genclients/predicate"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询人脸
func (that *FmClientApi) GetFaces(c *gin.Context) {

	profileId, _ := strconv.Atoi(c.DefaultQuery("profile_id", ""))
	isPrimary := c.DefaultQuery("is_primary", "")
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
	var conds = []predicate.Faces{}

	if profileId != 0 {
		conds = append(conds, faces.ProfileIDEQ(int64(profileId)))
	}
	if Ids != 0 {
		conds = append(conds, faces.IDEQ(int64(Ids)))
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apilog.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, faces.CreatedTimeGT(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apilog.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, faces.CreatedTimeLT(EndAtDate))
	}

	var ispm *int = nil
	if isPrimary != "" {
		ispms, err := strconv.Atoi(isPrimary)
		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
		ispm = &ispms
	}

	datas, err := SerEntORMApp.FacesList(&Page, &PageSize, &profileId, ispm, c, conds, tenantId)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)
}

// 删除
func (that *FmClientApi) UpdateFace(c *gin.Context) {

	var (
		PParm = seclients.FaceParam{}
	)

	err := c.ShouldBindJSON(&PParm)
	if err != nil {
		err = fmt.Errorf("%v:user:%v ", code.ZhCNText[code.ParamBindError], PParm)
		helpers.JSONs(c, code.ParamBindError, err)
		return
	}

	datas, err := SerEntORMApp.FacesUpdate(c, &PParm)
	apilog.Printf("query data:%#v PParm:%#v \n", datas, PParm)
	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": datas})

}
