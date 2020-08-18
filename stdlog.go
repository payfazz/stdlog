/*
Package stdlog provide log utilities.
*/
package stdlog

import (
	"errors"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

// Printer is interface for printing log
type Printer interface {
	Print(...interface{})
}

// OnelineLog is derived from "OnelineLog" env variable according to strconv.ParseBool
func OnelineLog() bool {
	onelineLog, _ := strconv.ParseBool(os.Getenv("OnelineLog"))
	return onelineLog
}

// TimestampLog is derived from "TimestampLog" env variable according to strconv.ParseBool
func TimestampLog() bool {
	noTimestampLog, _ := strconv.ParseBool(os.Getenv("TimestampLog"))
	return noTimestampLog
}

// NewFromEnv is same as New but timestamp and onelines are inherited from env
func NewFromEnv(w io.Writer, prefix string) *Logger {
	return New(w, prefix, TimestampLog(), OnelineLog())
}

var (
	outOnce sync.Once
	out     atomic.Value // *Logger

	errOnce sync.Once
	err     atomic.Value // *Logger
)

// Out return global out Logger
//
// if SetOut never called before, will return Logger that wrap os.Stdout and
// behave acording to OnelineLog and NoTimestampLog
func Out() *Logger {
	outOnce.Do(func() {
		SetOut(NewFromEnv(os.Stdout, ""))
	})

	return out.Load().(*Logger)
}

// SetOut set logger returned by Out
func SetOut(l *Logger) error {
	if l == nil {
		panic("l cannot be nil")
	}

	if out.Load() != nil {
		return errors.New("Out already set")
	}
	out.Store(l)

	return nil
}

// Err return global err Logger
//
// if SetErr never called before, will return Logger that wrap os.Stderr and
// behave acording to OnelineLog and NoTimestampLog
func Err() *Logger {
	errOnce.Do(func() {
		SetErr(NewFromEnv(os.Stderr, ""))
	})

	return err.Load().(*Logger)
}

// SetErr set logger returned by Err
func SetErr(l *Logger) error {
	if l == nil {
		panic("l cannot be nil")
	}

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
