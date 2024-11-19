# Étape de construction
FROM golang:alpine AS builder

# Installer les dépendances nécessaires
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copier les fichiers go.mod et go.sum
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le reste du code source
COPY . .

# Compiler l'application avec CGO activé
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main cmd/forum/main.go

# Étape finale
FROM alpine:latest

# Installer les bibliothèques nécessaires pour SQLite
RUN apk add --no-cache sqlite

WORKDIR /root/

COPY --from=builder /app .

EXPOSE 8080

CMD ["./main"]
