FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o my-application cmd/controller/main.go

CMD ["./my-application"]
