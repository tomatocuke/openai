### 说明
- 功能。该项目通过调用`ChatGPT`提供微信公众号自动回复服务 (需要你有`ChatGPT`账号)
- 限制。微信5s内收不到回复，会再重试2次，即单条消息最久15s，如果`ChatGPT`超时则没办法给出回复。
- 内容。`ChatGPT`的结果仅供参考，目测对中文的一些知识不友好，存在睁眼说瞎话的情况。
- 调节。可以自行查看代码和参考openai的文档进行`ChatGPT`调节参数


### 参数获取
- `ChatGPT`的`API_KEY`：翻墙，开启全局代理，[openai](https://beta.openai.com/account/api-keys) （如果访问被拒绝，注意全局代理，打开调试，Application清除LocalStorage后刷新，实测可以）
- 微信公众号令牌`TOKEN`：[微信公众平台](https://mp.weixin.qq.com/) -> 基本配置 -> 生成令牌 -> 按下边部署启动服务 -> 设置服务器地址 `http://x.x.x.x:8080/chatgpt`


### 项目部署
- Docker 
  ```bash
  git clone https://github.com/tomatocuke/chatgpt.git
  cd chatgpt
  docker build -t chatgpt-wechat .
  docker run -d --name=chatgpt-wechat -p 8080:8080 -e API_KEY=xxx -e WX_TOKEN=xxx chatgpt-wechat
  ```
- Golang运行
  ```bash 
  git clone https://github.com/tomatocuke/chatgpt.git
  cd chatgpt
  go run main.go -PORT=8080 -API_KEY=xxx -WX_TOKEN=xxx 
  ```


