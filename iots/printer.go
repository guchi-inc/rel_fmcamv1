// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package iots

import (
	"encoding/json"
	"fmcam/common/configs"
	"fmcam/ctrlib/utils"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models"
	"fmcam/models/erps"
	"fmcam/models/trade"
	"fmt"
)

/*

iots 设备操作层与 DAO(数据操作)层 位于同一层次.

1 对接第三方打印接口
2 封装打印相关功能
*/

// 打印执行接口
func (t *Iots) PrinterDo(ids string) (*trade.DCloudPrint, error) {
	task := allTask
	if ids == "" {
		ids = userId
	}
	repByte := task.PostPrinter(appKey, appSecret, "", ids, false)
	fmt.Printf("printer bytes:%T strings:%v\n", repByte, string(repByte))
	/*
		返回的结构
		strings:{
		"code":100,
		"data":"eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkdWRpYW4tYXBwLXNlY3VyaXR5Iiwic3ViIjoie1wiYXBwSWRcIjpcImRlbW8tYXBwLWtleVwiLFwidXNlcklkXCI6XCIxXCJ9IiwiYXVkIjoidXNlciIsImV4cCI6NDg0MDE2MzI3MSwiaWF0IjoxNjg2NTYzMjcxfQ.6g10EtPxagrE16bbGEpswax4K5aT1u3VxWyyAk0YYpo",
		"message":"成功",
		"attachment":null}
	*/
	newPrinter := trade.DCloudPrint{}
	err := json.Unmarshal(repByte, &newPrinter)
	if err != nil {
		return nil, err
	}
	if newPrinter.Code != 100 {
		logger.Printf("printer bytes:%v strings:%v\n", newPrinter, string(repByte))
		return nil, fmt.Errorf("DCloude print error:%v\n", newPrinter)
	}
	return &newPrinter, nil
}

// 打印执行接口, 获取token，执行打印
func (t *Iots) PostPrinterDoing(ids string, tmps []erps.Params) (*trade.DCloudPrint, error) {

	var (
		tmp_id string
		//"[{\"design_no\":\"130611230711179647\",\"date\": \"2023-07-30 17:13:21\",
		//\"length\":\"189\",\"color\":\"red\",\"fact_name\":\"New yang tech.\",
		//\"type\":\"JIng A0703\",\"grammage\":\"82de/298.2\",\"pre_batch\":\"230730\"}]"

		print_data string = ""
		datamap           = map[string]string{}
	)

	logger.Printf("system doing receive ids:%#v param:%#v \n", ids, tmps)

	task := allTask
	tNow := timeutil.CSTLayoutString(nil)
	//拼接两个打印样式相关的字段和值
	if ids != "backup" {
		tmp_id = configs.PrinterTmpId.Default

		d := erps.TmpDefault{}

		for sk, v := range tmps {
			logger.Printf("find sk:%#v and param value:%#v \n", sk, v)
			if v.Field == "design_no" {
				d.DesignNo = v.Value
			} else if v.Field == "date" {
				d.Date = v.Value
			} else if v.Field == "length" {
				d.Length = v.Value
			} else if v.Field == "color" {
				d.Color = v.Value
			}

		}

		d.Date = tNow
		print_byte, err := json.Marshal(d)
		logger.Printf("tmp_id:%#v json marsh tmp :%#v ,err:%#v \n", tmp_id, string(print_byte), err)
		if err != nil {
			return nil, err
		}
		print_data = string(print_byte)
	} else {

		tmp_id = configs.PrinterTmpId.Backup

		d := erps.TmpBackUp{}
		for sk, v := range tmps {
			logger.Printf("find sk:%#v and param value:%#v \n", sk, v)
			if v.Field == "design_no" {
				d.DesignNo = v.Value
			} else if v.Field == "date" {
				d.Date = v.Value
			} else if v.Field == "length" {
				d.Length = v.Value
			} else if v.Field == "color" {
				d.Color = v.Value
			} else if v.Field == "fact_name" {
				d.FactName = v.Value
			} else if v.Field == "type" {
				d.Type = v.Value
			} else if v.Field == "grammage" {
				d.Grammage = v.Value
			} else if v.Field == "pre_batch" {
				d.PreBatch = v.Value
			}
		}

		d.Date = tNow
		print_byte, err := json.Marshal(d)
		logger.Printf("tmp_id:%v json marsh tmp :%#v ,err:%#v \n", tmp_id, string(print_byte), err)
		if err != nil {
			return nil, err
		}
		print_data = string(print_byte)
	}

	//传递模板id 和 数据
	datamap["doc_no"] = tmp_id
	//数据格式需要匹配 数组形式，[], 里面才是键值对映射。
	datamap["print_data"] = fmt.Sprintf("[%v]", print_data)
	accessToken := "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkdWRpYW4tYXBwLXNlY3VyaXR5Iiwic3ViIjoie1wiYXBwSWRcIjpcImRlbW8tYXBwLWtleVwiLFwidXNlcklkXCI6XCIxXCJ9IiwiYXVkIjoidXNlciIsImV4cCI6NDgzOTY0Mjk4MiwiaWF0IjoxNjg2MDQyOTgyfQ.iqJ46KlYAz7om6M_gD6BArNrZgvpZ_ezyVXgnIxn6fc"
	repByte := task.PostPrinterDoing("", accessToken, datamap, false)
	logger.Printf("printer datamap:%#v resp strings:%v\n", datamap, string(repByte))
	/*
		返回的结构
		strings:{
		"code":100,
		"data":"eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkdWRpYW4tYXBwLXNlY3VyaXR5Iiwic3ViIjoie1wiYXBwSWRcIjpcImRlbW8tYXBwLWtleVwiLFwidXNlcklkXCI6XCIxXCJ9IiwiYXVkIjoidXNlciIsImV4cCI6NDg0MDE2MzI3MSwiaWF0IjoxNjg2NTYzMjcxfQ.6g10EtPxagrE16bbGEpswax4K5aT1u3VxWyyAk0YYpo",
		"message":"成功",
		"attachment":null}
	*/
	newPrinter := trade.DCloudPrint{}
	err := json.Unmarshal(repByte, &newPrinter)
	if err != nil {
		return nil, err
	}
	if newPrinter.Code != 100 {
		logger.Printf("printer bytes:%v strings:%v\n", newPrinter, string(repByte))
		return nil, fmt.Errorf("DCloude print error:%v\n", newPrinter)
	}
	return &newPrinter, nil
}

func MapListToStruct(tmps []erps.Params) []erps.TmpBackUp {
	backKeys := models.TP_printer["tmp_backup"]
	return utils.SliceBackUpWithMapZip(erps.TmpBackUp{}, tmps, backKeys, "tmp_backup", "design_no")
}
