# 作业
把我们的 httpserver 服务以 istio ingress gateway 的形式发布出来。
- 如何实现安全保证
- 七层路由规则
- 考虑 open tracing 的接入

# 实现
## 安装 istio
1. 安装 istioctl 命令
https://istio.io/latest/docs/setup/getting-started/#download
```
$ curl -L https://istio.io/downloadIstio | sh -
$ cd istio-1.13.2
$ sudo cp bin/istioctl /usr/bin/local
```
2. 安装 istio with jaegor
https://istio.io/latest/docs/ops/integrations/jaeger/#installation

```
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/addons/jaeger.yaml
istioctl install -f ./tracing.yaml
```
## create and enable istio-injection for namespace
```
kubectl apply -f namespace.yaml
```
## deploy gohttpserver application
```
kubectl apply -f deployment.yaml -n gohttpserver-istio
kubectl apply -f service.yaml -n gohttpserver-istio
```
## build gateway and virtualservice
参考

https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/

https://istio.io/latest/docs/tasks/traffic-management/ingress/secure-ingress/#generate-client-and-server-certificates-and-keys

### Determine the ingress IP and ports
```
kubectl get svc istio-ingressgateway -n istio-system
export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')
export TCP_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="tcp")].port}')
```
### generate client and server certificates and keys
1. create a root certificate and priviate key to sign the certificates for your service.
```
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cloudnative Inc./CN=istio.cloudnative-learn.com' -keyout istio.cloudnative-learn.com.key -out istio.cloudnative-learn.com.crt
```
2. create a certificate and a private key for gohttpserver.istio.cloudnative-learn.com
```
openssl req -out gohttpserver.istio.cloudnative-learn.com.csr -newkey rsa:2048 -nodes -keyout gohttpserver.istio.cloudnative-learn.com.key -subj "/CN=gohttpserver.istio.cloudnative-learn.com/O=gohttpserver.istio organization"

openssl x509 -req -sha256 -days 365 -CA istio.cloudnative-learn.com.crt -CAkey istio.cloudnative-learn.com.key -set_serial 0 -in gohttpserver.istio.cloudnative-learn.com.csr -out gohttpserver.istio.cloudnative-learn.com.crt
```
### create a secret for the ingress gateway
```
kubectl create -n istio-system secret tls gohttpserver-credential --key=gohttpserver.istio.cloudnative-learn.com.key --cert=gohttpserver.istio.cloudnative-learn.com.crt
```
### create gateway and virtualservice with TLS
istio gateway
```
kubectl apply -f gateway.yaml -n gohttpserver-istio
```

istio virtualservice
```
kubectl apply -f virtualservice.yaml -n gohttpserver-istio
```

## 测试
非 tls 的 
```
curl -s -I -HHost:gohttpserveristio.cloudnative-learn.com "http://$INGRESS_HOST:$INGRESS_PORT/metrics"
```
tls 的
```
# 命令行
curl -v -HHost:gohttpserver.istio.cloudnative-learn.com --resolve "gohttpserver.istio.cloudnative-learn.com:$SECURE_INGRESS_PORT:$INGRESS_HOST" \
--cacert gohttpserver.istio.cloudnative-learn.com.crt "https://gohttpserver.istio.cloudnative-learn.com:$SECURE_INGRESS_PORT/metrics"

# 或浏览器
https://gohttpserver.istio.cloudnative-learn.com/metrics
```

output
```
$ curl -v -HHost:gohttpserver.istio.cloudnative-learn.com --resolve "gohttpserver.istio.cloudnative-learn.com:$SECURE_INGRESS_PORT:$INGRESS_HOST" \
> --cacert gohttpserver.istio.cloudnative-learn.com.crt "https://gohttpserver.istio.cloudnative-learn.com:$SECURE_INGRESS_PORT/metrics"
* Added gohttpserver.istio.cloudnative-learn.com:443:20.63.208.219 to DNS cache
* Hostname gohttpserver.istio.cloudnative-learn.com was found in DNS cache
*   Trying 20.63.208.219:443...
* TCP_NODELAY set
* Connected to gohttpserver.istio.cloudnative-learn.com (20.63.208.219) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: gohttpserver.istio.cloudnative-learn.com.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: CN=gohttpserver.istio.cloudnative-learn.com; O=gohttpserver.istio organization
*  start date: Mar 21 05:16:24 2022 GMT
*  expire date: Mar 21 05:16:24 2023 GMT
*  common name: gohttpserver.istio.cloudnative-learn.com (matched)
*  issuer: O=cloudnative Inc.; CN=istio.cloudnative-learn.com
*  SSL certificate verify ok.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* Using Stream ID: 1 (easy handle 0x7fffdf11a840)
> GET /metrics HTTP/2
> Host:gohttpserver.istio.cloudnative-learn.com
> user-agent: curl/7.68.0
> accept: */*
>
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* old SSL session ID is stale, removing
* Connection state changed (MAX_CONCURRENT_STREAMS == 2147483647)!
< HTTP/2 200
< content-type: text/plain; version=0.0.4; charset=utf-8
< date: Mon, 21 Mar 2022 05:19:02 GMT
< x-envoy-upstream-service-time: 7
< server: istio-envoy
<
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
go_goroutines 11
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.16.14"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 3.55176e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 3.55176e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.44594e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 952
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 4.122512e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 3.55176e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.2283776e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.333568e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 10519
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.2283776e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6617344e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 11471
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 2400
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 45288
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 49152
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.473924e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 609604
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 491520
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 491520
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.3352456e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 7
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.05
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.4880768e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.64783971219e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.29346048e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 7
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
* Connection #0 to host gohttpserver.istio.cloudnative-learn.com left intact
```
![image](https://github.com/Yejing0609/CloudNativeLearning/blob/main/Homework/K8S/pics/istio-web-tls.PNG)

## Jeagor
```
$ istioctl dashboard jaeger
http://localhost:16686
```
![image](https://github.com/Yejing0609/CloudNativeLearning/blob/main/Homework/K8S/pics/jaegor.PNG)
