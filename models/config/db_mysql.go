// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package config

type MysqlDB struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *MysqlDB) Dsn() string {
	newDsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config

	// newLog := configs.DefaultSuggerLogger()
	// newLog.Infof("newDsn:%#v \n", newDsn)/
	return newDsn
}

func (m *MysqlDB) GetLogMode() string {
	return m.LogMode
}
