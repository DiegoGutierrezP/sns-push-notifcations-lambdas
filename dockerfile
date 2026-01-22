
# ------------------------------------------------------------
# Stage 1: Builder (compila todos los binarios de tus Lambdas)
# ------------------------------------------------------------
FROM golang:1.22 AS builder

# Variables de build
ARG TARGETOS=linux
ARG TARGETARCH=arm64
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /src

# Cache de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código
COPY . .

# Compilar cada lambda (ajusta la lista si cambian)
RUN mkdir -p /out/bin && \
    go build -ldflags="-s -w" -o /out/bin/send-notification ./cmd/send-notification && \
    go build -ldflags="-s -w" -o /out/bin/create-device     ./cmd/create-device && \
    go build -ldflags="-s - w" -o /out/bin/subscription     ./cmd/subscription && \
    go build -ldflags="-s -w" -o /out/bin/unsubscription    ./cmd/unsubscription && \
    go build -ldflags="-s -w" -o /out/bin/update-device     ./cmd/update-device

# (Opcional) Generar ZIPs para despliegue tradicional
RUN apt-get update && apt-get install -y zip && \
    (cd /out/bin && for f in send-notification create-device subscription unsubscription update-device; do zip -q "$f.zip" "$f"; done)

# ------------------------------------------------------------
# Stage 2: Export (artefactos listos en /out/bin)
# ------------------------------------------------------------
FROM alpine:3.20 AS export
WORKDIR /out/bin
COPY --from=builder /out/bin /out/bin

# Los binarios y ZIPs quedan disponibles al hacer:
# docker build -t push-lambdas . &&
