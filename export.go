package errors

import (
	"errors"

	"github.com/go-haru/field"
)

func With(err error, fields ...field.Field) *Exception {
	if err == nil {
		return nil
	}
	var exp Exception
	switch val := err.(type) {
	default:
		exp.cause, exp.field = val, fields
	case *Message:
		exp.message, exp.field = val, fields
	case *Exception:
		exp = *val
		exp.field = append(exp.field, fields...)
	}
	exp.tracer.trace(1)
	return &exp
}

func Wrap(msg *Message, cause error, fields ...field.Field) *Exception {
	if cause == nil {
		return nil
	}
	var exp = &Exception{message: msg, cause: cause, field: fields}
	exp.tracer.trace(1)
	return exp
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
