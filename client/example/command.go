package example

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

// Example FIFO
func CommandExampleFIFO() {
	examFlags := flag.NewFlagSet("example.fifo", flag.ExitOnError)
	serverAddress := examFlags.String("address", "localhost:50055", "the tube server address")
	examFlags.Parse(os.Args[2:])
	log.Infof("dial server %s", *serverAddress)
	err := ExampleFIFO(*serverAddress)
	if err != nil {
		panic(err)
	}
}
