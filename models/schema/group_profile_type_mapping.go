package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// chema 对应 group_profile_type_mapping 表
type GrouProfileTypeMapping struct {
	ent.Schema
}

// 自定义表名
func (GrouProfileTypeMapping) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "group_profile_type_mapping"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (GrouProfileTypeMapping) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.Int64("group_id").
			Comment("分组ID").
			StructTag(`json:"group_id" db:"group_id"`),
		field.Int64("profile_type_id").
			Comment("设备 ID").
			StructTag(`json:"profile_type_id" db:"profile_type_id"`),
		field.String("creator").
			Immutable().Comment("操作人").
			StructTag(`json:"creator" db:"creator"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("操作时间").StructTag(`json:"created_time" db:"created_time"`),
	}
}

func (GrouProfileTypeMapping) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
		index.Fields("group_id", "profile_type_id").Unique(),
	}
}

// func (Alerts) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		//定义反向边
// 		edge.From("Device", Devices.Type).
// 			Ref("DeviceToAlertsID").
// 			Field("device_id").
// 			Required(),
// 		edge.From("CaptureLogs", CaptureLogs.Type).
// 			Ref("CaptureToAlertsID").
// 			Field("capture_log_id").
// 			Required(),
// 		edge.From("fm_tenant", Tenants.Type).
// 			Ref("TenantToAlerts").
// 			Field("tenant_id").
// 			Required(),
// 	}
// }
