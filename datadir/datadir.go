package datadir

import (
        "os"
        "path/filepath"
	"log"
)

var dataDirPath string = ""

func SetDataDir(dir string) {
	dataDirPath = dir
}

func DataDir() string {
        if dataDirPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
                dataDirPath = filepath.Join(cwd, "rpctube")
        }
	return dataDirPath
}

func DataPath(path string) string {
	if path == "" {
		return DataDir()
	} else {
		return filepath.Join(DataDir(), path)
	}
}

func EnsureDataDir(dir string) string {
        dir = DataPath(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("make datadir %s", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
        return dir
}
