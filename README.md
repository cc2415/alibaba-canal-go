# alibaba-canal-go
阿里巴巴canal
proto要使用这个
"github.com/golang/protobuf/proto"

# 启动alibaba-canal服务
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
# 修改配置
``` golang
AppConfig.EsAddress = []string{"http://127.0.0.1:9200/", "http://127.0.0.1:9201/", "http://127.0.0.1:9202/"}
	AppConfig.AlibabaCanal.Address = "192.168.199.165" //alibaba-canal的服务的ip地址
	AppConfig.AlibabaCanal.Port = 11111 //是alibaba-canal 的端口
	AppConfig.AlibabaCanal.Username = "alibaba_canal" //是要被同步的数据库的账号
	AppConfig.AlibabaCanal.Password = "abcd" //是要被同步的数据库的密码
	AppConfig.AlibabaCanal.Destination = "destination" //是alibaba-canal的服务的名字，自定义的
	AppConfig.AlibabaCanal.Database = "database"
```

# 增加新的表同步到es
复制table目录下的chatMsg

修改表名和struct的内容数据就ok了
