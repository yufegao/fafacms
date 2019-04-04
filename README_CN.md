# 花花CMS(FaFa CMS)

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/fafacms.svg?style=social&label=Forks)](https://github.com/hunterhug/fafacms/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/fafacms.svg?style=social&label=Stars)](https://github.com/hunterhug/fafacms/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/fafacms)](https://goreportcard.com/report/github.com/hunterhug/fafacms)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms/issues)

开发中... 稳定分支为master. 

## 项目说明

此项目代号为 `fafacms`。花花拼音 `fafa`，名称来源于广东话发发，花花的谐音，听起来有诙谐，娱乐等感觉。

是一个使用 `Golang` 开发的前后端分离 --> 内容管理系统(CMS)，支持不同用户登录，并且可以发布文章，在首页可以看到不同用户的文章。本项目愿景是希望使用较简洁较规范的企业开发模式，可以给初学者一个示例。你可以把它当作弱化的博客园或简书。

后端主要返回JSON API，前端可用任意框架开发界面。 该项目架构可扩展为不同业务的Web项目。

依赖以下项目:

1. [Gin is a HTTP web framework written in Go (Golang)](https://github.com/gin-gonic/gin)
2. [The open source embeddable online markdown editor (component).](https://github.com/pandao/editor.md)
3. [Session management for Go 1.7+](https://github.com/alexedwards/scs)
4. ...

使用`Gin`是因为框架比较简洁， 主要使用到路由功能。使用 `Markdown` 编辑器是因为程序员朋友和很多媒体工作者也开始尝鲜...

代码结构:

```
├── config.json 配置文件
├── core    后端代码
│   ├── config      配置
│   ├── flog        日志
│   ├── controllers 控制器
│   ├── model       模型
│   ├── router      路由
│   ├── server      服务
│   └── util        工具
├── main.go 程序入口
└── web 前端代码(可用vue等开发)
```

后端API调试见: [/doc/README.md](/doc/README.md), 使用[insomnia](https://insomnia.rest/)

## 如何使用

### 后端部署(常规)

获取代码:

```
go get -v github.com/hunterhug/fafacms
```

代码就会保存在`Golang GOPATH`目录下.

运行:

```
fafacms -config=./config.json
```

其中`config.json`说明如下:

```
{
  "DefaultConfig": {
    "WebPort": ":8080", 				    # 程序运行端口(可改)
    "StoragePath": "./data/storage",  		# 本地文件保存地址(可改)
    "LogPath": "./data/log/fafacms_log.log", 	# 日志保存地址(可改)
    "Debug": true   					        # 打开调试(默认保持)
  },
  "DbConfig": {
    "DriverName": "mysql",  			# 关系型数据库驱动(默认保持)
    "Name": "fafa", 					# 关系型数据库名字(可改)
    "Host": "127.0.0.1", 				# 关系型数据库地址(可改)
    "User": "root", 					# 关系型数据库用户(可改)
    "Pass": "123456789", 				# 关系型数据库密码(可改)
    "Port": "3306", 					# 关系型数据库端口(可改)
    "Prefix": "fafa_", 					# 关系型数据库表前缀(可改)
    "MaxIdleConns": 20, 				# 关系型数据库池闲置连接数(默认保持)
    "MaxOpenConns": 20, 				# 关系型数据库池打开连接数(默认保持)
    "DebugToFile": true, 				# SQL调试是否输出到文件(默认保持)
    "DebugToFileName": "./data/log/fafacms_db.log", # SQL调试输出文件路径(默认保持)
    "Debug": true 					# SQL调试(默认保持)
  },
  "SessionConfig": {
    "RedisHost": "127.0.0.1:6379", 			# Redis地址(可改)
    "RedisMaxIdle": 64, 				# (默认保持)
    "RedisMaxActive": 0, 				# (默认保持)
    "RedisIdleTimeout": 120, 				# (默认保持)
    "RedisDB": 0, 					# Redis默认连接数据库(默认保持)
    "RedisPass": "123456789"   				# Redis密码(可为空,可改)
  }
}
```

本项目依赖`Mysql`,`Redis`和本地存储, 快速部署数据库环境请参考: [Docker easy use to run  Mysql/Redis](https://github.com/hunterhug/GoSpider-docker).

```
git clone https://github.com/hunterhug/GoSpider-docker
cd GoSpider-docker
chomd 777 build.sh
./build

sudo docker exec -it  GoSpider-mysqldb mysql -uroot -p123456789

> create database fafa default character set utf8mb4 collate utf8mb4_unicode_ci;

sudo docker exec -it GoSpider-redis redis-cli -a 123456789

> KEYS *
```

### 后端部署(Docker)

你也可以使用`docker`进行部署, 构建镜像(Docker版本必须大于17.06):

```
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh
````

先新建数据卷, 并且移动配置并修改:

```
mkdir $HOME/fafacms
cp docker_config.json $HOME/fafacms/config.json
```

启动容器:

```
sudo docker run -d --name fafacms --net=host -v $HOME/fafacms:/root/fafacms --env RUN_OPTS="-config=/root/fafacms/config.json" hunterhug/fafacms

sudo docker logs -f --tail 10 fafacms
```

其中`$HOME/fafacms`是挂载的持久化卷, 配置`config.json`放置在该文件夹下.

### 前端部署(常规)

路由:

```
var (
	HomeRouter = map[string]HttpHandle{
		"/":       {controllers.Home, GP},
		"/login":  {controllers.Login, GP},
		"/logout": {controllers.Logout, GP},
	}

	// /v1/user/create
	// need login group auth
	V1Router = map[string]HttpHandle{
		"/user/create": {controllers.CreateUser, POST},
		"/user/update": {controllers.UpdateUser, POST},
		"/user/delete": {controllers.DeleteUser, POST},
		"/user/take":   {controllers.TakeUser, GP},
		"/user/list":   {controllers.ListUser, GP},

		"/group/create": {controllers.CreateGroup, POST},
		"/group/update": {controllers.UpdateGroup, POST},
		"/group/delete": {controllers.DeleteGroup, POST},
		"/group/take":   {controllers.TakeGroup, GP},
		"/group/list":   {controllers.ListGroup, GP},

		"/resource/create": {controllers.CreateResource, POST},
		"/resource/update": {controllers.UpdateResource, POST},
		"/resource/delete": {controllers.DeleteResource, POST},
		"/resource/take":   {controllers.TakeResource, GP},
		"/resource/list":   {controllers.ListResource, GP},

		"/auth/update": {controllers.UpdateAuth, GP},

		"/node/create": {controllers.CreateNode, POST},
		"/node/update": {controllers.UpdateNode, POST},
		"/node/delete": {controllers.DeleteNode, POST},
		"/node/take":   {controllers.TakeNode, GP},
		"/node/list":   {controllers.ListNode, GP},

		"/content/create": {controllers.CreateContent, POST},
		"/content/update": {controllers.UpdateContent, POST},
		"/content/delete": {controllers.DeleteContent, POST},
		"/content/take":   {controllers.TakeContent, GP},
		"/content/list":   {controllers.ListContent, GP},

		"/comment/create": {controllers.CreateComment, POST},
		"/comment/update": {controllers.UpdateComment, POST},
		"/comment/delete": {controllers.DeleteComment, POST},
		"/comment/take":   {controllers.TakeComment, GP},
		"/comment/list":   {controllers.ListComment, GP},
	}

	// /b/upload
	// need login group auth
	BaseRouter = map[string]HttpHandle{
		"/upload": {controllers.Upload, POST},
	}
)
```

## 支持

微信支持:

![](/support/weixin.jpg)

支付宝支持:

![](/support/alipay.png)


## 界面

待开发