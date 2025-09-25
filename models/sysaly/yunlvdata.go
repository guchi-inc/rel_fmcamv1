package sysaly

//接收绿云 人员数据
type YunLvPostData struct {
	HotelId            int64  `json:"hotelId"`
	MasterId           int64  `json:"masterId"`
	Identifier         int64  `json:"Identifier"`
	GrpAccnt           int64  `json:"grpAccnt"`
	LinkId             int64  `json:"linkId"`
	Timestamp          int64  `json:"timestamp"`
	HotelGroupCode     string `json:"hotelGroupCode"`
	HotelGroupDescript string `json:"hotelGroupDescript"`
	HotelCode          string `json:"hotelCode"`
	HotelDescript      string `json:"hotelDescript"`
	BizType            string `json:"bizType"`
	Name               string `json:"name"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Language           string `json:"language"`
	Sex                string `json:"sex"`
	Vip                string `json:"vip"`
	IdCode             string `json:"idCode"`
	IdNo               string `json:"idNo"`
	Arr                string `json:"arr"`
	Dep                string `json:"dep"`
	Rmno               string `json:"rmno"`
	Mobile             string `json:"mobile"`
	Email              string `json:"email"`
	RsvClass           string `json:"rsvClass"`
	RsvMan             string `json:"rsvMan"`
	CrsNo              string `json:"crsNo"`
	RsvNo              string `json:"rsvNo"`
	RealRate           string `json:"realRate"`
	Date               string `json:"date"`
	Pay                string `json:"pay"`
	OldRmno            string `json:"oldRmno"`
	Memberno           string `json:"memberno"`
	Facepic            string `json:"facepic"`
	Scanpic            string `json:"scanpic"`
	Saleman            string `json:"saleman"`
	Remark             string `json:"remark"`
	CreditMan          string `json:"creditMan"`
	BuildingNo         string `json:"buildingNo"`
	FloorNo            string `json:"floorNo"`
	OldBuildingNo      string `json:"oldBuildingNo"`
	OldFloorNo         string `json:"oldFloorNo"`
}
