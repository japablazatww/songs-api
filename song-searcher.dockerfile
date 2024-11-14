# song-searcher.dockerfile
FROM golang:1.18-alpine AS builder

# Instala dependencias de compilación
RUN apk add --no-cache gcc musl-dev

# Crea un usuario no root
RUN adduser -D -g '' appuser

# Establecer directorio de trabajo
WORKDIR /app

# Copia archivos go.mod y go.sum primero para aprovechar el caché de Docker
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o songSearcherApp ./cmd/api

# Imagen final
FROM alpine:latest

# Actualiza e instala certificados CA
RUN apk --no-cache add ca-certificates && \
    update-ca-certificates

# Crea usuario no root
RUN adduser -D -g '' appuser

# Crea directorios necesarios
RUN mkdir /app && \
    chown -R appuser:appuser /app

# Copia el binario compilado mas archivos necesarios
COPY --from=builder --chown=appuser:appuser /app/songSearcherApp /app/
COPY --from=builder --chown=appuser:appuser /app/.env /app/
COPY --from=builder --chown=appuser:appuser /app/cmd/api/internal/infraestructure/config/origin_weights.json /app


# Cambia al usuario no root
USER appuser

WORKDIR /app

CMD ["./songSearcherApp"]