package users

import (
	"encoding/json"
	"fmcam/common/databases"
	"fmcam/models/seclients"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/predicate"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 执行 日志查询
func (u *UserRepository) SqlLogList(Page, PageSize int, Ctx *gin.Context, conds []predicate.SqlLog) (*seclients.SQLLogPages, error) {

	var (
		client    = databases.EntClient
		DBConn    = databases.DBMysql
		responses = seclients.SQLLogPages{Page: Page, PageSize: PageSize}
	)

	total, err := client.Debug().SqlLog.Query().Where(conds...).Count(Ctx)
	if err != nil {
		return nil, err
	}

	offset := PageSize * (Page - 1)
	logs, err := client.Debug().SqlLog.
		Query().Where(conds...).Order(genclients.Desc("id")).
		Limit(PageSize).
		Offset(offset).
		All(Ctx)
	if err != nil {
		return nil, err
	}

	//事务执行
	tx, err := DBConn.Beginx()
	if err != nil {
		return nil, err
	}
	tx.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")

	logger.Printf("logs of conds len:%v \n", len(logs))
	for _, log := range logs {
		creaTime := log.CreatedTime.Format("2006-01-02 15:04:05")

		var userPostion string
		if log.PkValue != 0 {
			baseSql := fmt.Sprintf("SELECT position FROM fm_user WHERE id = '%v' order by id desc limit 1 ", log.PkValue)
			err := tx.Get(&userPostion, baseSql)
			logger.Println("查询用户类型信息postion:", baseSql, "  ", err)
		}

		var userCreator string
		if log.PkValue != 0 {
			baseSql := fmt.Sprintf("SELECT login_name FROM fm_user WHERE id = '%v' or login_name = '%v' order by id desc limit 1 ", *log.Creator, *log.Creator)
			err := tx.Get(&userCreator, baseSql)
			logger.Println("查询操作者信息login_name:", baseSql, "  ", err)
		}

		resp := seclients.SQLLogResponse{
			ID:        log.ID,
			Action:    *log.Action,
			Table:     log.TableName,
			PK:        &log.PkValue,
			Args:      log.Args,
			CreatedAt: &creaTime,
			Position:  &userPostion,
			Creator:   &userCreator,
		}

		// // old_data
		// var old map[string]interface{}
		// if log.OldData != nil {
		// 	raw, _ := json.Marshal(log.OldData)
		// 	_ = json.Unmarshal(raw, &old)
		// }

		// // new_data
		var new map[string]interface{}
		if log.NewData != nil {
			raw, _ := json.Marshal(log.NewData)
			_ = json.Unmarshal(raw, &new)
		}

		// 组装变更对比
		var newCNAME string
		var qcnameField string
		for k, newVal := range new {
			// resp.Changes = append(resp.Changes, seclients.FieldDiff{
			// 	Field: k,
			// 	Old:   old[k],
			// 	New:   newVal,
			// })
			if *log.Action != "INSERT" {
				if len(qcnameField) == 0 {
					qcnameField += fmt.Sprintf("%v", "'"+k+"'")
				} else {
					qcnameField += fmt.Sprintf("%v", ","+"'"+k+"'")
				}
				if k == "enabled" {
					switch newVal.(type) {
					case float64:
						if newVal.(float64) == 1 {
							newCNAME += "启用"
						}
						if newVal.(float64) == 0 {
							newCNAME += "禁用"
						}
					case int:
						if newVal.(int) == 1 {
							newCNAME += "启用"
						}
						if newVal.(int) == 0 {
							newCNAME += "禁用"
						}
					case string:
						if newVal.(string) == "1" {
							newCNAME += "启用"
						}
						if newVal.(string) == "0" {
							newCNAME += "禁用"
						}
					}
				}

			}
		}

		//查询字段说明
		if qcnameField != "" {
			baseSql := fmt.Sprintf("SELECT cname FROM field_metadata WHERE table_name = '%v' AND name in (%v) order by id desc   ", log.TableName, qcnameField)
			var fcnames []string
			err = tx.Select(&fcnames, baseSql)
			logger.Println("查询字段说明cname:", baseSql, "  ", err)

			for _, fcname := range fcnames {
				newCNAME += fcname + ","
			}
		}

		if *log.Action == "INSERT" {
			var ac string
			if newCNAME != "" {
				ac = "新增" + newCNAME
			} else {
				if log.TableName == "fm_user" {
					ac = "新增用户"
				} else {
					ac = "新增操作"
				}
			}
			resp.ChangesCname = &ac
		}
		if *log.Action == "UPDATE" {
			ac := "修改" + newCNAME
			resp.ChangesCname = &ac
		}

		responses.Data = append(responses.Data, resp)
	}

	err = tx.Commit()
	if err != nil {
		logger.Println("查询组织日志失败:", responses, "  ", err)
	}
	responses.Total = total
	if len(responses.Data) > 0 {
		cc, err := databases.FmGlobalMap("sql_logs", nil)
		if err != nil {
			logger.Println("查询日志表字段说明失败:", cc, "  ", err)
		}
		responses.Columns = cc
	}
	return &responses, nil
}

// type DescList struct {
// 	ID        int      `json:"id" db:"id"`                 //所属sql_log的记录 id
// 	TabelName string   `json:"table_name" db:"table_name"` //该条记录所属表名
// 	Fields    []string `json:"fields" db:"fields"`         //要查询的字段
// 	Action    string   `json:"action" db:"action"`         //操作 UPDATE  INSERT 新增不需要查询字段解释
// 	CName     []string `json:"cname" db:"cname"`           // 字段解释
// 	Changes   string   `json:"changes" db:"changes"`       //修改内容说明
// }

// // 查询字段解释和变更说明
// func (u *UserRepository) SqlLogField(fields []DescList) ([]DescList, error) {

// 	var (
// 		DBConn  = databases.DBMysql
// 		baseSql = "SELECT name,cname,is_visible,is_searchable,is_editable,is_required,data_type,max_length,sort FROM field_metadata "
// 		filter  string
// 	)

// 	var NewDescList = []DescList{}
// 	for _, ques := range fields {
// 		var erpDescs = []erps.GcDesc{}
// 		if ques.TabelName == "" {
// 			return nil, fmt.Errorf("属性所属表不能为空")
// 		}

// 		filter += " WHERE table_name = " + "'" + ques.TabelName + "'"
// 		var fsField string
// 		if len(ques.Fields) > 0 {
// 			for _, f := range ques.Fields {
// 				if len(fsField) <= 0 {
// 					fsField += f
// 				} else {
// 					fsField += "," + f
// 				}
// 			}
// 		}

// 		filter += fmt.Sprintf(" AND name in (%v)", fsField)

// 		limits := " ORDER BY sort asc, id asc "
// 		baseSql += filter
// 		baseSql += limits
// 		err := DBConn.Select(&erpDescs, baseSql)
// 		logger.Printf("field select of table:%v sql: %v ,err :%v \n", ques.TabelName, baseSql, err)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if ques.Action == "INSERT" {
// 			ques.Changes = "新增"
// 		}
// 		if ques.Action == "UPDATE" {
// 			ques.Changes = "修改"
// 		}
// 		for _, desc := range erpDescs {
// 			if desc.CName != "" {
// 				ques.Changes += desc.CName + ","
// 			}
// 		}

// 		NewDescList = append(NewDescList, ques)

// 	}
// 	return NewDescList, nil
// }
