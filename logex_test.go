package logex

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	l := NewStd()
	l.Println("TestLog")
}

func TestLog2(t *testing.T) {
	l := NewStd()
	l.Fields(Fields{
		"Mirko": "Odora",
		"Age":   "69",
	}).Println("TestLog2")
}

func TestLog3(t *testing.T) {
	l := NewStd()
	l.Stack(2, 10).Debugln("TestLog3")
}

func TestLog4(t *testing.T) {
	l := NewStd()
	l.Caller(0).Println("TestLog4")
}

func TestLog5(t *testing.T) {
	l := New()
	l.AddOutput(os.Stdout, NewJSONFormatter())
	l.Fields(Fields{
		"Mirko": "Odora",
		"Age":   64,
	}).Println("TestLog5")
}

func TestLog6(t *testing.T) {
	l := New()
	l.AddOutput(os.Stdout, NewJSONFormatter())
	l.Stack(2, 10).Debugln("TestLog6")
}
