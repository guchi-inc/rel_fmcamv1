package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// fm_dedicated_services schema 对应 fm_dedicated_services 表
type FmDedicatedServices struct {
	ent.Schema
}

// 自定义表名
func (FmDedicatedServices) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_dedicated_services"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (FmDedicatedServices) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.Int64("work_id").
			Optional().Nillable().
			Comment("酒店id").
			StructTag(`json:"work_id" db:"work_id"`),
		field.String("contacts").
			Optional().Nillable().MaxLen(100).
			Validate(func(s string) error {
				if len(s) > 100 {
					return fmt.Errorf("contacts too long")
				}
				return nil
			}).
			Default("").Comment("服务客服名").
			StructTag(`json:"contacts"  db:"contacts"`),
		field.String("supplier").
			Optional().Nillable().MaxLen(255).
			Validate(func(s string) error {
				if len(s) > 255 {
					return fmt.Errorf("supplier too long")
				}
				return nil
			}).
			Default("").Comment("酒店名").
			StructTag(`json:"supplier"  db:"supplier"`),
		field.String("phonenum").
			Optional().Nillable().MaxLen(20).
			Validate(func(s string) error {
				if len(s) > 20 {
					return fmt.Errorf("phonenum too long")
				}
				return nil
			}).
			Default("phonenum").Comment("服务电话").
			StructTag(`json:"phonenum"  db:"phonenum"`),
		field.String("email").
			Optional().Nillable().MaxLen(255).
			Validate(func(s string) error {
				if len(s) > 255 {
					return fmt.Errorf("email too long")
				}
				return nil
			}).
			Default("email").Comment("服务电邮").
			StructTag(`json:"email"  db:"email"`),
		field.String("fax").
			Optional().Nillable().MaxLen(30).
			Validate(func(s string) error {
				if len(s) > 30 {
					return fmt.Errorf("fax too long")
				}
				return nil
			}).
			Default("fax").Comment("传真号").
			StructTag(`json:"fax"  db:"fax"`),
		field.String("description").
			Optional().Nillable().MaxLen(255).
			Validate(func(s string) error {
				if len(s) > 255 {
					return fmt.Errorf("description too long")
				}
				return nil
			}).
			Default("description").Comment("描述").
			StructTag(`json:"description"  db:"description"`),
		field.String("creator").
			Optional().Nillable().MaxLen(20).
			Validate(func(s string) error {
				if len(s) > 20 {
					return fmt.Errorf("name too long")
				}
				return nil
			}).
			Default("admin").Comment("操作人").
			StructTag(`json:"creator"  db:"creator"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
	}
}

func (FmDedicatedServices) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}
