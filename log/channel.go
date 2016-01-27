package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type (
	ch struct {
		mu        sync.Mutex
		name      string
		level     int
		out       io.Writer
		formatter Formatter
	}
)

func newChannel(name string, level int) Logger {
	return &ch{
		name:      name,
		level:     level,
		out:       os.Stderr,
		formatter: newFormatter(),
	}
}

func (c *ch) Panic(format string, v ...interface{}) {
	c.log(LevelPanic, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func (c *ch) Fatal(format string, v ...interface{}) {
	c.log(LevelFatal, format, v...)
	os.Exit(1)
}

func (c *ch) Error(format string, v ...interface{}) {
	c.log(LevelError, format, v...)
}

func (c *ch) Warning(format string, v ...interface{}) {
	c.log(LevelWarning, format, v...)
}

func (c *ch) Info(format string, v ...interface{}) {
	c.log(LevelInfo, format, v...)
}

func (c *ch) Debug(format string, v ...interface{}) {
	c.log(LevelDebug, format, v...)
}

func (c *ch) SetLevel(level int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.level = level
}

func (c *ch) Level() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.level
}

func (c *ch) SetName(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.name = name
}

func (c *ch) Name() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.name
}

func (c *ch) SetFormatter(formatter Formatter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.formatter = formatter
}

func (c *ch) Formatter() Formatter {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.formatter
}

func (c *ch) log(level int, format string, v ...interface{}) {
	if c.level < level {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.formatter.Format(c.out, level, c.name, fmt.Sprintf(format, v...))
}
