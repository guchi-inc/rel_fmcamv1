package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// City schema 对应 fm_gov_city 表
type GovCity struct {
	ent.Schema
}

// 自定义表名
func (GovCity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_gov_city"},
		entsql.WithComments(true), // false 全局关闭字段注释同步
	}
}

func (GovCity) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("code").
			MaxLen(20).Comment("编号").
			StructTag(`json:"code"  db:"code"`),
		field.String("name").
			MaxLen(200).
			Default("").Comment("名称").
			StructTag(`json:"name"  db:"name"`),
		field.String("province_code").
			MaxLen(200).
			Default("").Comment("省编号").
			StructTag(`json:"province_code"  db:"province_code"`),
		field.String("creator").
			MaxLen(20).
			Default("").Comment("操作人").
			StructTag(`json:"creator"  db:"creator"`),
		field.String("delete_flag").
			MaxLen(2).
			Default("0").Comment("删除标记").
			StructTag(`json:"delete_flag"  db:"delete_flag"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
	}
}

func (GovCity) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),

		// code  唯一
		index.Fields("code").Unique(),
	}
}

// func (GovCity) Hooks() []ent.Hook {
// 	//  自定义的日志记录钩子
// 	return []ent.Hook{
// 		LogHookStd(),
// 	}
// }
