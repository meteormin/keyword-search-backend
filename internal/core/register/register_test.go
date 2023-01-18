package register_test

import (
	"github.com/miniyus/keyword-search-backend/internal/core"
	"github.com/miniyus/keyword-search-backend/internal/core/register"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	c := core.New()
	register.Resister(c)

	var jobDispatcher worker.Dispatcher

	c.Resolve(&jobDispatcher)

	log.Print(jobDispatcher)
}
