package errors_test

import (
	"fmt"
	"github.com/Bowery/errors"
	"io"
	"strconv"
)

func ExampleNewStackError() {
	err := errors.NewStackError(io.EOF)
	stackErr := errors.IsStackError(err)

	fmt.Println(stackErr.Err == io.EOF) // The underlying error.
	fmt.Println(stackErr.Trace.Exception.Message)

	for _, frame := range stackErr.Trace.Frames {
		fmt.Println("File:", frame.File+":"+strconv.Itoa(frame.Line), "at", frame.Method)
	}

	// Or use Stack to print a nicely formatted stack trace.
	fmt.Println(stackErr.Stack())
}

func ExampleIsStackError() {
	err := errors.NewStackError(errors.New("github.com/Bowery/errors: Some error happened"))
	stackErr := errors.IsStackError(err)
	if stackErr != nil {
		fmt.Println("Yes it's a stack error.")
		fmt.Println(stackErr.Stack())
	}

	err = io.EOF
	if errors.IsStackError(err) == nil {
		fmt.Println("It is not a stack error.")
		fmt.Println(err)
	}
}
