package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// FmUserAccount schema 对应 fm_user 表
type FmUserAccount struct {
	ent.Schema
}

func (FmUserAccount) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_user"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (FmUserAccount) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("username").
			NotEmpty().
			MaxLen(25).Comment("姓名").
			StructTag(`json:"username"  db:"username"`),
		field.String("login_name").
			NotEmpty().
			MaxLen(50).
			Default("").Comment("登陆名").
			StructTag(`json:"login_name" db:"login_name"`),
		field.String("password").
			NotEmpty().
			MaxLen(50).
			Default("").Comment("登陆密码").
			StructTag(`json:"password" db:"password"`),
		field.String("leader_flag").
			Optional().Comment("领导标记").
			StructTag(`json:"leader_flag"  db:"leader_flag"`),
		field.String("position").
			Optional().Comment("职位").
			StructTag(`json:"position"  db:"position"`),
		field.String("department").
			Optional().Comment("所属部门").
			StructTag(`json:"department"  db:"department"`),
		field.String("email").
			Optional().Comment("电邮").
			StructTag(`json:"email"  db:"email"`),
		field.String("phonenum").
			Optional().Comment("手机号").
			StructTag(`json:"phonenum"  db:"phonenum"`),
		field.String("description").
			Optional().Comment("描述").
			StructTag(`json:"description"  db:"description"`),
		field.String("ethnicity").
			Optional().Comment("国籍民族").
			StructTag(`json:"ethnicity"  db:"ethnicity"`),
		field.String("gender").
			Optional().Comment("性别").
			StructTag(`json:"gender"  db:"gender"`),
		field.String("local").
			Optional().Comment("地区").
			StructTag(`json:"local"  db:"local"`),
		field.String("localhost").
			Optional().Comment("详细地址").
			StructTag(`json:"localhost"  db:"localhost"`),
		field.String("m2_localhost").
			Optional().Comment("酒店体量").
			StructTag(`json:"m2_localhost"  db:"m2_localhost"`),
		field.String("ismanager").
			Optional().Comment("是否管理").
			StructTag(`json:"ismanager"  db:"ismanager"`),
		field.String("isystem").
			Optional().Comment("是否内建").
			StructTag(`json:"isystem"  db:"isystem"`),
		field.Bool("enabled").
			Optional().Default(true).
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.Bool("is_sms").
			Optional().Comment("验证码登陆").
			StructTag(`json:"is_sms"  db:"is_sms"`),
		field.String("member_id").
			Optional().Comment("证件号").
			StructTag(`json:"member_id"  db:"member_id"`),
		field.String("leader_id").
			Optional().Comment("父账号id").
			StructTag(`json:"leader_id"  db:"leader_id"`),
		field.String("device_time").
			Optional().Comment("最近登陆").
			StructTag(`json:"device_time"  db:"device_time"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Optional().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.String("delete_flag").
			Optional().Comment("删除标记").
			StructTag(`json:"delete_flag"  db:"delete_flag"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			UpdateDefault(func() time.Time {
				return time.Now().Local().In(TimeLoc)
			}).
			Comment("更新时间").StructTag(`json:"updated_time" db:"updated_time"`),
		field.Time("deleted_time").
			Optional().Nillable().
			Comment("删除时间").StructTag(`json:"deleted_time" db:"deleted_time"`),
	}
}

func (FmUserAccount) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
		// login_name  phonenum 唯一
		index.Fields("login_name", "phonenum").Unique(),
	}
}
