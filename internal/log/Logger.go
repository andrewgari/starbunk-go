package log

import (
	"log"
	"os"
)

var (
	INFO  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
)

func SetupLogger() {
	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	INFO = log.New(file, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	WARN = log.New(file, "WARNING:  ", log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(file, "ERROR:  ", log.Ldate|log.Ltime|log.Lshortfile)
}
