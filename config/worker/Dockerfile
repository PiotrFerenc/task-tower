FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o dwas-worker cmd/worker/main.go

CMD ["./dwas-worker"]
