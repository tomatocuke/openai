### 03-03 被墙了
### 03-01更新，替换ChatGPT接口了！速度超快的！
### 持续优化中，喜欢的同学给个🌟关注一下

### 一、介绍
- 能干什么？ 
  - 通过调用`OpenAI`的接口智能回答问题。
  - 可以直接api调用
  - 可以用作公众号自动回复。（**声明**：仅用于个人娱乐体验，勿作商业用途）
- 免费吗？不算，`OpenAI`账号赠送18$，限期使用。 $0.002 / 1000 tokens，每次花费已经打印在日志里。（暂时没做上下文）
- 风险。依托于公众号存在一定风险。 我已加了[敏感词检测](https://github.com/tomatocuke/sieve) ，但不清楚微信具体机制此举是否有效。
- 体验。关注公众号`杠点杠`尝试提问，这仅是个人娱乐号，不推送。


### 公众号部署
> 如果仅用于api调用，忽略下文有关微信公众号的设置、忽略 WX_TOKEN 参数
1. 获取`API_KEY`。[OpenAI](https://beta.openai.com/account/api-keys) （如果访问被拒绝，注意全局代理，打开调试，Application清除LocalStorage后刷新，实测可以）
2. 获取微信公众号`令牌Token`：[微信公众平台](https://mp.weixin.qq.com/)->基本配置->服务器配置->令牌(Token)  (不使用公众号可调过)
3. 使用以上参数启动服务，以下两种方式选其一部署。(此处举例端口9001，如果用公众号且无域名须用80端口)
  - Docker
    ```bash
    docker run -p 9001:8080 -e API_KEY=xxx -e WX_TOKEN=xxx -d -v $PWD/log:/app/log tomatocuke/openai
    ```
  - Golang
    ```bash 
    git clone https://github.com/tomatocuke/openai.git
    cd openai
    go run main.go -PORT=9001 -API_KEY=xxx -WX_TOKEN=xxx 
    ```
4. 启动服务后简单测试 `curl 'http://127.0.0.1:9001/test?msg=怎么做锅包肉'` 
5. 查看日志 `tail ./log/data.log`
6. 公众号配置。 
  - 无域名。须用80端口部署，服务器地址(URL)填写 `http://服务器IP/wx`。
  - 有域名。nginx配置参考
    ```conf
    server {
      listen 80;
      server_name 域名; #你的域名，不带http，例: abc.com

      location /openai/ {
        proxy_pass http://127.0.0.1:9001/; # 服务端口号
      }
    }
    ```
    重新加载nginx配置`nginx -s reload`后，公众号服务器地址填写: `http://域名/openai/wx`。
    公众号设置明文方式传输，确认后点击「启用」。 （初次设置生效要等一会，过几分钟关闭再启用试试）
1. 直接调用api `curl 'http://域名/openai/test?msg=怎么做锅包肉'`
    

### 三、其他
- 有什么问题我github可能不及时查看，欢迎提问和交流，QQ:`772532526`
