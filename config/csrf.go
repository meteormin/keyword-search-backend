package config

import fCsrf "github.com/gofiber/fiber/v2/middleware/csrf"

func csrf() fCsrf.Config {
	return fCsrf.Config{}
}
