package main

import (
	"github.com/miniyus/keyword-search-backend/internal/app"
	"github.com/miniyus/keyword-search-backend/internal/create_admin"
)

func main() {
	create_admin.CreateAdmin(app.New())
}
