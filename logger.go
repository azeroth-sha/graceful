package graceful

import (
	"fmt"
	"log"
	"os"
)

var (
	flag       = log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix
	defaultLog = &logger{
		callDepth: 2,
		infoLog:   log.New(os.Stdout, `[graceful] `, flag),
		errorLog:  log.New(os.Stderr, `[graceful] `, flag),
	}
)

type logger struct {
	callDepth int
	infoLog   *log.Logger
	errorLog  *log.Logger
}

// Info calls l.infoLog to print to the logger.
// Arguments are handled in the manner of fmt.Sprintln.
func (l *logger) Info(a ...interface{}) {
	_ = l.infoLog.Output(l.callDepth, fmt.Sprintln(a...))
}

// Infof calls l.infoLog to print to the logger.
// Arguments are handled in the manner of fmt.Sprintf.
func (l *logger) Infof(format string, a ...interface{}) {
	_ = l.infoLog.Output(l.callDepth, fmt.Sprintf(format, a...))
}

// Error calls l.errorLog to print to the logger.
// Arguments are handled in the manner of fmt.Sprintln.
func (l *logger) Error(a ...interface{}) {
	_ = l.errorLog.Output(l.callDepth, fmt.Sprintln(a...))
}

// Errorf calls l.errorLog to print to the logger.
// Arguments are handled in the manner of fmt.Sprintf.
func (l *logger) Errorf(format string, a ...interface{}) {
	_ = l.errorLog.Output(l.callDepth, fmt.Sprintf(format, a...))
}

func SetDefaultLog(infoLog, errLog *log.Logger, callDepth ...int) {
	if len(callDepth) == 0 {
		callDepth = append(callDepth, 2)
	}
	defaultLog = &logger{
		callDepth: callDepth[0],
		infoLog:   infoLog,
		errorLog:  errLog,
	}
}
