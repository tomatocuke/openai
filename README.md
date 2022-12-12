
### 注意
- 本服务不需要翻墙，但是须翻墙获取api_key，注意全局模式，https://beta.openai.com/account/api-keys 。
- 返回的内容自行甄别。目测问它硬性知识有时候特别离谱。参数微调参考文档更改代码
- go版本根据你的环境直接更改go.mod即可，没引用外部包

### 两种使用方式
- 直接调用api
  `go run main.go -api_key $API_KEY`
- 公众号自动回复机器人
   1. https://mp.weixin.qq.com/ -> 基本配置
   2. 生成令牌Token
   3. 启动服务 `go run main.go -api_key $API_KEY -wx_token $TOKEN`
   4. 设置服务器地址 http://x.x.x.x:$PORT/chatgpt 

