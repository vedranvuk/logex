// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkLogEmpty(b *testing.B) {
	b.StopTimer()
	l := New(nil)
	l.SetLevel(LevelPrint)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Println(42)
	}
}

type fakewriter struct{}

func (fw *fakewriter) Write(p []byte) (n int, err error) { return len(p), nil }

func BenchmarkLogSimple(b *testing.B) {
	b.StopTimer()
	l := New(nil)
	l.AddOutput("out", &fakewriter{}, NewSimpleFormatter())
	l.SetLevel(LevelPrint)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Println(42)
	}
}

func BenchmarkLogJSON(b *testing.B) {
	b.StopTimer()
	l := New(nil)
	l.AddOutput("out", &fakewriter{}, NewJSONFormatter(true))
	l.SetLevel(LevelPrint)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Println(42)
	}
}

func TestConcurrent(t *testing.T) {

	l := New(nil)

	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := file.Sync(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	bufjs := bytes.NewBuffer(nil)
	buftxt := bytes.NewBuffer(nil)
	l.AddOutput("out", bufjs, NewJSONFormatter(true))
	l.AddOutput("out", buftxt, NewSimpleFormatter())
	// l.AddOutput(file, NewJSONFormatter(true))
	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func(threadid int) {
			for i := 0; i < 5; i++ {
				err := fmt.Errorf("Error number '%d' occured in thread '%d'", i, threadid)
				msg := "Erroring..."
				l.WithCaller(1).WithStack(1, 5).Errorf(err, msg)
			}
			done <- true
		}(i)
	}

	for total := 0; total < 5; total++ {
		<-done
	}

	// fmt.Printf("%s\n", string(buftxt.Bytes()))
	fmt.Printf("%s\n", string(bufjs.Bytes()))

	ioutil.WriteFile("logtest.json", bufjs.Bytes(), os.ModePerm)
	ioutil.WriteFile("logtest.txt", buftxt.Bytes(), os.ModePerm)
}

type writer struct {
	prefix string
}

func newwriter(prefix string) *writer { return &writer{prefix} }

func (w *writer) Write(p []byte) (int, error) {
	fmt.Printf("%s: %s", w.prefix, string(p))
	return len(p), nil
}

func TestOutputs(t *testing.T) {

	l := New(nil)
	l.AddOutput("1", newwriter("1"), NewSimpleFormatter())
	l.AddOutput("2", newwriter("2"), NewSimpleFormatter())
	l.AddOutput("3", newwriter("3"), NewSimpleFormatter())
	l.AddOutput("4", newwriter("4"), NewSimpleFormatter())
	l.AddOutput("5", newwriter("5"), NewSimpleFormatter())

	sub := l.ToOutputs("2", "4")
	sub.Println(LevelDebug, "test")
}
