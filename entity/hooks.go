package entity

import (
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gorm-extension/gormhooks"
)

func RegisterHooks(a app.Application) {
	gormhooks.Register(&Search{})
	searchHandler := newSearchHookHandler(a)
	searchHooks := gormhooks.GetHooks(&Search{})
	searchHooks.HandleAfterSave(searchHandler.AfterSave)
}
