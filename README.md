# Keyword Search Backend

## go lang with fiber framework

- gorm 사용
    - 현재 리포지토리에서는 postgreSQL만 사용중

## Install

```shell
git clone https://github.com/miniyus/gofiber.git

# 도커 컨테이너 이미지 및 gorm DB 드라이버는 postgreSQL을 사용 중입니다.
# 도커 컨테이너 생성 및 패키지 설치 명령
make deploy

# 로컬 실행 시
go mod download
make start # go run 명령 사용
# or
make build # go build 명령으로 빌드 후 실행
./build/gofiber
```

## Dot Env

```shell
APP_ENV=development
APP_NAME=go-fiber
APP_PORT=8080
LOCALE=ko_KR
TIME_ZONE=Asia/Seoul

GO_GROUP=20 # echo $GID / Mac OS인 경우..
GO_USER=501 # echo $UID / Mac OS인 경우..

DB_HOST=go-pgsql # docker compose로 구성한 경우
# DB_HOST=localhost
DB_DATABASE=go_fiber
DB_PORT=5432
DB_USERNAME=smyoo
DB_PASSWORD=smyoo
DB_AUTO_MIGRATE=true

REDIS_HOST=go-redis # docker compose로 구성한 경우
# REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

CREATE_ADMIN=true
CREATE_ADMIN_USERNAME=admin
CREATE_ADMIN_PASSWORD=admin
CREATE_ADMIN_EMAIL=admin@email.com
```

## Directory Structure

### miniyus/gofiber

- [miniyus/gofiber](https://github.com/miniyus/gofiber)

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

### Routes

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

### Middleware

```go
package main

import (
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
  "github.com/miniyus/gofiber"
  "github.com/miniyus/gofiber/app"
  "github.com/miniyus/gofiber/config"
)

func main() {
  cfg := config.GetConfigs()
  a := gofiber.New(cfg)
  a.Middleware(func(fiberApp *fiber.App, app app.Application) {
    // fiber middlewares...
    fiberApp.Use(logger.New())
  })
}

```

### IOC Container

```go
package main

import (
  "github.com/miniyus/gofiber"
  "github.com/miniyus/gofiber/app"
  "github.com/miniyus/gofiber/config"
)

func main() {
  cfg := config.GetConfigs()
  a := gofiber.New(cfg)
  // Register 메서드는 아래의 구조체 바인딩이 길어질 경우를 대비
  // 외부 패키지에서 작성하여 Register 메서드로 전달하여 실행 시킬 수 있다.
  a.Register(func(app app.Application) { 
    // Bind 메서드는 포인터 구조체의 타입을 기억 
    // 매핑된 클로저(함수)를 통해 구조체를 필요할 때 생성하여 사용 할 수 있다.  
    app.Bind(&cfg, func() *config.Configs {
      return cfg
    })
  })
  a.Singleton(&cfg) // 싱글톤 패턴으로 구조체를 바인딩하고 Resolve() 메서드로 가져올 수 있다.
  a.Resolve(&cfg)   // 바인딩된 구조체를 가져온다.

}

```

### 기타 패키지

- [miniyus/gollection](https://github.com/miniyus/gollection)
- [miniyus/gorm-extension](https://github.com/miniyus/gorm-extension)
- [miniyus/goworker](https://github.com/miniyus/goworker)