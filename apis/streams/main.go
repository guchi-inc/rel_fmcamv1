package streams

import (
	"fmcam/common/databases"
	"fmcam/systems"
	"log"
	"os"
)

var (
	apilog          = log.New(os.Stdout, "stream info -", 13)
	UserService     = systems.ServiceGroupApp.SysUserGroup
	SysStreamGroup  = systems.ServiceGroupApp.StreamGroup
	SysFileMan      = systems.ServiceGroupApp.SysFileGroup
	SysObsFileGroup = systems.ServiceGroupApp.SysObsFileGroup
	StrTask         = DeviceTask{}

	OBSFile = databases.OBFileTask
)

type StreamGroup struct {
	FRouter   *faceinfo
	StrRouter *DeviceTask
}

func BuildFacer() *StreamGroup {
	frouter := faceinfo{}

	return &StreamGroup{FRouter: &frouter, StrRouter: &StrTask}
}

// camera配置
type CameraConf struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	RtspUrl     string `json:"url"`
	Status      bool   `json:"status"`
	LastSuccess string `json:"last_success"`
	PlayUrl     string `json:"play_url"`
}
