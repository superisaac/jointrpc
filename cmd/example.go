package cmd

import (
	"flag"
	//log "github.com/sirupsen/logrus"
	client "github.com/superisaac/jointrpc/client"
	example "github.com/superisaac/jointrpc/client/example"
	"os"
)

// Example ARRAY
func CommandExampleArray() {
	examFlags := flag.NewFlagSet("example.array", flag.ExitOnError)
	serverFlag := client.NewServerFlag(examFlags)
	examFlags.Parse(os.Args[2:])

	err := example.ExampleArray(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}

func CommandExampleBenchmark() {
	examFlags := flag.NewFlagSet("example.benchmark", flag.ExitOnError)
	serverFlag := client.NewServerFlag(examFlags)
	pConcurrency := examFlags.Int("con", 1, "the number of concurreny workers")
	examFlags.Parse(os.Args[2:])

	err := example.ExampleBenchmark(serverFlag.Get(), *pConcurrency)
	if err != nil {
		panic(err)
	}
}
