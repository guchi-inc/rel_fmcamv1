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

// Faces schema 对应 face 表
type Faces struct {
	ent.Schema
}

// 自定义表名
func (Faces) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Face"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (Faces) Fields() []ent.Field {
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
			Immutable().
			StructTag(`json:"profile_id" db:"profile_id"`),
		field.Bytes("face_embedding").
			MaxLen(255).Comment("特征向量").
			StructTag(`json:"face_embedding"  db:"face_embedding"`),
		field.String("image_url").
			Default("").
			MaxLen(255).Comment("照片地址").
			StructTag(`json:"image_url"  db:"image_url"`),
		field.Bool("is_primary").
			Default(false).
			Comment("主照片").StructTag(`json:"is_primary" db:"is_primary"`),
		field.String("updated_location").
			Optional().Nillable().
			MaxLen(100).Comment("最后抓拍位置").
			StructTag(`json:"updated_location"  db:"updated_location"`),
		field.Int("capture_count").
			Comment("抓拍总次数").
			StructTag(`json:"capture_count" db:"capture_count"`),
		field.Time("expires_time").
			Comment("过期清理时间").StructTag(`json:"expires_time" db:"expires_time"`),
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

func (Faces) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}

// func (Faces) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		////定义反向边
// 		edge.From("Profile", Profiles.Type).
// 			Ref("ProfilesToFaces").
// 			Field("profile_id").
// 			Required(),
// 		edge.From("fm_tenant", Tenants.Type).
// 			Ref("TenantToFaces").
// 			Field("tenant_id").
// 			Required(),
// 	}
// }
