package main

import (
	"bufio"
	"encoding/csv"
	"github.com/joho/godotenv"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/internal/api/search"
	logger "github.com/miniyus/keyword-search-backend/logger"
	"github.com/miniyus/keyword-search-backend/utils"
	"log"
	"os"
	"path"
	"strconv"
)

type column int

const (
	hostId      column = 0
	queryKey    column = 1
	query       column = 2
	description column = 3
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
	service := search.NewService(repo)

	batchPath := path.Join(config.Path.DataPath, "/batch")

	files, err := os.ReadDir(batchPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			if path.Ext(file.Name()) != ".csv" {
				continue
			}
			log.Println(file.Name())

			f, err := os.Open(path.Join(batchPath, file.Name()))
			if err != nil {
				panic(err)
			}

			csvReader := csv.NewReader(bufio.NewReader(f))
			rows, err := csvReader.ReadAll()
			log.Println(rows)
			hId, err := strconv.Atoi(rows[1][hostId])
			if err != nil {
				panic(err)
			}

			utils.NewCollection(rows).Chunk(100, func(v [][]string, i int) {
				var searchSlice []*search.CreateSearch

				utils.NewCollection(v).For(func(v []string, j int) {
					if i*100+j == 0 {
						return
					}

					createSearch := csvToCreateSearch(v)

					if createSearch.HostId == uint(hId) {
						if createSearch.Description == "" {
							createSearch.Description = strconv.Itoa(i*100 + j)
						}

						searchSlice = append(searchSlice, createSearch)
					}
				})

				create, err := service.BatchCreate(uint(hId), searchSlice)

				log.Println(create)
				if err != nil {
					log.Println(err)
				}
			})
		}
	}

	if err != nil {
		panic(err)
	}

}

func csvToCreateSearch(data []string) *search.CreateSearch {
	hId, err := strconv.Atoi(data[hostId])
	if err != nil {
		panic(err)
	}

	q := data[query]
	qk := data[queryKey]
	desc := data[description]

	return &search.CreateSearch{
		HostId:      uint(hId),
		QueryKey:    qk,
		Query:       q,
		Description: desc,
		Publish:     true,
	}
}
