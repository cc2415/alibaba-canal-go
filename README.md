# alibaba-canal-go
阿里巴巴canal proto要使用这个 "github.com/golang/protobuf/proto"

# 功能
可设置数据库表是否需要初始化
自动跟踪表数据的增删改查同步es
设置副本数量（0就默认为es节点-1个）

# 配置
```yaml
esAddress: # es地址
  - "http://localhost:9200"
  - "http://localhost:9201"
  - "http://localhost:9202"

alibabaCanal: # alibaba-cnal服务地址
  address: 192.168.199.165 #alibaba-canal的服务的ip地址
  port: 11111 #alibaba-canal服务 的端口
  username: alibaba_canal #是要被同步的数据库的账号
  password: abcd #是要被同步的数据库的密码
  destination: destination #是alibaba-canal的服务的名字，自定义的
  database: databaseName #要被同步的数据库

databaseName: databaseName # 数据库名

esRolloverCondition: # es索引配置
  maxAge: 1d
  maxDocs: 10
  maxSize: 5gb
esNodeNumber: 3 #节点数量
esIndexShareReplicasNum: 2 # 副本数量
```

# 启动
## 启动alibaba-canal服务
要求数据库开启了binlog

bindlog使用的是ROW模式

创建一个账号，给这个账号对应库的权限

[可参考这个](https://www.codeccc.cn/index.php/2024/03/18/%e9%98%bf%e9%87%8c%e7%9a%84-binlog-%e7%9a%84%e5%a2%9e%e9%87%8f%e8%ae%a2%e9%98%85%e5%92%8c%e6%b6%88%e8%b4%b9%e7%bb%84%e4%bb%b6/)

---

Dockerfile
```
FROM canal/canal-server
```

docker-composer.yml
```yml
version: '3.1'
services:
  cc-alibaba-canal:
    build:
      context: ./
      dockerfile: Dockerfile
    container_name: alibaba-canal
    privileged: true
    environment:
      - canal.auto.scan=true
      - canal.destinations=destination # 这个是自定义的，随便写，后面链接服务的时候也是要和这个一样
      - canal.instance.master.address=127.0.0.1:3306 # 注意ip
      - canal.instance.dbUsername=alibaba_canal # 账号
      - canal.instance.dbPassword=abcd # 密码
      - canal.instance.connectionCharset=UTF-8
      - canal.instance.tsdb.enable=true
      - canal.instance.gtidon=false
    hostname: 192.168.199.165 # 本机ip地址
    ports:
      - "11110:11110"
      - "11111:11111"
      - "11112:11112"
      - "9100:9100"
    mem_limit: "4096m"
  cc-alibaba-admin: # 这个镜像暂时用不上，可以去掉
    image: canal/canal-admin
    container_name: canal-admin
    privileged: true
    environment:
      - server.port=8089
      - canal.adminUser=admin
      - canal.adminPasswd=admin
    ports:
      - "8089:8089"
    mem_limit: "1024m"
    hostname: 192.168.199.165 # 注意ip



```
## 配置env
复制 .env.template.yaml 到 .env 根据备注修改
## 修改env配置

## 设置表数据更新是否需要初始化和同步到es
复制table目录下的chatMsgModel
```go
package table

type ChatMsgModel struct { //自定义
	BaseModel
}

var chatMsgModel = ChatMsgModel{
	BaseModel{
		TableName:    "fa_chat_msg", //表名
		NeedInitData: true, //是否需要初始化表数据到es
		NeedInEs:     true, // 是否追踪增删改查到es
	}}

func init() { // 这个一定要有，不要动
	ModelList = append(ModelList, chatMsgModel)
}

```
