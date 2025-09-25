// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package sysfiles

import (
	"fmcam/common/databases"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type ObsFile struct{}

// 上传文件
func (S *ObsFile) UploadFile(fileName string, file *multipart.File) (*minio.UploadInfo, error) {
	return ObsApi.Upload(fileName, file)
}

/*
Get(fileName string) (string, error) //Base64
	Delete(fileName string) error
	List(prefix string) <-chan minio.ObjectInfo //指定文件列表
	Share(fileName string) (string, error)      //获取一次性访问地址
*/
// 获取文件（Base64）
func (S *ObsFile) GetFile(fileName string) (string, error) {
	return ObsApi.Get(fileName)
}

// 删除文件 从mysql 删除信息，从obs删除文件
func (S *ObsFile) DeleteFile(fileName string) error {

	DBConn := databases.DBMysql
	//查询  是否存在
	var exists int
	counts := fmt.Sprintf("SELECT id FROM fm_img WHERE img_parent = '%v'; ", fileName)
	err := DBConn.Get(&exists, counts)
	if err != nil {
		return ObsApi.Delete(fileName)
	}
	logger.Printf("fm_img img_parent:%v stime:  \n", fileName)

	imgUpdateSql := "DELETE FROM fm_img  WHERE id = ?"
	rst, err := DBConn.Exec(imgUpdateSql, exists)
	logger.Printf("fm_img delete img_parent:%v sql:%v rst:%#v err:%v \n", fileName, imgUpdateSql, rst, err)

	if err != nil {
		return ObsApi.Delete(fileName)
	}

	return ObsApi.Delete(fileName)
}

// 列出所有文件
func (S *ObsFile) ListFiles(prefix string) ([]map[string]interface{}, error) {
	return ObsApi.List(prefix)
}

// 单个文件的直接访问链接
func (S *ObsFile) ShareFiles(fileName string) (map[string]string, error) {
	return ObsApi.Share(fileName)
}
