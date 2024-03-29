<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [net 网络](#net-%E7%BD%91%E7%BB%9C)
  - [性能指标](#%E6%80%A7%E8%83%BD%E6%8C%87%E6%A0%87)
  - [网络配置](#%E7%BD%91%E7%BB%9C%E9%85%8D%E7%BD%AE)
  - [套接字信息](#%E5%A5%97%E6%8E%A5%E5%AD%97%E4%BF%A1%E6%81%AF)
  - [协议栈统计信息](#%E5%8D%8F%E8%AE%AE%E6%A0%88%E7%BB%9F%E8%AE%A1%E4%BF%A1%E6%81%AF)
  - [golang net包](#golang-net%E5%8C%85)
    - [Conn 接口：一个通用的面向流的网络连接](#conn-%E6%8E%A5%E5%8F%A3%E4%B8%80%E4%B8%AA%E9%80%9A%E7%94%A8%E7%9A%84%E9%9D%A2%E5%90%91%E6%B5%81%E7%9A%84%E7%BD%91%E7%BB%9C%E8%BF%9E%E6%8E%A5)
    - [PacketConn 接口：一种通用的面向数据包的网络连接](#packetconn-%E6%8E%A5%E5%8F%A3%E4%B8%80%E7%A7%8D%E9%80%9A%E7%94%A8%E7%9A%84%E9%9D%A2%E5%90%91%E6%95%B0%E6%8D%AE%E5%8C%85%E7%9A%84%E7%BD%91%E7%BB%9C%E8%BF%9E%E6%8E%A5)
  - [参考链接](#%E5%8F%82%E8%80%83%E9%93%BE%E6%8E%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# net 网络
![](.net_images/net_protocol.png)

- 应用层是网络应用程序和网络协议存放的分层，因特网的应用层包括许多协议，例如我们学 web 离不开的 HTTP，电子邮件传送协议 SMTP、端系统文件上传协议 FTP、还有为我们进行域名解析的 DNS 协议。应用层协议分布在多个端系统上，一个端系统应用程序与另外一个端系统应用程序交换信息分组，我们把位于应用层的信息分组称为 报文(message
- 因特网的运输层在应用程序断点之间传送应用程序报文，在这一层主要有两种传输协议 TCP和 UDP,把运输层的分组称为 报文段(segment)
- 因特网的网络层负责将称为 数据报(datagram) 的网络分层从一台主机移动到另一台主机。网络层一个非常重要的协议是 IP 协议，所有具有网络层的因特网组件都必须运行 IP 协议，IP 协议是一种网际协议，除了 IP 协议外，网络层还包括一些其他网际协议和路由选择协议
- 为了将分组从一个节点（主机或路由器）运输到另一个节点，网络层必须依靠链路层提供服务。链路层的例子包括以太网、WiFi 和电缆接入的 DOCSIS 协议，因为数据从源目的地传送通常需要经过几条链路，一个数据包可能被沿途不同的链路层协议处理，我们把链路层的分组称为 帧(frame)
- 物理层的作用是将帧中的一个个 比特 从一个节点运输到另一个节点

## 性能指标
通常用带宽、吞吐量、延时、PPS（Packet Per Second）等指标衡量网络的性能。

- 带宽，表示链路的最大传输速率，单位通常为 b/s （比特 / 秒）。在你为服务器选购网卡时，带宽就是最核心的参考指标。常用的带宽有 1000M、10G、40G、100G 等。
- 吞吐量，表示单位时间内成功传输的数据量，单位通常为 b/s（比特 / 秒）或者 B/s（字节 / 秒）。吞吐量受带宽限制，而吞吐量 / 带宽，也就是该网络的使用率
- 延时，表示从网络请求发出后，一直到收到远端响应，所需要的时间延迟。在不同场景中，这一指标可能会有不同含义。比如，它可以表示，建立连接需要的时间（比如 TCP 握手延时），或一个数据包往返所需的时间（比如 RTT）。
- PPS，是 Packet Per Second（包 / 秒）的缩写，表示以网络包为单位的传输速率。PPS 通常用来评估网络的转发能力，比如硬件交换机，通常可以达到线性转发（即 PPS 可以达到或者接近理论最大值）。而基于 Linux 服务器的转发，则容易受网络包大小的影响

网络的可用性（网络能否正常通信）、并发连接数（TCP 连接数量）、丢包率（丢包百分比）、重传率（重新传输的网络包比例）等也是常用的性能指标。

对 TCP 或者 Web 服务来说，更多会用并发连接数和每秒请求数（QPS，Query per Second）等指标，它们更能反应实际应用程序的性能

## 网络配置

查看网络接口的配置和状态。你可以使用 ifconfig 或者 ip 命令，来查看网络的配置。

> ifconfig 和 ip 分别属于软件包 net-tools 和 iproute2，iproute2 是 net-tools 的下一代。
> 通常情况下它们会在发行版中默认安装。但如果你找不到 ifconfig 或者 ip 命令，可以安装这两个软件包。

```shell
ubuntu@VM-16-12-ubuntu:~$ ifconfig eth0
eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.0.16.12  netmask 255.255.252.0  broadcast 10.0.19.255
        inet6 fe80::5054:ff:fe82:711f  prefixlen 64  scopeid 0x20<link>
        ether 52:54:00:82:71:1f  txqueuelen 1000  (Ethernet)
        RX packets 18861576  bytes 6778423168 (6.7 GB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 14497669  bytes 2177803591 (2.1 GB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

ubuntu@VM-16-12-ubuntu:~$ ip -s addr show dev eth0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether 52:54:00:82:71:1f brd ff:ff:ff:ff:ff:ff
    inet 10.0.16.12/22 brd 10.0.19.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::5054:ff:fe82:711f/64 scope link 
       valid_lft forever preferred_lft forever
    RX: bytes  packets  errors  dropped overrun mcast   
    6778428051 18861636 0       0       0       0       
    TX: bytes  packets  errors  dropped carrier collsns 
    2177813118 14497724 0       0       0       0  
```
解释：

1. 网络接口的状态标志。ifconfig 输出中的 RUNNING ，或 ip 输出中的 LOWER_UP ，都表示物理网络是连通的，即网卡已经连接到了交换机或者路由器中。如果你看不到它们，通常表示网线被拔掉了
2. MTU 的大小。MTU 默认大小是 1500，根据网络架构的不同（比如是否使用了 VXLAN 等叠加网络），你可能需要调大或者调小 MTU 的数值
3. 网络接口的 IP 地址、子网以及 MAC 地址。
4. 网络收发的字节数、包数、错误数以及丢包情况，特别是 TX 和 RX 部分的 errors、dropped、overruns、carrier 以及 collisions 等指标不为 0 时，通常表示出现了网络 I/O 问题。


## 套接字信息
ifconfig 和 ip 只显示了网络接口收发数据包的统计信息，但在实际的性能问题中，网络协议栈中的统计信息，我们也必须关注。你可以用 netstat 或者 ss ，来查看套接字、网络栈、网络接口以及路由表的信息。

> 使用 ss 来查询网络的连接信息，因为它比 netstat 提供了更好的性能（速度更快）

```shell
# head -n 3 表示只显示前面 3 行
# -l 表示只显示监听套接字
# -n 表示显示数字地址和端口 (而不是名字)
# -p 表示显示进程信息
root@VM-16-12-ubuntu:/home/ubuntu# netstat -nlp | head -n 3
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name    
tcp        0      0 0.0.0.0:9200            0.0.0.0:*               LISTEN      2390997/docker-prox 

# -l 表示只显示监听套接字
# -t 表示只显示 TCP 套接字
# -n 表示显示数字地址和端口 (而不是名字)
# -p 表示显示进程信息
root@VM-16-12-ubuntu:/home/ubuntu# ss -ltnp | head -n 3
State    Recv-Q   Send-Q                                   Local Address:Port                                   Peer Address:Port                               Process                                                                         
LISTEN   0        4096                                           0.0.0.0:9200                                        0.0.0.0:*                                   users:(("docker-proxy",pid=2390997,fd=4))                                      
LISTEN   0        10                                        192.168.80.1:53                                          0.0.0.0:*                                   users:(("named",pid=768,fd=62),("named",pid=768,fd=61),("named",pid=768,fd=60))
```

netstat 和 ss 的输出也是类似的，都展示了套接字的状态、接收队列、发送队列、本地地址、远端地址、进程 PID 和进程名称等

其中，接收队列（Recv-Q）和发送队列（Send-Q）需要你特别关注，它们通常应该是 0。当你发现它们不是 0 时，说明有网络包的堆积发生。当然还要注意，在不同套接字状态下，它们的含义不同。


当套接字处于连接状态（Established）时，

- Recv-Q 表示套接字缓冲还没有被应用程序取走的字节数（即接收队列长度）。

- 而 Send-Q 表示还没有被远端主机确认的字节数（即发送队列长度）。

当套接字处于监听状态（Listening）时，
- Recv-Q 表示 syn backlog 的当前值。

- 而 Send-Q 表示最大的 syn backlog 值。

而 syn backlog 是 TCP 协议栈中的半连接队列长度，相应的也有一个全连接队列（accept queue），它们都是维护 TCP 状态的重要机制.

所谓半连接，就是还没有完成 TCP 三次握手的连接，连接只进行了一半，而服务器收到了客户端的 SYN 包后，就会把这个连接放到半连接队列中，然后再向客户端发送 SYN+ACK 包

而全连接，则是指服务器收到了客户端的 ACK，完成了 TCP 三次握手，然后就会把这个连接挪到全连接队列中。这些全连接中的套接字，还需要再被 accept() 系统调用取走，这样，服务器就可以开始真正处理客户端的请求了


## 协议栈统计信息


```shell
$ netstat -s

Tcp:
    3244906 active connection openings
    23143 passive connection openings
    115732 failed connection attempts
    2964 connection resets received
    1 connections established
    13025010 segments received
    17606946 segments sent out
    44438 segments retransmitted
    42 bad segments received
    5315 resets sent
    InCsumErrors: 42
...
 
$ ss -s
Total: 186 (kernel 1446)
TCP:   4 (estab 1, closed 0, orphaned 0, synrecv 0, timewait 0/0), ports 0
 
Transport Total     IP        IPv6
*	  1446      -         -
RAW	  2         1         1
UDP	  2         2         0
TCP	  4         3         1

```
ss 只显示已经连接、关闭、孤儿套接字等简要统计，而 netstat 则提供的是更详细的网络协议栈信息


## golang net包

net包提供了可移植的网络I/O接口，包括TCP/IP、UDP、域名解析和Unix域socket等方式的通信。
其中每一种通信方式都使用 xxConn 结构体来表示，诸如IPConn、TCPConn等，这些结构体都实现了Conn接口，Conn接口实现了基本的读、写、关闭、获取远程和本地地址、设置timeout等功能。


### Conn 接口：一个通用的面向流的网络连接
多个goroutines可以同时调用Conn上的方法
```go
type Conn interface {
    // Read从连接中读取数据
    // Read方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
    Read(b []byte) (n int, err error)
    // Write从连接中写入数据
    // Write方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
    Write(b []byte) (n int, err error)
    // Close方法关闭该连接
    // 并会导致任何阻塞中的Read或Write方法不再阻塞并返回错误
    Close() error
    // 返回本地网络地址
    LocalAddr() Addr
    // 返回远端网络地址
    RemoteAddr() Addr
    // 设定该连接的读写deadline，等价于同时调用SetReadDeadline和SetWriteDeadline
    // deadline是一个绝对时间，超过该时间后I/O操作就会直接因超时失败返回而不会阻塞
    // deadline对之后的所有I/O操作都起效，而不仅仅是下一次的读或写操作
    // 参数t为零值表示不设置期限
    SetDeadline(t time.Time) error
    // 设定该连接的读操作deadline，参数t为零值表示不设置期限
    SetReadDeadline(t time.Time) error
    // 设定该连接的写操作deadline，参数t为零值表示不设置期限
    // 即使写入超时，返回值n也可能>0，说明成功写入了部分数据
    SetWriteDeadline(t time.Time) error
}
```
每种类型都是对应的结构体实现这些接口。

### PacketConn 接口：一种通用的面向数据包的网络连接
```go
type PacketConn interface {
    // ReadFrom方法从连接读取一个数据包，并将有效信息写入b
    // ReadFrom方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
    // 返回写入的字节数和该数据包的来源地址
    ReadFrom(b []byte) (n int, addr Addr, err error)
    // WriteTo方法将有效数据b写入一个数据包发送给addr
    // WriteTo方法可能会在超过某个固定时间限制后超时返回错误，该错误的Timeout()方法返回真
    // 在面向数据包的连接中，写入超时非常罕见
    WriteTo(b []byte, addr Addr) (n int, err error)
    // Close方法关闭该连接
    // 会导致任何阻塞中的ReadFrom或WriteTo方法不再阻塞并返回错误
    Close() error
    // 返回本地网络地址
    LocalAddr() Addr
    // 设定该连接的读写deadline
    SetDeadline(t time.Time) error
    // 设定该连接的读操作deadline，参数t为零值表示不设置期限
    // 如果时间到达deadline，读操作就会直接因超时失败返回而不会阻塞
    SetReadDeadline(t time.Time) error
    // 设定该连接的写操作deadline，参数t为零值表示不设置期限
    // 如果时间到达deadline，写操作就会直接因超时失败返回而不会阻塞
    // 即使写入超时，返回值n也可能>0，说明成功写入了部分数据
    SetWriteDeadline(t time.Time) error
}
```




## 参考链接
