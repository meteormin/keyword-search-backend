package short_url

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/utils"
	"log"
)

type Handler interface {
	Redirect(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &HandlerStruct{
		service: service,
	}
}

func (h *HandlerStruct) Redirect(c *fiber.Ctx) error {
	allParams := c.AllParams()
	code := allParams["code"]
	if code == "" {
		return fiber.ErrNotFound
	}
	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	realUrl, err := h.service.FindRealUrl(code, user.Id)
	if err != nil {
		return err
	}
	log.Print(realUrl)
	err = c.Redirect(realUrl)
	if err != nil {
		return c.JSON(fiber.Map{
			"url": realUrl,
		})
	}

	return err
}
