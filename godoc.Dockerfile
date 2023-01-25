FROM golang:1.19-alpine as build
WORKDIR /godoc

EXPOSE 6060

COPY . .
RUN rm -rf ./data

RUN CGO_ENABLED=0 go mod download
RUN CGO_ENABLED=0 go install golang.org/x/tools/cmd/godoc@latest

CMD ["godoc", "-http=:6060"]