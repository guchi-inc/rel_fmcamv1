package sysutils

import (
	"fmcam/common/databases"
	"fmcam/models/erps"
	"fmt"
)

type GovAddrUtil struct{}

// 获取运营管理 [省份] 列表
func (gt *GovAddrUtil) GAddrProvinceList(Page, PageSize int, name, codes string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT * FROM `fm_gov_province` "

		total    int
		useProvs = erps.FmGovProvinceList{}
		DBConn   = databases.DBMysql
	)

	//条件
	Filters := " WHERE id > 0 "
	if codes != "" {
		Filters += " AND code LIKE '" + codes + "%' "
	}
	if name != "" {
		Filters += " AND name LIKE '%" + name + "%' "
	}

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_gov_province`  " + Filters

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &useProvs, err
	}

	OffLimits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	//admin 查询全部角色数据，具有编辑权限
	if name == "admin" {
		baseSql = "SELECT  * FROM  `fm_gov_province`  "
		baseSql += Filters
		baseSql += OffLimits
		err = DBConn.Select(&useProvs.Data, baseSql)
		logsys.Printf("from db admin baseSql:%#v  err:%v \n", baseSql, err)
		if err != nil {
			return nil, err
		}

		useProvs.Total = uint(total)
		useProvs.Size = uint(len(useProvs.Data))
		return &useProvs, nil
	}

	baseSql += Filters
	baseSql += OffLimits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&useProvs.Data, baseSql)
	logsys.Printf("from limits select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useProvs.Total = uint(total)
	useProvs.Size = uint(PageSize)
	return &useProvs, nil
}

// 获取运营管理 [城市] 列表
func (gt *GovAddrUtil) GAddrCityList(Page, PageSize int, name, codes, provCode string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT * FROM `fm_gov_city` "

		total    int
		useProvs = erps.FmGovCityList{}
		DBConn   = databases.DBMysql
	)

	//条件
	Filters := " WHERE id > 0 "
	if codes != "" {
		Filters += " AND code LIKE '" + codes + "%' "
	}
	if name != "" {
		Filters += " AND name LIKE '%" + name + "%' "
	}

	if provCode != "" {
		Filters += " AND province_code LIKE '%" + provCode + "%' "
	}

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_gov_city`  " + Filters

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &useProvs, err
	}

	OffLimits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	//admin 查询全部角色数据，具有编辑权限
	if name == "admin" {
		baseSql = "SELECT  * FROM  `fm_gov_city`  "
		baseSql += Filters
		baseSql += OffLimits
		err = DBConn.Select(&useProvs.Data, baseSql)
		logsys.Printf("from db fm_gov_city baseSql:%#v  err:%v \n", baseSql, err)
		if err != nil {
			return nil, err
		}

		useProvs.Total = uint(total)
		useProvs.Size = uint(len(useProvs.Data))
		return &useProvs, nil
	}

	baseSql += Filters
	baseSql += OffLimits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&useProvs.Data, baseSql)
	logsys.Printf("from fm_gov_city select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useProvs.Total = uint(total)
	useProvs.Size = uint(PageSize)
	return &useProvs, nil
}

// 获取运营管理 [区县] 列表
func (gt *GovAddrUtil) GAddrAreaList(Page, PageSize int, name, codes, provCode, cityCode string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT * FROM `fm_gov_area` "

		total    int
		useProvs = erps.FmGovAreaList{}
		DBConn   = databases.DBMysql
	)

	//条件
	Filters := " WHERE id > 0 "
	if codes != "" {
		Filters += " AND code LIKE '" + codes + "%' "
	}
	if name != "" {
		Filters += " AND name LIKE '%" + name + "%' "
	}
	if provCode != "" {
		Filters += " AND province_code LIKE '%" + provCode + "%' "
	}
	if cityCode != "" {
		Filters += " AND city_code LIKE '%" + cityCode + "%' "
	}

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_gov_area`  " + Filters

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &useProvs, err
	}

	OffLimits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	//admin 查询全部角色数据，具有编辑权限
	if name == "admin" {
		baseSql = "SELECT  * FROM  `fm_gov_area`  "
		baseSql += Filters
		baseSql += OffLimits
		err = DBConn.Select(&useProvs.Data, baseSql)
		logsys.Printf("from db fm_gov_area baseSql:%#v  err:%v \n", baseSql, err)
		if err != nil {
			return nil, err
		}

		useProvs.Total = uint(total)
		useProvs.Size = uint(PageSize)
		return &useProvs, nil
	}

	baseSql += Filters
	baseSql += OffLimits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&useProvs.Data, baseSql)
	logsys.Printf("from limits select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useProvs.Total = uint(total)
	useProvs.Size = uint(PageSize)
	return &useProvs, nil
}

// 获取运营管理 [街道] 列表
func (gt *GovAddrUtil) GAddrStreetList(Page, PageSize int, name, codes, provCode, cityCode, areaCode string) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = "SELECT * FROM `fm_gov_street` "

		total    int
		useProvs = erps.FmGovStreetList{}
		DBConn   = databases.DBMysql
	)

	//条件
	Filters := " WHERE id > 0 "
	if codes != "" {
		Filters += " AND code LIKE '" + codes + "%' "
	}
	if name != "" {
		Filters += " AND name LIKE '%" + name + "%' "
	}
	if provCode != "" {
		Filters += " AND province_code LIKE '%" + provCode + "%' "
	}
	if cityCode != "" {
		Filters += " AND city_code LIKE '%" + cityCode + "%' "
	}
	if areaCode != "" {
		Filters += " AND area_code LIKE '%" + areaCode + "%' "
	}

	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_gov_street`  " + Filters

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &useProvs, err
	}

	OffLimits := fmt.Sprintf("  ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)
	//admin 查询全部角色数据，具有编辑权限
	if name == "admin" {
		baseSql = "SELECT  * FROM  `fm_gov_street`  "
		baseSql += Filters
		baseSql += OffLimits
		err = DBConn.Select(&useProvs.Data, baseSql)
		logsys.Printf("from db fm_gov_street baseSql:%#v  err:%v \n", baseSql, err)
		if err != nil {
			return nil, err
		}

		useProvs.Total = uint(total)
		useProvs.Size = uint(PageSize)
		return &useProvs, nil
	}

	baseSql += Filters
	baseSql += OffLimits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&useProvs.Data, baseSql)
	logsys.Printf("from limits select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useProvs.Total = uint(total)
	useProvs.Size = uint(PageSize)
	return &useProvs, nil
}
