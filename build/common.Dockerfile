FROM golang:1.21 AS builder

ARG TARGET

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/${TARGET}/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
