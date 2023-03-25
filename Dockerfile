FROM golang:1.20-alpine as build

RUN apk --no-cache add tzdata && \
	cp /usr/share/zoneinfo/Asia/Seoul /etc/localtime && \
	echo "${TIME_ZONE}" > /etc/timezone

RUN apk add --no-cache make

RUN mkdir /fiber

WORKDIR /fiber

COPY . .

RUN CGO_ENABLED=0 go mod download
RUN CGO_ENABLED=0 go install github.com/swaggo/swag/cmd/swag@latest
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} make build

FROM alpine AS base

LABEL maintainer="miniyu97@gmail.com"

RUN apk --no-cache add tzdata && \
	cp /usr/share/zoneinfo/Asia/Seoul /etc/localtime && \
	echo "Asia/Seoul" > /etc/timezone \
	apk del tzdata

RUN mkdir /home/gofiber

ARG SELECT_ENV
ARG GO_GROUP
ARG GO_VERSION

FROM base AS deploy

WORKDIR /home/gofiber

COPY --from=build /fiber/build ./build

COPY .env${SELECT_ENV} ./.env

EXPOSE $APP_PORT

CMD ["./build/gofiber"]