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

// ProfileType schema 对应 profile_type 表
type ProfileType struct {
	ent.Schema
}

// 自定义表名
func (ProfileType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ProfileType"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (ProfileType) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("type_name").
			NotEmpty().
			MaxLen(100).Comment("类型名").
			StructTag(`json:"type_name"  db:"type_name"`),
		field.Int("warning_level").
			Optional().Nillable().Comment("预警等级").
			StructTag(`json:"warning_level" db:"warning_level"`),
		field.Bool("warning_enabled").
			Optional().Nillable().
			Comment("预警启用").StructTag(`json:"warning_enabled" db:"warning_enabled"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Optional().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.String("description").
			Optional().
			MaxLen(255).Comment("描述").
			StructTag(`json:"description"  db:"description"`),
		field.Bool("deleteable").
			Optional().
			Comment("可删除").StructTag(`json:"deleteable" db:"deleteable"`),
		field.Bool("enabled").
			Optional().
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.String("type_code").
			Optional().Nillable().MaxLen(20).
			Default("陌生人").
			Comment("人员类型").
			StructTag(`json:"type_code"  db:"type_code"`),
		field.Int("face_validity_hours").
			Optional().Nillable().Default(0).
			Comment("有效时长").
			StructTag(`json:"face_validity_hours"  db:"face_validity_hours"`),
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

func (ProfileType) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
		// 联合唯一
		index.Fields("tenant_id", "type_code").Unique(),
		index.Fields("tenant_id", "type_name").Unique(),
	}
}

// func (ProfileType) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("TypeToProfiles", Profiles.Type),
// 		//定义反向边
// 		edge.From("fm_tenant", Tenants.Type).
// 			Ref("TenantToProfileType").
// 			Field("tenant_id").
// 			Required(),
// 	}
// }
