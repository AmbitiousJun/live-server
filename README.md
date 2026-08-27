<div align="center">
  <img height="150px" src="./assets/logo.png"></img>
</div>

<h1 align="center">live-server</h1>

<div align="center">
  <a href="https://github.com/AmbitiousJun/live-server/tree/v1.23.0"><img src="https://img.shields.io/github/v/tag/AmbitiousJun/live-server"></img></a>
  <a href="https://hub.docker.com/r/ambitiousjun/live-server/tags"><img src="https://img.shields.io/docker/image-size/ambitiousjun/live-server/v1.23.0"></img></a>
  <a href="https://hub.docker.com/r/ambitiousjun/live-server/tags"><img src="https://img.shields.io/docker/pulls/ambitiousjun/live-server"></img></a>
  <a href="https://goreportcard.com/report/github.com/AmbitiousJun/live-server"><img src="https://goreportcard.com/badge/github.com/AmbitiousJun/live-server"></img></a>
  <a href="https://github.com/AmbitiousJun/live-server/releases/latest"><img src="https://img.shields.io/github/downloads/AmbitiousJun/live-server/total"></img></a>
  <img src="https://img.shields.io/github/license/AmbitiousJun/live-server"></img>
</div>

<div align="center">
  一个 HTTP 电视直播服务
</div>

## DockerCompose 安装

docker-compose.yml:

```yaml
version: "3.1"
services:
  live-server:
    image: ambitiousjun/live-server:latest
    environment:
      - TZ=Asia/Shanghai
      - GIN_MODE=release
    container_name: live-server
    volumes:
      - ./data:/app/data
    ports:
      - 5666:5666
```

运行:

```shell
docker-compose up -d
```

查看运行日志:

```shell
docker logs -f live-server -n 1000
```

## 直接运行二进制

从 release 页下载对应的二进制文件

终端运行：

```shell
live-server -prod=true
```

程序默认运行在 `5666` 端口上，自定义端口：

```shell
live-server -prod=true -p 8880
```

开启守护进程，运行在后台：

```shell
nohup live-server -prod=true -p 8880 > "./app.log" 2>&1 &
```

停止 `8880` 端口进程：

```shell
lsof -ti:8880 | xargs -r kill -9
```

## 请我喝杯 9.9💰 的 Luckin Coffee☕️

<img height="500px" src="assets/2024-11-05-09-59-45.png" />
