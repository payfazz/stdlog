package stdlog

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/payfazz/go-oneliner"
)

const timeFormat = "2006-01-02T15:04:05.000Z07:00"

// Logger represent logging object.
type Logger struct {
	mu sync.Mutex

	prefix    string
	timestamp bool

	io.Writer
}

// New create new logger that write to b,
//
// if timestamp is true, every Print will prefixed with timestamp,
// if onelines is true, every Write will formated into JSON string,
func New(w io.Writer, prefix string, timestamp bool, onelines bool) *Logger {
	if onelines {
		w = oneliner.Wrap(w)
	}

	return &Logger{
		prefix:    prefix,
		timestamp: timestamp,

		Writer: w,
	}
}

// Write from io.Writer interface
//
// Write is safe called from multiple goroutine, it guarantees to serialize access to the Writer.
func (l *Logger) Write(p []byte) (n int, err error) {
	if l == nil {
		return len(p), nil
	}

	l.mu.Lock()
	n, err = l.Writer.Write(p)
	l.mu.Unlock()

	return
}

// Print from Printer interface
//
// Print the arguments. Arguments are handled in the manner of fmt.Print.
//
// Print is safe called from multiple goroutine, it guarantees to serialize access to the Writer.
func (l *Logger) Print(v ...interface{}) {
	if l == nil {
		return
	}

	buff := getBuffer()

	if l.prefix != "" {
		buff.WriteString(l.prefix)
	}

	if l.timestamp {
		buff.WriteString(time.Now().Format(timeFormat))
		buff.WriteByte(' ')
	}

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
