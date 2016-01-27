package log

import "fmt"

const (
	LevelPanic int = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var (
	levelNames = map[string]int{
		"panic":   LevelPanic,
		"fatal":   LevelFatal,
		"error":   LevelError,
		"warning": LevelWarning,
		"info":    LevelInfo,
		"debug":   LevelDebug,
	}
	levelOut = map[int]string{
		LevelPanic:   "PNC",
		LevelFatal:   "FTL",
		LevelError:   "ERR",
		LevelWarning: "WRN",
		LevelInfo:    "INF",
		LevelDebug:   "DBG",
	}
)

func ParseLevel(s string) (int, error) {
	level, ok := levelNames[s]
	if !ok {
		return 0, fmt.Errorf("unknown level %q", s)
	}
	return level, nil
}

func MustParseLevel(s string) int {
	level, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return level
}
