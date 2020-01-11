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
	ErrUnmarshalLevel = ErrLog.WrapFormat("error unmarshaling '%s' as loglevel")
)

const (
	// KeyTime specifies that field carries a timestamp value.
	KeyTime = "time"
	// KeyMessage specifies that field carries the log message.
	KeyMessage = "message"
	// KeyLogLevel specifies that field carries the log level value.
	KeyLogLevel = "loglevel"
	// KeyError specifies that field carries an error value.
	KeyError = "error"
	// KeyFrames
	KeyFrames = "frames"
	// KeyFile specifies that field carries name of file, possibly a caller.
	KeyFile = "file"
	// KeyLine specifies that field carries number of line, probably of caller.
	KeyLine = "line"
	// KeyFnc
	KeyFunc = "func"
)

// Fields maps keys to values in a log line.
type Fields map[string]interface{}

// Time returns Time field.
func (f Fields) Time() time.Time {
	if time, ok := (f[KeyTime]).(time.Time); ok {
		return time
	}
	return time.Time{}
}

// Message returns message field.
func (f Fields) Message() string {
	if message, ok := (f[KeyMessage]).(string); ok {
		return message
	}
	return ""
}

// Message returns log level field.
func (f Fields) LogLevel() LogLevel {
	if loglevel, ok := (f[KeyLogLevel]).(LogLevel); ok {
		return loglevel
	}
	return LevelNone
}

// Message returns error field.
func (f Fields) Error() error {
	if err, ok := (f[KeyError]).(error); ok {
		return err
	}
	return nil
}

// Frames returns frames.
func (f Fields) Frames() []Fields {
	if fields, ok := (f[KeyFrames]).([]Fields); ok {
		return fields
	}
	return nil
}

// File returns file field.
func (f Fields) File() string {
	if file, ok := (f[KeyFile]).(string); ok {
		return file
	}
	return ""
}

// Line returns line field.
func (f Fields) Line() int {
	if line, ok := (f[KeyLine]).(int); ok {
		return line
	}
	return 0
}

// Func returns func field.
func (f Fields) Func() string {
	if fun, ok := (f[KeyFunc]).(string); ok {
		return fun
	}
	return ""
}

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
	// Format must return a string representation of key/value pairs, such as JSON object, CSV, custom.
	Format(Fields) string
}

// LogLevel defines a logging level for a Logger.
type LogLevel byte

const (
	// LevelNone is undefined.
	LevelNone LogLevel = iota
	// LevelMute is the silent logging level.
	// It is used to silence the logger.
	LevelMute
	// LevelPrint is the print logging level that is always printed, unless LevelMute.
	LevelPrint
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
}

// NewLine returns a new Line instance that will output to Logger l.
func NewLine(l *Logger) *Line {
	return &Line{
		log:    l,
		fields: make(Fields),
	}
}

// flush outputs line fields to the Logger.
func (p *Line) flush(level LogLevel, message string) {
	p.fields[KeyTime] = time.Now()
	p.fields[KeyMessage] = message
	p.fields[KeyLogLevel] = level
	p.log.print(p.fields)
	p.fields = make(Fields)
}

func (p *Line) Debugf(format string, args ...interface{}) {
	p.flush(LevelDebug, fmt.Sprintf(format, args...))
}

func (p *Line) Debugln(args ...interface{}) {
	p.flush(LevelDebug, fmt.Sprint(args...)+"\n")
}

func (p *Line) Infof(format string, args ...interface{}) {
	p.flush(LevelInfo, fmt.Sprintf(format, args...))
}

func (p *Line) Infoln(args ...interface{}) {
	p.flush(LevelInfo, fmt.Sprint(args...)+"\n")
}

func (p *Line) Warningf(format string, args ...interface{}) {
	p.flush(LevelWarning, fmt.Sprintf(format, args...))
}

func (p *Line) Warningln(args ...interface{}) {
	p.flush(LevelWarning, fmt.Sprint(args...)+"\n")
}

func (p *Line) Errorf(err error, format string, args ...interface{}) {
	p.fields[KeyError] = err
	p.flush(LevelError, fmt.Sprintf(format, args...))
}

func (p *Line) Errorln(err error, args ...interface{}) {
	p.fields[KeyError] = err
	p.flush(LevelError, fmt.Sprint(args...)+"\n")
}

func (p *Line) Printf(format string, args ...interface{}) {
	p.flush(LevelPrint, fmt.Sprintf(format, args...))
}

func (p *Line) Println(args ...interface{}) {
	p.flush(LevelPrint, fmt.Sprint(args...)+"\n")
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
	for key, val := range fields {
		p.fields[key] = val
	}
	return p
}

// SimpleFormatter appends key/value pairs alphabetically to returned string.
type SimpleFormatter struct{}

// Format implements Formatter interface.
func (sf SimpleFormatter) Format(fields Fields) string {

	const TimeStampFormat = "2006-02-01 15:04:05"

	s := fields.Time().Format(TimeStampFormat)
	switch fields.LogLevel() {
	case LevelError:
		s += " EROR"
	case LevelWarning:
		s += " WARN"
	case LevelInfo:
		s += " INFO"
	case LevelDebug:
		s += " DEBG"
	}
	s += " "
	s += fields.Message()
	if err := fields.Error(); err != nil {
		s += fmt.Sprintf("\t%v\n", err)
	}
	if file := fields.File(); file != "" {
		s += fmt.Sprintf("\t%s (%d)\n", fields.File(), fields.Line())
	}
	if frames := fields.Frames(); frames != nil {
		for _, frame := range frames {
			s += fmt.Sprintf("\t%s (%d)\n\t\t%s\n", frame[KeyFile], frame[KeyLine], frame[KeyFunc])
		}
	}
	return s
}

// NewSimpleFormatter returns a new SimpleFormatter.
func NewSimpleFormatter() Formatter {
	return &SimpleFormatter{}
}

// JSONFormatter formats key/value pairs into a string of JSON object.
type JSONFormatter struct{ indent bool }

// Format implements Formatter interface.
func (jf *JSONFormatter) Format(fields Fields) string {
	var buf []byte
	var err error
	if jf.indent {
		buf, err = json.MarshalIndent(fields, "", "\t")
	} else {
		buf, err = json.Marshal(fields)
	}
	if err != nil {
		return err.Error()
	}
	return string(buf) + "\n"
}

// NewJSONFormatter returns a new JSONFormatter.
func NewJSONFormatter(indent bool) Formatter {
	return &JSONFormatter{indent}
}
