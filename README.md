# im中心项目

## 功能设计
### 用户设计
    采用用户中心支持的oauth2密码式，前端直接与用户中心通信

### 模块设计
#### 业务层
- 可以自己设计，实现房间、群组、系统内的其他业务长链接通知

#### 连接层
- 实现websocket服务、tcp服务
- 连接的管理：存储、查询、删除、定时清除无用的连接
- 长链接事件处理：上线、下线、广播
- 路由处理：心跳、ping
- 对外http：获取连接列表、获取连接详情、连接下线、获取在线用户列表、单发消息、多发消息、获取服务状态
- 对外rpc：获取连接列表、获取连接详情、连接下线、获取在线用户列表、单发消息、多发消息、获取服务状态
- 通过消息队列进行处理峰值流量（发消息的接口流量），rpc/http发消息的接口，请求达到一定阀值就批量把消息写入到消息队列
各节点连接管理器起一个协程，来轮训拉取消息队列的消息，进行下发

#### 依赖模块
- 用户中心，提供用户管理、oauth2等服务
- 文件中心，提供文件上传下载、管理等服务

### 连接关闭的情况
#### 客户端下线
- 读协程轮训连接对象报错，关闭连接对象，关闭读协程，关闭写协程，清理相关数据
#### 定时任务清理心跳超时连接下线
- 走下线事件流程
#### 请求接口主动下线
- 走下线事件流程

### 连接信息存储方案
* 节点是没那么容易挂掉的
* 如果节点挂掉，该节点所有连接也挂掉了；
* 方案1：其他节点长时间未收到该节点的心跳后，由一位节点（选举/redis分布式锁）来删除（后台任务批量删除）挂掉节点的所有连接缓存，并且所有未挂节点通知本节点所有连接（如果有必要）
* 方案2：其他节点长时间未收到该节点的心跳后，都同步该节点挂掉的状态，【并且通知本节点所有连接（如果有必要）】，该挂掉节点的连接缓存会在超时后自动删除（需实现心跳续超时时间），每次查到挂掉节点的连接时可以删了该连接
#### 方案1
- redis集群作为公共缓存，记录每个连接的信息，包括用户、设备、所在节点等
- 一个以（年/月/日，以体量来计算）为单位的用户操作日志表，记录用户操作的日志，包括用户操作的时间、操作类型、操作内容、操作结果、操作者等
- 以该表的单位生成缓存，规则是：只查每个用户的最近一次操作日志，如果为下线就是下线状态，反之就是在线，查不到就是下线状态
- 优点：把用户操作轨迹记录下来了，还利用了用户操作日志表，没有额外加表
- 缺点：业务逻辑比较复杂，需要处理一些特殊情况
#### 方案2
- 用户上线就缓存这个关系到redis集群，用户下线就删除这个缓存
- 优点：业务很清晰，所有节点直接判断是否存在该缓存就知道该用户是否上线

### 命令
    #构建网关代码-api
    #goctl api go -dir . -style goZero --home $HOME/.goctl/template -api xxx.api
    #构建rpc代码
    #goctl rpc proto -dir . -style goZero --home $HOME/.goctl/template -src xxx.proto
    #构建curl代码
    #goctl model mysql ddl -src="../../doc/init.sql" -dir="." -c --home $HOME/.goctl/template
    
    pkill -f  im-center
    
    # cd ../dataProcess
    # nohup go run dataprocess.go -f etc/dataprocess.yaml >> dataProcess.log &
    
    cd gateway
    goctl api plugin -plugin goctl-swagger="swagger -filename im-center.json" -api im-center.api -dir .
    # goctl api plugin -plugin goctl-swagger="swagger -filename im-center.json  -host 172.16.10.83:8888 " -api im-center.api -dir .
    nohup go run im-center.go -f etc/im-center.yaml >> gateway.log &