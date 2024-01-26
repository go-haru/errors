package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type Stack struct {
	Frame uintptr `json:"frame"`
	Step  uintptr `json:"step"`
	Func  string  `json:"func"`
	File  string  `json:"file"`
}

type Tracer interface {
	Frame() []uintptr
	Trace([]uintptr) []Stack
}

type tracer []uintptr

func (t *tracer) Frame() []uintptr { return *t }

func (t *tracer) trace(skip int) {
	pcs := make([]uintptr, 255)
	n := runtime.Callers(skip+2, pcs)
	*t = pcs[:n]
}

func (t *tracer) Trace(parentFrame []uintptr) []Stack {
	var currentStack = *t
	var sameFrames int
trimSameStack:
	for pi, pf := range parentFrame {
		for ti, tf := range currentStack[:len(currentStack)-1-(len(parentFrame)-1-pi)+1] {
			if pf == tf {
				sameFrames = len(currentStack) - (ti + 1) + 1
				break trimSameStack
			}
		}
	}
	var infoStack = make([]Stack, 0, len(currentStack)-sameFrames)
	frames := runtime.CallersFrames(currentStack)
	var more bool
	var frame runtime.Frame
	var hidden bool
	for i := len(currentStack) - 1 - sameFrames; i >= 0; i-- {
		switch frame, more = frames.Next(); {
		case frame.Function == "runtime.main",
			frame.Function == "runtime.goexit",
			strings.HasPrefix(frame.Function, "github.com/spf13/cobra."):
			if !hidden {
				infoStack, hidden = append(infoStack, Stack{Func: "// package or function omitted", File: "........"}), true
			}
		default:
			infoStack, hidden = append(infoStack, Stack{
				Frame: frame.PC,
				Step:  uintptr(sameFrames + i),
				Func:  strings.ReplaceAll(frame.Function, "%2e", "."),
				File:  fmt.Sprint(strings.Replace(frame.File, "@v0.0.0-00010101000000-000000000000", "", 1), ":", frame.Line),
			}), false
			if !more {
				break
			}
		}
	}
	return infoStack
}
