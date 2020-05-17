package stdlog

import (
	"bytes"
	"sync"
)

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
