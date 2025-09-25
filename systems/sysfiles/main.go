package sysfiles

import (
	"fmcam/ctrlib/utils/storeobs"
	"fmcam/systems/users"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	logger        = log.New(os.Stdout, "system-files: ", 13)
	MinFileObs    = ObsFile{}
	ObsApi        = storeobs.StorageObs
	UserService   = users.DefUserRepo //用户工具
	unSupportName = []string{" ", "\\", "+", "-", ":", ";", ",", "%", "@", "!", "#", "$", "%", "^", "&", "*", "(", ")"}
)

// 检查是否包括非法字符
func CheckNameVaild(distStr string) (bool, error) {

	for _, ss := range unSupportName {
		if strings.Contains(distStr, ss) {
			return false, fmt.Errorf("file name :%v error:%v", distStr, ss)
		}
	}

	if strings.Count(distStr, "/") > 5 {
		return false, fmt.Errorf("file name :%v error too much :%v", distStr, "/")
	}
	return true, nil
}
