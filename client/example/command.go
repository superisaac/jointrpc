package example

import (
	"flag"
	log "github.com/sirupsen/logrus"
	client "github.com/superisaac/rpctube/client"
	"os"
)

// Example ARRAY
func CommandExampleArray() {
	examFlags := flag.NewFlagSet("example.array", flag.ExitOnError)
	pAddress := examFlags.String("c", "localhost:50055", "the tube server address")
	pCertFile := examFlags.String("cert", "", "the tube cert files")

	examFlags.Parse(os.Args[2:])

	serverAddress, certFile := client.TryGetServerSettings(*pAddress, *pCertFile)
	log.Infof("dial server %s", serverAddress)

	err := ExampleArray(serverAddress, certFile)
	if err != nil {
		panic(err)
	}
}
