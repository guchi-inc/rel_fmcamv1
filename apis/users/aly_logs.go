package users

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/code"
	"fmcam/systems/genclients/predicate"
	"fmcam/systems/genclients/sqllog"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户操作日志
func (that *UsersApi) GetUserSQLLogs(c *gin.Context) {

	userId := c.DefaultQuery("id", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apidebug.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	var conds = []predicate.SqlLog{}

	conds = append(conds, sqllog.TableNameEQ("fm_user"))
	if userId != "" {
		userIds, _ := strconv.Atoi(userId)
		conds = append(conds, sqllog.PkValueEQ(userIds))
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apidebug.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, sqllog.CreatedTimeGTE(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apidebug.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, sqllog.CreatedTimeLTE(EndAtDate))
	}

	responses, err := UserService.SqlLogList(Page, PageSize, c, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"error": err})
		return
	}
	helpers.JSONs(c, code.Success, responses)

}

// 设备操作日志
func (that *UsersApi) GetDeviceSQLLogs(c *gin.Context) {

	deviceID := c.DefaultQuery("id", "")
	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	StartAt := c.DefaultQuery("start_time", "") //开始时间
	EndAt := c.DefaultQuery("end_time", "")     //结束时间

	if Page <= 0 {
		Page = 1
	}
	if PageSize < 5 {
		PageSize = 5
	}

	if StartAt != "" && EndAt != "" {
		StartAt, EndAt, err := timeutil.CheckDateRange(StartAt, EndAt)
		apidebug.Println("select un finished total err:", StartAt, "Fileter:", EndAt, err)

		if err != nil {
			helpers.JSONs(c, code.ParamError, err)
			return
		}
	}

	var conds = []predicate.SqlLog{}

	conds = append(conds, sqllog.TableNameEQ("Device"))
	if deviceID != "" {
		deviceIDs, _ := strconv.Atoi(deviceID)
		conds = append(conds, sqllog.PkValueEQ(deviceIDs))
	}

	if StartAt != "" && EndAt != "" {
		StartAtDate, err := timeutil.ParseCSTInLocation(StartAt)
		apidebug.Printf("  StartAtDate:%#v to StartAt:%#v err:%#v \n", StartAtDate, StartAt, err)

		conds = append(conds, sqllog.CreatedTimeGT(StartAtDate))
		EndAtDate, err := timeutil.ParseCSTInLocation(EndAt)
		apidebug.Printf("  EndAtDate:%#v to EndAt:%#v err:%#v \n", EndAtDate, EndAt, err)

		conds = append(conds, sqllog.CreatedTimeLT(EndAtDate))
	}

	responses, err := UserService.SqlLogList(Page, PageSize, c, conds)
	if err != nil {
		helpers.JSONs(c, code.NullData, gin.H{"error": err})
		return
	}
	helpers.JSONs(c, code.Success, responses)

}
