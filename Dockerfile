FROM golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main cmd/forum/main.go

FROM alpine:latest

RUN apk update && apk add --no-cache sqlite
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/internal/db internal/db
COPY --from=builder /app/static static/
COPY --from=builder /app/web web/

EXPOSE 8080
CMD ["./main"]