FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o main .

# 环境变量
ENV API_KEY ""
ENV WX_TOKEN ""

EXPOSE "$PORT"

CMD ["./main"]
