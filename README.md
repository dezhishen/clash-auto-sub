# CLASH-AUTO-SUB
## 说明
本项目是一个 Clash 订阅应用，可以将订阅的proxies部分替换本地配置文件中的proxies部分，调用接口,实现订阅的更新。
## 使用说明
              "CLASH_SUB_URL": "",
                "CLASH_CONF_PATH": "",
                "CLASH_URL": "",
                "CLASH_SECRET": "",
                "CLASH_CONF_PATH_IN_CLASH":""
```
docker run -d --restart=always \
  -e CLASH_SUB_URL="" -e CLASH_CONF_PATH="" \
  -e CLASH_URL="" -e CLASH_SECRET="" \
  -e CLASH_CONF_PATH_IN_CLASH="" \
  -e INTERVAL="3600" \
  dezhishen/clash-auto-sub
```
## 参数说明
| 参数 | 说明 |
| --- | --- |
| CLASH_SUB_URL | 订阅地址 |
| CLASH_CONF_PATH | 本地配置文件路径 |
| CLASH_URL | Clash API 地址 |
| CLASH_SECRET | Clash API 密钥 |
| CLASH_CONF_PATH_IN_CLASH | Clash中配置文件的路径,默认使用`CLASH_CONF_PATH`|