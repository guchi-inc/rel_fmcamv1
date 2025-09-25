package streams

import (
	"github.com/gin-gonic/gin"
)

// 用户资源查询
func (s *StreamRouter) InitFileObs(r *gin.RouterGroup) {

	// rstpCtl := ApiStreamCtl.DevStreamApis.FRouter
	// 路由 rtsp flv
	// menus := r.Group("/img")
	// // menus.Use(middleware.AuthorizeJWT())
	// menus.GET("/url/search", rstpCtl.StreamUrl)
	// menus.POST("/register", rstpCtl.InsertFace)
	// menus.POST("/compare", rstpCtl.CompareFace)

	streamCtl := ApiStreamCtl.DevStreamApis.StrRouter
	//视频设备管理
	stream := r.Group("/stream")
	// stream.Use(middleware.AuthorizeJWT())
	stream.GET("/device/query", streamCtl.StreamQuery)
	stream.POST("/device/new", streamCtl.StreamNew)
	stream.POST("/device/update", streamCtl.StreamUpdate)

	//TODO...
	stream.GET("/log/query", streamCtl.CheckLogSelect)
	stream.POST("/log/new", streamCtl.CheckLogNew)
	stream.POST("/log/update", streamCtl.CheckLogUpdate)

	//文件形式 form 表单 图形保存
	headers := r.Group("/header")
	//图片文件上传
	headers.POST("/file/upload", streamCtl.FileUploadBlob)
	//编码形式 json 格式图形保存
	headers.POST("/file/upload/img", streamCtl.FileUploadB64)
	headers.POST("/face/img", streamCtl.FileFaceB64)

	//查询存储桶文件
	obs := r.Group("/obs")
	obs.GET("/files/list", streamCtl.ListFiles)
	obs.GET("/files/name", streamCtl.GetBussFile)
	obs.GET("/files/del", streamCtl.DeleteBussFile)
	//1天期分享链接
	obs.GET("/files/share", streamCtl.ShareFiles)

}
