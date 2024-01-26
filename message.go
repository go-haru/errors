package errors

import (
	"fmt"

	"github.com/go-haru/field"
)

type Message struct {
	code    uint32
	message string
}

func NewMessage(code uint32, msg string) *Message { return &Message{code: code, message: msg} }

func (m *Message) Code() uint32 { return m.code }

func (m *Message) Message() string { return m.message }

func (m *Message) Is(target error) bool {
	if m == nil {
		return target == nil
	}
	if target == nil {
		return false
	}
	if m.code > 0 {
		switch t := target.(type) {
		case interface{ Code() uint32 }:
			return t.Code() == m.code
		case interface{ Code() int }:
			return t.Code() == int(m.code)
		}
	}
	switch t := target.(type) {
	case interface{ Message() string }:
		return t.Message() == m.message
	default:
		return t.Error() == m.message
	}
}

func (m *Message) Error() string {
	if m.code != 0 {
		return fmt.Sprintf("%d: %s", m.code, m.message)
	} else {
		return m.message
	}
}

func (m *Message) With(fields ...field.Field) *Exception {
	var exp = &Exception{message: m, field: fields}
	exp.tracer.trace(1)
	return exp
}

func (m *Message) Wrap(cause error, fields ...field.Field) *Exception {
	var exp = &Exception{message: m, cause: cause, field: fields}
	exp.tracer.trace(1)
	return exp
}
