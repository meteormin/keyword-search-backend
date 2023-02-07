package config

import (
	mConfig "github.com/miniyus/gofiber/config"
	"log"
	"os"
	"path"
)

type Path struct {
	BasePath string
	DataPath string
	LogPath  string
}

func getPath() mConfig.Path {
	getWd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed getwd %v", err)
	}
	dataPath := path.Join(getWd, "data")
	if _, err := os.Stat(dataPath); err != nil {
		log.Println("Not exists data directory")
		log.Printf("Create %s", dataPath)
		err = os.Mkdir(dataPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	logPath := path.Join(dataPath, "logs")
	if _, err := os.Stat(logPath); err != nil {
		log.Println("Not exists logs directory")
		log.Printf("Create %s", logPath)
		err = os.Mkdir(logPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return mConfig.Path{
		BasePath: getWd,
		DataPath: dataPath,
		LogPath:  logPath,
	}
}
