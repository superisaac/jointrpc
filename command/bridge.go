package command

import (
	"context"
	"flag"
	"fmt"
	client "github.com/superisaac/jointrpc/client"
	bridge "github.com/superisaac/jointrpc/cluster/bridge"
	"os"
)

func CommandStartBridge() {
	commandFlags := flag.NewFlagSet("bridge", flag.ExitOnError)
	pCertFile := commandFlags.String("cert", "", "cert file")

	commandFlags.Parse(os.Args[2:])
	if commandFlags.NArg() < 1 {
		fmt.Println("server addresses...")
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
	bridge.StartNewBridge(context.Background(), serverEntries)
}
