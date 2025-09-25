package resources

import (
	"fmcam/systems"
	"log"
	"os"
)

var (
	LimitsApis = LimitsApi{}

	apilog             = log.New(os.Stdout, "api limits INFO -", 13)
	UserService        = systems.ServiceGroupApp.SysUserGroup
	LimitsServiceGroup = systems.ServiceGroupApp.LimitsGroup
	SysCustomesGroup   = systems.ServiceGroupApp.CustomerGroup
	SerEntORMApp       = systems.ServiceGroupApp.SysClientsGroup
)

type LimitsApi struct{}
