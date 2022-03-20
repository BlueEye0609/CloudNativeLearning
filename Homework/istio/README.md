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
2. 安装 istio
```
$ istioctl install  --set profile=demo -y
```
## create and enable istio-injection for namespace
```
$ kubectl apply -f namespace.yaml
```
## deploy gohttpserver application
```
$ kubectl apply -f deployment.yaml -n gohttpserver-istio
$ kubectl apply -f service.yaml -n gohttpserver-istio
```
## build gateway and virtualservice
参考

https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/

https://istio.io/latest/docs/tasks/traffic-management/ingress/secure-ingress/#generate-client-and-server-certificates-and-keys

### Determine the ingress IP and ports
```
$ kubectl get svc istio-ingressgateway -n istio-system

$ export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

$ export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')

$ export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')

$ export TCP_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="tcp")].port}')

```
### generate client and server certificates and keys
1. create a root certificate and priviate key to sign the certificates for your service.
```
$ openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cloudnative Inc./CN=cloudnative-learn.com' -keyout cloudnative-learn.com.key -out cloudnative-learn.com.crt
```
2. create a certificate and a private key for gohttpserveristio.cloudnative-learn.com
```
$ openssl req -out gohttpserveristio.cloudnative-learn.com.csr -newkey rsa:2048 -nodes -keyout gohttpserveristio.cloudnative-learn.com.key -subj "/CN=gohttpserveristio.cloudnative-learn.com/O=gohttpserveristio organization"
$ openssl x509 -req -sha256 -days 365 -CA cloudnative-learn.com.crt -CAkey cloudnative-learn.com.key -set_serial 0 -in gohttpserveristio.cloudnative-learn.com.csr -out gohttpserveristio.cloudnative-learn.com.crt
```
### create a secret for the ingress gateway
```
kubectl create -n istio-system secret tls gohttpserver-credential --key=gohttpserveristio.cloudnative-learn.com.key --cert=gohttpserveristio.cloudnative-learn.com.crt
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
curl -v -HHost:gohttpserveristio.cloudnative-learn.com --resolve "gohttpserveristio.cloudnative-learn.com:$SECURE_INGRESS_PORT:$INGRESS_HOST" \
--cacert gohttpserveristio.cloudnative-learn.com.crt "https://gohttpserveristio.cloudnative-learn.com:$SECURE_INGRESS_PORT/metrics"
```
