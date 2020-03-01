// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"time"
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

var reservedkeys = map[string]struct{}{
	"time":     struct{}{},
	"message":  struct{}{},
	"loglevel": struct{}{},
	"error":    struct{}{},
	"frames":   struct{}{},
	"file":     struct{}{},
	"line":     struct{}{},
	"func":     struct{}{},
}

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

// Set sets a custom field under key to value.
// Set returns an error if a reserved key was is used.
func (f Fields) Set(key string, value interface{}) error {
	if _, ok := reservedkeys[key]; ok {
		return ErrReservedKey.WrapArgs(key)
	}
	f[key] = value
	return nil
}
