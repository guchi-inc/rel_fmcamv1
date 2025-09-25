package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// CaptureLogs schema 对应 CaptureLogs 表
type CaptureLogs struct {
	ent.Schema
}

// 自定义表名
func (CaptureLogs) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "CaptureLogs"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (CaptureLogs) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Immutable().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.Int64("device_id").
			Immutable().Comment("设备 ID").
			StructTag(`json:"device_id" db:"device_id"`),
		field.Int64("matched_profile_id").
			Immutable().Comment("比中人员ID").
			StructTag(`json:"matched_profile_id" db:"matched_profile_id"`),
		field.Int8("func_type").
			Optional().Default(0).
			Comment("是否抓拍").StructTag(`json:"func_type" db:"func_type"`),
		field.Float32("match_score").
			Optional().
			Default(0.0).
			Comment("比对得分").
			StructTag(`json:"match_score" db:"match_score"`),
		field.Bool("has_alert").
			Default(false).
			Comment("触发预警").StructTag(`json:"has_alert" db:"has_alert"`),
		field.String("device_name").
			Optional().Nillable().
			Comment("设备名称").
			StructTag(`json:"device_name"  db:"device_name"`),
		field.String("device_location").
			Optional().Nillable().
			Comment("设备位置").
			StructTag(`json:"device_location"  db:"device_location"`),
		field.String("content").
			Optional().Nillable().
			Comment("抓拍记录内容").
			StructTag(`json:"content"  db:"content"`),
		field.String("capture_image_url").
			Optional().
			Default("").
			MaxLen(64).Comment("抓拍的图片地址").
			StructTag(`json:"capture_image_url"  db:"capture_image_url"`),
		field.Time("capture_time").
			Optional().Nillable().
			Comment("抓拍时间").StructTag(`json:"capture_time" db:"capture_time"`),
	}
}

func (CaptureLogs) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}

// func (CaptureLogs) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("Alerts", Alerts.Type).
// 			Field("capture_log_id").
// 			StorageKey(edge.Column("capture_log_id")),
// 		edge.From("Profile", Profiles.Type).
// 			Ref("Profile").
// 			Field("matched_profile_id").
// 			Unique().
// 			Required(),
// 	}
// }

// edge.To("books", Book.Type)  // 一个作者 -> 多本书
// func (CaptureLogs) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("CaptureToAlertsID", Alerts.Type),
// 		//定义反向边
// 		edge.From("Device", Devices.Type).
// 			Ref("DeviceToCaptureLogsID").
// 			Field("device_id").
// 			Required(),
// 		edge.From("Profile", Profiles.Type).
// 			Ref("ProfilesToCaptureLogs").
// 			Field("matched_profile_id").
// 			Required(),
// 		edge.From("TemporaryFace", TemporaryFace.Type).
// 			Ref("TemporaryFaceToCaptureLogs").
// 			Field("matched_transient_id").
// 			Required(),
// 		edge.From("fm_tenant", Tenants.Type).
// 			Ref("TenantToCaptureLogs").
// 			Field("tenant_id").
// 			Required(),
// 	}
// }
