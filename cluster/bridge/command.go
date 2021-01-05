package bridge

import (
	"context"
	"flag"
	"fmt"
	client "github.com/superisaac/jointrpc/client"
	"os"
)

func printHelp() {
	fmt.Println("server addresses...")
}

func CommandStartBridge() {
	commandFlags := flag.NewFlagSet("bridge", flag.ExitOnError)
	pCertFile := commandFlags.String("cert", "", "cert file")

	commandFlags.Parse(os.Args[2:])
	if commandFlags.NArg() < 1 {
		printHelp()
		os.Exit(1)
	}

	args := commandFlags.Args()

	serverEntries := make([]client.ServerEntry, 0)

	for _, connect := range args {
		serverEntries = append(serverEntries, client.ServerEntry{
			ServerUrl: connect,
			CertFile:  *pCertFile,
		})
	}
	StartNewBridge(context.Background(), serverEntries)
}
