package main

import (
	"github.com/joho/godotenv"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/app/search"
	"github.com/miniyus/keyword-search-backend/internal/core/database"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("failed dotenv load")
	}

	config := configure.GetConfigs()
	db := database.DB(config.Database)
	loggerConfig := config.CustomLogger
	zLog := logger.New(logger.Config{
		TimeFormat: loggerConfig.TimeFormat,
		FilePath:   loggerConfig.FilePath,
		Filename:   loggerConfig.Filename,
		MaxAge:     loggerConfig.MaxAge,
		MaxBackups: loggerConfig.MaxBackups,
		MaxSize:    loggerConfig.MaxSize,
		Compress:   loggerConfig.Compress,
		TimeKey:    loggerConfig.TimeKey,
		TimeZone:   loggerConfig.TimeZone,
		LogLevel:   loggerConfig.LogLevel,
	})

	repo := search.NewRepository(db, zLog)
	buff, err := os.ReadFile(path.Join(config.Path.DataPath, "/batch/note.txt"))
	if err != nil {
		panic(err)
	}
	f := strings.Split(string(buff), "\n")

	service := search.NewService(repo)
	hostId := uint(1)

	var searchSlice []*search.CreateSearch

	for i, s := range f {
		query := strings.Trim(s, " ")
		query = strings.Trim(query, "\n")
		query = strings.Trim(query, "\t")

		searchSlice = append(searchSlice, &search.CreateSearch{
			HostId:      hostId,
			QueryKey:    "keyword",
			Query:       query,
			Description: strconv.Itoa(i),
			Publish:     true,
		})
		println(query)
	}

	//create, err := service.BatchCreate(hostId, searchSlice)
	//if err != nil {
	//	panic(err)
	//}
	//
	//println(create)
	if err == nil {
		_, err := service.BatchCreate(hostId, searchSlice)
		if err != nil {
			panic(err)
		}
	}
}
