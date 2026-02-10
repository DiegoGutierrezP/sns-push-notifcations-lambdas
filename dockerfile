
# ------------------------------------------------------------
# Stage 1: Builder (compila todos los binarios de tus Lambdas)
# ------------------------------------------------------------
FROM golang:1.25 AS builder

# Variables de build
ARG TARGETOS=linux
# ARG TARGETARCH=arm64
ARG TARGETARCH=amd64
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
    go build -ldflags="-s -w" -o /out/bin/device-subscribe ./cmd/device-subscribe && \
    go build -ldflags="-s -w" -o /out/bin/device-unsubscribe     ./cmd/device-unsubscribe && \
    go build -ldflags="-s -w" -o /out/bin/publish-notification ./cmd/publish-notification && \
    go build -ldflags="-s -w" -o /out/bin/register-device    ./cmd/register-device && \
    go build -ldflags="-s -w" -o /out/bin/update-device    ./cmd/update-device && \
    go build -ldflags="-s -w" -o /out/bin/delivery-log-processor    ./cmd/delivery-log-processor

RUN chmod +x /out/bin/*

# (Opcional) Generar ZIPs para despliegue tradicional
RUN apt-get update && apt-get install -y zip && \
    (cd /out/bin && for f in device-subscribe device-unsubscribe publish-notification register-device update-device delivery-log-processor; do zip -q "$f.zip" "$f"; done)

# ------------------------------------------------------------
# Stage 2: Export (artefactos listos en /out/bin)
# ------------------------------------------------------------

FROM public.ecr.aws/lambda/provided:al2

WORKDIR /var/task
COPY --from=builder /out/bin /var/task
# COPY .env .env

RUN printf '%s\n' \
    '#!/bin/sh' \
    'set -eu' \
    ': "${LAMBDA_BIN:?Debes definir LAMBDA_BIN (ej: register-device)}"' \
    'test -x "/var/task/${LAMBDA_BIN}" || { echo "No existe o no es ejecutable: /var/task/${LAMBDA_BIN}"; ls -la /var/task; exit 1; }' \
    'exec "/var/task/${LAMBDA_BIN}"' \
    > "${LAMBDA_RUNTIME_DIR}/bootstrap" && chmod +x "${LAMBDA_RUNTIME_DIR}/bootstrap"

CMD ["handler"]