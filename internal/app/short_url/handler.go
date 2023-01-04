package short_url

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"go.uber.org/zap"
)

type Handler interface {
	Redirect(c *fiber.Ctx) error
	FindUrlByCode(c *fiber.Ctx) error
	logger.HasLogger
}

type HandlerStruct struct {
	service Service
	logger.HasLoggerStruct
}

func NewHandler(service Service, log *zap.SugaredLogger) Handler {
	return &HandlerStruct{
		service:         service,
		HasLoggerStruct: logger.HasLoggerStruct{Logger: log},
	}
}

// FindUrlByCode
// @Summary find url by code
// @description find url by code
// @Tags ShortUrl
// @Param code path string true "short url code"
// @Success 200 {object} Response
// @Success 404
// @Accept json
// @Produce json
// @Router /api/short-url/{code} [get]
// @Security BearerAuth
func (h *HandlerStruct) FindUrlByCode(c *fiber.Ctx) error {
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

	return c.JSON(Response{
		Url: realUrl,
	})
}

// Redirect
// @Summary redirect real url from short-url
// @description redirect real url from short-url
// @Tags ShortUrl
// @Param code path string true "short url code"
// @Success 200 {object} Response
// @Success 302 {string} redirect
// @Accept json
// @Produce json
// @Router /api/short-url/{code}/redirect [get]
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

	err = c.Redirect(realUrl)
	if err != nil {
		return c.JSON(Response{
			Url: realUrl,
		})
	}

	return err
}
