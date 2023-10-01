FROM golang:1.21

WORKDIR /src

COPY go.mod ./
RUN go get github.com/gin-gonic/gin
RUN go get github.com/joho/godotenv
RUN go get golang.org/x/exp/slices

COPY * ./
RUN go mod tidy

RUN mkdir bin
RUN go build -o bin/api

CMD bin/api
