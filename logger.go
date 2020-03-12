// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"io"
	"os"
	"sync"
)

type output struct {
	w io.Writer
	f Formatter
}

type outputmap map[string]*output

type ErrorFunc func(err error)

// Logger is an implementation of Log.
type Logger struct {
	Log

	mu      sync.Mutex
	outputs outputmap
	lvl     LogLevel
	ef      ErrorFunc
}

// print prints fields to registered writers using associated formatters.
func (l *Logger) print(fields *Fields, outputnames ...string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if fields.LogLevel() > l.lvl {
		return
	}
	var err error
	var out *output
	var ok bool
	if len(outputnames) > 0 {
		for _, name := range outputnames {
			if out, ok = l.outputs[name]; ok {
				if _, err = out.w.Write([]byte(out.f.Format(fields))); err != nil && l.ef != nil {
					l.ef(err)
				}
			}
		}
	} else {
		for _, out = range l.outputs {
			if _, err = out.w.Write([]byte(out.f.Format(fields))); err != nil && l.ef != nil {
				l.ef(err)
			}
		}
	}
	l.Log = NewLine(l)
}

// AddOutput registers an output writer with formatter f unser specified
// name which must be unique and not empty or returns an error.
func (l *Logger) AddOutput(name string, w io.Writer, f Formatter) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if name == "" {
		return ErrInvalidName
	}
	if _, exists := l.outputs[name]; exists {
		return ErrDuplicateName.WrapArgs(name)
	}
	l.outputs[name] = &output{w, f}
	return nil
}

// SetLevel sets Logger's LogLevel.
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lvl = level
}

// New returns a new Logger with no defined outputs.
// Initial logging level is set to LevelDebug.
func New(ef ErrorFunc) *Logger {
	p := &Logger{
		mu:      sync.Mutex{},
		outputs: make(outputmap),
		lvl:     LevelDebug,
		ef:      ef,
	}
	p.Log = NewLine(p)
	return p
}

// NewStd returns a new Logger initialized to stdout using a default formatter.
func NewStd(ef ErrorFunc) *Logger {
	p := New(ef)
	p.AddOutput("stdout", os.Stdout, NewSimpleFormatter())
	return p
}
