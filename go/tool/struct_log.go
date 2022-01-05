package tool

import (
	"io/ioutil"
	"log"
	"os"
)

// log
var flag = log.Lshortfile | log.LstdFlags
var loggerError = log.New(os.Stderr, "[ERROR] ", flag)
var loggerInfo = log.New(os.Stdout, "[INFO] ", flag)
var loggerWarn = log.New(os.Stdout, "[WARN] ", flag)
var loggerDiscard = log.New(ioutil.Discard, "", 0)

func LoggerError() *log.Logger {
	return loggerError
}
func LoggerInfo() *log.Logger {
	return loggerInfo
}
func LoggerWarn() *log.Logger {
	return loggerWarn
}
func LoggerDiscard() *log.Logger {
	return loggerDiscard
}

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
}
type Logge struct {
}

func (l *Logge) Info(msg string) {
	LoggerInfo().Println(msg)
}

func (l *Logge) Warn(msg string) {
	LoggerWarn().Println(msg)
}

func (l *Logge) Error(msg string) {
	LoggerError().Println(msg)
}

func NewLogger() Logger {
	return &Logge{}
}

type Log struct {
	Typ     int
	Message string
}

type SequenceLogger struct {
	logger Logger
	seq    *chan *Log
}

func NewSequenceLogger(logger Logger) *SequenceLogger {
	seq := make(chan *Log, 100)
	sl := &SequenceLogger{
		logger: logger,
		seq:    &seq,
	}
	sl.out()
	return sl
}

func (s *SequenceLogger) out() {
	go func() {
		for {
			select {
			case l := <-*s.seq:
				switch l.Typ {
				case 0:
					s.logger.Info(l.Message)
				case 1:
					s.logger.Warn(l.Message)
				case 2:
					s.logger.Error(l.Message)
				}
			}
		}
	}()
}

func (s *SequenceLogger) Info(msg string) {
	*s.seq <- &Log{
		Typ:     0,
		Message: msg,
	}
}

func (s *SequenceLogger) Warn(msg string) {
	*s.seq <- &Log{
		Typ:     1,
		Message: msg,
	}
}

func (s *SequenceLogger) Error(msg string) {
	*s.seq <- &Log{
		Typ:     2,
		Message: msg,
	}
}
