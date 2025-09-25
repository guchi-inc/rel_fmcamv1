package handlers

import (
	"fmcam/common/configs"
	"fmcam/systems"
	"fmcam/systems/sysclients"
	"log"
	"os"
)

var (
	apilog      = log.New(os.Stdout, "handlers INFO -", 13)
	zapLog      = configs.GZapLog
	GovApis     = GovAddrApi{}
	UserService = systems.ServiceGroupApp.SysUserGroup

	SerEntORMApp   = systems.ServiceGroupApp.SysClientsGroup
	SysGovAddGroup = systems.ServiceGroupApp.GovAddrGroup
	//迁入和签出登记
	SysCheckOutIn = sysclients.SysCleint{}
)

type GovAddrApi struct{}

type FmClientApi struct{}
