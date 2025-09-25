package sysutils

import (
	"fmcam/common/databases"
	"fmcam/models/seclients"
	"fmt"

	"github.com/google/uuid"
)

// 根据租户号 查人员类型信息
func (ut *UserUtil) UserTypeList(Page, PageSize int, TenantId string) (*seclients.ProfileTypeList, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		useType = seclients.ProfileTypeList{}
		DBConn  = databases.DBMysql
		err     error
		Filters = " WHERE enabled = 1  "
	)

	if TenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}

	limits := fmt.Sprintf(" ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	baseSql := "SELECT id,BIN_TO_UUID(tenant_id) as tenant_id,type_name,type_code,warning_level,warning_enabled,description,enabled FROM ProfileType " + Filters
	baseSql += limits
	err = DBConn.Select(&useType.Data, baseSql)
	logsys.Printf("get profile type:%v db get user:%#v err:%v \n", baseSql, useType, err)
	if err != nil {
		return nil, err
	}

	useType.Page = Page
	useType.PageSize = PageSize
	useType.Columns, _ = databases.FmGlobalMap("ProfileType", nil)
	return &useType, nil

}
