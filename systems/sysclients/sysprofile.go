package sysclients

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils"
	"fmcam/models/erps"
	"fmcam/models/seclients"
	"fmcam/models/sysaly"
	"fmcam/systems/genclients"
	"fmcam/systems/genclients/profiles"
	"fmcam/systems/sysutils"
	"fmt"
	"io"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 新增 人员信息
func (oc *ORMCleint) ProfileNew(ctx *gin.Context, PParm *seclients.ProfileInfos) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	Ctx := databases.WithOperatorFromGin(ctx)

	// 3. 存在 => 执行更新
	build := client.Profiles.Create()

	if PParm.TypeID != nil {
		build.SetTypeID(*PParm.TypeID)
	}
	if PParm.Name != nil {
		build.SetName(*PParm.Name)
	}
	if PParm.PhoneNumber != nil {
		build.SetPhoneNumber(*PParm.PhoneNumber)
	}
	if PParm.IDCardNumber != nil {
		build.SetIDCardNumber(*PParm.IDCardNumber)
	}
	if PParm.RoomID != nil {
		build.SetRoomID(*PParm.RoomID)
	}
	if PParm.TmpURL != nil {
		build.SetTmpURL(*PParm.TmpURL)
	}
	if PParm.Enabled != nil {
		build.SetEnabled(*PParm.Enabled)
	}

	upProfile, err := build.Save(Ctx)
	sysdebug.Printf("new Profile:%#v err:%#v \n", upProfile, err)
	if err != nil {
		return 0, err
	}
	if PParm.TenantID != nil {
		upSql := fmt.Sprintf(` UPDATE Profile Set tenant_id = %v `, fmt.Sprintf("(UUID_TO_BIN('%v'))", uuid.MustParse(*PParm.TenantID)))
		_, err := DBConn.Exec(upSql)
		if err != nil {
			return 0, err
		}
	}

	return upProfile.ID, nil

}

// 更新 人员信息
func (oc *ORMCleint) ProfileUpdate(ctx *gin.Context, PParm *seclients.ProfileInfos) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	Ctx := databases.WithOperatorFromGin(ctx)
	// 1. 查找是否存在
	existing, err := client.Profiles.
		Query().
		Where(profiles.ID(PParm.ID)).
		Order(profiles.ByID(entsql.OrderDesc())).
		Only(Ctx)
	sysdebug.Printf("query Profiles check in existing?:%#v err:%#v \n", existing, err)

	if err != nil {
		return 0, err
	}
	if existing != nil && existing.ID != 0 {

		//查 人员类型 guest id
		if PParm.DeleteFlag != nil && *PParm.DeleteFlag == "1" {
			DelSql := fmt.Sprintf(" DELETE FROM Profile WHERE id = '%v' ", PParm.ID)
			rst, err := DBConn.Exec(DelSql)
			sysdebug.Printf(" Profile  DelSql:%#v err:%#v \n", DelSql, err)
			if err != nil {
				err := fmt.Errorf("del Profile:%v fail:%v", DelSql, err)
				return 0, err
			}
			return rst.LastInsertId()
		}

		// 3. 存在 => 执行更新
		updater := client.Profiles.UpdateOneID(PParm.ID)

		if PParm.TypeID != nil {
			updater.SetTypeID(*PParm.TypeID)
		}
		if PParm.RoomID != nil {
			updater.SetRoomID(*PParm.RoomID)
		}
		if PParm.TmpURL != nil {
			updater.SetTmpURL(*PParm.TmpURL)
		}
		if PParm.Enabled != nil {
			updater.SetEnabled(*PParm.Enabled)
		}

		upProfile, err := updater.Save(ctx)
		sysdebug.Printf("upProfile:%#v err:%#v \n", upProfile, err)
	}
	return 1, nil

}

// 仅用于 住客 信息解码
func (uc *SysCleint) ProfileJsonParam(ctx *gin.Context) (*seclients.CheckoutProfile, error) {

	var (
		profileUser seclients.CheckoutProfile
	)

	//接收json参数
	var body []byte
	body, err := io.ReadAll(ctx.Request.Body)
	sysdebug.Printf("before dencry byte:%v \n", string(body))
	if err != nil {
		sysdebug.Printf("ioutil.ReadAll:%v failed:%v", body, err)
		return nil, err
	}

	//解密
	if configs.PostEncrypt {
		// body = utils.XDencryptByte(body, []byte(configs.AccessKey))
		strBody, err := utils.NewXDecryptString(string(body), configs.AccessKey)
		sysdebug.Printf("dencry strBody:%v err:%v", strBody, err)

		if strBody == "" {
			sysdebug.Printf("dencry str:%v strBody:%v", string(body), string(strBody))
			err = fmt.Errorf("after dencry:%v map error:%v", string(strBody), err)
			return nil, err
		}
		body = []byte(strBody)

		sysdebug.Println("dencry byte str:", string(body))
	}

	err = json.Unmarshal(body, &profileUser)
	if err != nil {
		err := fmt.Errorf("after dencry byte:%v to map error:%v", string(body), err)
		return nil, err
	}

	if profileUser.LoginName != nil || profileUser.LoginPhonenum != nil || profileUser.IdCardGuest != nil {
		return &profileUser, nil
	}

	return nil, fmt.Errorf("关键数据为空")

}

// 更新 人员签入信息
func (oc *SysCleint) ProfileCheckin(ctx *gin.Context, fuser *erps.FmUser, PParm *seclients.CheckoutProfile) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	Ctx := databases.WithOperatorFromGin(ctx)
	// 1. 查找是否存在
	existing, err := client.Profiles.
		Query().
		Where(profiles.IDCardNumber(*PParm.IdCardGuest)).
		Order(profiles.ByID(entsql.OrderDesc())).
		Only(Ctx)
	sysdebug.Printf("query Profiles check in existing?:%#v err:%#v \n", existing, err)

	if existing == nil || existing.ID == 0 {

		//查 人员类型 guest id
		guestId := genclients.ProfileType{}
		err := DBConn.Get(&guestId, "SELECT id,type_name,type_code FROM ProfileType WHERE type_code = 'guest' ")
		sysdebug.Printf("query ProfileType guestId:%#v err:%#v \n", guestId, err)

		//新增陌生人 离店
		accouTenant := fuser.TenantId
		baseSql := fmt.Sprintf(`INSERT INTO Profile (name, tenant_id,type_id,phone_number,id_card_number,enabled,room_id,tmp_url) VALUES ('%v', %v, %v, '%v', '%v', %v, '%v', '%v')
			`, guestId.TypeName, fmt.Sprintf("(UUID_TO_BIN('%v'))", accouTenant), guestId.ID, PParm.PhoneGuest, PParm.IdCardGuest, true, PParm.RoomId, PParm.TmpUrl)
		rst, errec := DBConn.Exec(baseSql)
		sysdebug.Printf("new Profile:%#v err:%#v \n", baseSql, errec)
		if errec != nil {

			return 0, errec
		}
		return rst.LastInsertId()
	}

	// 3. 存在 => 执行更新
	upProfile, err := client.Debug().Profiles.
		UpdateOne(existing).
		SetRoomID(*PParm.RoomId).
		SetTmpURL(*PParm.TmpUrl).
		SetEnabled(true). // ...更新字段
		Save(ctx)
	sysdebug.Printf("upProfile:%#v err:%#v \n", upProfile, err)
	if err != nil {
		return 0, err
	}
	return upProfile.ID, nil

}

// 更新 人员签出 离店信息
func (oc *SysCleint) ProfileCheckout(ctx *gin.Context, fuser *erps.FmUser, PParm *seclients.CheckoutProfile) (int64, error) {

	var (
		client = databases.EntClient
		DBConn = databases.DBMysql
	)

	Ctx := databases.WithOperatorFromGin(ctx)
	// 1. 查找是否存在
	existing, err := client.Profiles.
		Query().
		Where(profiles.IDCardNumber(*PParm.IdCardGuest)).
		Order(profiles.ByID(entsql.OrderDesc())).
		Only(Ctx)

	sysdebug.Printf("query  check out existing:%#v err:%#v \n", existing, err)

	if existing == nil || existing.ID == 0 {

		//查 人员类型 guest id
		guestId := genclients.ProfileType{}
		getTenant := fmt.Sprintf("SELECT id,type_name,type_code FROM ProfileType WHERE type_code = 'guest' AND BIN_TO_UUID(tenant_id) = '%v'", configs.DefTenantId)
		err := DBConn.Get(&guestId, getTenant)
		sysdebug.Printf("query ProfileType guestId:%#v err:%#v \n", guestId, err)

		//新增陌生人 离店
		accouTenant := fuser.TenantId
		baseSql := fmt.Sprintf(`INSERT INTO Profile (name, tenant_id,type_id,phone_number,id_card_number,enabled) VALUES ('%v', %v, %v, '%v', '%v', %v)
			`, guestId.TypeName, fmt.Sprintf("(UUID_TO_BIN('%v'))", accouTenant), guestId.ID, PParm.PhoneGuest, PParm.IdCardGuest, false)
		rst, errec := DBConn.Exec(baseSql)
		sysdebug.Printf("new Profile:%#v err:%#v \n", baseSql, errec)
		if errec != nil {
			return 0, errec
		}
		return rst.LastInsertId()
	}

	// 3. 存在 => 执行更新
	upProfile, err := client.Debug().Profiles.
		UpdateOne(existing).
		SetEnabled(false). // ...更新字段
		Save(ctx)
	sysdebug.Printf("upProfile:%#v err:%#v \n", upProfile, err)
	if err != nil {
		return 0, err
	}

	return upProfile.ID, nil

}

// 绿云 数据记录和处理
func (oc *SysCleint) YunLvPostDeal(ctx *gin.Context, req *sysaly.YunLvPostData) error {

	var (
		client = databases.EntClient
		// DBConn = databases.DBMysql
	)
	bizType := req.BizType // req["bizType"].(string)
	idNo := req.IdNo       // req["idNo"].(string)
	sysdebug.Printf("param all:%T bizType:%v idno:%v", req.IdCode, bizType, idNo)
	switch bizType {
	case "CHECKIN", "FACECHG":
		return YunLvProfile(ctx, req)
	case "CHECKOUT":
		//更新已有的
		p, err := client.Debug().Profiles.Query().Where(profiles.IDCardNumber(idNo)).Order(genclients.Desc("created_time")).First(ctx)
		sysdebug.Printf("param all:%#v bizType:%v idno:%v", req.HotelDescript, bizType, idNo)

		if err == nil {
			return client.Debug().Profiles.UpdateOne(p).SetEnabled(false).Exec(ctx)
		} else {
			//新增一个签出人员记录
			return YunLvProfile(ctx, req)
		}
	default:
		// do nothing
		sysdebug.Printf("receive other bizType:%#v \n", bizType)
		return nil
	}

}

// 云绿检查登记
func YunLvProfile(ctx *gin.Context, req *sysaly.YunLvPostData) error {

	var (
		DBConn     = databases.DBMysql
		tenantInfo = erps.FmTenant{}
		customUt   = sysutils.CustomerUtil{}
	)

	/*
			A checkin流程
			查询是否存在相同的人员。
				修改已有相同id_card_number 人员enabled为false

		    checkin新增流程
			1 检查fm_tenant是否存在 type对应的 hotelCode

			根据hotelCode 关联 fm_tenant表的 type字段，
				如果没有该酒店信息，则自动新增一个酒店并关联新增的tenant_id到人员。
				如果type字段存在该酒店信息，返回 对应的tenant_id 用于新增人员.

			2 新增人员信息
				如果有 face 图像 或 scanpic 证件图像，存储图像并返回本地共享链接
		        新增人员信息到Profile表
	*/

	//查是否存在 type 对应 绿云的 hotelCode
	//读已提交
	zkId := genclients.ProfileType{}
	//使用默认的人员类型
	profileTypeSql := fmt.Sprintf("SELECT id,type_name,type_code FROM ProfileType WHERE type_name = '住客' AND  BIN_TO_UUID(tenant_id) = '%v' ", uuid.MustParse(configs.DefTenantId))
	err := DBConn.Get(&zkId, profileTypeSql)
	sysdebug.Printf("query ProfileType guestId:%#v err:%#v \n", zkId, err)
	if err != nil {
		return err
	}

	errTx := DBConn.MustExec("SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
	var logSteps []gin.H

	tx, err := DBConn.Beginx()
	if err != nil {
		return err
	}

	tenantSql := fmt.Sprintf("SELECT id,BIN_TO_UUID(tenant_id)  AS tenant_id  FROM fm_tenant WHERE type = '%v' order by id desc limit 1 ", req.HotelCode)
	err = tx.Get(&tenantInfo, tenantSql)
	sysdebug.Printf("已有租户 %v 已有租户号 %v 个 .err:%v \n", tenantInfo.Type, tenantInfo.TenantId, err)
	logSteps = append(logSteps, gin.H{"exist tenant": tenantInfo.TenantId})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sysdebug.Println("没有查询到数据，不算错误，可以继续逻辑")
		} else {
			return err
		}
	}

	//新增租户信息
	if tenantInfo.TenantId == "" {
		//系统新增数据 使用 admin 作为创建者 hotelCode
		newtid, err := customUt.NewCustomer("admin",
			&seclients.TenantFull{Enabled: true,
				Supplier:    req.HotelDescript,
				Description: req.HotelGroupDescript,
				Type:        req.HotelCode})

		newTenant := fmt.Sprintf("new tenant:%#v err:%#v \n", newtid, err)
		logSteps = append(logSteps, gin.H{"new tenant": newTenant})
		if err != nil {
			tx.Rollback()
			return err
		}
		tenantSql := fmt.Sprintf("SELECT BIN_TO_UUID(tenant_id) AS tenant_id FROM fm_tenant WHERE id = '%v'  ", newtid)
		err = tx.Get(&tenantInfo, tenantSql)
		if err != nil {
			return err
		}
	}

	var idIndetNumber string
	if req.IdNo != "" {
		idIndetNumber = req.IdNo
	} else {
		idIndetNumber = fmt.Sprintf("%v", req.MasterId)
	}

	//更新已有人员信息
	profileSql := fmt.Sprintf("UPDATE Profile SET enabled = '0' WHERE enabled = '1' AND id_card_number = '%v' ", idIndetNumber)
	rst, err := tx.Exec(profileSql)
	profileEdit := fmt.Sprintf("  已有记录  %v tenantSql %v .profileSql:%v err:%v errTx:%#v \n", tenantInfo, tenantSql, profileSql, err, errTx)
	logSteps = append(logSteps, gin.H{"UPDATE profile": profileEdit})
	if err != nil {
		tx.Rollback()
		return err
	}

	Affected, _ := rst.RowsAffected()
	sysdebug.Printf("已更新人员%v 已有记录 %v 个为非住店状态.err:%v \n", req.IdNo, Affected, err)

	//checkout 类型的信息变动，不新增人员记录
	if req.BizType == "CHECKOUT" {
		return nil
	}

	// 图像处理 使用证件号作为 图像文件名
	var base64Data string
	if req.Facepic != "" {
		//面部图
		base64Data = req.Facepic
	} else {
		//证件图
		base64Data = req.Scanpic
	}

	var shareUrl string
	if base64Data != "" {
		shareUrl, err = ImgFileDeal(ctx, &erps.ImageUploadRequest{ImgParent: idIndetNumber, Base64Data: base64Data})
		sysdebug.Printf("query img file shareUrl:%#v err:%#v \n", shareUrl, err)
		if err != nil {
			return err
		}
	}

	// 新增入住记录 包括图像信息链接
	accouTenant := tenantInfo.TenantId

	baseSql := fmt.Sprintf(`INSERT INTO Profile (name, tenant_id,type_id,phone_number,id_card_number,enabled,room_id,tmp_url) VALUES ('%v', %v, %v, '%v', '%v', %v, '%v', '%v')
		`, req.Name, fmt.Sprintf("(UUID_TO_BIN('%v'))", accouTenant), zkId.ID, req.Mobile, idIndetNumber, true, req.Rmno, shareUrl)
	rst, errec := tx.Exec(baseSql)
	insertProfile := fmt.Sprintf("new Profile:%#v err:%#v \n", baseSql, errec)
	logSteps = append(logSteps, gin.H{"new profile": insertProfile})
	if errec != nil {
		tx.Rollback()
		sysdebug.Printf("INSERT Profile:%#v errec:%#v \n", baseSql, errec)
		return errec
	}

	newCheckIn, errs := rst.LastInsertId()
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %v", err)
	}

	sysdebug.Printf("commit err:%#v newCheckIn:%v logSteps:%#v \n", err, newCheckIn, logSteps)
	return errs
}

// 返回 存储到 obs 后的分享链接url 和可能的err
func ImgFileDeal(c *gin.Context, reqImg *erps.ImageUploadRequest) (string, error) {

	//图像处理
	header, file, err := OBSFile.ImgDealBase64(c, reqImg)
	sysdebug.Printf("图像处理后 header:%#v info:%#v ImgDealBase64 erropr:%#v \n ", header != nil, file != nil, err)
	if err != nil {
		return "", fmt.Errorf("文件处理失败 %v", err)
	}

	//对象存储服务
	info, err := ObsApi.Upload(header.Filename, file)
	sysdebug.Println("info: ", info != nil, "err:", err)
	if err != nil {
		return "", fmt.Errorf("文件上传失败 %v", err)
	}

	//返回的键值对 img_parent  url
	shareUrl, err := ObsApi.Share(header.Filename)
	if err != nil {
		return "", fmt.Errorf("文件链接获取失败 %v", err)
	}
	return shareUrl["url"], nil
}
