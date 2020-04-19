// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"encoding/json"
	"sync"
	"time"
)

// FieldKey is a Fields key.
type FieldKey string

const (
	// KeyTime specifies that field carries a timestamp value.
	KeyTime FieldKey = "time"
	// KeyMessage specifies that field carries the log message.
	KeyMessage FieldKey = "message"
	// KeyLogLevel specifies that field carries the log level value.
	KeyLogLevel FieldKey = "loglevel"
	// KeyError specifies that field carries an error value.
	KeyError FieldKey = "error"
	// KeyFrames
	KeyFrames FieldKey = "frames"
	// KeyFile specifies that field carries name of file, possibly a caller.
	KeyFile FieldKey = "file"
	// KeyLine specifies that field carries number of line, probably of caller.
	KeyLine FieldKey = "line"
	// KeyFnc
	KeyFunc FieldKey = "func"
)

// list of reserved keys.
var reservedkeys = map[FieldKey]struct{}{
	KeyTime:     struct{}{},
	KeyMessage:  struct{}{},
	KeyLogLevel: struct{}{},
	KeyError:    struct{}{},
	KeyFrames:   struct{}{},
	KeyFile:     struct{}{},
	KeyLine:     struct{}{},
	KeyFunc:     struct{}{},
}

// keyreserved returns if a key is reserved.
func keyreserved(key FieldKey) (reserved bool) {
	_, reserved = reservedkeys[key]
	return
}

type fieldsMap map[FieldKey]interface{}

// Fields maps keys to values in a log line.
type Fields struct {
	mu sync.Mutex
	fieldsMap
}

// NewFields creates new Fields.
func NewFields() *Fields {
	return &Fields{
		mu:        sync.Mutex{},
		fieldsMap: make(fieldsMap),
	}
}

func (f *Fields) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &f.fieldsMap) }
func (f *Fields) MarshalJSON() ([]byte, error)    { return json.Marshal(f.fieldsMap) }

// set sets a field under key to value.
func (f *Fields) set(key FieldKey, value interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fieldsMap[key] = value
}

// Set sets a custom field under key to value.
// Set returns an error if a reserved key was is used.
func (f *Fields) Set(key FieldKey, value interface{}) error {
	if keyreserved(key) {
		return ErrReservedKey.WrapArgs(key)
	}
	f.set(key, value)
	return nil
}

// Get gets a field by key and returns it and a truth if it exists.
func (f *Fields) Get(key FieldKey) (val interface{}, exists bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	val, exists = f.fieldsMap[key]
	return
}

// len returns number of fields.
func (f *Fields) Len() int {
	return len(f.fieldsMap)
}

// Custom returns custom fields
func (f *Fields) Custom() *Fields {
	cf := NewFields()
	for key, val := range f.fieldsMap {
		if !keyreserved(key) {
			cf.set(key, val)
		}
	}
	return cf
}

// WalkFunc is a prototype of a func Walk calls.
type WalkFunc = func(key FieldKey, val interface{}) bool

// Walk walks the fields and calls f for each field.
// f should return true to continue the walk.
// Walk returns an error if f is invalid.
func (f *Fields) Walk(wf WalkFunc) error {
	if wf == nil {
		return ErrInvalidWalkFunc
	}
	for key, val := range f.fieldsMap {
		if !wf(key, val) {
			break
		}
	}
	return nil
}

// Time returns Time field.
func (f *Fields) Time() time.Time {
	t, ok := f.Get(KeyTime)
	if !ok {
		return time.Time{}
	}
	return t.(time.Time)
}

// Message returns message field.
func (f *Fields) Message() string {
	msg, ok := f.Get(KeyMessage)
	if !ok {
		return ""
	}
	return msg.(string)
}

// Message returns log level field.
func (f *Fields) LogLevel() LogLevel {
	lvl, ok := f.Get(KeyLogLevel)
	if !ok {
		return LevelNone
	}
	return lvl.(LogLevel)
}

// Message returns error field.
func (f *Fields) Error() error {
	err, ok := f.Get(KeyError)
	if !ok {
		return nil
	}
	if err == nil {
		return nil
	}
	return err.(error)
}

// Frames returns frames.
func (f *Fields) Frames() []*Fields {
	fields, ok := f.Get(KeyFrames)
	if !ok {
		return nil
	}
	return fields.([]*Fields)
}

// File returns file field.
func (f *Fields) File() string {
	file, ok := f.Get(KeyFile)
	if !ok {
		return ""
	}
	return file.(string)
}

// Line returns line field.
func (f *Fields) Line() int {
	line, ok := f.Get(KeyLine)
	if !ok {
		return -1
	}
	return line.(int)
}

// Func returns func field.
func (f *Fields) Func() string {
	fun, ok := f.Get(KeyFunc)
	if !ok {
		return ""
	}
	return fun.(string)
}
