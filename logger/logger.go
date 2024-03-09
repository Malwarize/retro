package logger

import (
	"bufio"
	"log"
	"os"

	"github.com/Malwarize/retro/config"
)

var (
	INFOLogger  *log.Logger
	WARNLogger  *log.Logger
	ERRORLogger *log.Logger
)

func init() {
	logFile, err := os.OpenFile(
		config.GetConfig().LogFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o666,
	)
	if err != nil && !os.IsNotExist(err) {
		logFile = os.Stdout
	}
	INFOLogger = log.New(
		logFile,
		"INFO: ",
		log.Ldate|log.Ltime,
	)
	WARNLogger = log.New(
		logFile,
		"WARN: ",
		log.Ldate|log.Ltime,
	)
	ERRORLogger = log.New(
		logFile,
		"ERROR: ",
		log.Ldate|log.Ltime,
	)
}

func LogError(err RetroError, extra ...any) error {
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
	WARNLogger.Println(warn, extra)
}

func GetLogs() ([]string, error) {
	logFile, err := os.Open(
		config.GetConfig().LogFile,
	)

	if err != nil {
		return nil, err
	}

	var lastLines []string
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		lastLines = append(lastLines, scanner.Text())
		if len(lastLines) > 200 {
			lastLines = lastLines[1:]
		}
	}
	return lastLines, nil
}
