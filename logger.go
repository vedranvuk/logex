// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"io"
	"os"
	"sync"
)

// Logger is an implementation of Log.
type Logger struct {
	Log

	mu  sync.Mutex
	wrs map[io.Writer]Formatter
	lvl LogLevel
}

// print prints fields to registered writers using associated formatters.
func (l *Logger) print(fields *Fields) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if fields.LogLevel() > l.lvl {
		return
	}
	for writer, formatter := range l.wrs {
		writer.Write([]byte(formatter.Format(fields)))
	}
}

// AddOutput adds an io.Writer to Logger to be written to using the specified formatter f.
// Duplicates are replaced, if found.
func (l *Logger) AddOutput(w io.Writer, f Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.wrs[w] = f
}

// SetLevel sets Logger's LogLevel.
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lvl = level
}

// New returns a new Logger with no defined outputs.
// Initial logging level is set to LevelDebug.
func New() *Logger {
	p := &Logger{
		mu:  sync.Mutex{},
		wrs: make(map[io.Writer]Formatter),
		lvl: LevelDebug,
	}
	p.Log = NewLine(p)
	return p
}

// NewStd returns a new Logger initialized to stdout using a default formatter.
func NewStd() *Logger {
	p := New()
	p.AddOutput(os.Stdout, NewSimpleFormatter())
	return p
}
