# 多阶段构建：Go 1.22-alpine 编译，运行时内置 python3/nodejs/openjdk17 支撑在线评测沙箱
FROM golang:1.22-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata python3 nodejs openjdk17-jdk
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["/app/server"]
