package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type SqlLog struct {
	ent.Schema
}

// 自定义表名
func (SqlLog) Table() string {
	return "sql_logs"
}

func (SqlLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sql_logs"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (SqlLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("table_name").
			MaxLen(64).Comment("表名").
			StructTag(`json:"table_name"  db:"table_name"`),
		field.String("query").
			Optional().Nillable().
			MaxLen(500).Comment("执行操作").
			StructTag(`json:"query"  db:"query"`),
		field.String("args").
			Optional().Nillable().
			MaxLen(200).Comment("参数").
			StructTag(`json:"args"  db:"args"`),
		field.String("action").
			Optional().Nillable().
			Comment("操作").
			StructTag(`json:"action" db:"action"`),
		field.String("db_name").
			Optional().Nillable().MaxLen(64).
			Default("").Comment("数据库").
			StructTag(`json:"db_name"  db:"db_name"`),
		field.Int("pk_value").
			Optional().Comment("操作键").
			StructTag(`json:"pk_value" db:"pk_value"`),
		field.JSON("old_data", map[string]interface{}{}).
			Optional().Comment("旧值").
			StructTag(`json:"old_data" db:"old_data"`),
		field.JSON("new_data", map[string]interface{}{}).
			Optional().Comment("新值").
			StructTag(`json:"new_data" db:"new_data"`),
		field.String("creator").
			Optional().Nillable().MaxLen(20).
			Default("").Comment("操作人").
			StructTag(`json:"creator"  db:"creator"`),
		field.Time("created_time").
			Optional().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
	}
}
func (SqlLog) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}
