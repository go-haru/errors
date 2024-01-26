package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-haru/field"
)

type Exception struct {
	message *Message
	cause   error
	tracer  tracer
	field   field.Fields
}

func (e *Exception) Code() uint32 { return e.message.code }

func (e *Exception) Unwrap() error { return e.cause }

func (e *Exception) Is(target error) bool {
	if e == nil {
		return target == nil
	}
	if target == nil {
		return false
	}
	if e.cause != nil {
		if errors.Is(e.cause, target) {
			return true
		}
	}
	if e.message != nil {
		if errors.Is(e.message, target) {
			return true
		}
	}
	return false
}

func (e *Exception) Tracer() Tracer { return &e.tracer }

func (e *Exception) Error() string {
	var buf strings.Builder
	e.encodeText(&buf, nil)
	return buf.String()
}

func (e *Exception) encodeText(buf *strings.Builder, parentFrame []uintptr) {
	if e.cause != nil {
		var parentExp interface {
			error
			encodeText(buf *strings.Builder, parentFrame []uintptr)
		}
		if errors.As(e.cause, &parentExp) {
			parentExp.encodeText(buf, e.tracer.Frame())
		} else {
			buf.WriteString(e.cause.Error())
		}
	}
	if e.message != nil {
		buf.WriteRune('\n')
		buf.WriteString(e.message.Error())
	}
	if len(e.field) > 0 {
		buf.WriteString(" # ")
		_ = e.field.EncodeJSON(buf)
	}
	for _, stack := range e.Tracer().Trace(parentFrame) {
		buf.WriteString("\n\x20\x20")
		if stack.Step > 0 {
			buf.WriteString(fmt.Sprintf("[%d] ", stack.Step))
		} else {
			buf.WriteString("[*] ") // omitted
		}
		buf.WriteString(stack.Func)
		if stack.File != "" {
			buf.WriteString("\n\x20\x20\x20\x20")
			buf.WriteString(stack.File)
		}
	}
}
