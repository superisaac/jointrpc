package example

import (
	"flag"
	"log"
	"os"
	//client "github.com/superisaac/rpctube/client"
)

// Example FIFO
func CommandExampleFIFO() {
	examFlags := flag.NewFlagSet("example.fifo", flag.ExitOnError)
	serverAddress := examFlags.String("address", "localhost:50055", "the tube server address")
	examFlags.Parse(os.Args[2:])
	log.Printf("dial server %s", *serverAddress)
	err := ExampleFIFO(*serverAddress)
	if err != nil {
		panic(err)
	}
}
