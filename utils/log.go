package utils

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

func InitLog() {
	//log_level := os.GetEnv("TUBE_LOG_LEVEL")
	f := os.Stdout
	flag := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	WarningLogger = log.New(f, "WARNING: ", flag)
	ErrorLogger = log.New(f, "ERROR: ", flag)
	InfoLogger = log.New(f, "INFO: ", flag)
	DebugLogger = log.New(f, "DEBUG: ", flag)
}
