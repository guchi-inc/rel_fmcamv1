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

// ProfileType schema 对应 fm_alert_group 表
type FmAlertGroup struct {
	ent.Schema
}

// 自定义表名
func (FmAlertGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_alert_group"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (FmAlertGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("group_name").
			NotEmpty().
			MaxLen(100).Comment("组名").
			StructTag(`json:"group_name"  db:"group_name"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Immutable().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.Int8("group_type").
			Optional().Default(1).
			Comment("组类别").StructTag(`json:"group_type" db:"group_type"`),
		field.Bool("enabled").
			Optional().Nillable().Default(true).
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.String("customization").
			Optional().NotEmpty().
			MaxLen(255).Comment("定制化").
			StructTag(`json:"customization"  db:"customization"`),
		field.String("description").
			Optional().MaxLen(20).
			Comment("描述").
			StructTag(`json:"description"  db:"description"`),
		field.String("creator").
			Optional().MaxLen(50).
			Comment("操作员").
			StructTag(`json:"creator"  db:"creator"`),
		field.String("uniq_enabled_group").
			Optional().Nillable().
			Immutable().
			Comment("启用户").StructTag(`json:"uniq_enabled_group,omitempty" db:"uniq_enabled_group"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time,omitempty" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			UpdateDefault(func() time.Time {
				return time.Now().Local().In(TimeLoc)
			}).
			Comment("更新时间").StructTag(`json:"updated_time,omitempty" db:"updated_time"`),
	}
}

func (FmAlertGroup) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}
