// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/common/databases"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
)

// 用户更新 地址 和 密码信息
func (u *UserRepository) BussUpdateUser(userInfo *erps.FmUser) (any, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		DBConn = databases.DBMysql
	)

	userOrigin, err := u.GetUserById(int(userInfo.Id))

	if err != nil {
		err = fmt.Errorf("数据不存在:%v", err)
		return nil, err
	}

	//只允许 管理员 根据 id 更新
	if userOrigin.LoginName == "admin" {
		err = fmt.Errorf("%v:正在修改超级管理员密码", code.ZhCNText[code.AuthorizedPowerError])
		return nil, err
	}
	upTotal := 0

	if userInfo.Password != "" && (userOrigin.Password != userInfo.Password) {
		userOrigin.Password = userInfo.Password
		upTotal += 1
	}

	if upTotal == 0 {
		err = fmt.Errorf("数据无需修改。")
		return nil, err
	}

	personStructs := []erps.FmUser{
		*userOrigin,
	}

	rst, err := DBConn.NamedExec(`UPDATE fm_user SET 
	password=:password WHERE id=:id`, personStructs)
	logger.Printf("from db select:%#v  err:%v \n", rst, err)
	if err != nil {
		return nil, err
	}

	logger.Printf("from db get:%v user:%#v err:%v \n", userInfo, rst, err)

	return rst.RowsAffected()
}
