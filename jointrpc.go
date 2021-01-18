package main

import (
	"fmt"
	"sort"
	//"context"
	log "github.com/sirupsen/logrus"
	//client "github.com/superisaac/jointrpc/client"
	//example "github.com/superisaac/jointrpc/client/example"
	//server "github.com/superisaac/jointrpc/server"
	command "github.com/superisaac/jointrpc/command"	
	"os"
	//"strings"
	//"strings"
	//bridge "github.com/superisaac/jointrpc/cluster/bridge"
)

var commands map[string]string = map[string]string{
	"server": "start join server",
	"rpc": "call jsonrpc method, the same as call",
	"call": "call jsonrpc method",
	"watch": "watch notifies and print them",
	"notify": "send notify",
	"watchstate": "watch server state changes",
	"methods": "list served methods",
	"delegates": "list delegated methods",
	"bridge": "run as a bridge between servers",
	"example.array": "start an example array service",
	"help": "print this methods",
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
	var arr []string
	for cmd, _ := range commands {
		arr = append(arr, cmd)
	}
	sort.Strings(arr)

	//fmt.Printf("commands are: %s\n", strings.Join(arr, ","))
	fmt.Printf("commands:\n")
	//for cmd, help := range commands {
	for _, cmd := range arr {
		help, _ := commands[cmd]
		if len(cmd) > 5 {
			fmt.Printf("  %s\t%s\n", cmd, help)			
		} else {
			fmt.Printf("  %s\t\t%s\n", cmd, help)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		command.CommandStartServer()
	case "methods":
		setupClientSideLogger("")
		command.CommandListMethods()
	case "delegates":
		setupClientSideLogger("")
		command.CommandListDelegates()
	case "watchstate":
		setupClientSideLogger("")
		command.CommandWatchState()
	case "rpc":
		setupClientSideLogger("")
		command.CommandCallRPC("rpc")
	case "call":
		setupClientSideLogger("")
		command.CommandCallRPC("call")
	case "notify":
		setupClientSideLogger("")
		command.CommandSendNotify()
	case "watch":
		setupClientSideLogger("")
		command.CommandWatch()
	case "bridge":
		setupClientSideLogger("")
		command.CommandStartBridge()
	case "example.array":
		setupClientSideLogger("")
		command.CommandExampleArray()
	case "help":
		showHelp()
	default:
		//fmt.Println("expect subcommands")
		showHelp()
		os.Exit(1)
	}

}
