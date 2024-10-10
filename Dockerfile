FROM golang:1.23.1-alpine

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]