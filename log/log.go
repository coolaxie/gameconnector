package log

var (
	DebugReporter func(string, ...interface{})
	ReleaseReporter func(string, ...interface{})
	ErrorReporter func(string, ...interface{})
	FatalReporter func(string, ...interface{})
)

func Debug(format string, v ...interface{}) {
	if nil != DebugReporter {
		DebugReporter(format, v...)
	}
}

func Release(format string, v ...interface{}) {
	if nil != ReleaseReporter {
		ReleaseReporter(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if nil != ErrorReporter {
		ErrorReporter(format, v...)
	}
}

func Fatal(format string, v ...interface{}) {
	if nil != FatalReporter {
		FatalReporter(format, v...)
	}
}

