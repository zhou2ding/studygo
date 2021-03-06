# day06

## 数据库

- 常见数据库`SQLlite`、`MySQL`、`SQLServer`、`Oracle`、`postgreSQL`

- 不同数据库的占位符不同

  |   数据库   |        占位符语法         |
  | :--------: | :-----------------------: |
  |   MySQL    |            `?`            |
  | PostgreSQL |       `$1`, `$2`等        |
  |   SQLite   |        `?` 和`$1`         |
  |   Oracle   | `:name`（name是字段名字） |

## MySQL

> 主流的关系型数据库

### 知识点

- SQL语句：`DQL`、`DML`、`DDL`、`DCL`、`TCL`
- 存储引擎：`InnoDB`、`MyISAM`，支持插件式的存储引擎

### database/sql包

- 原生支持连接池，是并发安全的

- 并没有具体的实现，只是列出了一些需要第三方库实现的具体内容

- 使用方法

  - `go get -u github.com/go-sql-driver/mysql`

  - `dsn := username:password@tcp(ipAddress)/database`

  - `sql.Open("msyql", dsn)`，返回一个`sql.DB`指针，不会校验用户名和密码

  - 在导入`mysql`的包时自动调用了其`init`方法，此方法向`database/sql`包中注册了`"mysql"`这个驱动

    ```go
    import (
    	"database/sql"
        
    	_ "github.com/go-sql-driver/mysql"
    )
    func main() {
    	dsn := "root:564710@tcp(localhost:3306)/bjpowernode"
    	_, err := sql.Open("mysql", dsn)
    	if err != nil {
    		fmt.Printf("open %s failed, error:%v\n", dsn, err)
    		return
    	}
    }
    ```

- `select SUBSTRING_INDEX(host,':',1) as ip , count(*) from information_schema.processlist group by ip;`获取本机连接mysql的IP地址

- `db.SetMaxOpenConns()`设置数据库连接池的最大建立连接的数量

- `db.SetMaxIdleConns()`设置数据库连接池的最大闲置连接数

#### 查询

- 单行查询

  1. 写查询单条记录的sql语句
  2. 执行查询(`QueryRow()`方法，接收一个字符串和可变长度的任意类型变量，返回`*sql.Row`对象)
  3. 拿到结果（用`Scan()`方法，且必须使用，因为此方法会释放数据库连接），不释放的话就会卡住，等待把连接归还给连接池

  ```go
  type user struct {
      id int
      name string
      age int
  }
  var u1 user
  sqlStr := `select id, name, age from user where id=?;` // ?是占位符
  db.QueryRow(sqlstr, 1).Scan(&u1.id, &u1.name, &u1.age)
  ```

- 多行查询

  1. SQL语句
  2. 执行`db.Query()`，返回`*db.Rows`对象和一个`error`
  3. `defer`关闭rows
  4. 循环取值`rows.Next()`

  ```go
  sqlStr := `select id, name, age from user where id > ?;`
  rows, _ := db.Query(sqlStr, 0)
  for rows.Next() {
      var u1 user
      _ = rows.Scan(&u1.id, &u1.name, &u1.age)
      fmt.Println(u1)
  }
  ```

#### 增删改

- `ret, err := db.Exec()`接收一个字符串和一个可变长度任意类型的变量，返回`db.Result`接口类型变量和`error`
  - 如果是插入操作，会拿到插入数据的`id`，`ret.LastInsertId()`，返回`id`和`error`
  - 如果是修改操作，会拿到受影响的行数
  - 如果是删除操作，会拿到受影响的行数

#### 预处理

>  好处

1. 优化MySQL服务器重复执行SQL的方法，可以提升服务器性能，提前让服务器编译，一次编译多次执行，节省后续编译的成本。
2. 避免SQL注入问题。

>  执行过程

1. 把SQL语句分成两部分，命令部分与数据部分。
2. 先把命令部分发送给MySQL服务端，MySQL服务端进行SQL预处理。
3. 然后把数据部分发送给MySQL服务端，MySQL服务端对SQL语句进行占位符替换。
4. MySQL服务端执行完整的SQL语句并将结果返回给客户端。

>  使用

- `db.Prepare()`接收一个字符串，返回准备好的状态`*db.Stmt`和`error`
- `defer stmt.Cloes()`
- 调用`stmt`的`QueryRow()`、`Query()`、`Exec()`方法执行操作

#### 事务

- `db.Begin()`，无参，返回一个`db.Tx`（transaction）和一个`error`，后续执行sql语句就调动`tx`的`Query()`、`QueryRow()`、`Exec()`方法
- `db.Commit()`，无参，返回一个`error`
- `db.RollBack()`，无参，返回一个`error`

#### sqlx

> 第三方库，更方便地使用mysql
>
> 结构体的字段必须大写，因为sqlx是通过反射获取字段信息的

- `sqlx.Connect("mysql", dsn)`，`open`数据库并`ping`数据库

- `db.Get(&user,sqlStr,id)`：查询单条，不用一个字段一个字段去修改了，直接修改整个结构体变量
- `db.Select(&userlist, sqlStr)`：查询多条，`userlist`是个切片，虽然是引用类型，但也要传它的引用，因为`sqlx`只对指针类型进行了处理，其他引用类型没管

#### sql注入

``xxx or 1=1 #`

``xxx union select * from user #`

## Redis

> KV数据库，支持master/slave模式

> 场景

- cache缓存
- 计数场景
- 简单的队列
- 排行榜

> 使用

```go
var resisdb *redis.Client
resisdb = redis.NewClient(&redis.Options{
    Addr: "xxx",
    Password: "",
    DB: 0, //0~15共16个DB
})
```

## 消息队列

常用的消息队列

- RabbitMQ
- Kafka
-  ActiveMQ 
- RocketMQ
- NSQ

### NSQ

#### 概述

> [NSQ](https://nsq.io/)是Go语言编写的一个开源的实时分布式内存消息队列，其性能十分优异。 NSQ的优势有以下优势：
>
> 1. NSQ提倡分布式和分散的拓扑，没有单点故障，支持容错和高可用性，并提供可靠的消息交付保证
> 2. NSQ支持横向扩展，没有任何集中式代理。
> 3. NSQ易于配置和部署，并且内置了管理界面。

![1618646437997](D:\资料\Go\src\studygo\Golang学习笔记\golang笔记_进阶.assets\1618646437997.png)

#### 使用

1. 启动NSQ三个组件
   1. `nsqlookupd.exe`
   2. `nsqd.exe -broadcast-address=127.0.0.1 -lookupd-tcp-address=127.0.0.1:4160`
   3. `nsqadmin.exe -lookupd-http-address=127.0.0.1:4161`
   4. 浏览器输入`http://127.0.0.1:4171`

# Beego

## 命令

```powershell
#先main.go中添加orm.RunCommand()，且要把结构体注册一下
#常用：go run main.go orm -v
go run main.go orm syncdb		#根据结构体建表，不带参数时，只根据结构体创建不存在的表
				==>-db=string	#指定数据库的别名，默认“default”
				==>-force		#建表前把已有的表先删除，慎用！！！不带此参数或指定为false时不执行此操作
				==>-v			#查看详情（verbose，打印sql语句等）
					
go run main.go orm sqlall		#打印建表语句
go run main.go orm help			#查看帮助，orm后面不带参数时默认带的help

#修改源码：cmd_utils.go中getColumnAddQuery()的return中添加
fi.description,
```

```bash
#打包成windows：打包完成后，安装nssm，然后nssm install <服务名>，然后启动服务即可
bee pack -be GOOS=windows
#如果报错，先执行如下命令再bee pack
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
```

```bash
#打包成linux：
SET CGO_ENABLED=0  // 禁用CGO
SET GOOS=linux  // 目标平台是linux
SET GOARCH=amd64  // 目标处理器架构是amd64
bee pack -be GOOS=linux
#linux上运行(&表示在后台运行)
#方式一
nohup ./要执行的文件 &
#方式二：
./nginx -c ../conf/nginx.conf
supervisord -c /etc/supervisord.conf
supervisorctl
start ions
```

## 日志模板

> 设计模板步骤：
>
> 1. 确定是常用的目标
> 2. 确定需要哪些字段
> 3. 将这些字段拼成想要的格式

```json
过程进度：starttime:%v | current/count:%v/%v | use_time:%v | endtime:%v
请求接口：url:%v | method:%v | body:%v | header:%v
数据库连接：ip:%v | port:%v | database:%v
抽象对象：parentStruct:%v | childStruct%v | params:%v
耗时：starttime:%v | use_time:%v | spend:%v | endtime:%v
```

==日志引擎==

- console
- file
- **multifile**
- smtp：日志警告，邮件发送
- conn：网络输出
- ElasticSearch：输出到ES
- ......

==cache引擎==

- memory
- file（不常用）
- redis
- memcache

## 设计思想

> 多利用\<input type="hidden" id="xx" value="{{.Id}}"/>标签来获取要操作的id值

- 删除：切记把`is_del`标志位改为1就行，不要直接删
- 批量删除：前端通过哪些框被选中这个属性来得到批量id，前端把这些id通过`,`拼接成字符串，后端再split成byte切片，再遍历改`is_del`
- 新增：
- 修改：可以在前端urlfor时就把id通过`"?id={{.Id}}"`传给后端，后端再通过id把此条数据的值回传给前端
- 搜索：触发搜索时就把关键字通过URL传给后端`"?kw="+kw`；后端查完后，要给前端传个kw，前端拿到后在换页的href中加上`&kw={{$.kw}}`避免换页后搜索条件被重置
- 封装json对象：先封装成servejson

# 微服务

> 单体架构：MVC（model、view、controller），在同一个服务（进程）中运行
>
> 缺点：逻辑复杂、需求变更成本高、不好扩展、更新技术成本高、新人融入成本高、维护成本高、代码腐化越来越严重

> 微服务：微服务就是一些具有足够小的粒度、能够相互协作而且自治的服务体系。
>
> 微服务设计：将复杂系统使用组件化的方式进行拆分，并通过通讯的方式进行整合的一种设计方法
>
> 特点：独立性，职能小，进程隔离，轻量型通信，灵活性高，和语言无关
>
> 缺点：技术难度大，运维要求高，不好调试，接口调整成本高，重复性代码多

## protoBuf

> 是一种语言无关、平台无关、可扩展的数据结构，是序列化的结构。
>
> 体积小，传输效率高，支持多平台跨语言，兼容性好，序列化和反序列化速度快

### 安装

- 首先下载[protoc-3.17.1-win64.zip](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.1/protoc-3.17.1-win64.zip)解压后，把`protoc.exe`放到`%gopath%/bin`
- 然后在`%gopath%/src`下面新建`google.golang.org\protobuf`
  - 先`git clone https://github.com/protocolbuffers/protobuf-go`
  - 再把`protobuf-go`目录下全部文件移动到`protobuf`下
  - 最后`go get -u github.com/golang/protobuf/protoc-gen-go`在`%gopath%/bin`下生成`protoc-gen-go.exe`

### 使用

```powershell
protoc my.proto --go_out=.	#把my.proto转成.go文件，转完放到当前目录下；支持通配符：*.proto；两个参数谁先谁后都行
```

### 语法

```protobuf
syntax = "proto3";					//protobuf的版本
option go_package = "./;protoDemo";	//指定生成的.go文件的package；./必须有
package protoDemo;					//指定package

//消息类型，属于messages，另外两种类型是enum和service；一个文件能有多个，但最好是同一种类型，不要多种类型混用，否则会导致依赖性膨胀
//可以嵌套，最外层的是level1，往里level依次增加，不同level的字段编号互不影响
message Hello {
    string name = 1;				//字段值只能是数字：表示字段的编号；19000~19999之间的不能用，是预留的
    int32 age = 2;					//字段类型有string、int32、int64等，或其他message类型
    optional string addr = 3;
    repeated addr hobbies = 4;
    RetHello ret = 6;
}
message RetHello {
	int32 code = 1;
	string msg = 2;
}

//枚举类型，字段没有类型，值必须从0开始
enum HelloEnum {
  RunStatus = 0;
  StopStatus = 1;
  DeleteStatus = 2;
}
```
```go
func main() {
    //protoc生成的.go文件，包名是go_package中指定的pbf
    phone := []int64{1,3,8,10}
    name := "zhangsan"
    helloData := &pbf.Hello{	//必须是指针
        Name:&name,				//optional规则的字段类型是指针
        Phone:phone,
        Addr: "河北",
    }
    fmt.Println(helloData.Addr)
    fmt.Println(*helloData.Name)
    fmt.Println(helloData.GetName())	//其他字段都有对应的Get方法
    
    //编解码：Marshal和Unmarshal
    ret,_ := proto.Marshal(helloData)
    helloData2 := pbf.Hello{}
    _ = proto.Unmarshal(ret,&helloData2)
}
```

> 字段规
>
> - optional:字段可出现0次或1次，为空可以指定默认值 [default=10]，否则使用语言的默认值
>
>   - optional int32 result_per_page = 3 [default = 10];
>
>     字符串默认为空字符串
>
>     数字默认为0
>
>     bool默认为false
>
>     枚举默认为第一个列出的值，一定要注意枚举的顺序，容易有坑
>
> - repeated：字段可出现任意多次（包括 0），数组或列表要使用这种
>
> - required：字段只能也必须出现1次，proto3中已废弃

### 字段类型

| .proto   | 说明     | Go语言  |
| :------- | :------- | :------ |
| double   | 浮点类型 | float64 |
| float    | 浮点类型 | float32 |
| int32    |          | int32   |
| int64    |          | int64   |
| uint32   |          | uint32  |
| uint64   |          | uint64  |
| sint32   |          | int32   |
| sint64   |          | int64   |
| fixed32  |          | uint32  |
| fixed64  |          | uint64  |
| sfixed32 |          | int32   |
| sfixed34 |          | int64   |
| bool     |          | bool    |
| string   |          | string  |
| bytes    |          | []byte  |

## grpc

### rpc基础

- remote procedure call：一种远程通讯协议，也是发送请求和接收响应，使用TCP

- 组成和过程：发过去四步，回来四步
  - 客户端：把调用请求发到客户端存根
  - 客户端存根(client stud)：收到请求后把信息序列化，发到保存的服务端的地址
  - 网络传输模块
  - 服务端存根(server stub)：解码消息（反序列化），调用服务
  - 服务端：提供真正的服务

- go实现rpc

  ```go
  //客户端
  func main() {
  	//与服务端创建连接
  	client,err := rpc.DialHTTP("tcp","127.0.0.1:8080")
  	if err != nil { return }
  	defer func() {
  		_ = client.Close()
  	}()
  	//调用server端提供的服务，三个参数分别是提供服务的方法，方法的入参，方法的返回
  	var output string
  	err = client.Call("User.SayHello","张三",&output)
  	fmt.Println("结果为：",output)
  }
  ```
  
  ```go
  //服务端
  type User struct {}
  
  func (u *User) SayHello(input string, output *string) error {
  	*output =  "hello, " + input
  	return nil
  }
  
  func main() {
  	usr := new(User)
  	//把usr对象注册到rpc服务中
  	_ = rpc.Register(usr)
  	//把usr提供的服务注册到HTTP协议上
  	rpc.HandleHTTP()
  	//监听tcp连接
  	listener,err := net.Listen("tcp","127.0.0.1:8080")
  	if err != nil { return }
  	//接收到监听器后服务启动
  	_ = http.Serve(listener, nil)
  }
  ```
  
  > 注意：对外暴露的服务方法定义标准
  >
  > 1、对外暴露的方法有且只能有两个参数，这两个参数只能是输出类型或内建类型，两种类型中的一种。
  > 2、方法的第二个参数必须是指针类型。
  > 3、方法的返回类型为error。

### grpc安装使用

1. `git clone https://github.com/grpc/grpc-go.git`，然后把文件夹grpc-go的名字改为grpc

2. `protoc hello.proto --go_out=plugins=grpc:. `，proto文件如下

   ```protobuf
   service TestService {								//服务端注册的就是这个service
       rpc SayHello(HelloReq) returns (HelloResp) {}; //接收请求返回响应，固定写法，服务端实现的就是这个方法
   }
   message HelloReq {
       string name = 1;
   }
   message HelloResp {
       string ret = 1;
   }
   ```

3. 服务端

   ```go
   //主要注意点：8行，15 16行，23 24行，30行
   package main
   
   import (
   	"context"
   	"fmt"
   	"google.golang.org/grpc"
   	hello "grpc_protobuf"
   	"net"
   	"net/http"
   )
   
   type User struct {}
   
   func (s *User) SayHello(ctx context.Context, input *hello.HelloReq) (output *hello.HelloResp,err error){
   	output = &hello.HelloResp{
   		Ret: "hello, " + input.Name,
   	}
   	return
   }
   
   func main() {
   	server := grpc.NewServer()
   	hello.RegisterTestServiceServer(server,&User{})
   	listener,err := net.Listen("tcp","127.0.0.1:8080")
   	if err != nil { return }
       _ = server.Serve(listener,nil)	//不能用http.Serve()，否则会报connection closed错误
   }
   ```

4. 客户端

   ```go
   func main() {
   	//WithInsecure：跳过证书的验证
   	conn,err := grpc.Dial("127.0.0.1:8080",grpc.WithInsecure())
   	if err != nil { return }
   	client := hello.NewTestServiceClient(conn)
   	resp, err := client.SayHello(context.Background(),&hello.HelloReq{
   		Input: "张三",
   	})
   	if err != nil { return }
   	fmt.Println(resp.Output)
   }
   ```

## consul

> 介绍：一个开源工具，是一个用来实现分布式系统的服务发现与配置的开源工具
>
> 特性：服务发现、健康检查、键值对存储、多数据中心等
>
> 术语：代理是consul集群中每个成员的守护进程；数据中心是一个私有、低延迟和高带宽的网络环境

### 安装使用

1. 下载后把exe放到`gopath/bin`下即可

2. ==配置conf.json==（服务注册）

   >  checks配置完成后启动consul，然后启动server.go即可把此server注册到consul的service中了

   ```json
   {
       "service":{
           id:节点id，name相同时靠id区分
           name:节点name，UI界面最外层的区分依据
           address:节点IP
           port:节点端口
           tags:节点标签，是个列表
           checks:服务检查
           //如下是TCP+ Interval方式，还有HTTP+ Interval和Script+ Interval的方式
           [{
               {	//一般只配置tcp和interval（每隔多少秒检查一次）就行
                   "id": "ssh",
                   "name": "SSHTCP on port 22",
                   "tcp": "192.168.0.105:8080",
                   "interval": "10s",
                   "timeout": "1s"
               }
           }
       }
   }
   
   //datacenter：同命令行参数-datacenter；data_dir：同命令行参数-data_dir；bootstrap：同命令行参数-bootstrap；bootstrap_expect：同命令行参数-bootstrap_expect；bind_addr：同命令行参数-bind；enable_syslog：同命令行参数-syslog；log_level：同命令行参数-log_level；node_name：同命令行参数node；server：是否是server节点
   ```

3. 启动consul：要使用consul必须运行agent，它可以运行为server或client模式，每个数据中心至少一个server，一般一个集群3~5个server

   ```powershell
   #启动服务端，可以多起几个，其中一个server的ip和集群的ip相同
   consul agent -server -bootstrap-expect 5 -data-dir /tmp/consul -node=s1 -bind=192.168.0.110 -ui -rejoin -client 0.0.0.0 -join 192.168.0.105
   
   #启动客户端：不用加-bootstrap、-server和-client
   consul agent -data-dir /tmp/consul -node=c1 -bind=192.168.0.116 -ui -rejoin -join 192.168.0.105
   
   consul members		#查看成员
   consul leave		#停止agent
   consul kv get <key>	#查看键值对
   ```

- ==命令行参数介绍==
    - server ：定义agent运行在server模式
    - bootstrap-expect ：在一个datacenter中期望提供的server节点数目，当该值提供的时候，consul一直等到达到指定sever数目的时候才会引导整个集群
    - data-dir=/tmp/consul：数据存放目录，windows是在c:\tmp\consul
    - node=s1 节点在集群中的名称，在一个集群中必须是唯一的，默认是该节点的主机名
    - bind：该服务器节点名,该地址用来在集群内部的通讯，集群内的所有节点到地址都必须是可达的，默认是0.0.0.0
    - -ui：能使用ui界面
    - rejoin：使consul忽略先前的离开，在再次启动后仍旧尝试加入集群中。
    - `config-dir`：配置文件目录，里面所有以.json结尾的文件都会被加载。一般都要按配置文件启动
    - client：consul服务侦听地址，这个地址提供HTTP、DNS、RPC等服务，默认是127.0.0.1所以不对外提供服务，如果你要对外提供服务改成0.0.0.0
    - join：节点要加到哪个集群中，后面跟集群的ip（如本机ip）（可以在server都启动后再统一consul join把几个节点的ip都加进来，但要先启动和集群ip一样的那个server）

- ui界面介绍(127.0.0.1:8500)

  - Services：所有的服务（以配置文件中的name区分）
  - Nodes:所有的节点，包括服务端和客户端，以及健康情况
  - Key/Value：进行Consul KV的管理，使用命令consul kv get username可以获取key
  - ACL：访问控制列表，对UI、API、CLI以及服务通信进行安全上的保证，用于生产环境
  - Intentions：可以通过此页面对Intention进行增删改查的操作

- consul members查看到的字段介绍

  | Node | Address        | Status | Type   | Build | Protocol | DC       | Segment |
  | ---- | -------------- | ------ | ------ | ----- | -------- | -------- | ------- |
  | 节点 | 地址           | 状态   | 类型   | 版本  | 协议     | 数据中心 | 分管    |
  | s1   | 127.0.0.1:8301 | alive  | server | 1.9.5 | 2        | dc1      | all     |

- ==leader职责==

  - 同步注册信息给其他server
  - 负责各节点的健康检查

4. 各个服务器上的consul启动后，继续启动各个服务器上的server.go。把client.go放到以client方式启动consul的服务器上，并在client.go中添加服务发现的代码

   > client端能获取所有服务地址，当健康检查发现一个service挂掉后，client能切换到健康的节点上去

     ~~实现方法一：每个服务器都启动一个server，把所有节点的ip放进数组中，client中遍历数组进行拨号连接~~

     ==实现方法二==：用consul提供的方法，过滤掉不健康的（还可以加其他过滤条件，如使用量等），client只连接健康的

   ```go
   //在grpc章节的client实现的基础上新增如下
   config := api.DefaultConfig()
   config.Address = "127.0.0.1:8500"
   
   var waitIndex uint64
   cliApi, _ := api.NewClient(config)
   //获取符合筛选条件的所有服务实体。
   //前三个参数分别是配置文件中的name、tag，和"是否只保留通过筛选的服务"
   services,_,_ := cliApi.Health().Service("sayHello","sayhello",true,&api.QueryOptions{
       WaitIndex: waitIndex,
   })
   //services[0]还可以获得Node、Checks的具体信息
   url := fmt.Sprintf("%v:%v",services[0].Service.Address,services[0].Service.Port)
   
   //WithInsecure：跳过证书的验证
   conn,err := grpc.Dial(url,grpc.WithInsecure())
   ```

5. 最后把consul设为自启动即可

## micro

> 一个采用了微服务体系结构模式的开源框架，隐藏了分布式系统的复杂性

### 命令

```powershell
 new命令：
   --namespace "go.micro"	Namespace for the service e.g com.example
   --type "srv"				Type of service e.g api, fnc, srv, web
   --fqdn					全限定域名
   --alias					给项目名称起别名
   --plugin [--plugin option --plugin option]
```

### 使用

1. `git clone https://github.com/micro/micro.git`，然后`go build`，把exe放到`gopath/bin`下
2. `micro new micro_project`，然后去//.
3. 

## docker部署

## 实战

# 第三方库

```bash
#生成验证码图片
go get -u github.com/mojocn/base64Captcha

#go操作kafka ES etcd
go get -u github.com/Shopify/sarama	
go get -u github.com/olivere/elastic
go get -u github.com/etcd-io/etcd/clientv3
#tail用来跟踪日志文件，每新增一行日志就往kafka写一条消息
go get -u github.com/hpcloud/tail

#操作表格；解析ini配置文件
go get -u github.com/360EntSecGroup-Skylar/excelize
go get -u github.com/go-ini/ini

#goland的三个工具，goimports常用
go get -u golang.org/x/tools/cmd/goimports
go get -u golang.org/x/lint/golint
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint


#论坛项目
#zap日志，luberjac分割日志文件，，
go get -u go.uber.org/zap
go get -u github.com/natefinch/lumberjack
#viper读取配置
go get -u github.com/spf13/viper

#validator校验,locales和universal-translator把校验结果翻译成中文
go get -u github.com/go-playground/validator/v10
go get -u github.com/go-playground/locales/zh
go get -u github.com/go-playground/locales/en
go get -u github.com/go-playground/universal-translator

#jwt（json web token），代替session的用户认证和鉴权，详见web基础.md
go get -u github.com/dgrijalva/jwt-go
#雪花算法
go get -u github.com/bwmarrin/snowflake

#swag生成接口文档
go get -u github.com/swaggo/swag/cmd/swag

#支持跨域请求
go get -u github.com/gin-contrib/cors

#单元测试
go get -u github.com/stretchr/testify/assert
```

