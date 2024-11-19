FROM golang:alpine

RUN apk update && apk add --no-cache gcc musl-dev sqlite

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/forum/main.go

RUN apk del gcc musl-dev

EXPOSE 8080

CMD ["./main"]