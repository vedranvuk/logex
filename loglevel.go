// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logex

import (
	"fmt"
	"strconv"
	"strings"
)

// LogLevel is Log logging level.
type LogLevel byte

const (
	// LevelNone is undefined logging level. It prints nothing.
	LevelNone LogLevel = iota
	// LevelMute is the silent logging level used to silence the logger.
	LevelMute
	// LevelError is the error logging level that prints errors only.
	LevelError
	// LevelWarning is the warning logging level that prints warnings and errors.
	LevelWarning
	// LevelInfo is the info logging level that prints information, warnings and errors.
	LevelInfo
	// LevelDebug is the debug logging level that prints debug messages, information, warnings and errors.
	LevelDebug
	// LevelCustom and levels up to LevelPrint are custom logging levels.
	// To define a custom logging level use: MyLevel := LogLevel(LevelCustom +1).
	LevelCustom
	// LevelPrint is the print logging level that prints everything that gets logged.
	LevelPrint = LogLevel(255)
)

// String implements the Stringer interface.
func (ll LogLevel) String() string {
	switch ll {
	case LevelNone:
		return "None"
	case LevelMute:
		return "Mute"
	case LevelError:
		return "Error"
	case LevelWarning:
		return "Warning"
	case LevelInfo:
		return "Info"
	case LevelDebug:
		return "Debug"
	case LevelPrint:
		return "Print"
	default:
		if ll >= LevelCustom && ll < LevelPrint {
			return fmt.Sprintf("Custom(%d)", byte(ll))
		}
	}
	return ""
}

// MarshalText implements the TextMarshaler interface.
func (ll *LogLevel) MarshalText() ([]byte, error) {
	return []byte(ll.String()), nil
}

// UnmarshalText implements the TextUnmarshaler interface.
func (ll *LogLevel) UnmarshalText(text []byte) error {
	switch s := strings.ToLower(string(text)); s {
	case "none":
		*ll = LevelNone
	case "mute":
		*ll = LevelMute
	case "error":
		*ll = LevelError
	case "warning":
		*ll = LevelWarning
	case "info":
		*ll = LevelInfo
	case "debug":
		*ll = LevelDebug
	case "print":
		*ll = LevelPrint
	default:
		if strings.HasPrefix(s, "custom") {
			s = strings.TrimPrefix(s, "custom")
			if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
				s = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "("), ")"))
				lvl, err := strconv.Atoi(s)
				if err != nil {
					return ErrUnmarshalLevel.WrapArgs(text)
				}
				*ll = LogLevel(lvl)
				return nil
			}
		}
		return ErrUnmarshalLevel.WrapArgs(string(text))
	}
	return nil
}
