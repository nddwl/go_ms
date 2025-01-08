# 微服务

## 1.微服务框架图

![API架构图](https://github.com/nddwl/zhihu/blob/master/png/apiGateway.png)

## 2.项目配置

~~~ssss
tree -L 2 ./
./
├── README.md
├── application
│   ├── applet
│   ├── article
│   ├── chat
│   ├── concerned
│   ├── member
│   ├── message
│   ├── qa
│   └── user
├── db
│   └── user.sql
├── go.mod
├── go.sum
└── pkg
    ├── encrypt
    ├── jwt
    └── util
~~~

### 1.生成api项目

~~~shell
cd ~/beyond/application/applet

goctl api go --dir=./ --api applet.api
~~~

### 2.生成user项目

~~~shell
cd ~/beyond/application/user/rpc

goctl rpc protoc ./user.proto --go_out=. --go-grpc_out=. --zrpc_out=./
~~~

### 3.生成model

~~~shell
cd ~/beyond/application/user/rpc

goctl model mysql datasource --dir ./internal/model --table user --cache true --url "root:Zsg123456@tcp(127.0.0.1:3306)/beyond_user"
~~~



## 3.服务注册与发现

![服务注册与发现](https://github.com/nddwl/zhihu/blob/master/png/%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0.png)

