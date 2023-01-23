package main

import (
	"github.com/miniyus/keyword-search-backend/app"
	"github.com/miniyus/keyword-search-backend/create_admin"
)

func main() {
	create_admin.CreateAdmin(app.New())
}
