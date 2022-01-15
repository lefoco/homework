
#编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化
sudo docker build -t lefoco/httpserver:${tag} .

#将镜像推送至 docker 官方镜像仓库
sudo docker push lefoco/httpserver:${tag}

#通过 docker 命令本地启动 httpserver(pull from 官方镜像仓库)
sudo docker run --name lefoco-httpserver -p 80:80 lefoco/httpserver

#通过 nsenter 进入容器查看 IP 配置
cadmin@k8snode:~/httpserver$ sudo docker ps
CONTAINER ID   IMAGE                                               COMMAND                  CREATED          STATUS          PORTS                                   NAMES
e448c57a0a17   lefoco/httpserver:v1.0                              "/bin/sh -c /httpser…"   18 minutes ago   Up 18 minutes   0.0.0.0:80->80/tcp, :::80->80/tcp       lefoco-httpserver

cadmin@k8snode:~/httpserver$ sudo docker inspect e448c57a0a17 |grep -i pid
"Pid": 687685,
"PidMode": "",
"PidsLimit": null,

#nsenter查看ip
cadmin@k8snode:~/httpserver$ sudo nsenter -t 687685 -n ip a

1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
inet 127.0.0.1/8 scope host lo
valid_lft forever preferred_lft forever
67: eth0@if68: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
link/ether 02:42:ac:11:00:05 brd ff:ff:ff:ff:ff:ff link-netnsid 0
inet 172.17.0.5/16 brd 172.17.255.255 scope global eth0
valid_lft forever preferred_lft forever

#访问容器ip
cadmin@k8snode:~/httpserver$ curl 172.17.0.5
hello httpserver...

#访问主机ip
cadmin@k8snode:~/httpserver$ curl 192.168.34.3
hello httpserver...

#访问主机ip/healthz
cadmin@k8snode:~/httpserver$ curl 192.168.34.3/healthz
