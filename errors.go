package xErr

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Defined severities
type sev int
const (
	severe sev = iota
	high
	low
	info
)

type extendedErr interface {
	GetSeverity() sev
	SetSeverity(severity sev)
	Error() string
}

type xErr struct {
	err error
	severity	sev
}

func (x *xErr) Error() string {
	return x.err.Error()
}

func (x *xErr) GetSeverity() sev {
	return x.severity
}

func (x *xErr) SetSeverity(severity sev)  {
	x.severity = severity
}

var (
	osExit           = os.Exit
	output io.Writer = os.Stdout // modified during testing
)

func UnwrapError(err error) error {
	nestedError := errors.Unwrap(err)
	if nestedError == nil {
		return err
	} else {
		fmt.Println(nestedError.Error())
		return UnwrapError(nestedError)
	}
}

func WrapError(id, new error) error {
	return fmt.Errorf("%w:\n    %v", id, new)
}

func HandleError(err extendedErr) {
	_, _ = fmt.Fprintf(output, "Error:\n  "+err.Error()+"\n")

	//sourceError := UnwrapError(xErr{})

	// severe, high and low severities break the execution of the code
	switch e := err.GetSeverity(); e {
	case severe: os.Exit(2)
	case high: os.Exit(1)
	case low: os.Exit(0)
	}
}
