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

// Redirect
// @Summary create search
// @description create search
// @Tags Redirect
// @Param code path string true "short url code"
// @Success 200 {object} RedirectResponse
// @Success 302 {string} redirect
// @Accept json
// @Produce json
// @Router /api/redirect/{code} [get]
// @Security BearerAuth
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
		return c.JSON(RedirectResponse{
			Url: realUrl,
		})
	}

	return err
}
