// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package code

import (
	"fmcam/ctrlib/utils/errors"
	"fmt"

	ut "github.com/go-playground/universal-translator" //使用 CLDR 数据 + 复数规则的 Go/Golang 的 i18n 翻译器
	"github.com/go-playground/validator/v10"           //Go Struct 和 Field 验证，包括 Cross Field、Cross Struct、Map、Slice 和 Array
)

var _ BusinessError = (*businessError)(nil)
var trans ut.Translator

type BusinessError interface {
	// i 为了避免被其他包实现
	i()

	// WithError 设置错误信息
	WithError(err error) BusinessError

	// WithAlert 设置告警通知
	WithAlert() BusinessError

	// BusinessCode 获取业务码
	BusinessCode() int

	// HTTPCode 获取 HTTP 状态码
	HTTPCode() int

	// Message 获取错误描述
	Message() string

	// StackError 获取带堆栈的错误信息
	StackError() error

	// IsAlert 是否开启告警通知
	IsAlert() bool
}

type businessError struct {
	httpCode     int    // HTTP 状态码
	businessCode int    // 业务码
	message      string // 错误描述
	stackError   error  // 含有堆栈信息的错误
	isAlert      bool   // 是否告警通知
}

func Error(httpCode, businessCode int, message string) BusinessError {
	return &businessError{
		httpCode:     httpCode,
		businessCode: businessCode,
		message:      message,
		isAlert:      false,
	}
}

func (be *businessError) StructError() error {
	return fmt.Errorf("%#v", be)
}

func (e *businessError) i() {}

func (e *businessError) WithError(err error) BusinessError {
	e.stackError = errors.WithStack(err)
	return e
}

func (e *businessError) WithAlert() BusinessError {
	e.isAlert = true
	return e
}

func (e *businessError) HTTPCode() int {
	return e.httpCode
}

func (e *businessError) BusinessCode() int {
	return e.businessCode
}

func (e *businessError) Message() string {
	return e.message
}

func (e *businessError) StackError() error {
	return e.stackError
}

func (e *businessError) IsAlert() bool {
	return e.isAlert
}

func init() {
	// lang := configs.GCONFIG.Language.Local

	// if lang == configs.ZhCN {
	// 	trans, _ = ut.New(zh.New()).GetTranslator("zh")
	// 	if err := zhTranslation.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), trans); err != nil {
	// 		Logger.Println("validator zh translation error", err)
	// 	}
	// }

	// if lang == configs.EnUS {
	// 	trans, _ = ut.New(en.New()).GetTranslator("en")
	// 	if err := enTranslation.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), trans); err != nil {
	// 		Logger.Println("validator en translation error", err)
	// 	}
	// }
}

func ErrorSer(err error) (message string) {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return err.Error()
	} else {
		for _, e := range validationErrors {
			message += e.Translate(trans) + ";"
		}
	}
	return message
}
