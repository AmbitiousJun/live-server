<h1 align="center">live-server</h1>

<div align="center">
  <a href="https://github.com/AmbitiousJun/live-server/tree/v1.2.0"><img src="https://img.shields.io/github/v/tag/AmbitiousJun/live-server"></img></a>
  <a href="https://hub.docker.com/r/ambitiousjun/live-server/tags"><img src="https://img.shields.io/docker/image-size/ambitiousjun/live-server/v1.2.0"></img></a>
  <a href="https://hub.docker.com/r/ambitiousjun/live-server/tags"><img src="https://img.shields.io/docker/pulls/ambitiousjun/live-server"></img></a>
  <a href="https://goreportcard.com/report/github.com/AmbitiousJun/live-server"><img src="https://goreportcard.com/badge/github.com/AmbitiousJun/live-server"></img></a>
  <img src="https://img.shields.io/github/license/AmbitiousJun/live-server"></img>
</div>

<div align="center">
  一个 HTTP 电视直播服务
</div>

## DockerCompose 安装

docker-compose.yml:

```yaml
version: '3.1'
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