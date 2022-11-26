package config

import (
	"log"
	"os"
	"path"
)

type Path struct {
	BasePath string
	DataPath string
}

func GetPath() Path {
	getWd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed getwd %v", err)
	}

	return Path{
		BasePath: getWd,
		DataPath: path.Join(getWd, "data"),
	}
}
