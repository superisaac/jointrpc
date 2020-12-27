package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/rpctube/client"
	example "github.com/superisaac/rpctube/client/example"
	server "github.com/superisaac/rpctube/server"
	"os"
	"strings"
	//"strings"
	//tube "github.com/superisaac/rpctube/tube"
)

var commands []string = []string{
	"server", "rpc", "call",
	"watch", "notify",
	"listmethods", "watchmethods",
	"example.array",
	"help",
}

func setupClientSideLogger(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})

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

	if logLevel == "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func showHelp() {
	fmt.Printf("usage: %s <command> [<args>]\n", os.Args[0])
	fmt.Printf("commands are: %s\n", strings.Join(commands, ","))
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		server.CommandStartServer()
	case "listmethods":
		setupClientSideLogger("")
		client.CommandListMethods()
	case "watchmethods":
		setupClientSideLogger("")
		client.CommandWatchMethods()
	case "rpc":
		setupClientSideLogger("")
		client.CommandCallRPC("rpc")
	case "call":
		setupClientSideLogger("")
		client.CommandCallRPC("call")
	case "notify":
		setupClientSideLogger("")
		client.CommandSendNotify()
	case "watch":
		setupClientSideLogger("")
		client.CommandWatch()
	case "example.array":
		setupClientSideLogger("")
		example.CommandExampleArray()
	case "help":
		showHelp()
	default:
		//fmt.Println("expect subcommands")
		showHelp()
		os.Exit(1)
	}

}
