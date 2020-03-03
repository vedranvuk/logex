// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"errors"
	"os"
	"testing"
)

func TestLogLevel(t *testing.T) {
	testfunc := func(level LogLevel) {
		b, err := level.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		level = LevelMute
		if err := level.UnmarshalText(b); err != nil {
			t.Fatal(err)
		}
	}
	testfunc(LevelNone)
	testfunc(LevelMute)
	testfunc(LevelError)
	testfunc(LevelWarning)
	testfunc(LevelInfo)
	testfunc(LevelDebug)
	testfunc(LevelPrint)
	testfunc(LogLevel(42))
}

func DoLog(l Log) {
	lvl := LogLevel(42)
	l.Printf(lvl, "%s\n", "Printf")
	l.Println(lvl, "Println")
	l.Debugf("%s\n", "Debug")
	l.Debugln("Debugln")
	l.Infof("%s\n", "Info")
	l.Infoln("Infoln")
	l.Warningf("%s\n", "Warningf")
	l.Warningln("Warningln")
	l.Errorf(errors.New("error: catastrophic failure"), "%s\n", "Errorf")
	l.Errorln(errors.New("error: success"), "Errorln")
	l.Caller(1).Errorln(errors.New("error: caller"), "ERROR")
	l.Stack(0, 10).Errorln(errors.New("error: stack"), "ERROR")
}

func TestLog(t *testing.T) {
	l := NewStd()
	l.SetLevel(LevelPrint)
	DoLog(l)
}

func TestLog2(t *testing.T) {
	l := New()
	l.AddOutput(os.Stdout, NewJSONFormatter(true))
	DoLog(l)
}

func TestLog3(t *testing.T) {
	l := NewStd()
	l.SetLevel(LevelError)
	l.Println(42, "testis")
}

func TestLog4(t *testing.T) {
	l := NewStd()
	l.SetLevel(LevelPrint)
	f := NewFields()
	f.Set("mirko", "odora")
	l.Fields(f).Println(64, "test")
}

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
