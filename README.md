### 说明
- 功能。通过调用`OpenAI`提供微信公众号自动回复服务。（内容仅供参考）
- 注意。有别于网页版`ChatGPT`基于GPT-3.5，本项目是调用GPT-3。(网页版速度无法满足微信时间限制)
- 限制。公众号为被动回复，微信5s内收不到回复，会再重试2次，即单条消息最久15s，超时则没办法给出回复

### 准备
- [微信公众号](https://mp.weixin.qq.com/)
- [OpenAI账号](https://beta.openai.com) (需要梯子)
- `Docker`或`Golang`环境

### 操作流程
1. 获取`API_KEY`。[OpenAI](https://beta.openai.com/account/api-keys) （如果访问被拒绝，注意全局代理，打开调试，Application清除LocalStorage后刷新，实测可以）
2. 获取微信公众号令牌`TOKEN`：[微信公众平台](https://mp.weixin.qq.com/) -> 基本配置 -> 生成令牌 
3. 使用以上两个参数按照↓项目部署。（需要nginx代理到443或者80端口）
4. 继续设置公众号服务器地址 `http(s)://xxx/chatgpt` (接口路由为`chatgpt`)


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


