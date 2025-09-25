package schema

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"entgo.io/ent"
	"go.uber.org/zap"
)

var (
	schemaLog  = log.New(os.Stdout, "INFO -", 13)
	TimeLoc, _ = time.LoadLocation("Asia/Shanghai")
	CSTLayout  = "2006-01-02 15:04:05"
)

// LogHook 使用标准库 log 打印操作信息（包括调用位置）
func LogHookStd() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			start := time.Now()

			// 执行变更
			val, err := next.Mutate(ctx, m)

			// 获取调用行号
			_, file, line, ok := runtime.Caller(2)
			callerInfo := "unknown"
			if ok {
				callerInfo = fmt.Sprintf("%s:%d", file, line)
			}

			// 获取实体类型、操作类型
			entity := m.Type()
			op := m.Op().String()
			// 尝试获取主键 ID（如果存在）
			var id interface{} = "(not set)"
			// 获取 ID：通过 Mutation 类型断言
			switch mut := m.(type) {
			case interface{ ID() (int, bool) }:
				if i, ok := mut.ID(); ok {
					id = i
				}
			case interface{ ID() (int32, bool) }:
				if i, ok := mut.ID(); ok {
					id = i
				}
			case interface{ ID() (int64, bool) }:
				if i, ok := mut.ID(); ok {
					id = i
				}
			case interface{ ID() (string, bool) }:
				if i, ok := mut.ID(); ok {
					id = i
				}
			default:
				// 尝试从返回值获取 ID
				if entVal, ok := val.(interface{ ID() any }); ok {
					id = entVal.ID()
				}
			}

			// zap 记录
			schemaLog.Println("Ent Operation Executed",
				zap.String("operation", "log this."),
				zap.Error(err),
			)

			schemaLog.Printf("[EntLog] %s | %s %s | ID=%v | Duration=%v | Caller=%s | Err=%v",
				time.Now().Format("2006-01-02 15:04:05"),
				entity, op,
				id,
				time.Since(start),
				callerInfo,
				err,
			)

			return val, err
		})
	}
}

// 1. 创建一个自定义的 NullTime 类型
type NullTime struct {
	sql.NullTime
}

// 2. 实现 ent.Field 类型的接口
func (nt NullTime) Value() (driver.Value, error) {
	if nt.Valid {
		return nt.Time, nil
	}
	return nil, nil
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	// 尝试将值转换为 time.Time 类型
	switch v := value.(type) {
	case time.Time:
		nt.Time = v
		nt.Valid = true
		return nil
	case []byte:
		// 如果返回的是 []byte 类型（如 MySQL 返回的时间格式）
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		nt.Time = parsedTime
		nt.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported scan type for NullTime: %T", v)
	}
}
