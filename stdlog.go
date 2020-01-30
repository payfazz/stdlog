/*
Package stdlog provide log utilities.

It parse environment variable "OnelineLog", if it true according to strconv.ParseBool,
then every call to Print is encoded to JSON String to make sure that every log is one line.

*/
package stdlog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/payfazz/go-oneliner"
)

// Printer is interface for printing log
//
// NOTE: log.Logger implement this interface
type Printer interface {
	Print(...interface{})
}

// Logger represent logging object.
type Logger struct {
	m sync.Mutex
	b io.Writer
}

// static type check
var (
	_ io.Writer = (*Logger)(nil)
	_ Printer   = (*Logger)(nil)
)

var (
	// Onelines is derived from OnelineLog env variable according to strconv.ParseBool
	Onelines bool

	// Out is wrapper of os.Stdout.
	Out *Logger

	// Err is wrapper of os.Stderr.
	Err *Logger

	// Discard is nop Logger
	Discard Printer = New(ioutil.Discard, false)
)

func init() {
	Onelines, _ = strconv.ParseBool(os.Getenv("OnelineLog"))
	Out = New(os.Stdout, Onelines)
	Err = New(os.Stderr, Onelines)
}

// New create new logger that write to b, if onelines
func New(b io.Writer, onelines bool) *Logger {
	if onelines && b != ioutil.Discard {
		b = oneliner.Wrap(b)
	}
	return &Logger{
		b: b,
	}
}

// PrintOut wrap a call to Out.Print() method
func PrintOut(v ...interface{}) {
	Out.Print(v)
}

// PrintErr wrap a call to Err.Print() method
func PrintErr(v ...interface{}) {
	Err.Print(v)
}

// Print from Printer interface
//
// Print the arguments. Arguments are handled in the manner of fmt.Print.
//
// Print is safe called from multiple goroutine, it guarantees to serialize access to the Writer.
//
// Print always ended with newline.
func (l *Logger) Print(v ...interface{}) {
	if l.b == ioutil.Discard {
		return
	}

	buff := getBuffer()

	for _, value := range v {
		if s, ok := value.(string); ok {
			buff.WriteString(s)
		} else {
			buff.WriteString(fmt.Sprint(value))
		}
	}

	if buff.Bytes()[buff.Len()-1] != '\n' {
		buff.WriteByte('\n')
	}

	l.Write(buff.Bytes())

	putBuffer(buff)
}

// Write from io.Writer interface
//
// Write is safe called from multiple goroutine, it guarantees to serialize access to the Writer.
func (l *Logger) Write(p []byte) (int, error) {
	if l.b == ioutil.Discard {
		return len(p), nil
	}

	l.m.Lock()
	n, err := l.b.Write(p)
	l.m.Unlock()
	return n, err
}

var pool sync.Pool

func getBuffer() *bytes.Buffer {
	if x := pool.Get(); x != nil {
		b := x.(*bytes.Buffer)
		b.Reset()
		return b
	}
	return &bytes.Buffer{}
}

func putBuffer(b *bytes.Buffer) {
	pool.Put(b)
}
