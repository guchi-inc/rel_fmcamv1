package handlers

import (
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmcam/systems/genclients/devices"
	"fmcam/systems/genclients/predicate"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询设备信息
func (that *FmClientApi) GetDevices(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	location := c.DefaultQuery("location", "")
	url := c.DefaultQuery("url", "")
	ids, _ := strconv.Atoi(c.DefaultQuery("id", "0"))

	tenantId := c.DefaultQuery("tenant_id", "")

	PageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	Page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var conds = []predicate.Devices{}

	if name != "" {
		conds = append(conds, devices.NameContains(name))
	}
	if location != "" {
		conds = append(conds, devices.LocationContains(location))
	}
	if url != "" {
		conds = append(conds, devices.URLEQ(url))
	}
	if ids != 0 {
		conds = append(conds, devices.IDEQ(int64(ids)))
	}

	datas, err := SerEntORMApp.DevicesList(Page, PageSize, c, conds, tenantId)
	apilog.Printf("query data:%#v conds:%#v \n", datas, conds)

	if err != nil {
		helpers.JSONs(c, code.NullData, err)
		return
	}
	helpers.JSONs(c, code.Success, datas)

}
