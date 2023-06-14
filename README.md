### 基于GPT-3.5的公众号自动回复机器人
### 持续优化中，喜欢的同学给个🌟关注一下

### 一、介绍
- 说明
  - 这是一个用于**公众号自动回复机器人**的项目。需要你有 OpenAI 账号、公众号、海外服务器或代理。
  - 花费。`OpenAI`账号赠送18$，限期使用。[按字算钱](https://platform.openai.com/tokenizer)，0.002刀 /1000 tokens 。[价格参考](https://openai.com/pricing)
  - 观点。我觉得公众号不是一个好的使用场景，订阅号入口麻烦，服务号需要资质且风险更大。所以如果你只是玩玩可以部署。
  - 体验。关注公众号`杠点杠`尝试提问，这仅是个人娱乐号，不推送。   别问预测和实事问题，它不会。

### 二、Feature
- [x] 解决微信被动回复限制问题。(设定超时时间，滚动返回)  
- [x] 支持用户语音输入。（要主动开启，设置与开发->接口权限->接收语音识别结果。已关注用户可能24小时内生效，可重新关注尝试）
- [x] 设置代理
- [x] prompt 提示、max_tokens、temperature 参数调节
- [x] [敏感词](https://github.com/tomatocuke/sieve)检测及自定义添加。(不清楚这样是否降低风险。代码内置隐藏了一些敏感词，你也可以启动时在根目录添加`keyword.txt`自定义敏感词。  如有敏感词误杀，你可以向我反馈)
- [ ] 上下文。(其实开发也不算难。主要是OpenAI不记录会话，上下文的本质是把之前的QA都作为新的参数传过去，这会叠加消耗token)
- [ ] 用户身份验证。(待开发)

### 三、部署
1. 获取`API_KEY`。[OpenAI](https://beta.openai.com/account/api-keys) （如果访问被拒绝，注意全局代理，打开调试，Application清除LocalStorage后刷新，实测可以）
2. 获取微信公众号`令牌Token`：[微信公众平台](https://mp.weixin.qq.com/)->基本配置->服务器配置->令牌(Token) 
3. 克隆项目，修改配置文件 `config.yaml`
4. 两种方式部署。（简单举例占用80端口，如果需要别的端口自己配置nginx等）

  - 直接二进制启动 (Linux amd64)
      ```sh
      mkdir log
      # 尝试启动
      ./openaiBin 
      # 守护进程启动
      nohup ./openaiBin >> log/data.log 2>&1 &
      ```
  - 使用Docker启动服务
      ```bash
      # 注意这里会拷贝配置到容器里，如果修改配置，需到容器内修改，或者启用新的容器
      docker run -d -p 80:80 -v $PWD/log:/app/log -v $PWD/config.yaml:/app/config.yaml tomatocuke/openai
      # 查看状况
      docker logs 容器ID 
      ```
5. 服务器地址(URL)填写 `http://服务器IP/wx`，设置明文方式传输，提交后，点击「启用」。
  
### 四、QA
- 日志出现 `openai请求超时` <br>
答：对openai的请求发不过去

- 出现报错  `connection reset by peer` 或 `Post "https://api.openai.com/v1/chat/completions": EOF`  <br>
答：是否使用了代理呢？ 大概率是IP被多人使用的结果，换个IP，但是部署不建议使用代理的方式，不稳定。

- 服务正常，但是公众号无响应？ <br>
答：初次设置生效要等一会，过几分钟把公众号的服务器设置按钮关闭再启用试试。

- 文档真特么烂，我部署不成功！ <br>
答：别github提问，我很少看。或者你有什么好的建议，欢迎加我QQ:`772532526`，不收费。
