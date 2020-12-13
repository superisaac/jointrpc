package example

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

// Example ARRAY
func CommandExampleArray() {
	examFlags := flag.NewFlagSet("example.array", flag.ExitOnError)
	serverAddress := examFlags.String("address", "localhost:50055", "the tube server address")
	examFlags.Parse(os.Args[2:])
	log.Infof("dial server %s", *serverAddress)
	err := ExampleArray(*serverAddress)
	if err != nil {
		panic(err)
	}
}
