/*
Package stdlog provide log utilities.

It follow https://12factor.net/logs.

It parse environment variable "OnelineLog", if it true according to strconv.ParseBool,
then every call to Print is encoded to JSON String.

*/
package stdlog

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/payfazz/go-oneliner"
)

// Logger represent logging object.
type Logger struct {
	m sync.Mutex
	b io.Writer
}

var (
	// Out is wrapper of os.Stdout.
	Out = new(os.Stdout)

	// Err is wrapper of os.Stderr.
	Err = new(os.Stderr)
)

func new(b io.Writer) *Logger {
	onelines, _ := strconv.ParseBool(os.Getenv("OnelineLog"))
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
	var s string
	if len(v) == 1 {
		if s2, ok := v[0].(string); ok {
			s = s2
			if s[len(s)-1] != '\n' {
				s += "\n"
			}
		}
	}
	if s == "" {
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprint(v...))
		if sb.String()[sb.Len()-1] != '\n' {
			sb.WriteByte('\n')
		}
		s = sb.String()
	}

	// peform unsafe zero-copy conversion from string to byte slice
	sHeader := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sHeader.Data,
		Len:  sHeader.Len,
		Cap:  sHeader.Len,
	}))

	l.m.Lock()
	l.b.Write(bytes) // safe because Write only *read* the content
	l.m.Unlock()

	// make sure s live until here
	runtime.KeepAlive(s)
}

// O is shortcut to Out.Print.
func O(v ...interface{}) {
	Out.Print(v...)
}

// E is shortcut to Err.Print.
func E(v ...interface{}) {
	Err.Print(v...)
}
