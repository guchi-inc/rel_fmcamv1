package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// FmAlertDefinition schema 对应 fm_alert_definition 表
type FmAlertDefinition struct {
	ent.Schema
}

// 自定义表名
func (FmAlertDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_alert_definition"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (FmAlertDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.Int64("alert_group_id").
			Optional().Comment("关联的分组id").
			StructTag(`json:"alert_group_id"  db:"alert_group_id"`),
		field.Int("level").
			Optional().
			Comment("预警等级").
			StructTag(`json:"level"  db:"level"`),
		field.Int64("profile_type_id").
			Optional().Comment("人员类型id").
			StructTag(`json:"profile_type_id"  db:"profile_type_id"`),
		field.String("action").
			Optional().MaxLen(20).
			Comment("处理操作").
			StructTag(`json:"action"  db:"action"`),
		field.String("alarm_sound").
			Optional().MaxLen(20).
			Comment("预警声音").
			StructTag(`json:"alarm_sound"  db:"alarm_sound"`),
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

func (FmAlertDefinition) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}
