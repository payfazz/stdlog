package stdlog

import (
	"bytes"
	"testing"
)

func assert(t *testing.T, got, expected string) {
	if got != expected {
		t.Fatalf("got %s, expected %s", got, expected)
	}
}

func TestLogger1(t *testing.T) {
	buff := &bytes.Buffer{}
	a := New(buff, "", false, false)
	a.Print("test")
	a.Print("test", "test")
	a.Print("test\n")
	a.Print("test", "test\n")
	assert(t, buff.String(), "test\ntesttest\ntest\ntesttest\n")
}

func TestLogger2(t *testing.T) {
	buff := &bytes.Buffer{}
	a := New(buff, "", false, true)
	a.Print("test")
	a.Print("test", "test")
	a.Print("test\n")
	a.Print("test", "test\n")
	assert(t, buff.String(), "\"test\\n\"\n\"testtest\\n\"\n\"test\\n\"\n\"testtest\\n\"\n")
}

func BenchmarkWithoutNewline(b *testing.B) {
	buff := &bytes.Buffer{}
	a := New(buff, "", false, false)
	for i := 0; i < b.N; i++ {
		buff.Reset()
		a.Print("test")
		a.Print("test", "test")
	}
}

func BenchmarkWithNewline(b *testing.B) {
	buff := &bytes.Buffer{}
	a := New(buff, "", false, false)
	for i := 0; i < b.N; i++ {
		buff.Reset()
		a.Print("test\n")
		a.Print("test", "test\n")
	}
}
