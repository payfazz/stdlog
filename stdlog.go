/*
Package stdlog provide log utilities.
*/
package stdlog

import (
	"errors"
	"io"
	"os"
	"strconv"
	"sync/atomic"
)

// Printer is interface for printing log
type Printer interface {
	Print(...interface{})
}

var (
	_OnelineLog     bool
	_NoTimestampLog bool
)

func init() {
	_OnelineLog, _ = strconv.ParseBool(os.Getenv("OnelineLog"))
	_NoTimestampLog, _ = strconv.ParseBool(os.Getenv("NoTimestampLog"))
}

// New2 is same as New but timestamp and onelines are inherited from env
func New2(w io.Writer, prefix string) *Logger {
	return New(os.Stdout, prefix, !_NoTimestampLog, _OnelineLog)
}

// OnelineLog is derived from "OnelineLog" env variable according to strconv.ParseBool
func OnelineLog() bool {
	return _OnelineLog
}

// NoTimestampLog is derived from "NoTimestampLog" env variable according to strconv.ParseBool
func NoTimestampLog() bool {
	return _NoTimestampLog
}

var (
	out atomic.Value // *Logger
	err atomic.Value //*Logger
)

// Out is wrapper of os.Stdout.
//
// Wrapped logger will behave acording to OnelineLog and NoTimestampLog env
func Out() *Logger {
	if out.Load() == nil {
		out.Store(New2(os.Stdout, ""))
	}

	return out.Load().(*Logger)
}

// SetOut set logger returned by Out, only can be called once before Out is called
func SetOut(l *Logger) error {
	if out.Load() != nil {
		return errors.New("Out already set")
	}

	out.Store(l)
	return nil
}

// Err is wrapper of os.Stderr.
//
// Wrapped logger will behave acording to OnelineLog and NoTimestampLog env
func Err() *Logger {
	if err.Load() == nil {
		err.Store(New2(os.Stderr, ""))
	}

	return err.Load().(*Logger)
}

// SetErr set logger returned by Err, only can be called once before Err is called
func SetErr(l *Logger) error {
	if err.Load() != nil {
		return errors.New("Out already set")
	}

	err.Store(l)
	return nil
}

// PrintOut is shortcut to Out().Print
func PrintOut(v ...interface{}) {
	Out().Print(v...)
}

// PrintErr is shortcut to Err().Print
func PrintErr(v ...interface{}) {
	Err().Print(v...)
}
