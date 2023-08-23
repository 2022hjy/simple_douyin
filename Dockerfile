
FROM golang:1.21-alpine
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
WORKDIR /app
COPY . .

# 构建后端和安装环境
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o simple_douyin main.go \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add ffmpeg

# 暴露端口
EXPOSE 8080

CMD ["/app/simple_douyin"]