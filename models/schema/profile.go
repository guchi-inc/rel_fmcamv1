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

// Profiles schema 对应 profile 表
type Profiles struct {
	ent.Schema
}

// 自定义表名
func (Profiles) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Profile"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (Profiles) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.Int64("type_id").
			Optional().Nillable().Comment("人员类型").
			StructTag(`json:"type_id" db:"type_id"`),
		field.String("name").
			Optional().Nillable().
			MaxLen(100).Comment("姓名").
			StructTag(`json:"name"  db:"name"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Optional().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.String("id_card_number").
			Optional().Nillable().
			MaxLen(20).Comment("身份证号").
			StructTag(`json:"id_card_number"  db:"id_card_number"`),
		field.String("phone_number").
			Optional().Nillable().
			Default("").
			Comment("手机号").
			StructTag(`json:"phone_number"  db:"phone_number"`),
		field.Bool("enabled").
			Optional().Nillable().Default(true).
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.String("room_id").
			Optional().Nillable().
			Comment("房间号").StructTag(`json:"room_id" db:"room_id"`),
		field.String("tmp_url").
			Optional().Nillable().
			Comment("面部临时地址").StructTag(`json:"tmp_url" db:"tmp_url"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			UpdateDefault(func() time.Time {
				return time.Now().Local().In(TimeLoc)
			}).
			Comment("更新时间").StructTag(`json:"updated_time" db:"updated_time"`),
	}
}

func (Profiles) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}

// func (Profiles) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		// 一个类型 -> 多个作者
// 		edge.To("Profile", CaptureLogs.Type).
// 			Field("id").
// 			StorageKey(edge.Column("matched_profile_id")),
// 	}
// }
