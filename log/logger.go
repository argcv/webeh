package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"github.com/argcv/webeh/repl"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
	DISABLED
)

func (l LogLevel) String() (s string) {
	s = ""
	switch l {
	case DEBUG:
		s = repl.NewColoredText("DEBUG").SetFG(repl.Cyan).String()
	case INFO:
		s = repl.NewColoredText("INFO").SetFG(repl.Green).String()
	case WARN:
		s = repl.NewColoredText("WARN").SetFG(repl.Magenta).String()
	case ERROR:
		s = repl.NewColoredText("ERROR").SetFG(repl.BrightRed).String()
	case FATAL:
		s = repl.NewColoredText("FATAL").SetBG(repl.Red).SetFG(repl.BrightWhite).String()
		//default:
		//s = repl.NewColoredText(fmt.Sprintf("??:%v", l)).SetFG(repl.Cyan).String()
	}
	s = fmt.Sprintf("[%v] ", s)
	return
}

var (
	loggers = map[LogLevel]*log.Logger{
		DEBUG: nil,
		INFO:  nil,
		WARN:  nil,
		ERROR: nil,
		FATAL: nil,
	}
	loggersMtx          = sync.Mutex{}
	cLevel     LogLevel = INFO
)

func SetLevel(level LogLevel) {
	cLevel = level
}

func Verbose() {
	SetLevel(DEBUG)
}

func Quiet() {
	SetLevel(DISABLED)
}

func Output(level LogLevel, msg string, calldepth int) {
	if cLevel <= level {
		loggers[level].Output(calldepth+1, msg)
	}
}

func Debug(v ...interface{}) {
	Output(DEBUG, fmt.Sprintln(v...), 2)
}

func Info(v ...interface{}) {
	Output(INFO, fmt.Sprintln(v...), 2)
}

func Warn(v ...interface{}) {
	Output(WARN, fmt.Sprintln(v...), 2)
}

func Error(v ...interface{}) {
	Output(ERROR, fmt.Sprintln(v...), 2)
}

func Fatal(v ...interface{}) {
	Output(FATAL, fmt.Sprintln(v...), 2)
}

// calldepth == 0 == current
func Debugd(calldepth int, v ...interface{}) {
	Output(DEBUG, fmt.Sprintln(v...), calldepth + 2)
}

func Infod(calldepth int,v ...interface{}) {
	Output(INFO, fmt.Sprintln(v...), calldepth + 2)
}

func Warnd(calldepth int,v ...interface{}) {
	Output(WARN, fmt.Sprintln(v...), calldepth + 2)
}

func Errord(calldepth int,v ...interface{}) {
	Output(ERROR, fmt.Sprintln(v...), calldepth + 2)
}

func Fatald(calldepth int,v ...interface{}) {
	Output(FATAL, fmt.Sprintln(v...), calldepth + 2)
}

func Debugf(f string, v ...interface{}) {
	Output(DEBUG, fmt.Sprintf(f, v...), 2)
}

func Infof(f string, v ...interface{}) {
	Output(INFO, fmt.Sprintf(f, v...), 2)
}

func Warnf(f string, v ...interface{}) {
	Output(WARN, fmt.Sprintf(f, v...), 2)
}

func Errorf(f string, v ...interface{}) {
	Output(ERROR, fmt.Sprintf(f, v...), 2)
}

func Fatalf(f string, v ...interface{}) {
	Output(FATAL, fmt.Sprintf(f, v...), 2)
}

func IfEligible(level LogLevel, f func()) bool {
	if cLevel <= level {
		f()
		return true
	} else {
		return false
	}
}

func IfDebug(f func()) bool {
	return IfEligible(DEBUG, f)
}

func SetLogger(out io.Writer, flag int) {
	loggersMtx.Lock()
	defer loggersMtx.Unlock()
	for k, _ := range loggers {
		loggers[k] = log.New(out, k.String(), flag)
	}
}

func init() {
	SetLogger(os.Stderr, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}