FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/main ./cmd/api/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /out/main ./main

EXPOSE 8080

ENTRYPOINT ["/app/main"]
