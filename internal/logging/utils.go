package logging

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// Sets logger in "info" mode for one print operation.
func Info(logger *log.Logger, format string, v ...any) {
	printfWith(logger, "INFO: ", log.Ldate|log.Ltime, format, v...)
}

// Sets logger in "debug" mode for one print operation.
func Debug(logger *log.Logger, format string, v ...any) {
	printfWith(logger, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile, format, v...)
}

// Sets logger in "error" mode for one print operation.
func Error(logger *log.Logger, format string, v ...any) {
	printfWith(logger, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile, format, v...)
}

func printfWith(logger *log.Logger, prefix string, flags int, format string, v ...any) {
	if flags&log.Lshortfile != 0 || flags&log.Llongfile != 0 {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		if flags&log.Lshortfile != 0 {
			flags &^= log.Lshortfile
			file = filepath.Base(file)

		} else if flags&log.Llongfile != 0 {
			flags &^= log.Llongfile
			f, err := filepath.Abs(file)
			if err == nil {
				file = f
			}
		}

		format = fmt.Sprintf("%s:%d: %s", file, line, format)

	}

	orgPrefix := logger.Prefix()
	orgFlags := logger.Flags()

	logger.SetPrefix(prefix)
	logger.SetFlags(flags)

	logger.Printf(format, v...)

	logger.SetPrefix(orgPrefix)
	logger.SetFlags(orgFlags)
}
