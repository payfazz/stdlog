package stdlog

import (
	"fmt"
	"io"
	"log"
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
	inner     io.Writer
}

// New create new logger that write to b,
//
// if timestamp is true, every Print will prefixed with timestamp.
//
// if onelines is true, every Write will formated into JSON string,
// this option is ignored when w is *Logger
func New(w io.Writer, prefix string, timestamp bool, onelines bool) *Logger {
	if logger, ok := w.(*Logger); ok {
		// already Logger, inherit the Writer
		w = logger.inner
	} else if onelines {
		w = oneliner.Wrap(w)
	}

	return &Logger{
		prefix:    prefix,
		timestamp: timestamp,
		inner:     w,
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
	defer l.mu.Unlock()

	n, err = l.inner.Write(p)

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

	for i, value := range v {
		if i != 0 {
			buff.WriteByte(' ')
		}
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

// AsLogger return log.Logger instance with l as inner
func (l *Logger) AsLogger() *log.Logger {
	return log.New(asLoggerWrapper{inner: l}, "", 0)
}

type asLoggerWrapper struct {
	inner *Logger
}

func (w asLoggerWrapper) Write(p []byte) (n int, err error) {
	w.inner.Print(string(p))
	return len(p), nil
}
