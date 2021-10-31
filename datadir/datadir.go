package datadir

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var datadirPath string = ""

func SetDatadir(dir string) {
	datadirPath = dir
}

func Datadir() string {
	if datadirPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		datadirPath = filepath.Join(cwd, ".jointrpc")
	}
	return datadirPath
}

func Datapath(path string) string {
	if path == "" {
		return Datadir()
	} else {
		return filepath.Join(Datadir(), path)
	}
}

func EnsureDatadir(dir string) string {
	dir = Datapath(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Debugf("make datadir %s", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
	return dir
}
