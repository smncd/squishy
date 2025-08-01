package logging

import "log"

func SetToInfo(logger *log.Logger) {
	logger.SetPrefix("INFO: ")
	logger.SetFlags(log.Ldate | log.Ltime)
}

func SetToDebug(logger *log.Logger) {
	logger.SetPrefix("DEBUG: ")
	logger.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func SetToFatal(logger *log.Logger) {
	logger.SetPrefix("FATAL: ")
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
