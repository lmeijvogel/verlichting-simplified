FROM golang:1.21-alpine

WORKDIR /src

COPY go.mod ./
RUN go get github.com/gin-gonic/gin
RUN go get github.com/joho/godotenv
RUN go get golang.org/x/exp/slices

RUN go mod tidy

COPY * ./

RUN mkdir bin
RUN go build -o bin/api

CMD bin/api
