FROM golang:1.19-alpine as build

RUN apk --no-cache add tzdata && \
	cp /usr/share/zoneinfo/Asia/Seoul /etc/localtime && \
	echo "Asia/Seoul" > /etc/timezone

RUN apk add --no-cache make

RUN mkdir /fiber

WORKDIR /fiber

COPY ../go.mod .
COPY ../go.sum .
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY .. .

#RUN swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/gofiber/main.go --output ./api/gofiber
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} make build

FROM ubuntu:22.04

LABEL maintainer="miniyu97@gmail.com"

RUN mkdir /home/gofiber

ARG SELECT_ENV
ARG GO_GROUP
ARG GO_VERSION

RUN apt -y upgrade \
    && apt update \
    && apt install -y tzdata \
    && apt install -y supervisor \
    && apt install -y vim \
    && apt install -y curl \
    && apt install -y wget \
    && apt install -y make

RUN ln -snf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && echo ${TIME_ZONE} > /etc/timezone
RUN groupadd --force -g $GO_GROUP gofiber
RUN useradd -ms /bin/bash --no-user-group -g $GO_GROUP -u 1337 gofiber

WORKDIR /home/gofiber

#RUN curl -O -L "https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz"
#RUN tar xvzf go${GO_VERSION}.${GOOS}-${GOARCH}.tar.gz -C /usr/local  &&  rm "go${GO_VERSION}.linux-${GOARCH}.tar.gz"

COPY --from=build /fiber/build ./build

COPY .env${SELECT_ENV} ./.env
COPY supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ../scripts/start-container.sh /usr/local/bin/start-container

RUN chown -R gofiber:$GO_GROUP /home/gofiber
RUN chgrp -R $GO_GROUP /home/gofiber

RUN chmod +x /usr/local/bin/start-container

#RUN export GOPATH=/usr/local/go/bin
#RUN export PATH=$PATH:/usr/local/go/bin

EXPOSE $APP_PORT

CMD ["start-container"]