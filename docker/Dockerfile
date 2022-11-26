FROM ubuntu:22.04

LABEL maintainer="miniyu97@gmail.com"

RUN mkdir /home/gofiber

ARG GO_GROUP
ARG GO_VERSION
ARG ARCH

RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

RUN apt -y upgrade \
    && apt update \
    && apt install -y supervisor \
    && apt install -y vim \
    && apt install -y curl \
    && apt install -y wget \
    && apt install -y make

RUN groupadd --force -g $GO_GROUP gofiber
RUN useradd -ms /bin/bash --no-user-group -g $GO_GROUP -u 1337 gofiber

WORKDIR /home/gofiber

RUN curl -O -L "https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz"
RUN tar xvzf go${GO_VERSION}.linux-${ARCH}.tar.gz -C /usr/local  &&  rm "go${GO_VERSION}.linux-${ARCH}.tar.gz"

COPY supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY start-container.sh /usr/local/bin/start-container

RUN chown -R gofiber:$GO_GROUP /home/gofiber
RUN chgrp -R $GO_GROUP /home/gofiber

RUN chmod +x /usr/local/bin/start-container

RUN export GOPATH=/usr/local/go/bin
RUN export PATH=$PATH:/usr/local/go/bin

EXPOSE $APP_PORT

CMD ["start-container"]