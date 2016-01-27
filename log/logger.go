package log

import (
	"io"
	"sync"
)

type (
	Logger interface {
		Panic(format string, v ...interface{})
		Fatal(format string, v ...interface{})
		Error(format string, v ...interface{})
		Warning(format string, v ...interface{})
		Info(format string, v ...interface{})
		Debug(format string, v ...interface{})

		SetLevel(level int)
		Level() int

		SetName(name string)
		Name() string

		SetFormatter(formatter Formatter)
		Formatter() Formatter
	}

	Formatter interface {
		Format(out io.Writer, level int, channel string, msg string)
	}
)

var (
	defaultLevel = LevelWarning
	defaultName  = "main"
	logger       = newChannel(defaultName, defaultLevel)
	channels     = map[string]Logger{
		defaultName: logger,
	}

	me sync.Mutex
)
