package sysfiles

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils/timeutil"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FilesDB struct{}

// 上传 系统文件接口页
func (that *FilesDB) UpFileBlob(ctx *gin.Context, file_parent, time_nums string) (string, error) {

	b, err := that.upSysFileBlob(ctx, file_parent, time_nums)
	if err != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return "", err
	}
	return b, nil
}

// 上传 系统文件接口页
func (that *FilesDB) upSysFileBlob(ctx *gin.Context, img_name, time_num string) (string, error) {

	//图像名称 bright, dark
	// 接收文件
	file, header, err := ctx.Request.FormFile("file") // ctx.FormFile("file")  返回 header 和 err

	newFileName := header.Filename
	logger.Printf("get file size:%v, filePath:%v \n", header.Size, newFileName)

	if err != nil {
		return "", fmt.Errorf("file_path not exist err:%v", err)
	}

	//文件名替换为 fparent
	var storeFileName string = fmt.Sprintf("%v%v", time.Now().Day(), rand.IntN(10))
	if strings.Contains(newFileName, ".") {
		newFileNameStrings := strings.Split(newFileName, ".")
		if !strings.Contains(img_name, ".") {
			storeFileName += fmt.Sprintf("%v.%v", img_name, newFileNameStrings[1])
		} else {
			storeFileName = img_name
		}
	} else {
		storeFileName += fmt.Sprintf("%v.%v", img_name, "png")
	}

	// //路径存取
	// err = ctx.SaveUploadedFile(header, filepath.Join(file_path, storeFileName))
	// if err != nil {
	// 	return "", err
	// }

	//如果文件名为空 .png 使用订单号和日期做为文件名
	if len(storeFileName) == 4 {
		storeFileName = fmt.Sprintf("%v%v%v", time_num, time.Now().Day(), rand.IntN(10)) + storeFileName
	}
	// 检查文件名称，上传文件至指定目录
	if t, err := CheckNameVaild(storeFileName); !t {
		return "", fmt.Errorf("名称非法%v", err)
	}

	//存数据库 二进制
	// Read the file's content
	// fileBytes, err := io.ReadAll(file)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
	// 	return "", err
	// }

	fmt.Println("file size:", header.Size)
	if header.Size > 350*1024*1024 {
		return "", fmt.Errorf("文件不得超过3M")
	}

	//路径 obs 存取
	uploadInfo, err := MinFileObs.UploadFile(storeFileName, &file)
	logger.Printf("upload file:%v to %v file uploadInfo:%#v \n", header.Filename, storeFileName, uploadInfo)
	if err != nil {
		return "", err
	}

	//用户信息
	user, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusNonAuthoritativeInfo, gin.H{"error": err})
		return "", err
	}

	login_name := user.LoginName //"admin" //
	fileBytes := []byte{}        //mysql数据库不存储用户上传的文件。只存储预装文件。
	_, err = that.SysImgInsert(login_name, storeFileName, time_num, uploadInfo.Bucket, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, gin.H{"error": err})
		return "", err
	}
	return storeFileName, err
}

// 系统相关图像 写入
func (fr *FilesDB) SysImgInsert(creator, imgName, order_number, file_path string, fileBytes []byte) (int64, error) {

	DBConn := databases.DBMysql
	//存png名称
	//查询 图像名称信息是否存在
	var existImgId int64
	typeIdSql := fmt.Sprintf("SELECT id FROM fm_sys_img WHERE img_parent = '%v'; ", imgName)
	err := DBConn.Get(&existImgId, typeIdSql)

	if err != nil || existImgId == 0 {
		fmt.Printf("系统图像信息:%v 不存在. err:%v ", imgName, err)
	}

	//如果已经存在,则更新,否在新增
	tNow := time.Now().In(configs.Loc).Local()
	stime := fmt.Sprintf("%v", timeutil.CSTLayoutString(&tNow)) //.Format(timeutil.CSTLayout)

	fmt.Printf("fm_img imgName:%v stime:%v \n", imgName, stime)

	if existImgId == 0 {
		imgSql := "INSERT INTO fm_sys_img (img_parent,order_number, local_path, picture, enabled, creator, created_time) VALUES (?,?,?,?, ?,?,?)"
		rst, err := DBConn.Exec(imgSql, imgName, order_number, file_path, fileBytes, "1", creator, stime)
		fmt.Printf("fm_sys_img new insert name:%v sql:%v err:%v \n", imgName, imgSql, err)

		if err != nil {
			return 0, err
		}

		return rst.LastInsertId()
	}

	imgUpdateSql := "Update fm_sys_img SET img_parent = ?, order_number = ?, local_path = ?, picture = ?, creator = ?, updated_time = ? WHERE id = ?"
	rst, err := DBConn.Exec(imgUpdateSql, imgName, order_number, file_path, fileBytes, creator, stime, existImgId)
	fmt.Printf("fm_sys_img update imgname:%v sql:%v rst:%#v err:%v \n", imgName, imgUpdateSql, rst, err)

	if err != nil {
		return 0, err
	}

	return int64(existImgId), nil
}
