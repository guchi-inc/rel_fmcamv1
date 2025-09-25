package resources

import "fmcam/apis"

var (
	ApiLimitCtl    = apis.ApiGroupApp.LimitsApis
	ApiGovAddrCtl  = apis.ApiGroupApp.GAddrsApis
	ApiCustomerCtl = apis.ApiGroupApp.CustomerTenantApis
)

type ResourcesRouter struct{}
