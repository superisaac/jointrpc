package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	logsyslog "github.com/sirupsen/logrus/hooks/syslog"
	client "github.com/superisaac/rpctube/client"
	example "github.com/superisaac/rpctube/client/example"
	server "github.com/superisaac/rpctube/server"
	"log/syslog"
	"os"
	//"strings"
	//tube "github.com/superisaac/rpctube/tube"
)

func setupLogger() {
	envLogOutput := os.Getenv("LOG_OUTPUT")
	if envLogOutput == "" || envLogOutput == "console" || envLogOutput == "stdout" {
		log.SetOutput(os.Stdout)
	} else if envLogOutput == "stderr" {
		log.SetOutput(os.Stderr)
	} else {
		file, err := os.OpenFile(envLogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}

	envLogSyslog := os.Getenv("LOG_SYSLOG")
	if envLogSyslog != "disabled" && envLogSyslog != "no" && envLogSyslog != "false" {
		hook, err := logsyslog.NewSyslogHook("", envLogSyslog, syslog.LOG_INFO, "")
		if err != nil {
			panic(err)
		}
		log.AddHook(hook)
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	switch envLogLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expect subcommands")
		os.Exit(1)
	}

	setupLogger()

	switch os.Args[1] {
	case "server":
		server.CommandStartServer()
	case "listmethods":
		client.CommandListMethods()
	case "rpc":
		client.CommandCallRPC("rpc")
	case "call":
		client.CommandCallRPC("call")
	case "example.fifo":
		example.CommandExampleFIFO()
	default:
		fmt.Println("expect subcommands")
		os.Exit(1)
	}

}
