package alysis

import (
	"fmcam/systems"
	"fmcam/systems/alysys"
	"fmcam/systems/sysutils"
	"log"
	"os"
)

var (
	UserUtil       = sysutils.NewUserUtil()
	CustomerUtil   = sysutils.CustomerUtil{}
	SerEntORMApp   = systems.ServiceGroupApp.SysClientsGroup
	AlysysOperator = alysys.AlysysOperator
	logapi         = log.New(os.Stdout, "INFO -", 13)
)
