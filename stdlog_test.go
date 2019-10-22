package stdlog

import (
	"bytes"
	"testing"
)

func TestLogger1(t *testing.T) {
	buff := &bytes.Buffer{}
	a := New(buff, false)
	a.Print("test")
	a.Print("test", "test")
	a.Print("test\n")
	a.Print("test", "test\n")
	if buff.String() != "test\ntesttest\ntest\ntesttest\n" {
		t.FailNow()
	}
}

func TestLogger2(t *testing.T) {
	buff := &bytes.Buffer{}
	a := New(buff, true)
	a.Print("test")
	a.Print("test", "test")
	a.Print("test\n")
	a.Print("test", "test\n")
	if buff.String() != "\"test\\n\"\n\"testtest\\n\"\n\"test\\n\"\n\"testtest\\n\"\n" {
		t.FailNow()
	}
}

func BenchmarkWithoutNewline(b *testing.B) {
	buff := &bytes.Buffer{}
	a := New(buff, false)
	for i := 0; i < b.N; i++ {
		buff.Reset()
		a.Print("test")
		a.Print("test", "test")
	}
}

func BenchmarkWithNewline(b *testing.B) {
	buff := &bytes.Buffer{}
	a := New(buff, false)
	for i := 0; i < b.N; i++ {
		buff.Reset()
		a.Print("test\n")
		a.Print("test", "test\n")
	}
}
