FROM golang:1.20

WORKDIR /app
COPY go/src/myapp/go.mod go/src/myapp/go.sum ./
RUN go mod download

COPY go/src/myapp .
COPY go/src/myapp/.env .

RUN go build -o main . && chmod +x main

EXPOSE 8080
CMD ["./main"]

