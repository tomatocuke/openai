### 最近持续优化中，喜欢的同学给个🌟关注一下

### 一、介绍
- 能干什么？ 通过调用`OpenAI`的接口智能回答问题。(API调用 或 用作公众号自动回复)
- 是`ChatGPT`吗？  不是。`ChatGPT`基于GPT-3.5，本项目是调用GPT-3，有很大差距。现在`ChatGPT`还没开放接口，安全限制很高，现有市面基本都是此类冒充的。
- 免费吗？不算，`OpenAI`账号赠送18$，限期使用。 消耗根据问题和回复长度计算。
- 有什么不足？ 
  - 回复内容准确度仅供参考，更适合开放性问题。 
  - 不支持上下文。 (也方便做，但是花费更多的tokens）
  - 速度和回复长度很难兼得。如果是订阅号，只能被动回复，[限制最久15s做出回复](https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Passive_user_reply_message.html)，回复可能超时或者是截断的(做了缓存优化，可稍等再次提问直接获得答案)  。 非订阅号，48小时内有20条主动回复额度，这个版本也开发了，在dev分支，还不稳定。 
- 内容安全。我做了[敏感词检测](https://github.com/tomatocuke/sieve)
- 体验。关注公众号`杠点杠`尝试提问，这仅是个人娱乐号，不推送。

### 二、部署
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
4. 启动服务后简单测试 `curl 'http://127.0.0.1:9001/test?msg=中国在哪个洲'` 
5. 查看日志 `tail ./log/data.log`
6. 公众号配置。 
  - 无域名。须用80端口部署，服务器地址(URL)填写 `http://x.x.x.x/` 你的服务器IP。
  - 有域名。nginx配置参考
    ```conf
    server {
      listen 80;
      server_name xxx.com; #你的域名

      location / {
        proxy_pass http://127.0.0.1:9001/; # 服务端口号
      }
    }
    ```
    重新加载nginx配置`nginx -s reload`后，公众号服务器地址填写: `http://xxx.com/`。(设置失败的话，`curl 'http://xxx.com/test?msg=中国在哪个洲'` 看看公网能不能访问)
    启用公众号服务器配置  (初次设置可能要等待2分钟生效）


### 三、其他
- 有什么问题我github可能不及时查看，加QQ:`772532526`
- 寻求一份杭州的Golang开发工作，HTTP方面。
