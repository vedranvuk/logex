// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"encoding/json"
	"fmt"
)

// Formatter formats Fields to a custom format.
type Formatter interface {
	// Format must return a string representation of key/value pairs, such as JSON object, CSV, custom.
	Format(Fields) string
}

// SimpleFormatter sorts Fields alphabetically and appends them as "key"="value" pairs separated by space.
type SimpleFormatter struct{}

// NewSimpleFormatter returns a new SimpleFormatter.
func NewSimpleFormatter() Formatter { return &SimpleFormatter{} }

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

// JSONFormatter formats Fields into a JSON object.
type JSONFormatter struct{ indent bool }

// NewJSONFormatter returns a new JSONFormatter.
func NewJSONFormatter(indent bool) Formatter { return &JSONFormatter{indent} }

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
