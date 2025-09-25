package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Alerts schema 对应 Alerts 表
type Alerts struct {
	ent.Schema
}

// 自定义表名
func (Alerts) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Alerts"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (Alerts) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Immutable().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.Int64("capture_log_id").
			Comment("抓拍记录 ID").
			StructTag(`json:"capture_log_id" db:"capture_log_id"`),
		field.Int64("device_id").
			Comment("设备 ID").
			StructTag(`json:"device_id" db:"device_id"`),
		field.Int8("alert_level").
			Default(0).
			Comment("预警等级: 1=一级预警, 2=二级预警, 3=三级预警, 4=四级预警").
			StructTag(`json:"alert_level" db:"alert_level"`),
		field.Int8("status").
			Default(0).
			Comment("处理状态: 0=未处理, 1=已处理, 2=已忽略").
			StructTag(`json:"status" db:"status"`),
		field.Int64("fm_user_id").
			Optional().Comment("处理人ID").
			StructTag(`json:"fm_user_id" db:"fm_user_id"`),
		field.Time("handled_time").
			Optional().Nillable().
			Comment("开始处理的时间").StructTag(`json:"handled_time" db:"handled_time"`),
		field.Text("remarks").
			Optional().Default("").
			Comment("处理意见或备注").
			StructTag(`json:"remarks,omitempty" db:"remarks"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("预警产生时间").StructTag(`json:"created_time" db:"created_time"`),
	}
}

func (Alerts) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}

// func (Alerts) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		//定义反向边
// 		edge.From("once", CaptureLogs.Type).
// 			Ref("once"). //指向 另一边的 edge 名
// 			Field("capture_log_id").
// 			Unique().   /// 与 CaptureLogs 一对一关系
// 			Required(), //表示这个 edge 不能为空（Alerts 必须有抓拍记录）
// 	}
// }
