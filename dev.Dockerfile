FROM golang:alpine

WORKDIR /api
COPY . /api

RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/gin-gonic/gin

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main