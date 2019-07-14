/*
Package stdlog provide log utilities.

It follow https://12factor.net/logs.

It parse environment variable "OnelineLog", if it true according to strconv.ParseBool,
then every call to Print is encoded to JSON String.

*/
package stdlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/payfazz/go-oneliner"
)

// Logger represent logging object.
type Logger struct {
	m sync.Mutex
	b io.Writer
}

var (
	// Out is wrapper of os.Stdout.
	Out *Logger

	// Err is wrapper of os.Stderr.
	Err *Logger
)

func init() {
	onelines, _ := strconv.ParseBool(os.Getenv("OnelineLog"))
	Out = new(os.Stdout, onelines)
	Err = new(os.Stderr, onelines)
}

func new(b io.Writer, onelines bool) *Logger {
	if onelines {
		b = oneliner.Wrap(b)
	}
	return &Logger{
		b: b,
	}
}

// Print the arguments. Arguments are handled in the manner of fmt.Print.
//
// Print is safe called from multiple goroutine, it guarantees to serialize access to the Writer.
//
// Print always ended with newline.
func (l *Logger) Print(v ...interface{}) {
	buff := getBuffer()
	buff.WriteString(fmt.Sprint(v...))
	if buff.Bytes()[buff.Len()-1] != '\n' {
		buff.WriteByte('\n')
	}

	l.m.Lock()
	l.b.Write(buff.Bytes())
	l.m.Unlock()

	putBuffer(buff)
}

// O is shortcut to Out.Print.
func O(v ...interface{}) {
	Out.Print(v...)
}

// E is shortcut to Err.Print.
func E(v ...interface{}) {
	Err.Print(v...)
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
