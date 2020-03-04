// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package logex implements a logger.
package logex

import (
	"github.com/vedranvuk/errorex"
)

// Log defines an interface for a logger.
type Log interface {

	// Debugf will log a debug message formed from format string and parameters.
	Debugf(string, ...interface{})
	// Debugln will log parameters as a debug message.
	Debugln(...interface{})
	// Infof will log an info message formed from format string and parameters.
	Infof(string, ...interface{})
	// Infoln will log parameters as an info message.
	Infoln(...interface{})
	// Warningf will log a warning message formed from format string and parameters.
	Warningf(string, ...interface{})
	// Warningln will log parameters as a warning message.
	Warningln(...interface{})
	// Errorf will log an error and an error message formed from format string and parameters.
	Errorf(error, string, ...interface{})
	// Errorln will log an error and parameters as a warning message.
	Errorln(error, ...interface{})

	// Printf will log a message with a custom logging level formed from format string and parameters.
	Printf(LogLevel, string, ...interface{})
	// Println will log parameters as a message with custom logging level.
	Println(LogLevel, ...interface{})

	// Caller will append the caller field to the next logged line.
	Caller(skip int) Log
	// Stack will append the stack field to the next logged line.
	Stack(skip int, depth int) Log
	// Fields will append the specified fields to the next logged line.
	Fields(*Fields) Log
}

var (
	// ErrLogex is the base error of logex package.
	ErrLogex = errorex.New("logex")
	// ErrUnmarshalLevel is returned when unmarshaling an invalid value from text as LogLevel.
	ErrUnmarshalLevel = ErrLogex.WrapFormat("error unmarshaling '%s' as loglevel")
	// ErrReservedKey is returned when a reserved key is being set to Fields.
	ErrReservedKey = ErrLogex.WrapFormat("cannot set field '%s', key is reserved")
	// ErrInvalidWalkFunc is returned when an invalid func was passed to Fields.Walk().
	ErrInvalidWalkFunc = ErrLogex.Wrap("invalid walk func")
)
