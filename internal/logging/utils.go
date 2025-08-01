package logging

import "log"

// Sets logger in "info" mode for one print operation.
func Info(logger *log.Logger, format string, v ...any) {
	printfWith(logger, "INFO: ", log.Ldate|log.Ltime, format, v...)
}

// Sets logger in "debug" mode for one print operation.
func Debug(logger *log.Logger, format string, v ...any) {
	printfWith(logger, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile, format, v...)
}

// Sets logger in "fatal" mode for one print operation.
func Fatal(logger *log.Logger, format string, v ...any) {
	printfWith(logger, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile, format, v...)
}

func printfWith(logger *log.Logger, prefix string, flags int, format string, v ...any) {
	orgPrefix := logger.Prefix()
	orgFlags := logger.Flags()

	logger.SetPrefix(prefix)
	logger.SetFlags(flags)

	logger.Printf(format, v...)

	logger.SetPrefix(orgPrefix)
	logger.SetFlags(orgFlags)
}
