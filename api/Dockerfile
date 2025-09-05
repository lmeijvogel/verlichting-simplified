FROM golang:1.24-alpine AS builder

WORKDIR /src

# First copy go.mod and go.sum, then download deps (cached)
COPY go.mod ./
RUN go get github.com/gin-gonic/gin
RUN go get github.com/joho/godotenv
RUN go get golang.org/x/exp/slices

RUN go mod download
RUN go mod tidy

# Now copy the rest of the source
COPY * ./

# Build the binary
RUN go build -o api .

# Stage 2: Run
FROM golang:1.24-alpine

RUN apk add --no-cache tini

WORKDIR /app

COPY --from=builder /src/api /bin/api
COPY .env ./.env

ENTRYPOINT ["tini", "--"]

CMD /bin/api
