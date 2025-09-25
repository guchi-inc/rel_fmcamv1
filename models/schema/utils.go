package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

// Decimal 是对 shopspring/decimal 的封装，支持 JSON 和数据库转换及基本运算。
type Decimal struct {
	decimal.Decimal
}

// IsZero 判断 Decimal 是否为零。
func (d Decimal) IsZero() bool {
	return d.Decimal.IsZero()
}

// MarshalJSON 实现 JSON 编码接口（以字符串格式输出，保持精度）。
func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON 实现 JSON 解码接口（兼容字符串和裸数字）。
func (d *Decimal) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		dec, err := decimal.NewFromString(str)
		if err != nil {
			return err
		}
		d.Decimal = dec
		return nil
	}
	dec, err := decimal.NewFromString(string(b))
	if err != nil {
		return err
	}
	d.Decimal = dec
	return nil
}

// Value 实现 driver.Valuer，用于数据库写入。
func (d Decimal) Value() (driver.Value, error) {
	return d.Decimal.String(), nil
}

// Scan 实现 sql.Scanner，用于数据库读取。
func (d *Decimal) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		dec, err := decimal.NewFromString(string(v))
		if err != nil {
			return err
		}
		d.Decimal = dec
		return nil
	case string:
		dec, err := decimal.NewFromString(v)
		if err != nil {
			return err
		}
		d.Decimal = dec
		return nil
	case float64:
		d.Decimal = decimal.NewFromFloat(v)
		return nil
	case nil:
		d.Decimal = decimal.Zero
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Decimal", v)
	}
}

// String 实现 fmt.Stringer。
func (d Decimal) String() string {
	return d.Decimal.String()
}

// -----------------------------
// 构造函数 & 转换函数
// -----------------------------

// NewDecimalFromFloat64 创建 Decimal。
func NewDecimalFromFloat64(f float64) Decimal {
	return Decimal{decimal.NewFromFloat(f)}
}

// NewDecimalFromString 创建 Decimal。
func NewDecimalFromString(s string) (Decimal, error) {
	dec, err := decimal.NewFromString(s)
	if err != nil {
		return Decimal{}, err
	}
	return Decimal{dec}, nil
}

// ToFloat64 转为 float64。
func (d Decimal) ToFloat64() (float64, bool) {
	return d.Decimal.Float64()
}

// ToString 转为字符串。
func (d Decimal) ToString() string {
	return d.Decimal.String()
}

// -----------------------------
// 运算函数
// -----------------------------

// Add 加法。
func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{d.Decimal.Add(other.Decimal)}
}

// Sub 减法。
func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{d.Decimal.Sub(other.Decimal)}
}

// Mul 乘法。
func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{d.Decimal.Mul(other.Decimal)}
}

// Div 除法（如除以零，返回 NaN）。
func (d Decimal) Div(other Decimal) Decimal {
	if other.IsZero() {
		return Decimal{decimal.NewFromFloat(0).Div(decimal.NewFromFloat(0))} // NaN
	}
	return Decimal{d.Decimal.Div(other.Decimal)}
}

// -----------------------------
// 比较函数
// -----------------------------

// GT 判断是否大于另一个 Decimal。
func (d Decimal) GT(other Decimal) bool {
	return d.Decimal.GreaterThan(other.Decimal)
}

// LT 判断是否小于另一个 Decimal。
func (d Decimal) LT(other Decimal) bool {
	return d.Decimal.LessThan(other.Decimal)
}

// EQ 判断是否等于另一个 Decimal。
func (d Decimal) EQ(other Decimal) bool {
	return d.Decimal.Equal(other.Decimal)
}

// -----------------------------
// 可空 Decimal（支持 JSON null）
// -----------------------------

// NullableDecimal 支持 JSON null 的 Decimal。
type NullableDecimal struct {
	Valid bool     // 是否有效
	Value *Decimal // 实际值
}

// MarshalJSON 实现 JSON 编码接口。
func (nd NullableDecimal) MarshalJSON() ([]byte, error) {
	if !nd.Valid || nd.Value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(nd.Value)
}

// UnmarshalJSON 实现 JSON 解码接口。
func (nd *NullableDecimal) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nd.Valid = false
		nd.Value = nil
		return nil
	}
	var d Decimal
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	nd.Valid = true
	nd.Value = &d
	return nil
}
