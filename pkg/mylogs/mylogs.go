package mylogs

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Lgr struct {
	*logrus.Logger
}

func New() *Lgr {
	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetReportCaller(true)

	if err := os.MkdirAll("logs", 0644); err != nil {
		panic(err)
	}

	abs, err := filepath.Abs("/logs")
	if err != nil {
		panic(err)
	}

	p := filepath.Join(abs, "all.log")
	fl, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	l.AddHook(&HookForFile{
		Writer:    []io.Writer{fl},
		LogLevels: logrus.AllLevels,
	})

	l.AddHook(&HookForAny{
		Writer:    []io.Writer{os.Stdout},
		LogLevels: []logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})

	return &Lgr{
		Logger: l,
	}
}

type HookForFile struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *HookForFile) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *HookForFile) Fire(entry *logrus.Entry) error {
	entry.Logger.Level = logrus.TraceLevel
	entry.Logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fName := path.Base(f.File)
			function = fmt.Sprintf("%s()", f.Function)
			file = fmt.Sprintf("%s:%d", fName, f.Line)
			return function, file
		},
		FullTimestamp: true,
	}

	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(str))
	}
	return err
}

type HookForAny struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *HookForAny) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *HookForAny) Fire(entry *logrus.Entry) error {
	entry.Logger.Formatter = &logrus.TextFormatter{
		ForceColors: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fName := path.Base(f.File)
			function = fmt.Sprintf("%s()\nMessage:", f.Function)
			file = fmt.Sprintf(" | %s:%d |", fName, f.Line)
			return function, file
		},
		FullTimestamp: true,
	}

	str, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range h.Writer {
		_, err = w.Write([]byte(str + "\n"))
	}
	return err
}
