package entity

import (
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gorm-extension/gormhooks"
)

func RegisterHooks(a app.Application) {
	gormhooks.Register(&Host{})
	hostHandler := newHostHookHandler(a)
	hostHooks := gormhooks.GetHooks(&Host{})
	hostHooks.HandleAfterSave(hostHandler.HostAfterSave)
}
