## 模块三：Docker 核心技术
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
# docker build -t gohttpserver:v2 .
```

[Docker hub 镜像](https://hub.docker.com/repository/docker/yejing0609/gohttpserver)

docker 命令本地启动 httpserver
```
# docker run -d yejing0609/gohttpserver:v2
```

通过 nsenter 进入容器查看 IP 配置
```
# PID=$(docker inspect --format "{{.State.Pid}}" <container>)
# nsenter --target $PID --net ip a
```

## 模块十
### 作业内容
1. 为 HTTPServer 添加 0-2 秒的随机延时；
1. 为 HTTPServer 项目添加延时 Metric；
1. 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；
1. 从 Promethus 界面中查询延时指标数据；

### 修改
- 在 Go 的代码里添加 prometheus metrics 相关代码。[main.go](httpServer/main.go), [metrics.go](httpServer/metrics/metrics.go)
- 创建新的 docker image: [yejing0609/gohttpserver:v5](https://hub.docker.com/repository/docker/yejing0609/gohttpserver)
- 修改 [service.yaml](../K8S/Files/service.yaml) 的 annotaiton、image、port 使其能够被 prometheus 使用
- 部署到测试集群
- 安装 promethues，配置 prometheus。create [promethues-additional.yaml](../K8S/Files/prometheus/prometheus-additional.yaml) as secrets, [rbac.yaml](../K8S/Files/rbac.yaml)
- 安装 grafana
- 从 Promethues 和 grafana 界面中获取结果

### 测试结果
```
root@master-node:~# kubectl get pods -n httpserver -o wide
NAME                         READY   STATUS    RESTARTS   AGE   IP                NODE          NOMINATED NODE   READINESS GATES
httpserver-f66d596cf-5pm82   1/1     Running   0          25s   192.168.168.177   worker-node   <none>           <none>
httpserver-f66d596cf-vfj9f   1/1     Running   0          25s   192.168.168.176   worker-node   <none>           <none>
root@master-node:~# curl 192.168.168.177
root@master-node:~# curl 192.168.168.177
root@master-node:~# curl 192.168.168.177
root@master-node:~# curl 192.168.168.177
```

```
root@master-node:~# curl 192.168.168.177/metrics
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.16.14"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 791288
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 791288
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 4172
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 203
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 4.017656e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 791288
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.5347584e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.368064e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 3858
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.5347584e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6715648e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 4061
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 1200
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 19176
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 32768
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.473924e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 470476
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 393216
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 393216
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.165032e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 7
# HELP httpserver_execution_latency_seconds Time spent.
# TYPE httpserver_execution_latency_seconds histogram
httpserver_execution_latency_seconds_bucket{step="total",le="0.001"} 0
httpserver_execution_latency_seconds_bucket{step="total",le="0.002"} 0
httpserver_execution_latency_seconds_bucket{step="total",le="0.004"} 0
httpserver_execution_latency_seconds_bucket{step="total",le="0.008"} 0
httpserver_execution_latency_seconds_bucket{step="total",le="0.016"} 0
httpserver_execution_latency_seconds_bucket{step="total",le="0.032"} 0
httpserver_execution_latency_seconds_bucket{step="total",le="0.064"} 1
httpserver_execution_latency_seconds_bucket{step="total",le="0.128"} 2
httpserver_execution_latency_seconds_bucket{step="total",le="0.256"} 2
httpserver_execution_latency_seconds_bucket{step="total",le="0.512"} 2
httpserver_execution_latency_seconds_bucket{step="total",le="1.024"} 2
httpserver_execution_latency_seconds_bucket{step="total",le="2.048"} 4
httpserver_execution_latency_seconds_bucket{step="total",le="4.096"} 4
httpserver_execution_latency_seconds_bucket{step="total",le="8.192"} 4
httpserver_execution_latency_seconds_bucket{step="total",le="16.384"} 4
httpserver_execution_latency_seconds_bucket{step="total",le="+Inf"} 4
httpserver_execution_latency_seconds_sum{step="total"} 3.876411071
httpserver_execution_latency_seconds_count{step="total"} 4
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.02
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 9
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 7.278592e+06
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.64687917156e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.27117824e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0

```


参考

https://cloud.tencent.com/document/product/1416/56033

https://prometheus.io/docs/guides/go-application/

