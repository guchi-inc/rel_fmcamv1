package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// FMPMSApi schema 对应 fm_pms_apis 表
type FMPMSApi struct {
	ent.Schema
}

// 自定义表名
func (FMPMSApi) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_pms_apis"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (FMPMSApi) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("pms_name").
			MaxLen(50).
			Validate(func(s string) error {
				if len(s) > 50 {
					return fmt.Errorf("pms_name too long")
				}
				return nil
			}).
			Default("").Comment("PMS商名").
			StructTag(`json:"pms_name"  db:"pms_name"`),
		field.String("pms_api").
			MaxLen(500).
			Validate(func(s string) error {
				if len(s) > 500 {
					return fmt.Errorf("pms_api too long")
				}
				return nil
			}).
			Default("").Comment("对接地址").
			StructTag(`json:"pms_api"  db:"pms_api"`),
		field.Bool("enabled").
			Optional().Nillable().Default(true).
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.String("contact").
			MaxLen(50).
			Validate(func(s string) error {
				if len(s) > 50 {
					return fmt.Errorf("contact too long")
				}
				return nil
			}).
			Default("").Comment("联系人").
			StructTag(`json:"contact"  db:"contact"`),
		field.String("phonenum").
			Optional().Nillable().MaxLen(20).
			Validate(func(s string) error {
				if len(s) > 20 {
					return fmt.Errorf("phonenum too long")
				}
				return nil
			}).
			Default("admin").Comment("联系电话").
			StructTag(`json:"phonenum"  db:"phonenum"`),
		field.String("description").
			Optional().Nillable().MaxLen(100).
			Validate(func(s string) error {
				if len(s) > 100 {
					return fmt.Errorf("description too long")
				}
				return nil
			}).
			Default("0").Comment("功能描述").
			StructTag(`json:"description"  db:"description"`),
		field.String("delete_flag").
			Optional().Nillable().MaxLen(1).
			Validate(func(s string) error {
				if len(s) > 1 {
					return fmt.Errorf("delete_flag too long")
				}
				return nil
			}).
			Default("0").Comment("删除标记").
			StructTag(`json:"delete_flag"  db:"delete_flag"`),
		field.String("creator").
			Optional().Nillable().MaxLen(25).
			Validate(func(s string) error {
				if len(s) > 25 {
					return fmt.Errorf("creator too long")
				}
				return nil
			}).
			Default("0").Comment("操作员").
			StructTag(`json:"creator"  db:"creator"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			Comment("更新时间").StructTag(`json:"updated_time" db:"updated_time"`),
	}
}

func (FMPMSApi) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),

		// pms_name  pms_api 唯一
		index.Fields("pms_name", "pms_api").Unique(),
	}
}

// func (FMPMSApi) Hooks() []ent.Hook {
// 	//  自定义的日志记录钩子
// 	return []ent.Hook{
// 		LogHookStd(),
// 	}
// }
