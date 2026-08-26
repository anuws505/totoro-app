# Stage 1: Build
FROM golang:1.27-alpine AS builder

ARG CA_CERT_VER="20260611-r0"
ARG TZDATA_VER="2026c-r0"

RUN apk add --no-cache \
    ca-certificates=${CA_CERT_VER} \
    tzdata=${TZDATA_VER}

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/api/main.go

# Stage 2: Run
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

COPY --from=builder /app/main /app/main
RUN chmod 755 /app/main

USER 10001:10001

EXPOSE 8080

HEALTHCHECK --interval=60s --timeout=2s --start-period=15s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/healthz || exit 1

CMD ["./main"]
