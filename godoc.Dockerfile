FROM golang:1.19-alpine AS build
WORKDIR /godoc

EXPOSE 6060

COPY go.mod go.mod
COPY go.sum go.sum

RUN CGO_ENABLED=0 go mod download
RUN CGO_ENABLED=0 go install golang.org/x/tools/cmd/godoc@latest

FROM build AS godoc

COPY . .
RUN rm -rf ./data

CMD ["godoc", "-http=:6060"]