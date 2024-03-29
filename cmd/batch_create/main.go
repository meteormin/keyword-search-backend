package main

import (
	"bufio"
	"encoding/csv"
	"github.com/joho/godotenv"
	"github.com/miniyus/gofiber/database"
	"github.com/miniyus/gollection"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"github.com/miniyus/keyword-search-backend/repo"
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
	db := database.New(config.Database["default"])

	repository := repo.NewSearchRepository(db)
	service := search.NewService(repository)

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

			gollection.NewCollection(rows).Chunk(100, func(v [][]string, i int) {
				var searchSlice []*search.CreateSearch

				gollection.NewCollection(v).For(func(v []string, j int) {
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

				create, err := service.BatchCreate(uint(hId), 1, searchSlice)

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
