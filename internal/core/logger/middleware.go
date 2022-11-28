package logger

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"strconv"
	"time"
)

func Middleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	elapsed := time.Since(start).Milliseconds()
	cu, ok := c.Locals(config.AuthUser).(auth.User)
	userID := ""
	if !ok {
		userID = "guest"
	} else {
		userID = strconv.Itoa(int(cu.Id))
	}
	req := c.Path()
	method := c.Method()

	errString := ""
	if err != nil {
		errString = fmt.Sprintf("| %s", err.Error())
	}

	GetLogger().Info(fmt.Sprintf("user: %4s | IP: %15s | %6dms | %s | %s %s",
		userID, c.IP(), elapsed, method, req, errString))

	c.Locals(config.Logger, GetLogger())

	return err
}
