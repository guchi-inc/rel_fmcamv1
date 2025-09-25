package streams

import (
	"fmcam/ctrlib/helpers"
	"fmcam/ctrlib/utils"
	"fmcam/models/code"
	"fmcam/models/erps"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 人员信息面部信息图，检查和创建路径
func (st *DeviceTask) FileFaceB64(c *gin.Context) {

	var req erps.ImageUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.JSONs(c, code.ParamBindError, gin.H{"error": err})
		return
	}

	//人员信息图名称
	if req.ImgParent == "" {
		helpers.JSONs(c, code.ParamError, gin.H{"error": code.ZhCNText[code.ParamError] + "图名不能为空"})
	}
	req.ImgParent = "profile_" + utils.PinMain(req.ImgParent)
	//图像处理
	header, file, err := OBSFile.ImgDealBase64(c, &req)
	apilog.Printf("图像处理后 header:%#v info:%#v err:%#v \n ", header, file != nil, err)
	if err != nil {
		helpers.JSONs(c, code.Failed, gin.H{"error": fmt.Errorf("文件处理失败 %v", err)})
		return
	}

	//对象存储服务
	info, err := SysObsFileGroup.UploadFile(header.Filename, file)
	apilog.Println("info: ", info != nil, "err:", err)
	if err != nil {
		helpers.JSONs(c, code.Failed, gin.H{"error": fmt.Errorf("文件上传失败 %v", err)})
		return
	}

	//图 访问链接 7 天有效
	imgShareUrl, err := SysObsFileGroup.ShareFiles(header.Filename)
	apilog.Println("get buss share url img_parent:", header.Filename, err)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件链接查询失败" + err.Error()})
		return
	}

	helpers.JSONs(c, code.Success, gin.H{"message": info.Key, "data": imgShareUrl})
}

// 上传 系统文件，检查和创建路径
func (st *DeviceTask) FileUploadB64(c *gin.Context) {

	var req erps.ImageUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.JSONs(c, code.ParamBindError, gin.H{"error": err})
		return
	}

	//图像处理
	header, file, err := OBSFile.ImgDealBase64(c, &req)
	apilog.Printf("图像处理后 header:%#v info:%#v err:%#v \n ", header, file != nil, err)
	if err != nil {
		helpers.JSONs(c, code.Failed, gin.H{"error": fmt.Errorf("文件处理失败 %v", err)})
		return
	}

	//对象存储服务
	info, err := SysObsFileGroup.UploadFile(header.Filename, file)
	apilog.Println("info: ", info != nil, "err:", err)
	if err != nil {
		helpers.JSONs(c, code.Failed, gin.H{"error": fmt.Errorf("文件上传失败 %v", err)})
		return
	}

	helpers.JSONs(c, code.Success, gin.H{"message": "success", "data": info})
}

// 上传 系统文件，检查和创建路径
func (st *DeviceTask) FileUploadBlob(ctx *gin.Context) {
	// 上传单个文件
	// 测试 curl -X POST http://localhost:1818/v1/stream//file/update   -F "file=./idCard1.jpg"  -H "Content-Type: multipart/form-data"
	// 单文件

	// 单文件，每次post的接口调用都重置路径为当天路径

	//名称  bright, dark
	file_parent := ctx.DefaultPostForm("img_parent", "")
	//序号
	time_nums := ctx.DefaultPostForm("time_num", "")

	//必传 img 名称
	if file_parent == "" && time_nums == "" {
		helpers.JSONs(ctx, http.StatusBadRequest, fmt.Errorf("param parent and order_number must be needed"))
		return
	}
	apilog.Printf("img_parent: %v save file to order_number:%v \n", file_parent, time_nums)

	b, err := SysFileMan.UpFileBlob(ctx, file_parent, time_nums)
	if err != nil {
		helpers.JSONs(ctx, http.StatusInternalServerError, err)
		return
	}

	helpers.JSONs(ctx, http.StatusOK, gin.H{"data": b, "message": fmt.Sprintf("'%s' uploaded at %v!", file_parent, b)})
}

// 上传文件
func (BA *DeviceTask) UploadBussFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件参数错误收取失败"})
		return
	}
	defer file.Close()
	info, err := SysObsFileGroup.UploadFile(header.Filename, &file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"message": "success", "data": info})

}

// 获取文件（Base64）
func (BA *DeviceTask) GetBussFile(c *gin.Context) {
	filename := c.DefaultQuery("img_parent", "")
	fileBase64, err := SysObsFileGroup.GetFile(filename)
	apilog.Println("get buss file img_parent:", filename, "fileBase64:", fileBase64, "err:", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件查询失败"})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": fileBase64})

}

// 删除文件
func (BA *DeviceTask) DeleteBussFile(c *gin.Context) {
	filename := c.DefaultQuery("img_parent", "")

	err := SysObsFileGroup.DeleteFile(filename)
	if err != nil {
		helpers.JSONs(c, code.CacheNotExist, gin.H{"error": fmt.Sprintf("失败:%v", err)})
		return
	}

	helpers.JSONs(c, code.Success, gin.H{"message": "文件删除成功", "data": ""})
}

// 列出所有文件
func (BA *DeviceTask) ListFiles(c *gin.Context) {
	filename := c.DefaultQuery("img_parent", "")
	fileList, err := SysObsFileGroup.ListFiles(filename)
	apilog.Println("get file by filePrefix:", filename, "fileList:", len(fileList), "err:", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件链接查询失败"})
		return
	}
	helpers.JSONs(c, code.Success, gin.H{"data": fileList})

}

// 单个文件的直接访问链接
func (BA *DeviceTask) ShareFiles(c *gin.Context) {
	filename := c.DefaultQuery("img_parent", "")
	fileShare, err := SysObsFileGroup.ShareFiles(filename)
	apilog.Println("get buss share url img_parent:", filename, err)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件链接查询失败"})
		return
	}

	helpers.JSONs(c, code.Success, gin.H{"data": fileShare})
}
