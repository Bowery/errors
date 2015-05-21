// Copyright 2015 Bowery, Inc.

package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// New creates a new error, this solves issue of name collision with
// errors pkg.
func New(args ...interface{}) error {
	return errors.New(strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Newf creates a new error, from an existing error template.
func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// StackError is an error with stack information.
type StackError struct {
	Err   error
	Trace *Trace
}

// IsStackError returns the error as a StackError if it's a StackError, nil
// otherwise.
func IsStackError(err error) *StackError {
	se, ok := err.(*StackError)
	if ok {
		return se
	}

	return nil
}

// NewStackError creates a stack error including the stack.
func NewStackError(err error) error {
	se := &StackError{
		Err: err,
		Trace: &Trace{
			Frames:    make([]*Frame, 0),
			Exception: &Exception{Message: err.Error(), Class: errClass(err)},
		},
	}

	// Get stack frames excluding the current one.
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			// Couldn't get another frame, so we're finished.
			break
		}

		f := &Frame{File: file, Line: line, Method: routineName(pc)}
		se.Trace.Frames = append(se.Trace.Frames, f)
	}

	return se
}

func (se *StackError) Error() string {
	return se.Err.Error()
}

// Stack prints the stack trace in a readable format.
func (se *StackError) Stack() string {
	stack := ""

	for i, frame := range se.Trace.Frames {
		stack += strconv.Itoa(i+1) + ": File \"" + frame.File + "\" line "
		stack += strconv.Itoa(frame.Line) + " in " + frame.Method + "\n"
	}
	stack += se.Trace.Exception.Class + ": " + se.Trace.Exception.Message

	return stack
}

// Trace contains the stack frames, and the exception information.
type Trace struct {
	Frames    []*Frame   `json:"frames"`
	Exception *Exception `json:"exception"`
}

// Exception contains the error message and it's class origin.
type Exception struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}

// Frame contains line, file and method info for a stack frame.
type Frame struct {
	File   string `json:"filename"`
	Line   int    `json:"lineno"`
	Method string `json:"method"`
}

// errClass retrieves the string representation for the errors type.
func errClass(err error) string {
	class := strings.TrimPrefix(reflect.TypeOf(err).String(), "*")
	if class == "" {
		class = "panic"
	}

	return class
}

// routineName returns the routines name for a given program counter.
func routineName(pc uintptr) string {
	fc := runtime.FuncForPC(pc)
	if fc == nil {
		return "???"
	}

	return fc.Name() // Includes the package info.
}
