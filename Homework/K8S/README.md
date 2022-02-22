# 第一部分
作业要求：编写 Kubernetes 部署脚本将 httpserver 部署到 Kubernetes 集群，以下是你可以思考的维度。

- 优雅启动
- 优雅终止
- 资源需求和 QoS 保证
- 探活
- 日常运维需求，日志等级
- 配置和代码分离

# 第二部分
除了将 httpServer 应用优雅的运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。

来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。

在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：
- Service
- Ingress

可以考虑的细节
- 如何确保整个应用的高可用。
- 如何通过证书保证 httpServer 的通讯安全。