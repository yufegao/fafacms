# FaFa CMS

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/fafacms.svg?style=social&label=Forks)](https://github.com/hunterhug/fafacms/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/fafacms.svg?style=social&label=Stars)](https://github.com/hunterhug/fafacms/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/fafacms)](https://goreportcard.com/report/github.com/hunterhug/fafacms)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/fafacms.svg)](https://github.com/hunterhug/fafacms/issues)

Developing--[中文](/README_CN.md)

## Project description

 `fafa` -- means `flower` in Cantonese.

A content management system in go, frontend and backend is highly splited. Support multi-users, post blogs, view blogs. We want to bring out an example 

Backend returns JSON API. Feel free to use any mainstream frameworks to develop frontend. This project framework is scalable.

Dependencies:

1. [Gin is a HTTP web framework written in Go (Golang)](https://github.com/gin-gonic/gin)
2. [The open source embeddable online markdown editor (component).](https://github.com/pandao/editor.md)
3. [Session management for Go 1.7+](https://github.com/alexedwards/scs)
4. ...

`Gin` is mainly used for server functions.

Structure:

```
├── config.json 
├── core    	# backend files
│   ├── config      
│   ├── flog        
│   ├── controllers 
│   ├── model       
│   ├── router     
│   ├── server      
│   └── util        
├── main.go 	# entrance
└── web  		# frontend files
```

How to run and debug backend: [/doc/README.md](/doc/README.md), use [insomnia](https://insomnia.rest/)

## Instruction

### Backend deployment(normal)

get codes:

```
go get -v github.com/hunterhug/fafacms
```



run:

```
fafacms -config=/root/config.json
```

description of`config.json`:

```
{
  "DefaultConfig": {
    "WebPort": ":8080", 				    	# port for project(changeable)
    "StoragePath": "/root/data/data",  			# Path for file saving(changeable)
    "LogPath": "/root/data/log/fafacms.log", 	# flog saving path(changeable)
    "Debug": true   					        # Debug(default)
  },
  "DbConfig": {
    "DriverName": "mysql",  			# Relational DB driver(default)
    "Name": "fafa", 					# DB name(changeable)
    "Host": "127.0.0.1", 				# DB host(changeable)
    "User": "root", 					# DB user(changeable)
    "Pass": "123456789", 				# DB password(changeable)
    "Port": "3306", 					# DB port(changeable)
    "Prefix": "fafa_", 					# DB prefix(changeable)
    "MaxIdleConns": 20, 				# Max Idle connections(default)
    "MaxOpenConns": 20, 				# Max Idle connections(default)
    "DebugToFile": true, 				# Debug output files(default)
    "DebugToFileName": "/root/data/log/fafacms_db.log", # SQL output file path(default)
    "Debug": true 										# sql Debug(default)
  },
  "SessionConfig": {
    "RedisHost": "127.0.0.1:6379", 						# RedisHost(changeable)
    "RedisMaxIdle": 64, 								# (default)
    "RedisMaxActive": 0, 								# (default)
    "RedisIdleTimeout": 120, 							# (default)
    "RedisDB": 0, 										# Redis connect database(default)
    "RedisPass": "123456789"   							# Redis password(optional, changeable)
  }
}
```

The project use`Mysql`,`Redis` and local storage, to deploy database, please refer : [Docker easy use to run  Mysql/Redis](https://github.com/hunterhug/GoSpider-docker).

```
git clone https://github.com/hunterhug/GoSpider-docker
chomd 777 build.sh
./build

sudo docker exec -it  GoSpider-mysqldb mysql -uroot -p123456789

> create database fafa default character set utf8mb4 collate utf8mb4_unicode_ci;

sudo docker exec -it GoSpider-redis redis-cli -a 123456789

> KEYS *
```

### Backend deployment(Docker)

We can also use `docker` to deploy, construct the mirror(Docker version must later than 17.06):

```
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh
````

Build data volume and config:

```
mkdir $HOME/data
cp docker_config.json $HOME/data/config.json
```

Initialize container:

```
sudo docker run -d --name fafacms -p 8080:8080 -v $HOME/data:/root/data hunterhug/fafacms fafacms -config=/root/data/config.json

sudo docker logs -f --tail 10 fafacms
```

`/root/data` is durable volume, please put `config.json` under the folder.

### Frontend deployment(Developing)

## Sponser

wechat:

![](/support/weixin.jpg)

alipay:

![](/support/alipay.png)

