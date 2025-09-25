package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Province schema 对应 fm_gov_province 表
type Province struct {
	ent.Schema
}

// 自定义表名
func (Province) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_gov_province"},
		entsql.WithComments(true), // 全局关闭字段注释同步
	}
}

func (Province) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("code").
			MaxLen(20).
			Validate(func(s string) error {
				if len(s) > 20 {
					return fmt.Errorf("code too long")
				}
				return nil
			}).
			Default("").Comment("编号").
			StructTag(`json:"code"  db:"code"`),
		field.String("name").
			MaxLen(200).
			Validate(func(s string) error {
				if len(s) > 200 {
					return fmt.Errorf("name too long")
				}
				return nil
			}).
			Default("").Comment("名称").
			StructTag(`json:"name"  db:"name"`),
		field.String("creator").
			MaxLen(20).
			Validate(func(s string) error {
				if len(s) > 20 {
					return fmt.Errorf("name too long")
				}
				return nil
			}).
			Default("admin").Comment("操作人").
			StructTag(`json:"creator"  db:"creator"`),
		field.String("delete_flag").
			MaxLen(2).
			Validate(func(s string) error {
				if len(s) > 2 {
					return fmt.Errorf("name too long")
				}
				return nil
			}).
			Default("0").Comment("删除标记").
			StructTag(`json:"delete_flag"  db:"delete_flag"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
	}
}

func (Province) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),

		//  code 联合唯一
		index.Fields("code").Unique(),
	}
}

// func (Province) Hooks() []ent.Hook {
// 	//  自定义的日志记录钩子
// 	return []ent.Hook{
// 		LogHookStd(),
// 	}
// }
