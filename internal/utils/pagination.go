package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Paginator struct {
	Page
	TotalCount int64       `json:"total_count"`
	Data       interface{} `json:"data"`
}

func GetPageFromCtx(c *fiber.Ctx) (Page, error) {
	pageQuery := c.Query("page", "1")
	pageSizeQuery := c.Query("page_size", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return Page{1, 10}, err
	}

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil {
		return Page{1, 10}, err
	}

	return Page{
		Page: page, PageSize: pageSize,
	}, nil
}

func Paginate(pageInfo Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pageInfo.Page
		if page == 0 {
			page = 1
		}

		pageSize := pageInfo.PageSize
		if pageSize <= 0 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}
