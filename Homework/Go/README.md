## 模块三：Docker 核心计划
### 作业内容
- 构建本地镜像
- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化
- 将镜像推送至 docker 官方镜像仓库
- 通过 docker 命令本地启动 httpserver
- 通过 nsenter 进入容器查看 IP 配置

### 作业查看
[Dockerfile](Dockerfile)

构建本地镜像
```
# docker build -t gohttpserver:v1 .
```

[Docker hub 镜像](https://hub.docker.com/repository/docker/yejing0609/gohttpserver)

docker 命令本地启动 httpserver
```
# docker run -d yejing0609/gohttpserver:v1
```

通过 nsenter 进入容器查看 IP 配置
```
# PID=$(docker inspect --format "{{.State.Pid}}" <container>)
# nsenter --target $PID --net ip a
```