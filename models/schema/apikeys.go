package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Apikeys schema 对应 Apikeys 表
type Apikeys struct {
	ent.Schema
}

// 自定义表名
func (Apikeys) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ApiKeys"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (Apikeys) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.Int64("user_id").
			Optional().
			StructTag(`json:"user_id" db:"user_id"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Immutable().
			Comment("酒店号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.Int64("usage_count").
			Comment("使用计数").
			StructTag(`json:"usage_count" db:"usage_count"`),
		field.String("api_key").
			Comment("APIKey字符串").
			StructTag(`json:"api_key" db:"api_key"`),
		field.String("key_name").
			Comment("Key名称").
			StructTag(`json:"key_name" db:"key_name"`),
		field.Bool("enabled").
			Optional().
			Comment("启用").
			StructTag(`json:"enabled" db:"enabled"`),
		field.Time("expires_time").
			Immutable().
			Comment("到期日期").StructTag(`json:"expires_time" db:"expires_time"`),
		field.Time("created_time").
			Immutable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().
			Comment("创建时间").StructTag(`json:"updated_time" db:"updated_time"`),
		field.Time("last_used_time").
			Optional().
			Comment("最近使用时间").StructTag(`json:"last_used_time" db:"last_used_time"`),
		field.Int("type").
			Optional().Comment("使用场景").
			StructTag(`json:"type" db:"type"`),
	}
}

func (Apikeys) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}
