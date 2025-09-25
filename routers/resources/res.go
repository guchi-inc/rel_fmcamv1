package resources

import (
	"fmcam/common/databases"
	"fmcam/ctrlib/middleware"

	"github.com/gin-gonic/gin"
)

// 资源查询
func (s *ResourcesRouter) InitResDevice(r *gin.RouterGroup) {

	r.Use(middleware.AuthorizeJWT())

	resCtl := ApiLimitCtl
	res := r.Group("/resource")
	res.GET("/limits/list", resCtl.LimitsList)
	res.POST("/limits/new", resCtl.LimitsInsert)
	res.POST("/limits/update", resCtl.LimitsUpdate)

	//fm_role_limits 表更新，决定用户能使用哪些功能菜单

	res.GET("/link/role/list", resCtl.RoleLimitsList)
	res.POST("/link/role/new", resCtl.RoleLimitsInsert)
	res.POST("/link/role/update", resCtl.RoleLimitsOneUpdate)

	res.GET("/link/role/history", resCtl.RoleLimitsListHistory)
	//单个删除
	res.POST("/link/role/delete", resCtl.DelRoleLimits)

	//省市区地址选择器
	addrCtl := ApiGovAddrCtl
	govaddr := r.Group("/govaddr")
	govaddr.GET("/province/list", addrCtl.GetProvinces)
	govaddr.POST("/province/new", addrCtl.NewProvince)
	govaddr.POST("/province/update", addrCtl.UpdateProvince)

	govaddr.GET("/city/list", addrCtl.GetCitys)
	govaddr.POST("/city/update", addrCtl.UpdateCity)
	govaddr.POST("/city/new", addrCtl.NewCity)

	govaddr.GET("/area/list", addrCtl.GetAreas)
	govaddr.POST("/area/update", addrCtl.UpdateArea)
	govaddr.POST("/area/new", addrCtl.NewArea)

	govaddr.GET("/street/list", addrCtl.GetStreets)
	govaddr.POST("/street/update", addrCtl.UpdateStreet)
	govaddr.POST("/street/new", addrCtl.NewStreet)

	/// 客服人员维护
	ser := r.Group("/service")
	ser.GET("/dedicate/list", addrCtl.GetDedicatedService)
	ser.POST("/dedicate/update", addrCtl.UpdateDedicatedService)
	ser.POST("/dedicate/new", addrCtl.NewDedicatedService)

	//酒店留言建议列表
	ser.GET("/demands/list", addrCtl.GetDemands)
	ser.POST("/demands/new", addrCtl.NewDemands)

	//字段属性查询
	res.Use(databases.OperatorMiddleware())
	res.GET("/fields/states", resCtl.GetFieldsState)
	res.POST("/fields/states", resCtl.UpdateFieldMetadata)

}
