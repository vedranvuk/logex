// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Line implements the Log interface.
// It forms a log line from standard properties like timestamp, message, stack
// and optional user defined values.
type Line struct {
	mu     sync.Mutex
	fields *Fields
	log    *Logger
	cloned bool
}

// NewLine returns a new Line instance that will output to Logger l.
func NewLine(l *Logger) *Line {
	return &Line{
		mu:     sync.Mutex{},
		fields: NewFields(),
		log:    l,
	}
}

// flush outputs line fields to the Logger.
func (p *Line) flush(level LogLevel, message string) {
	p.fields.set(KeyLogLevel, level)
	p.fields.set(KeyMessage, message)
	p.fields.set(KeyTime, time.Now())
	p.log.print(p.fields)
}

func (p *Line) Debugf(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelDebug, fmt.Sprintf(format, args...))
}

func (p *Line) Debugln(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelDebug, fmt.Sprint(args...)+"\n")
}

func (p *Line) Infof(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelInfo, fmt.Sprintf(format, args...))
}

func (p *Line) Infoln(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelInfo, fmt.Sprint(args...)+"\n")
}

func (p *Line) Warningf(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelWarning, fmt.Sprintf(format, args...))
}

func (p *Line) Warningln(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelWarning, fmt.Sprint(args...)+"\n")
}

func (p *Line) Errorf(err error, format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.fields.set(KeyError, err)
	p.flush(LevelError, fmt.Sprintf(format, args...))
}

func (p *Line) Errorln(err error, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.fields.set(KeyError, err)
	p.flush(LevelError, fmt.Sprint(args...)+"\n")
}

func (p *Line) Printf(level LogLevel, format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelPrint, fmt.Sprintf(format, args...))
}

func (p *Line) Println(level LogLevel, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelPrint, fmt.Sprint(args...)+"\n")
}

// lazyclone returns a clone of self if not already cloned.
func (p *Line) lazyclone() *Line {
	if p.cloned {
		return p
	}
	nl := NewLine(p.log)
	nl.cloned = true
	p.fields.Walk(func(key FieldKey, val interface{}) bool {
		nl.fields.set(key, val)
		return true
	})
	return nl
}

// Caller appends the caller fields to the Line.
func (p *Line) Caller(skip int) Log {
	p.mu.Lock()
	defer p.mu.Unlock()
	l := p.lazyclone()
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		l.fields.set(KeyFile, file)
		l.fields.set(KeyLine, line)
	}
	return l
}

// Stack appends the stack to the Line.
func (p *Line) Stack(skip, depth int) Log {
	p.mu.Lock()
	defer p.mu.Unlock()
	l := p.lazyclone()
	callers := make([]uintptr, depth)
	if runtime.Callers(skip, callers) > 0 {
		frames := runtime.CallersFrames(callers)
		frameslice := []*Fields{}
		for i := 0; i < depth; i++ {
			frame, more := frames.Next()
			f := NewFields()
			f.set(KeyFile, frame.File)
			f.set(KeyLine, frame.Line)
			f.set(KeyFunc, frame.Func.Name())
			frameslice = append(frameslice, f)
			if !more {
				break
			}
		}
		l.fields.set(KeyFrames, frameslice)
	}
	return l
}

// Fields appends custom fields to Line.
func (p *Line) Fields(fields *Fields) Log {
	p.mu.Lock()
	defer p.mu.Unlock()
	l := p.lazyclone()
	fields.Walk(func(key FieldKey, val interface{}) bool {
		l.fields.Set(key, val)
		return true
	})
	return l
}
