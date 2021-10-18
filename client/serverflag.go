package client

import (
	"flag"
	"os"
)

// misc functions
func NewServerFlag(flagSet *flag.FlagSet) *ServerFlag {
	seflag := new(ServerFlag)
	seflag.pAddress = flagSet.String("c", "", "the tube server address")
	seflag.pCertFile = flagSet.String("cert", "", "the cert file, default empty")
	return seflag
}

func (self *ServerFlag) ptrValue() ServerEntry {
	return ServerEntry{
		ServerUrl: *self.pAddress,
		CertFile:  *self.pCertFile,
	}
}

func (self *ServerFlag) Get() ServerEntry {
	value := self.ptrValue()
	if value.ServerUrl == "" {
		value.ServerUrl = os.Getenv("TUBE_CONNECT")
	}

	if value.ServerUrl == "" {
		value.ServerUrl = "h2c://localhost:50055"
	}

	if value.CertFile == "" {
		value.CertFile = os.Getenv("TUBE_CONNECT")
	}
	return value
}
