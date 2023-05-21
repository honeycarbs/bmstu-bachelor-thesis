package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("build/logs", 0755)

	if err != nil || os.IsExist(err) {
		panic("can't create log dir. no configured logging to files")
	} else {
		allFile, err := os.OpenFile("build/logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			panic(fmt.Sprintf("[Message]: %s", err))
		}

		l.SetOutput(ioutil.Discard)

		l.AddHook(&writerHook{
			Writer:    []io.Writer{allFile, os.Stdout},
			LogLevels: logrus.AllLevels,
		})
	}

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
