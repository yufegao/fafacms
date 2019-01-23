# FaFa Blog

开发中...

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/fafa.svg?style=social&label=Forks)](https://github.com/hunterhug/fafa/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/fafa.svg?style=social&label=Stars)](https://github.com/hunterhug/fafa/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/fafa.svg)](https://github.com/hunterhug/fafa)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/fafa)](https://goreportcard.com/report/github.com/hunterhug/fafa)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/fafa.svg)](https://github.com/hunterhug/fafa/issues)


## 项目说明

使用`Golang`开发的多用户简易博客系统,依赖以下项目:

1. [Gin is a HTTP web framework written in Go (Golang)](https://github.com/gin-gonic/gin)
2. [The open source embeddable online markdown editor (component).](https://github.com/pandao/editor.md)
3. [Session management for Go 1.7+](https://github.com/alexedwards/scs)

项目结构:

```
├── config.json 配置文件
├── core    后端代码
│   ├── config      配置
│   ├── controllers 控制器
│   ├── model       模型
│   ├── router      路由
│   ├── server      服务
│   └── util        工具
├── main.go 程序入口
└── web 前端代码
```

本项目依赖`Mysql`,`Redis`和本地存储, 快速部署数据库环境请参考: [Docker easy use to run  Mysql/Redis](https://github.com/hunterhug/GoSpider-docker).

```
git clone https://github.com/hunterhug/GoSpider-docker
chomd 777 build.sh
./build

sudo docker exec -it  GoSpider-mysqldb mysql -uroot -p123456789

> create database blog default character set utf8mb4 collate utf8mb4_unicode_ci;

sudo docker exec -it GoSpider-redis redis-cli -a 123456789

> KEYS *
```

## 如何使用

获取代码:

```
go get -v github.com/hunterhug/fafa
```

运行如下:

```
fafa -config=/root/config.json
```

其中`config.json`说明如下:

```
{
  "DefaultConfig": {
    "WebPort": ":8080", # 程序运行端口(可改)
    "StoragePath": "/home/hunterhug/data",  # 本地文件保存地址(可改)
    "LogPath": "/home/hunterhug/data/log/fafa.log", # 日志保存地址(可改)
    "Debug": true   # 打开调试(默认保持)
  },
  "DbConfig": {
    "DriverName": "mysql",  # 关系型数据库驱动(默认保持)
    "Name": "blog", # 关系型数据库名字(可改)
    "Host": "127.0.0.1", # 关系型数据库地址(可改)
    "User": "root", # 关系型数据库用户(可改)
    "Pass": "123456789", # 关系型数据库密码(可改)
    "Port": "3306", # 关系型数据库端口(可改)
    "Prefix": "fafa_", # 关系型数据库表前缀(可改)
    "MaxIdleConns": 20, # 关系型数据库池闲置连接数(默认保持)
    "MaxOpenConns": 20, # 关系型数据库池打开连接数(默认保持)
    "DebugToFile": true, # SQL调试是否输出到文件(默认保持)
    "DebugToFileName": "/home/hunterhug/data/log/fafa_db.log", # SQL调试输出文件路径(默认保持)
    "Debug": true # SQL调试(默认保持)
  },
  "SessionConfig": {
    "RedisHost": "127.0.0.1:6379", # Redis地址(可改)
    "RedisMaxIdle": 64, # (默认保持)
    "RedisMaxActive": 0, # (默认保持)
    "RedisIdleTimeout": 120, # (默认保持)
    "RedisDB": 0, # Redis默认连接数据库(默认保持)
    "RedisPass": "123456789"   # Redis密码(可为空,可改)
  }
}
```

你也可以使用`docker`进行部署:

构建镜像(Docker版本必须大于17.06):

```
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh
````

新建数据卷, 并且移动配置并修改:

```
mkdir $HOME/data
cp docker_config.json $HOME/data/config.json
```

启动容器:

```
sudo docker run -d --name fafablog -p 8080:8080 -v $HOME/data:/root/data hunterhug/fafa fafa -config=/root/data/config.json

sudo docker logs -f --tail 10 fafablog
```

其中`/root/data`是挂载的持久化卷, 配置`config.json`放置在该文件夹下.

## 支持

微信支持:

![](/support/weixin.jpg)

支付宝支持:

![](/support/alipay.png)


## 界面

待开发