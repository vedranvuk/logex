package logex

import (
	"errors"
	"os"
	"testing"
)

func DoLog(l Log) {
	l.Printf("%s\n", "Printf")
	l.Println("Println")
	l.Debugf("%s\n", "Debug")
	l.Debugln("Debugln")
	l.Infof("%s\n", "Info")
	l.Infoln("Infoln")
	l.Warningf("%s\n", "Warningf")
	l.Warningln("Warningln")
	l.Errorf(errors.New("error: catastrophic failure"), "%s\n", "Errorf")
	l.Errorln(errors.New("error: success"), "Errorln")
	l.Caller(1).Errorln(errors.New("error: caller"), "ERROR")
	l.Stack(0, 10).Errorln(errors.New("error: stack"), "ERROR")
}

func TestLog(t *testing.T) {
	l := NewStd()
	DoLog(l)
}

func TestLog2(t *testing.T) {
	l := New()
	l.AddOutput(os.Stdout, NewJSONFormatter(true))
	DoLog(l)
}

func TestLog3(t *testing.T) {
	l := NewStd()
	l.SetLevel(LevelError)
	l.Println("testis")
}
