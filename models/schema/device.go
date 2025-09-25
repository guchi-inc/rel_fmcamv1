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

// Devices schema 对应 face 表
type Devices struct {
	ent.Schema
}

// 自定义表名
func (Devices) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Device"},
		entsql.WithComments(false), // 全局关闭字段注释同步
	}
}

func (Devices) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique().
			Immutable().
			StructTag(`json:"id" db:"id"`),
		field.UUID("tenant_id", uuid.UUID{}).
			Immutable().
			Comment("租户号").
			StructTag(`json:"tenant_id"  db:"tenant_id"`),
		field.String("name").
			Optional().Nillable().
			Default("").
			MaxLen(64).Comment("设备名称").
			StructTag(`json:"name"  db:"name"`),
		field.String("url").
			Default("").
			MaxLen(255).Comment("设备管理地址").
			StructTag(`json:"url"  db:"url"`),
		field.String("location").
			Default("").
			MaxLen(128).Comment("安装地址").
			StructTag(`json:"location"  db:"location"`),
		field.Int("func_type").
			Optional().Default(0).
			Comment("是否监控").StructTag(`json:"func_type" db:"func_type"`),
		field.Uint16("display_width").
			Optional(). // 对应 NULLABLE
			Nillable(). // 允许 Go 层使用 nil 来表示 NULL
			Comment("图像宽度").
			StructTag(`json:"display_width"  db:"display_width"`),
		field.Uint16("display_height").
			Optional(). // 对应 NULLABLE
			Nillable(). // 允许 Go 层使用 nil 来表示 NULL
			Comment("图像高度").
			StructTag(`json:"display_height"  db:"display_height"`),
		field.Uint16("roi_x").
			Optional(). // 对应 NULLABLE
			Nillable(). // 允许 Go 层使用 nil 来表示 NULL
			Comment("采集区域左上角X坐标").
			StructTag(`json:"roi_x"  db:"roi_x"`),
		field.Uint16("roi_y").
			Optional().
			Nillable().
			Comment("采集区域左上角y坐标").
			StructTag(`json:"roi_y"  db:"roi_y"`),
		field.Uint16("roi_width").
			Optional().
			Nillable().
			Comment("采集区域宽度").
			StructTag(`json:"roi_width"  db:"roi_width"`),
		field.Uint16("roi_height").
			Optional().
			Nillable().
			Comment("采集区域高度").
			StructTag(`json:"roi_height"  db:"roi_height"`),
		field.Float("roi_rotation_angle").
			Optional().
			Nillable().
			Comment("采集区域旋转角度").
			StructTag(`json:"roi_rotation_angle"  db:"roi_rotation_angle"`),
		field.Bool("roi_enabled").
			Default(true).
			Comment("已设置ROI").StructTag(`json:"roi_enabled" db:"roi_enabled"`),
		field.Bool("enabled").
			Default(true).
			Comment("启用").StructTag(`json:"enabled" db:"enabled"`),
		field.Int64("target_fps").
			Default(5).Comment("目标帧率").
			StructTag(`json:"target_fps" db:"target_fps"`),
		field.Int64("dwell_duration").
			Optional().Nillable().
			Comment("逗留时间").
			StructTag(`json:"dwell_duration" db:"dwell_duration"`),
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

func (Devices) Indexes() []ent.Index {
	return []ent.Index{
		// id 唯一
		index.Fields("id").Unique(),
		index.Fields("tenant_id", "url").Unique(),
		index.Fields("tenant_id", "name").Unique(),
	}
}

// func (Devices) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("DeviceToCaptureLogsID", CaptureLogs.Type),
// 		edge.To("DeviceToAlertsID", Alerts.Type),
// 		//定义反向边
// 		edge.From("fm_tenant", Tenants.Type).
// 			Ref("TenantToDevices").
// 			Field("tenant_id").
// 			Required(),
// 	}
// }
