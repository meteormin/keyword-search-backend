package config

import fCors "github.com/gofiber/fiber/v2/middleware/cors"

func cors() fCors.Config {
	return fCors.Config{}
}
