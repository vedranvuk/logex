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
	l := New()
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
	l := New()
	l.AddOutput(&fakewriter{}, NewSimpleFormatter())
	l.SetLevel(LevelPrint)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Println(42)
	}
}

func BenchmarkLogJSON(b *testing.B) {
	b.StopTimer()
	l := New()
	l.AddOutput(&fakewriter{}, NewJSONFormatter(true))
	l.SetLevel(LevelPrint)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Println(42)
	}
}

func TestConcurrent(t *testing.T) {

	l := New()

	bufjs := bytes.NewBuffer(nil)
	buftxt := bytes.NewBuffer(nil)
	l.AddOutput(bufjs, NewJSONFormatter(false))
	l.AddOutput(buftxt, NewSimpleFormatter())
	done := make(chan bool)

	for i := 0; i < 16; i++ {
		go func(threadid int) {
			for i := 0; i < 100; i++ {
				err := fmt.Errorf("Error number '%d' occured in thread '%d'", i, threadid)
				msg := "Erroring..."
				l.Caller(1).Stack(1, 5).Errorf(err, msg)
			}
			done <- true
		}(i)
	}

	for total := 0; total < 16; total++ {
		<-done
	}

	ioutil.WriteFile("logtest.json", bufjs.Bytes(), os.ModePerm)
	ioutil.WriteFile("logtest.txt", buftxt.Bytes(), os.ModePerm)
}
