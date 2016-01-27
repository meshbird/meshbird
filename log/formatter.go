package log

import (
	"io"
	"time"
)

type (
	stdFormatter struct {
		buf []byte
	}
)

var (
	chLeft = []byte{' ', '['}
	chRight = []byte{']', ' '}
)

func newFormatter() Formatter {
	return &stdFormatter{}
}

func (f *stdFormatter) Format(out io.Writer, level int, channel string, msg string) {
	now := time.Now()

	f.buf = f.buf[:0]
	f.formatHeader(&f.buf, now, level, channel)
	f.buf = append(f.buf, msg...)

	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		f.buf = append(f.buf, '\n')
	}

	out.Write(f.buf)
}

func (f *stdFormatter) formatHeader(buf *[]byte, t time.Time, level int, channel string) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')

	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	*buf = append(*buf, ' ')

	*buf = append(*buf, levelOut[level]...)

	*buf = append(*buf, chLeft...)
	*buf = append(*buf, channel...)
	*buf = append(*buf, chRight...)
}
