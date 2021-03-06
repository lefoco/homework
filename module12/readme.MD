把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：
如何实现安全保证；
七层路由规则；
考虑 open tracing 的接入。

## istio 配置
```
kubectl create ns module12

kubectl create -f myhttpserver.yaml -n module12

kubectl create -f nginx.yaml -n module12

openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt

kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt

kubectl create -f istio-myhttpserver-specs.yaml -n module12

curl --resolve httpsserver.cncamp.io:443:10.121.3.32 https://httpsserver.cncamp.io/hello -v -k
```
