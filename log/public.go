package log

func L(name string) Logger {
	me.Lock()
	defer me.Unlock()

	log, ok := channels[name]
	if !ok {
		log = newChannel(name, defaultLevel)
		channels[name] = log
	}

	return log
}

func RemoveLogger(name string) {
	me.Lock()
	defer me.Unlock()
	delete(channels, name)
}

func Panic(format string, v ...interface{}) {
	logger.Panic(format, v...)
}

func Fatal(format string, v ...interface{}) {
	logger.Fatal(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

func Warning(format string, v ...interface{}) {
	logger.Warning(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func SetLevel(level int) {
	me.Lock()
	defer me.Unlock()

	defaultLevel = level

	for _, log := range channels {
		log.SetLevel(level)
	}
}
