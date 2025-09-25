package resources

import (
	"fmcam/ctrlib/middleware"

	"github.com/gin-gonic/gin"
)

// 资源查询
func (s *ResourcesRouter) InitCustomerRouter(r *gin.RouterGroup) {

	resCtl := ApiCustomerCtl
	r.Use(middleware.AuthorizeJWT())
	res := r.Group("/customer")

	res.GET("/list", resCtl.GetTenantsList)
	res.POST("/new", resCtl.NewTenantsInfo)
	res.POST("/update", resCtl.UpdateTenant)

	//pms管理api器
	pms := r.Group("/pms")
	pms.GET("/list", resCtl.GetPMSApiList)
	pms.POST("/update", resCtl.UpdatePMSApi)
	pms.POST("/new", resCtl.NewPMSApi)

}
