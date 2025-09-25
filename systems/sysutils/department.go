// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package sysutils

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/models/erps"
	"fmt"
	"time"
)

// 根据部门名称获取部门信息
func (ut *DepartmentUtil) GeDepartmentByName(dep_name string) (*erps.GcDepartmentStruct, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow

	var (
		DBConn = databases.DBMysql
		usedep = erps.GcDepartmentStruct{}
	)

	err := DBConn.Get(&usedep, "SELECT * FROM fm_department  WHERE dep_name=?", dep_name)

	logsys.Printf("from db get:%v name:%#v err:%v \n", usedep, dep_name, err)
	if err != nil {
		return nil, err
	}

	return &usedep, nil
}

// 根据部门id获取部门信息
func (ut *DepartmentUtil) GeDepartmentById(id int64) (*erps.GcDepartmentStruct, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow

	var (
		DBConn = databases.DBMysql
		usedep = erps.GcDepartmentStruct{}
	)

	err := DBConn.Get(&usedep, "SELECT * FROM fm_department  WHERE id=?", id)

	logsys.Printf("from db get:%v name:%#v err:%v \n", usedep, id, err)
	if err != nil {
		return nil, err
	}

	return &usedep, nil
}

// 获取部门 [名称] 列表  full = false
func (ut *DepartmentUtil) GeDepartmentNameList(full bool) (*erps.GcDepartmentStructList, error) {
	// 查询单个数据 You can also get a single result, a la QueryRow

	var (
		DBConns = databases.DBMysql
		baseSql = `SELECT id, 
		IFNULL(dep_type, '') AS dep_type , 
		IFNULL(dep_name,'') AS dep_name, 
		IFNULL(manager_name,'') AS manager_name,created_time,creator FROM fm_department WHERE delete_flag != '1' ;`
	)

	usedep := erps.GcDepartmentStructList{}

	if !full {
		baseSql = "SELECT dep_name FROM fm_department WHERE delete_flag != '1' ;"
	}

	logsys.Printf("DBConn:%#v  is nil:%#v\n", DBConns, DBConns == nil)
	err := DBConns.Select(&usedep.Data, baseSql)

	logsys.Printf("from db select:%#v  err:%v \n", usedep, err)
	if err != nil {
		return nil, err
	}

	usedep.Total = uint(len(usedep.Data))
	usedep.Size = uint(len(usedep.Data))

	usedep.Columns = FiledOper.DepartmentMap(true)
	return &usedep, nil
}

// 新增后台管理 [部门管理]
func (ut *DepartmentUtil) InsertDepartment(role_record *erps.GcDepartmentStruct) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	role_record.CreatedAt = time.Now().Local().In(configs.Loc).Format("2006-01-02 15:04:05")
	recordStructs := []erps.GcDepartmentStruct{
		*role_record,
	}

	rst, err := DBConn.NamedExec(`INSERT INTO fm_department ( 
		dep_type,
		manager_name,
		created_time,
		dep_name,
		creator)
        VALUES (:dep_type, :manager_name,:created_time,:dep_name, :creator)`, recordStructs)

	logsys.Printf("new fm_department insert:%#v, with db rst:%v err:%v\n", recordStructs, rst, err)
	if err == nil {
		total, err = rst.LastInsertId()
		return total, err
	}

	return 0, fmt.Errorf("Insert fm_department failed:%#v\n", err)
}

// 新增后台管理 [部门管理]
func (ut *DepartmentUtil) UpdateDepartment(recordStructs *erps.GcDepartmentStruct) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		total  int64
		DBConn = databases.DBMysql
	)

	recordStructs.CreatedAt = time.Now().Local().In(configs.Loc).Format("2006-01-02 15:04:05")

	currentDepartMent, err := ut.GeDepartmentById(recordStructs.ID)

	if currentDepartMent == nil || err != nil {
		return 0, fmt.Errorf("Department %s not exist", recordStructs.DepName)
	}

	if recordStructs.DeleteFlag != "" {
		rst, err := DBConn.NamedExec(`UPDATE fm_department SET 
		delete_flag =:delete_flag
        WHERE id =:id`, recordStructs)
		logsys.Printf("删除:%v,err:%v", recordStructs.DepName, err)
		if err != nil {
			return 0, fmt.Errorf("删除失败: %s err:%v", recordStructs.DepName, err)
		}
		return rst.RowsAffected()
	}

	var totalUp int
	if recordStructs.DepName != "" && recordStructs.DepName != currentDepartMent.DepName {
		currentDepartMent.DepName = recordStructs.DepName
		totalUp += 1
	}
	if recordStructs.DepType != "" && recordStructs.DepType != currentDepartMent.DepType {
		currentDepartMent.DepType = recordStructs.DepType
		totalUp += 1
	}
	if recordStructs.ManagerName != "" && recordStructs.ManagerName != currentDepartMent.ManagerName {
		currentDepartMent.ManagerName = recordStructs.ManagerName
		totalUp += 1
	}
	if recordStructs.Creator != "" && recordStructs.Creator != currentDepartMent.Creator {
		currentDepartMent.Creator = recordStructs.Creator
		totalUp += 1
	}

	if totalUp == 0 {
		return 0, fmt.Errorf("no update")
	}

	rst, err := DBConn.NamedExec(`UPDATE fm_department SET 
		dep_type =:dep_type,
		manager_name =:manager_name,
		created_time =:created_time,
		dep_name =:dep_name,
		creator =:creator
        WHERE id =:id`, currentDepartMent)

	logsys.Printf("new fm_department insert:%#v, with db rst:%v err:%v\n", recordStructs, rst, err)
	if err == nil {
		total, err = rst.LastInsertId()
		return total, err
	}

	return 0, fmt.Errorf("Insert fm_department failed:%#v\n", err)
}
