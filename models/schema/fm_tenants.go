package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Area schema 对应 fm_gov_area 表
type Tenants struct {
	ent.Schema
}

// 自定义表名
func (Tenants) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fm_tenant"},
		entsql.WithComments(true), // 全局关闭字段注释同步
	}
}

func (Tenants) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.String("supplier").
			MaxLen(255).Comment("租户名").
			StructTag(`json:"supplier"  db:"supplier"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.String("contacts").
			MaxLen(100).
			Default("").Comment("联系人").
			StructTag(`json:"contacts"  db:"contacts"`),
		field.String("email").
			MaxLen(50).
			Default("").Comment("电子邮箱").
			StructTag(`json:"email"  db:"email"`),
		field.String("description").
			MaxLen(100).
			Default("").Comment("备注").
			StructTag(`json:"description"  db:"description"`),
		field.String("type").
			MaxLen(20).
			Default("").Comment("类型").
			StructTag(`json:"type"  db:"type"`),
		field.String("province").
			MaxLen(50).
			Default("").Comment("省份").
			StructTag(`json:"province"  db:"province"`),
		field.String("city").
			MaxLen(50).
			Default("").Comment("城市").
			StructTag(`json:"city"  db:"city"`),
		field.String("area").
			MaxLen(50).
			Default("").Comment("区县").
			StructTag(`json:"area"  db:"area"`),
		field.String("street").
			MaxLen(50).
			Default("").Comment("街道乡镇").
			StructTag(`json:"street"  db:"street"`),
		field.String("address").
			MaxLen(50).
			Default("").Comment("门牌地址").
			StructTag(`json:"address"  db:"address"`),
		field.String("addr_code").
			MaxLen(30).
			Default("").Comment("地址编码").
			StructTag(`json:"addr_code"  db:"addr_code"`),
		field.String("fax").
			MaxLen(30).
			Default("").Comment("传真").
			StructTag(`json:"fax"  db:"fax"`),
		field.String("phone_num").
			MaxLen(30).
			Default("").Comment("固定电话").
			StructTag(`json:"phone_num"  db:"phone_num"`),
		field.String("telephone").
			MaxLen(30).
			Default("").Comment("手机号").
			StructTag(`json:"telephone"  db:"telephone"`),
		field.String("tax_num").
			MaxLen(50).
			Default("").Comment("纳税人,识别号").
			StructTag(`json:"tax_num"  db:"tax_num"`),
		field.String("bank_name").
			MaxLen(50).
			Default("").Comment("开户行").
			StructTag(`json:"bank_name"  db:"bank_name"`),
		field.String("account_number").
			MaxLen(50).
			Default("").Comment("账号").
			StructTag(`json:"account_number"  db:"account_number"`),
		field.String("sort").
			MaxLen(10).
			Default("").Comment("排序").
			StructTag(`json:"sort"  db:"sort"`),
		field.Bool("enabled").
			Optional().Default(true).
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.String("delete_flag").
			Default("0").
			Comment("已删除").StructTag(`json:"delete_flag" db:"delete_flag"`),
		field.Bool("isystem").
			Default(false).
			Comment("系统内置").StructTag(`json:"isystem" db:"isystem"`),
		field.Other("tax_rate", Decimal{}).
			SchemaType(map[string]string{
				"mysql": "decimal(24,6)",
			}).
			Default(NewDecimalFromFloat64(0.0)).
			Comment("税率"),
		field.Other("advance_in", Decimal{}).
			SchemaType(map[string]string{
				"mysql": "decimal(24,6)",
			}).
			Default(NewDecimalFromFloat64(0.0)).
			Comment("预收款"),
		field.Other("begin_need_get", Decimal{}).
			SchemaType(map[string]string{
				"mysql": "decimal(24,6)",
			}).
			Default(NewDecimalFromFloat64(0.0)).
			Comment("期初应收"),
		field.Other("begin_need_pay", Decimal{}).
			SchemaType(map[string]string{
				"mysql": "decimal(24,6)",
			}).
			Default(NewDecimalFromFloat64(0.0)).
			Comment("期初应付"),
		field.Other("all_need_get", Decimal{}).
			SchemaType(map[string]string{
				"mysql": "decimal(24,6)",
			}).
			Default(NewDecimalFromFloat64(0.0)).
			Comment("累计应收"),
		field.Other("all_need_pay", Decimal{}).
			SchemaType(map[string]string{
				"mysql": "decimal(24,6)",
			}).
			Default(NewDecimalFromFloat64(0.0)).
			Comment("累计应付"),
		field.String("creator").
			MaxLen(20).
			Default("").Comment("操作人").
			StructTag(`json:"creator"  db:"creator"`),
		field.Time("created_time").
			Optional().Nillable().
			Comment("创建时间").StructTag(`json:"created_time" db:"created_time"`),
		field.Time("updated_time").
			Optional().Nillable().
			UpdateDefault(func() time.Time {
				return time.Now().Local().In(TimeLoc)
			}).
			Comment("更新时间").StructTag(`json:"updated_time" db:"updated_time"`),
	}
}

func (Tenants) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),

		// tenant_id  唯一
		index.Fields("tenant_id").Unique(),
	}
}

// func (Tenants) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("TenantToCaptureLogs", CaptureLogs.Type),
// 		edge.To("TenantToAlerts", Alerts.Type),
// 		edge.To("TenantToDevices", Devices.Type),
// 		edge.To("TenantToFaces", Faces.Type),
// 		edge.To("TenantToTemporaryFace", TemporaryFace.Type),
// 		edge.To("TenantToProfileType", ProfileType.Type),
// 		edge.To("TenantToProfiles", Profiles.Type),
// 	}
// }

/*
client.Authors.

	Query().
	Select("id", "name", "full_address").
	Scan(ctx, &results)
*/
// func AddGeneratedColumnHook() schema.Hook {
// 	return func(next schema.Applier) schema.Applier {
// 		return schema.ApplyFunc(func(ctx context.Context, conn *schema.Conn) error {
// 			// 调用下一步迁移
// 			if err := next.Apply(ctx, conn); err != nil {
// 				return err
// 			}

// 			// 添加 full_address 虚拟列
// 			_, err := conn.ExecContext(ctx, `
//                 ALTER TABLE authors
//                 ADD COLUMN full_address VARCHAR(255)
//                 GENERATED ALWAYS AS (
//                     CONCAT(province, city, area, IFNULL(street, ''), IFNULL(address, ''))
//                 ) STORED
//                 COMMENT '详细地址';
//             `)
// 			return err
// 		})
// 	}
// }
