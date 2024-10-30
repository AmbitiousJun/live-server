# live-server

一个 http 电视直播服务

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