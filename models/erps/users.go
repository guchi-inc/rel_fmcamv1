// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package erps

type userRepository struct {
	// DB *sqlx.DB
}

type UserRepository interface {
	GetUser(int) (*FmUser, error)
	// GetByName(string) (User, error)
	// AddUser(User) (User, error)
	GetByEmail(string) (*FmUser, error)
	// GetAllUser() ([]User, error)
	// UpdateUser(User) (User, error)
	// DeleteUser(User) (User, error)
	// Migrate() error
}

// 部门列表
type GcDepartmentStructList struct {
	Total   uint                 `json:"total" db:"total" form:"total"`
	Size    uint                 `json:"size" db:"size" form:"size"`
	Data    []GcDepartmentStruct `json:"data" db:"data" form:"data"`
	Columns []GcDesc             `json:"columns" db:"columns"`
}

//用户单人信息 描述
type UserCardInfo struct {
	T       uint     `json:"total" db:"total"`
	S       uint     `json:"size" db:"size"`
	Data    []FmUser `json:"data" db:"data"`
	Columns []GcDesc `json:"columns" db:"columns"`
}

type UserList struct {
	T    uint      `json:"total" db:"total"`
	S    uint      `json:"size" db:"size"`
	Data []*FmUser `json:"data" db:"data"`
}

type UserListStruct struct {
	T       uint     `json:"total" db:"total"`
	S       uint     `json:"size" db:"size"`
	Data    []FmUser `json:"data" db:"data"`
	Columns []GcDesc `json:"columns" db:"columns"`
}

// 订阅消息
type MsgSubscribe struct {
	Type      string   `json:"type"`
	SubName   string   `json:"sub_name"`
	TableList []string `json:"table_list"`
}

// 订阅列表
type MsgSubscribeList struct {
	Total     uint           `json:"total"`
	LoginName string         `json:"login_name"`
	Data      []MsgSubscribe `json:"data"`
}

// 重置消息
type MsgReback struct {
	PhoneNum string `json:"phonenum"`
	Email    string `json:"email"`
	Codes    string `json:"code"`
	Password string `json:"password"`
}
