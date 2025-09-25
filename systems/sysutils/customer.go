package sysutils

import (
	"fmcam/common/configs"
	"fmcam/common/databases"
	"fmcam/ctrlib/utils/timeutil"
	"fmcam/models/seclients"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CustomerUtil struct{}

// 获取后台管理 [租户] 列表
func (mt *CustomerUtil) GetCustomerList(Page, PageSize int, supplier, Telephone, Contacts, fullAddress, StartAt, EndAt, TenantId string) (*seclients.TenantsPage, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = BaseTenantSql

		total      int
		useTenants = seclients.TenantsPage{}
		DBConn     = databases.DBMysql
	)

	//条件 supplier, province, city, area, street, fullAddress, TenantId
	Filters := " WHERE delete_flag != '1' "

	if supplier != "" {
		Filters += " AND supplier LIKE '" + supplier + "%' "
	}
	if Telephone != "" {
		Filters += " AND telephone LIKE '" + Telephone + "%' "
	}
	if fullAddress != "" {
		Filters += " AND full_address LIKE '%" + fullAddress + "%' "
	}
	if TenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}
	if Contacts != "" {
		Filters += " AND contacts LIKE '%" + Contacts + "%' "
	}

	if StartAt != "" && EndAt != "" {
		Filters += fmt.Sprintf(" AND (created_time between '%v' AND '%v' )", StartAt, EndAt)
	}
	totalSql := "SELECT COUNT(*) AS Total FROM  `fm_tenant`  " + Filters

	res := DBConn.QueryRow(totalSql)
	err := res.Scan(&total)
	if err != nil {
		logsys.Printf(" err:%#v totalSql:%v\n", err, totalSql)

		return &useTenants, err
	}

	Limits := fmt.Sprintf(" ORDER BY id ASC LIMIT %v OFFSET %v ", PageSize, (Page-1)*PageSize)

	baseSql += Filters
	baseSql += Limits
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err = DBConn.Select(&useTenants.Data, baseSql)
	logsys.Printf("from limits select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	useTenants.Total = total
	useTenants.PageSize = PageSize
	useTenants.Page = Page
	if len(useTenants.Data) > 0 {
		cm, _ := databases.FmGlobalMap("fm_tenant", nil)
		useTenants.Columns = cm
	}

	return &useTenants, nil
}

// 后台管理 [租户] 新增
func (mt *CustomerUtil) NewCustomer(login_name string, suppliers *seclients.TenantFull) (int64, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		DBConn = databases.DBMysql
	)

	tn := time.Now().Local().In(configs.Loc)
	ts := timeutil.CSTLayoutString(&tn)

	suppliers.Creator = login_name
	suppliers.CreatedAt = ts

	// 插入 SQL
	newuuid := uuid.New().String()
	query := fmt.Sprintf(`
		INSERT INTO fm_tenant (supplier,contacts,tenant_id, type, description, creator) VALUES ('%v', '%v', %v, '%v', '%v', '%v')
		`, suppliers.Supplier, suppliers.Contacts, fmt.Sprintf("(UUID_TO_BIN('%v'))", newuuid), suppliers.Type, suppliers.Description, suppliers.Creator)

	//写入 执行
	rst, err := DBConn.Exec(query)
	logsys.Printf("new fm_tenant:%v insert:%#v, with db rst:%#v err:%v\n", newuuid, query, rst, err)

	if err != nil {
		logsys.Println("插入失败: ", err)
		return 0, err
	}

	return rst.LastInsertId()

}

// 获取后台管理 [租户] 列表
func (mt *CustomerUtil) GetCustomerInfo(tid int, supplier, Telephone, TenantId string) (*seclients.TenantFull, error) {

	// 查询单个数据 You can also get a single result, a la QueryRow
	var (
		baseSql = `SELECT id,
		supplier,
		contacts,
		email,
		BIN_TO_UUID(tenant_id) as tenant_id,
		type,
		province,
		city,
		area,
		street,
		address,
		addr_code,
		full_address,
		fax,
		phone_num,
		telephone,
		tax_num,
		bank_name,
		account_number,
		sort,
		description,
		enabled,
		delete_flag,
		isystem,
		advance_in,
		begin_need_get,
		begin_need_pay,
		all_need_get,
		all_need_pay,
		tax_rate,
		creator,
		created_time,
		updated_time FROM fm_tenant `

		useTenants = seclients.TenantFull{}
		DBConn     = databases.DBMysql
	)

	//条件 supplier, province, city, area, street, fullAddress, TenantId
	Filters := " WHERE delete_flag != '1' "
	if tid != 0 {
		Filters += fmt.Sprintf(" AND id = '%v' ", tid)
	}
	if supplier != "" {
		Filters += " AND supplier LIKE '" + supplier + "%' "
	}
	if Telephone != "" {
		Filters += " AND telephone LIKE '" + Telephone + "%' "
	}
	if TenantId != "" {
		Filters += fmt.Sprintf(" AND BIN_TO_UUID(tenant_id)  = '%v' ", uuid.MustParse(TenantId))
	}

	baseSql += Filters
	// 筛选菜单名称列表,先查所在权限组，然后查菜单列表
	err := DBConn.Get(&useTenants, baseSql)
	logsys.Printf("from limits select:%#v  err:%v \n", baseSql, err)
	if err != nil {
		return nil, err
	}

	return &useTenants, nil
}
