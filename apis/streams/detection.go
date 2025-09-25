package streams

import (
	"bytes"
	"encoding/json"
	"fmcam/ctrlib/helpers"
	"fmcam/models/code"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type faceinfo struct{}

/*
检测图片是否存在面部信息
*/
var (
	Collection   = "stream_collection"
	CollectionDB = "stream_db"
	VectorDim    = 128
)

// --- 处理接口 ---

type RegisterRequest struct {
	ImageBase64 string `json:"image_base64"`
	UserID      string `json:"user_id"`
}

type CompareRequest struct {
	ImageBase64 string `json:"image_base64"`
}

// --- FRS 检测 & 向量生成 ---

// 定义 MediaPipe 响应格式
type mediaPipeResponse struct {
	Vector []float32 `json:"vector"`
}

// 综合用户信息和面部数据 写入mysql和faiss，关联id
func (fi *faceinfo) InsertFace(ctx *gin.Context) {

}

// 根据 id 查 flv播放地址信息
func (fi *faceinfo) StreamUrl(ctx *gin.Context) {

	camId := ctx.DefaultQuery("page_id", "")
	apilog.Printf("receive  camId:%#v \n", camId)

	user, err := UserService.GetUserByCtxToken(ctx)
	if err != nil {
		helpers.JSONs(ctx, code.AuthorizationError, gin.H{"error": err, "message": code.ZhCNText[code.AuthorizationError]})
		return
	}

	if camId == "" {
		helpers.JSONs(ctx, code.ParamError, gin.H{"error": "page cam id error", "message": code.ZhCNText[code.ParamError]})
		return
	}

	datas, err := SysStreamGroup.StreamSelect(1, 10, user.TenantId, camId, "")
	apilog.Printf(" cam info datas:%v, err:%v \n", datas, err)

	if err != nil {
		helpers.JSONs(ctx, code.NullData, gin.H{"error": err, "message": code.ZhCNText[code.NullData]})
		return
	}

	helpers.JSONs(ctx, code.Success, datas)

}

// 寻找和匹配面部数据
func (fi *faceinfo) CompareFace(ctx *gin.Context) {

	var vector []float32
	err := ctx.ShouldBind(&vector)
	if err != nil {
		helpers.JSONs(ctx, code.ParamBindError, err)
	}
	searchRst, searchCode, err := fi.SearchFace(vector)
	apilog.Println("searchRst, searchCode, err:", searchRst, searchCode, err)

	helpers.JSONs(ctx, code.Success, searchRst)
}

func (fi *faceinfo) InsertVectorToFaiss(vector []float32, userId string) error {
	payload := map[string]interface{}{
		"user_id": userId,
		"vector":  vector,
	}
	jsonBody, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:5001/faiss/add", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FAISS 添加失败，状态码 %d", resp.StatusCode)
	}
	return nil
}

// 调用ai服务
func (fi *faceinfo) SearchFace(vector []float32) (string, float32, error) {
	payload := map[string]interface{}{
		"vector": vector,
	}
	jsonBody, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:5001/faiss/search", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", 0, err
	}

	if res["user_id"] == nil {
		return "", 0, fmt.Errorf("未找到匹配人脸")
	}

	return res["user_id"].(string), float32(res["score"].(float64)), nil
}
