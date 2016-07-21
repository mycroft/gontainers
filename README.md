Gontainers: POC on how to build a simplified docker version
===========================================================

Create all mandatory stuff (bridge, network endpoints, overlay fs to launch a process isolated from host environment.
It will save container state into /tmp/rootfs.tar after the process dies.

Initialization
--------------

```
sudo iptables -t nat -A POSTROUTING -s 172.19.80.0/24 -j MASQUERADE
```
Where 172.19.80.0/24 is a custom network created for the container.

```
go build
mkdir -p fs/orig fs/rootfs
cd fs/orig && wget https://github.com/jpetazzo/docker-busybox/blob/master/tarmaker-buildroot/rootfs.tar && sudo tar xvf rootfs.tar
```


Usage
-----

```
sudo ./gontainers start
sudo ./gontainers run /bin/bash
```

Result
------

```
# ./gontainers run /bin/bash
Pid: 6475
Pid: 1
Adding route {Ifindex: 15 Dst: <nil> Src: <nil> Gw: 172.19.80.1 Flags: [] Table: 0}
/ # /bin/ps auxw
PID   USER     COMMAND
    1 root     /proc/self/exe child /bin/bash
    4 root     /bin/bash
    6 root     /bin/ps auxw
/ # /sbin/ip link show
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN mode DEFAULT group default qlen 1
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
15: myveth1@if16: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP mode DEFAULT group default 
    link/ether 5a:e0:7e:00:60:c5 brd ff:ff:ff:ff:ff:ff
/ # /bin/mount
overlay on / type overlay (rw,seclabel,relatime,lowerdir=fs/orig,upperdir=/tmp/container087141550/upper,workdir=/tmp/container087141550/work)
proc on /proc type proc (rw,relatime)

/ # ping -c 1 8.8.8.8
PING 8.8.8.8 (8.8.8.8): 56 data bytes
64 bytes from 8.8.8.8: seq=0 ttl=127 time=20.354 ms

--- 8.8.8.8 ping statistics ---
1 packets transmitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 20.354/20.354/20.354 ms
/ # echo "nameserver 8.8.8.8" > /etc/resolv.conf
/ # ping -c 1 google.com
PING google.com (172.217.18.238): 56 data bytes
64 bytes from 172.217.18.238: seq=0 ttl=127 time=16.277 ms

--- google.com ping statistics ---
1 packets transmitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 16.277/16.277/16.277 ms

/ # echo "this is my file" > /testaroo
/ # exit
Removing route {Ifindex: 21 Dst: <nil> Src: <nil> Gw: 172.19.80.1 Flags: [] Table: 0}
fs/rootfs

$ tar tvf /tmp/rootfs.tar |grep testaroo
-rw-r--r-- root/root            16 2016-07-21 09:11 ./testaroo
```



Notes
-----

http://stackoverflow.com/questions/18274088/how-can-i-make-my-own-base-image-for-docker
https://github.com/jpetazzo/docker-busybox
https://github.com/jpetazzo/docker-busybox/raw/master/rootfs.tar

https://www.infoq.com/articles/build-a-container-golang
https://godoc.org/github.com/vishvananda/netlink
https://godoc.org/github.com/vishvananda/netns
https://github.com/opencontainers/runc
http://lk4d4.darth.io/posts/unpriv1/
http://crosbymichael.com/creating-containers-part-1.html

http://stackoverflow.com/questions/22889241/linux-understanding-the-mount-namespace-clone-clone-newns-flag

