// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

var logger = New()

// Debugf logs a debug message formed from format string and args using the default logger.
func Debugf(format string, args ...interface{}) { logger.Debugf(format, args...) }

// Debugln logs args as a debug message using the default logger.
func Debugln(args ...interface{}) { logger.Debugln(args...) }

// Infof logs an info message formed from format string and args using the default logger.
func Infof(format string, args ...interface{}) { logger.Infof(format, args...) }

// Infoln logs args as an info message using the default logger.
func Infoln(args ...interface{}) { logger.Infoln(args...) }

// Warningf logs a warning message formed from format string and args using the default logger.
func Warningf(format string, args ...interface{}) { logger.Warningf(format, args...) }

// Warningln logs args as a warning message using the default logger.
func Warningln(args ...interface{}) { logger.Warningln(args...) }

// Errorf logs an error and an error message formed from format string and args using the default logger.
func Errorf(err error, format string, args ...interface{}) { logger.Errorf(err, format, args...) }

// Errorln logs an error and args as a warning message using the default logger.
func Errorln(err error, args ...interface{}) { logger.Errorln(err, args...) }

// Printf logs a message with a custom logging level formed from format string and args using the default logger.
func Printf(level LogLevel, format string, args ...interface{}) {
	logger.Printf(level, format, args...)
}

// Println logs args as a message with custom logging level using the default logger.
func Println(level LogLevel, args ...interface{}) { logger.Println(level, args...) }

// Caller appends the caller field to the next logged line using the default logger.
func WithCaller(skip int) Log { return logger.WithCaller(skip) }

// Stack appends the stack field to the next logged line using the default logger.
func WithStack(skip int, depth int) Log { return logger.WithStack(skip, depth) }

// Fields appends the specified fields to the next logged line using the default logger.
func WithFields(f *Fields) Log { return logger.WithFields(f) }
