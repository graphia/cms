package main

import (
	"io"
	"log"
	"os"
)

var (
	// Debug level logging (defaults to stdout)
	Debug *log.Logger

	// Info level logging (defaults to stdout)
	Info *log.Logger

	// Warning level logging (defaults to stdout)
	Warning *log.Logger

	// Error level logging (defaults to stderr)
	Error *log.Logger
)

func setupLogging(file bool) {

	if file {
		logFile, err := os.OpenFile(
			config.Logfile,
			os.O_RDWR|os.O_CREATE|os.O_APPEND,
			0666,
		)

		if err != nil {
			panic(err)
		}

		initializeLoggers(
			logFile,
			logFile,
			logFile,
			logFile,
		)

	} else {

		initializeLoggers(
			os.Stdout,
			os.Stdout,
			os.Stdout,
			os.Stderr,
		)

	}
}

func initializeLoggers(
	debugHandle,
	infoHandle,
	warningHandle,
	errorHandle io.Writer,
) {

	Debug = log.New(
		debugHandle,
		"[debug] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Info = log.New(
		infoHandle,
		"[info] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Warning = log.New(
		warningHandle,
		"[warning] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Error = log.New(
		errorHandle,
		"[error] ",
		log.Ltime|log.Lshortfile,
	)

}
