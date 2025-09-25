package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// schema 对应 field_metadata 表
type FieldMetadata struct {
	ent.Schema
}

// 自定义表名
func (FieldMetadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "field_metadata"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

// Fields of the FieldMetadata.
func (FieldMetadata) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			Comment("主键").
			Positive().
			StructTag(`json:"id" db:"id"`),
		field.String("table_name").
			MaxLen(64).
			NotEmpty().
			Comment("表名").StructTag(`json:"table_name" db:"table_name"`),
		field.String("name").
			MaxLen(64).
			NotEmpty().
			Comment("字段名").StructTag(`json:"name" db:"name"`),
		field.String("cname").
			MaxLen(128).
			NotEmpty().
			Comment("中文名").StructTag(`json:"cname" db:"cname"`),
		field.String("data_type").
			MaxLen(32).
			NotEmpty().
			Comment("字段类型").StructTag(`json:"data_type" db:"data_type"`),
		field.Bool("is_visible").
			Default(true).
			Comment("可访问").StructTag(`json:"is_visible" db:"is_visible"`),
		field.Bool("is_searchable").
			Default(true).
			Comment("可搜索").StructTag(`json:"is_searchable" db:"is_searchable"`),
		field.Bool("is_editable").
			Default(true).
			Comment("可编辑").StructTag(`json:"is_editable" db:"is_editable"`),
		field.Bool("is_required").
			Default(false).
			Comment("必需").StructTag(`json:"is_required" db:"is_required"`),
		field.Int("max_length").
			Default(20).
			Comment("最长").StructTag(`json:"max_length" db:"max_length"`),
		field.String("default_value").
			MaxLen(20).
			Default("").
			Comment("默认值").StructTag(`json:"default_value" db:"default_value"`),
		field.String("description").
			MaxLen(20).
			Default("").
			Comment("描述").StructTag(`json:"description" db:"description"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").
			StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			UpdateDefault(func() time.Time {
				return time.Now().Local().In(TimeLoc)
			}).
			Comment("更新时间").StructTag(`json:"updated_time" db:"updated_time"`),
	}
}

func (FieldMetadata) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),

		// code 联合唯一
		index.Fields("table_name", "name").Unique(),
	}
}

// func (FieldMetadata) Hooks() []ent.Hook {
// 	//  自定义的日志记录钩子
// 	return []ent.Hook{
// 		LogHookStd(),
// 	}
// }
