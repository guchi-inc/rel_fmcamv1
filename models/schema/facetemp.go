package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// TemporaryFace schema 对应 face 表
type TemporaryFace struct {
	ent.Schema
}

// 自定义表名
func (TemporaryFace) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "TemporaryFace"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (TemporaryFace) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Immutable().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.Int64("profile_id").
			Optional().Nillable().
			StructTag(`json:"profile_id" db:"profile_id"`),
		field.Bytes("face_embedding").
			MaxLen(255).Comment("特征向量").
			StructTag(`json:"face_embedding"  db:"face_embedding"`),
		field.String("img_url").
			Optional().Nillable().
			MaxLen(255).Comment("照片地址").
			StructTag(`json:"img_url"  db:"img_url"`),
		field.String("updated_location").
			Optional().Nillable().
			MaxLen(100).Comment("最后抓拍位置").
			StructTag(`json:"updated_location"  db:"updated_location"`),
		field.Int("capture_count").
			StructTag(`json:"capture_count" db:"capture_count"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("首次被抓拍时间").StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			Comment("最后被抓拍时间").StructTag(`json:"updated_time" db:"updated_time"`),
		field.Time("expires_time").
			Comment("过期清理时间").StructTag(`json:"expires_time" db:"expires_time"`),
	}
}

func (TemporaryFace) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}

// func (TemporaryFace) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("TemporaryFaceToCaptureLogs", CaptureLogs.Type),

// 		// 定义反向边
// 		edge.From("fm_tenant", Tenants.Type).
// 			Ref("TenantToTemporaryFace").
// 			Field("tenant_id").
// 			Required(),
// 	}
// }
