package logger

import "log"

var (
	INFOLogger  *log.Logger
	WARNLogger  *log.Logger
	ERRORLogger *log.Logger
)

func init() {
	INFOLogger = log.New(
		log.Writer(),
		"INFO: ",
		log.Ldate,
	)
	WARNLogger = log.New(
		log.Writer(),
		"WARN: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	ERRORLogger = log.New(
		log.Writer(),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}

func LogError(err GoPlayError, extra ...any) error {
	ERRORLogger.Println(
		err,
		extra,
	)
	return err
}

func LogInfo(info string, extra ...any) {
	INFOLogger.Println(info, extra)
}

func LogWarn(warn string, extra ...any) {
	WARNLogger.Println(warn)
}
