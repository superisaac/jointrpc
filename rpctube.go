package main

import (
	"os"
	"fmt"
	"log"
	server "github.com/superisaac/rpctube/server"
	client "github.com/superisaac/rpctube/client"
	//tube "github.com/superisaac/rpctube/tube"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expect subcommands")
		os.Exit(1)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	switch os.Args[1] {
	case "node":
		server.StartEntrypoint()
	case "rpc":
		client.CallRPC()
	default:
		fmt.Println("expect subcommands")
		os.Exit(1)
	}

}
