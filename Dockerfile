FROM alpine

WORKDIR /app
# 需要先本地编译，手动 GOOS=linux GOARCH=amd64 go build -o openai
COPY openaiBin .
COPY keyword.txt .

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 环境变量
ENV API_KEY ""
ENV WX_TOKEN ""

EXPOSE "$PORT"

CMD ["./openaiBin"]
