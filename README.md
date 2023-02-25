# Keyword Search Backend

## go lang with fiber framework

- gorm 사용
    - 현재 리포지토리에서는 postgreSQL만 사용중

## Install

```shell
git clone https://github.com/miniyus/gofiber.git

# 도커 컨테이너 및 gorm DB 드라이버는 postgreSQL을 사용 중입니다.
docker compose up -d --build 

# 로컬 실행 시
go mod download
make start
# or
make build
./build/gofiber
```

## Dot Env

```shell
APP_ENV=develop
APP_NAME=go-fiber
APP_PORT=8080

TIME_ZONE=Asia/Seoul

GO_GROUP=1000
GO_USER=1000

DB_HOST=go-pgsql
DB_DATABASE=go_fiber
DB_PORT=5432
DB_USERNAME=?
DB_PASSWORD=?
DB_AUTO_MIGRATE=true

```

## Directory Structure

### miniyus/gofiber

- [gofiber](https://github.com/miniyus/gofiber)

### config

- 설정 관리
- go 언어의 구조체를 활용하여 관리

```go
package main

import "github.com/miniyus/keyword-search-backend/config"

func main() {
	config.GetConfigs()
}

```

### internal

- api 구현

### routes

- 라우팅

```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
)

const ApiPrefix = "/api"

func Api(router app.Router, a app.Application) {
	router.Route(
		"test",
		func(r fiber.Router) {
			r.Get("/", func(ctx *fiber.Ctx) error {
				return ctx.JSON("test")
			})
		},
	).Name("api.test")
}
```

```go
package main

import (
	"github.com/miniyus/gofiber"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/gofiber/config"
	"github.com/miniyus/keyword-search-backend/routes"
)

func main() {
	cfg := config.GetConfigs()
	a := gofiber.New(cfg)
	a.Route("prefix", func(r app.Router) {
		routes.Api(r, a)
	}, "router group name")
}

```

### 기타 외부 패키지

- [miniyus/gollection](https://github.com/miniyus/gollection)
- [miniyus/gorm-extension](https://github.com/miniyus/gorm-extension)
- [miniyus/goworker](https://github.com/miniyus/goworker)