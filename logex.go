// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package logex implements a logger.
package logex

import (
	"encoding/json"
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
	ErrUnmarshalLevel = ErrLog.Wrap("error unmarshaling '%s' as loglevel")
)

// Fields maps keys to values in a log line.
type Fields map[string]interface{}

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

// Formatter is an interface to a type that formats a map of key/value pairs to a log message.
type Formatter interface {
	// Format must return a string representation of key/value pairs.
	Format(Fields) string
}

// LogLevel defines a logging level for a Logger.
type LogLevel byte

const (
	// LevelDebug is the debug logging level.
	// It should contain verbose debugging information that is useful in debugging.
	LevelDebug LogLevel = iota
	// LevelInfo is the info logging level.
	// It should include information notices that might be useful to user.
	LevelInfo
	// LevelWarning is the warning logging level that might be important to user.
	LevelWarning
	// LevelError is the error logging level that contains important error information.
	LevelError
	// LevelPrint is the print logging level that is always printed, unless LevelMute.
	LevelPrint
	// LevelMute is the silent logging level.
	// It is used to silence the logger.
	LevelMute
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
		return ErrUnmarshalLevel.WithArgs(string(text))
	}
	return nil
}

// Logger logs.
type Logger struct {
	Log

	mu  sync.Mutex
	wrs map[io.Writer]Formatter
}

// print prints fields to registered writers using associated formatters.
func (l *Logger) print(fields Fields) {
	l.mu.Lock()
	defer l.mu.Unlock()
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

// New returns a new Logger with no defined outputs.
func New() *Logger {
	p := &Logger{
		mu:  sync.Mutex{},
		wrs: make(map[io.Writer]Formatter),
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

const (
	// KeyTime specifies that field carries a timestamp value.
	KeyTime = "time"
	// KeyMessage specifies that field carries the log message.
	KeyMessage = "message"
	// KeyLogLevel specifies that field carries the log level value.
	KeyLogLevel = "loglevel"
	// KeyError specifies that field carries an error value.
	KeyError = "error"
	// KeyFile specifies that field carries name of file, possibly a caller.
	KeyFile = "file"
	// KeyLine specifies that field carries number of line, probably of caller.
	KeyLine = "line"
)

// Line defines a log line of fields gained from logging calls which
// are passed to a Logger that converts them to log lines using a Formatter.
type Line struct {
	log    *Logger
	fields Fields
}

// NewLine returns a new Line instance that will output to Logger l.
func NewLine(l *Logger) *Line {
	return &Line{
		log:    l,
		fields: make(Fields),
	}
}

// output outputs line fields to the Logger.
func (p *Line) output(level LogLevel, message string) {
	p.fields[KeyTime] = time.Now()
	p.fields[KeyMessage] = message
	p.fields[KeyLogLevel] = level
	p.log.print(p.fields)
	p.fields = make(Fields)
}

func (p *Line) Debugf(format string, args ...interface{}) {
	p.output(LevelDebug, fmt.Sprintf(format, args...))
}

func (p *Line) Debugln(args ...interface{}) {
	p.output(LevelDebug, fmt.Sprint(args...)+"\n")
}

func (p *Line) Infof(format string, args ...interface{}) {
	p.output(LevelInfo, fmt.Sprintf(format, args...))
}

func (p *Line) Infoln(args ...interface{}) {
	p.output(LevelInfo, fmt.Sprint(args...)+"\n")
}

func (p *Line) Warningf(format string, args ...interface{}) {
	p.output(LevelWarning, fmt.Sprintf(format, args...))
}

func (p *Line) Warningln(args ...interface{}) {
	p.output(LevelWarning, fmt.Sprint(args...)+"\n")
}

func (p *Line) Errorf(err error, format string, args ...interface{}) {
	p.fields[KeyError] = err
	p.output(LevelError, fmt.Sprintf(format, args...))
}

func (p *Line) Errorln(err error, args ...interface{}) {
	p.fields[KeyError] = err
	p.output(LevelError, fmt.Sprint(args...)+"\n")
}

func (p *Line) Printf(format string, args ...interface{}) {
	p.output(LevelPrint, fmt.Sprintf(format, args...))
}

func (p *Line) Println(args ...interface{}) {
	p.output(LevelPrint, fmt.Sprint(args...)+"\n")
}

// Caller appends the caller fields to the Line.
func (p *Line) Caller(skip int) Log {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		p.fields[KeyFile] = file
		p.fields[KeyLine] = line
	}
	return p
}

// Stack appends the stack to the Line.
func (p *Line) Stack(skip, depth int) Log {
	callers := make([]uintptr, depth)
	fn := runtime.Callers(skip, callers)
	fmt.Println(fn)
	frames := runtime.CallersFrames(callers)
	for i := 0; i < depth; i++ {
		frame, more := frames.Next()
		p.fields[fmt.Sprintf("frame.%0d.file", i)] = frame.File
		p.fields[fmt.Sprintf("frame.%0d.line", i)] = frame.Line
		p.fields[fmt.Sprintf("frame.%0d.func", i)] = frame.Function
		if !more {
			break
		}
	}
	return p
}

// Fields appends specified fields to Line.
func (p *Line) Fields(fields Fields) Log {
	for key, val := range fields {
		p.fields[key] = val
	}
	return p
}

// SimpleFormatter appends key/value pairs alphabetically to returned string.
type SimpleFormatter struct{}

// Format implements Formatter interface.
func (sf SimpleFormatter) Format(fields Fields) string {
	s := ""
	for key, val := range fields {
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("\"%s\"=\"%v\"", key, val)
	}
	return s
}

// NewSimpleFormatter returns a new SimpleFormatter.
func NewSimpleFormatter() Formatter {
	return &SimpleFormatter{}
}

// JSONFormatter formats key/value pairs into a JSON string.
type JSONFormatter struct{}

// Format implements Formatter interface.
func (jf *JSONFormatter) Format(fields Fields) string {
	buf, err := json.Marshal(fields)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

// NewJSONFormatter returns a new JSONFormatter.
func NewJSONFormatter() Formatter {
	return &JSONFormatter{}
}

// CSVFormatter formats key/value pairs into CSV data.
type CSVFormatter struct{}

// Format implements Formatter interface.
func (cf *CSVFormatter) Format(fields Fields) string {
	return ""
}

// LogFormatter formats key/value pairs into generalized log lines.
type LogFormatter struct{}

// Format implements Formatter interface.
func (tf *LogFormatter) Format(fields Fields) string {

	const DefTimestamp = "2006-02-01 15:04:05"

	return ""
}
