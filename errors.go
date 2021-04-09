package xErr

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Defined severities
type Sev int

const (
	Severe Sev = iota
	High
	Low
	Info
)

type ExtendedErr interface {
	GetSeverity() Sev
	SetSeverity(severity Sev)
	Error() string
	Handle()
}

type xErr struct {
	err      error
	severity Sev
}

func (x *xErr) Error() string {
	return x.err.Error()
}

func NewError(newErr error, severity Sev) xErr {
	e := xErr{err: newErr, severity: severity}
	return e
}

func (x *xErr) GetSeverity() Sev {
	return x.severity
}

func (x *xErr) SetSeverity(severity Sev) {
	x.severity = severity
}

func (x *xErr) Handle() {
	_, _ = fmt.Fprintf(output, "Error:\n  "+x.Error()+"\n")

	//sourceError := UnwrapError(xErr{})

	// severe, high and low severities break the execution of the code
	switch e := x.GetSeverity(); e {
	case Severe:
		os.Exit(2)
	case High:
		os.Exit(1)
	case Low:
		os.Exit(0)
	}
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


