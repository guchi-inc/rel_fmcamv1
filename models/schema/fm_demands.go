package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// fm_demands schema 对应 fm_demands 表
type FmDemands struct {
	ent.Schema
}

// 自定义表名
func (FmDemands) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_demands"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (FmDemands) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("supplier").
			Optional().Nillable().MaxLen(255).
			Validate(func(s string) error {
				if len(s) > 255 {
					return fmt.Errorf("supplier too long")
				}
				return nil
			}).
			Default("").Comment("公司名").
			StructTag(`json:"supplier"  db:"supplier"`),
		field.String("username").
			Optional().Nillable().MaxLen(200).
			Validate(func(s string) error {
				if len(s) > 200 {
					return fmt.Errorf("username too long")
				}
				return nil
			}).
			Default("").Comment("姓名").
			StructTag(`json:"username"  db:"username"`),
		field.String("phonenum").
			Optional().Nillable().MaxLen(20).
			Validate(func(s string) error {
				if len(s) > 20 {
					return fmt.Errorf("phonenum too long")
				}
				return nil
			}).
			Default("phonenum").Comment("联系电话").
			StructTag(`json:"phonenum"  db:"phonenum"`),
		field.String("email").
			Optional().Nillable().MaxLen(255).
			Validate(func(s string) error {
				if len(s) > 255 {
					return fmt.Errorf("email too long")
				}
				return nil
			}).
			Default("email").Comment("联系电邮").
			StructTag(`json:"email"  db:"email"`),
		field.String("province").
			Optional().Nillable().MaxLen(100).
			Validate(func(s string) error {
				if len(s) > 100 {
					return fmt.Errorf("province too long")
				}
				return nil
			}).
			Default("province").Comment("省份").
			StructTag(`json:"province"  db:"province"`),
		field.String("city").
			Optional().Nillable().MaxLen(100).
			Validate(func(s string) error {
				if len(s) > 100 {
					return fmt.Errorf("city too long")
				}
				return nil
			}).
			Default("city").Comment("城市").
			StructTag(`json:"city"  db:"city"`),
		field.String("area").
			Optional().Nillable().MaxLen(100).
			Validate(func(s string) error {
				if len(s) > 100 {
					return fmt.Errorf("area too long")
				}
				return nil
			}).
			Default("area").Comment("区县").
			StructTag(`json:"area"  db:"area"`),
		field.String("street").
			Optional().Nillable().MaxLen(255).
			Validate(func(s string) error {
				if len(s) > 255 {
					return fmt.Errorf("street too long")
				}
				return nil
			}).
			Default("street").Comment("街道").
			StructTag(`json:"street"  db:"street"`),
		field.String("message").
			Optional().Nillable().MaxLen(500).
			Validate(func(s string) error {
				if len(s) > 500 {
					return fmt.Errorf("message too long")
				}
				return nil
			}).
			Default("message").Comment("留言").
			StructTag(`json:"message"  db:"message"`),
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

func (FmDemands) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
	}
}
