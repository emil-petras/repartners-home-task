FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./main"]
