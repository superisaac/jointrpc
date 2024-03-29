package cmd

import (
	//"context"
	"flag"
	"fmt"
	"github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/playbook"
	"os"
)

func CommandPlaybook() {
	pbFlags := flag.NewFlagSet("playbook", flag.ExitOnError)
	serverFlag := client.NewServerFlag(pbFlags)

	pbFlags.Parse(os.Args[2:])

	if pbFlags.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "need playbook.yml\n")
		os.Exit(1)
	}
	pbYamlPath := pbFlags.Args()[0]

	if _, err := os.Stat(pbYamlPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "playbook yaml %s does not exist\n", pbYamlPath)
		os.Exit(1)
	}

	pb := playbook.NewPlaybook()
	err := pb.Config.ReadConfig(pbYamlPath)
	if err != nil {
		panic(err)
	}
	err = pb.Run(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}
