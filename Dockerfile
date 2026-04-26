FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY web /app/web

RUN chmod +x /app/server

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

EXPOSE 7540

CMD ["/app/server"]