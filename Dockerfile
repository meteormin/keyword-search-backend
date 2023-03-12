FROM golang:1.19-alpine as build

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

FROM ubuntu:22.04 AS base

LABEL maintainer="miniyu97@gmail.com"

RUN mkdir /home/gofiber

ARG SELECT_ENV
ARG GO_GROUP
ARG GO_VERSION

RUN apt-get -y upgrade \
    && apt-get update \
    && apt-get install -y tzdata supervisor vim curl wget make

RUN ln -snf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && echo ${TIME_ZONE} > /etc/timezone
RUN groupadd --force -g $GO_GROUP gofiber
RUN useradd -ms /bin/bash --no-user-group -g $GO_GROUP -u 1337 gofiber

#RUN curl -O -L "https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz"
#RUN tar xvzf go${GO_VERSION}.${GOOS}-${GOARCH}.tar.gz -C /usr/local  &&  rm "go${GO_VERSION}.linux-${GOARCH}.tar.gz"

FROM base AS deploy

WORKDIR /home/gofiber

COPY --from=build /fiber/build ./build

COPY .env${SELECT_ENV} ./.env
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY start-container.sh /usr/local/bin/start-container

RUN chown -R gofiber:$GO_GROUP /home/gofiber
RUN chgrp -R $GO_GROUP /home/gofiber

RUN chmod +x /usr/local/bin/start-container

#RUN export GOPATH=/usr/local/go/bin
#RUN export PATH=$PATH:/usr/local/go/bin

EXPOSE $APP_PORT

CMD ["start-container"]