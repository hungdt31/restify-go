########################################
# Stage 1: Builder - build static binary
########################################
FROM golang:1.25-alpine AS builder

# Cài đặt các dependencies cần thiết
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Tách bước dependency để tận dụng cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn
COPY . .

# Build binary tĩnh cho Linux
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags="-s -w -extldflags '-static'" -o /out/main .

########################################
# Stage 2: Runtime - tối giản để chạy
########################################
FROM alpine:latest AS runtime

# Cài đặt ca-certificates cho HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary từ stage builder
COPY --from=builder /out/main /app/main

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Tạo user không phải root để chạy ứng dụng
RUN adduser -D -s /bin/sh appuser && \
    chown -R appuser:appuser /app

# Chuyển sang user không phải root
USER appuser

# Cổng lắng nghe của ứng dụng
EXPOSE 8080

# Health check để đảm bảo ứng dụng đang chạy
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Lệnh khởi động container
ENTRYPOINT ["/app/main"]
