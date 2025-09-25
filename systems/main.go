// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package systems

import (
	"fmcam/systems/streams"
	"fmcam/systems/sysclients"
	"fmcam/systems/sysfiles"
	"fmcam/systems/sysutils"
	"fmcam/systems/users"
)

var (
	UserDef       = users.DefUserRepo
	DepUtils      = sysutils.NewDepartmentUtil()
	CustomerUtils = sysutils.CustomerUtil{}
)

type ServiceGroup struct {
	SysUserGroup    *users.UserRepository
	StreamGroup     *streams.StreamRep
	SysFileGroup    *sysfiles.FilesDB
	SysObsFileGroup *sysfiles.ObsFile
	LimitsGroup     *sysutils.LimitsUtil
	GovAddrGroup    *sysutils.GovAddrUtil
	CustomerGroup   *sysutils.CustomerUtil
	SysClientsGroup *sysclients.ORMCleint
}

func NewServiceGroup(url string) *ServiceGroup {

	sg := &ServiceGroup{
		SysUserGroup:    UserDef,
		StreamGroup:     streams.NewStreamRepo(),
		SysFileGroup:    &sysfiles.FilesDB{},
		SysObsFileGroup: &sysfiles.MinFileObs,
		LimitsGroup:     &sysutils.LimitsUtil{},
		GovAddrGroup:    &sysutils.GovAddrUtil{},
		CustomerGroup:   &sysutils.CustomerUtil{},

		SysClientsGroup: &sysclients.ORMCleint{},
	}

	return sg
}

var ServiceGroupApp = NewServiceGroup("")
