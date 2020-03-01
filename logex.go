// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package logex implements a logger.
package logex

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/vedranvuk/errorex"
)

var (
	// ErrLog is a base error of log package.
	ErrLog = errorex.New("log")
	// ErrUnmarshalLevel is returned when unmarshaling an invalid value from text as LogLevel.
	ErrUnmarshalLevel = ErrLog.WrapFormat("error unmarshaling '%s' as loglevel")
	// ErrReservedKey
	ErrReservedKey = ErrLog.WrapFormat("cannot set field '%s', key is reserved")
)

// Log is the public interface to Logger.
type Log interface {
	Debugf(string, ...interface{})
	Debugln(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})
	Warningf(string, ...interface{})
	Warningln(...interface{})
	Errorf(error, string, ...interface{})
	Errorln(error, ...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	Caller(int) Log
	Stack(int, int) Log
	Fields(Fields) Log
}

// LogLevel defines a logging level for a Logger.
type LogLevel byte

const (
	// LevelNone is undefined.
	LevelNone LogLevel = iota
	// LevelMute is the silent logging level.
	// It is used to silence the logger.
	LevelMute
	// LevelError is the error logging level that contains important error information.
	LevelError
	// LevelWarning is the warning logging level that might be important to user.
	LevelWarning
	// LevelInfo is the info logging level.
	// It should include information notices that might be useful to user.
	LevelInfo
	// LevelDebug is the debug logging level.
	// It should contain verbose debugging information that is useful in debugging.
	LevelDebug
	// LevelCustom and above is a custom logging level.
	LevelCustom
	// LevelPrint is the print logging level that is always printed, unless LevelMute.
	LevelPrint = LogLevel(255)
)

// String implements the Stringer interface.
func (ll LogLevel) String() string {
	switch ll {
	case LevelDebug:
		return "Debug"
	case LevelInfo:
		return "Info"
	case LevelWarning:
		return "Warning"
	case LevelError:
		return "Error"
	case LevelPrint:
		return "Print"
	}
	return ""
}

// MarshalText implements the TextMarshaler interface.
func (ll *LogLevel) MarshalText() ([]byte, error) {
	return []byte(ll.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface.
func (ll *LogLevel) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "debug":
		*ll = LevelDebug
	case "info":
		*ll = LevelInfo
	case "warning":
		*ll = LevelWarning
	case "error":
		*ll = LevelError
	case "print":
		*ll = LevelPrint
	default:
		return ErrUnmarshalLevel.WrapArgs(string(text))
	}
	return nil
}

// Logger is an implementation of Log.
type Logger struct {
	Log

	mu  sync.Mutex
	wrs map[io.Writer]Formatter
	lvl LogLevel
}

// print prints fields to registered writers using associated formatters.
func (l *Logger) print(fields Fields) {
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

// Line defines a log line consisting of fields populted by logging calls which
// are ultimately passed to a Logger that converts them to log lines using a Formatter.
type Line struct {
	log    *Logger
	fields Fields
	mu     sync.Mutex
}

// NewLine returns a new Line instance that will output to Logger l.
func NewLine(l *Logger) *Line {
	return &Line{
		log:    l,
		fields: make(Fields),
		mu:     sync.Mutex{},
	}
}

// flush outputs line fields to the Logger.
func (p *Line) flush(level LogLevel, message string) {
	p.fields[KeyLogLevel] = level
	p.fields[KeyMessage] = message
	p.fields[KeyTime] = time.Now()
	p.log.print(p.fields)
	p.fields = make(Fields)
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
	p.fields[KeyError] = err
	p.flush(LevelError, fmt.Sprintf(format, args...))
}

func (p *Line) Errorln(err error, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.fields[KeyError] = err
	p.flush(LevelError, fmt.Sprint(args...)+"\n")
}

func (p *Line) Printf(format string, args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelPrint, fmt.Sprintf(format, args...))
}

func (p *Line) Println(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flush(LevelPrint, fmt.Sprint(args...)+"\n")
}

// Caller appends the caller fields to the Line.
func (p *Line) Caller(skip int) Log {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		p.fields[KeyFile] = file
		p.fields[KeyLine] = line
	}
	return p
}

// Stack appends the stack to the Line.
func (p *Line) Stack(skip, depth int) Log {
	p.mu.Lock()
	defer p.mu.Unlock()
	callers := make([]uintptr, depth)
	if runtime.Callers(skip, callers) > 0 {
		frames := runtime.CallersFrames(callers)
		frameslice := []Fields{}
		for i := 0; i < depth; i++ {
			frame, more := frames.Next()
			f := Fields{
				KeyFile: frame.File,
				KeyLine: frame.Line,
				KeyFunc: frame.Func.Name(),
			}

			frameslice = append(frameslice, f)
			if !more {
				break
			}
		}
		p.fields[KeyFrames] = frameslice
	}
	return p
}

// Fields appends custom fields to Line.
func (p *Line) Fields(fields Fields) Log {
	p.mu.Lock()
	defer p.mu.Unlock()
	for key, val := range fields {
		p.fields[key] = val
	}
	return p
}
