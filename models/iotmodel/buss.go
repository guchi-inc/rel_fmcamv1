// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package iotmodel

import "fmcam/models/erps"

type PMStatus struct {
	ID        int    `json:"id" db:"id"`
	StateName string `json:"state_name" db:"state_name"`
	StateType string `json:"state_type" db:"state_type"`

	Description string `json:"description" db:"description"`
	IsFinal     string `json:"is_final" db:"is_final"`
}

type PMStateTrans struct {
	ID            int    `json:"id" db:"id"`
	FromStateName string `json:"from_state_name" db:"from_state_name"`
	ToStateName   string `json:"to_state_name" db:"to_state_name"`
	Description   string `json:"description" db:"description"`
}

type PMStatusList struct {
	Total int64      `json:"total" db:"total"`
	Size  int64      `json:"size" db:"size"`
	Data  []PMStatus `json:"data" db:"data"`

	Columns []erps.GcDesc `db:"columns" json:"columns" form:"columns"` //字段标识

}
