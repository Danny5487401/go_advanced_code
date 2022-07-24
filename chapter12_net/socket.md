# socket
![](.socket_images/socket_buffer.png)

socket 在操作系统层面，可以理解为一个文件。 我们可以对这个文件进行一些方法操作。

用listen方法，可以让程序作为服务器监听其他客户端的连接。

用connect，可以作为客户端连接服务器。

用send或write可以发送数据，recv或read可以接收数据

## socket缓冲区

在建立好连接之后，这个 socket 文件就像是远端机器的 "代理人" 一样。比如，如果我们想给远端服务发点什么东西，那就只需要对这个文件执行写操作就行了。 那写到了这个文件之后，剩下的发送工作自然就是由操作系统内核来完成了。

既然是写给操作系统，那操作系统就需要提供一个地方给用户写。同理，接收消息也是一样。 这个地方就是 socket 缓冲区。
    
- 用户发送消息的时候写给 send buffer（发送缓冲区）
    
- 用户接收消息的时候写给 recv buffer（接收缓冲区）
    
也就是说一个socket ，会带有两个缓冲区，一个用于发送，一个用于接收。因为这是个先进先出的结构，有时候也叫它们发送、接收队列
## 观察 socket 缓冲区

在linux环境下执行 netstat -nt 
```shell
# netstat -nt
Active Internet connections (w/o servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State      
tcp        0     60 172.22.66.69:22         122.14.220.252:59889    ESTABLISHED

```
这上面表明了，这里有一个协议（Proto）类型为 TCP 的连接，同时还有本地（Local Address）和远端（Foreign Address）的IP信息，状态（State）是已连接。

还有Send-Q 是发送缓冲区，下面的数字60是指，当前还有60 Byte在发送缓冲区中未发送。
而 Recv-Q 代表接收缓冲区，此时是空的，数据都被应用进程接收干净了

## tcp四次握手
## 执行 send 发送的字节，会立马发送吗
![](.socket_images/socket_send.png)

答案是不确定！执行 send 之后，数据只是拷贝到了socket 缓冲区。至 于什么时候会发数据，发多少数据，全听操作系统安排

## 如果缓冲区满了会怎么办
首先，socket在创建的时候，是可以设置是阻塞的还是非阻塞的。
![](.socket_images/send_block.png)
![](.socket_images/send_nonblock.png)

## 如果接收缓冲区为空，执行 recv 会怎么样？
如果此时 socket 是阻塞的，那么程序会在那干等，直到接收缓冲区有数据，就会把数据从接收缓冲区拷贝到用户缓冲区，然后返回
![](.socket_images/recv_block.png)
![](.socket_images/recv_nonblock.png)

## 如果接收缓冲区有数据时，执行close了，会怎么样？
![](.socket_images/recvbuf_nonEmpty.png)
socket close 时，主要的逻辑在 tcp_close() 里实现。

先说结论，关闭过程主要有两种情况：

- 如果接收缓冲区还有数据未读，会先把接收缓冲区的数据清空，然后给对端发一个RST。
    
- 如果接收缓冲区是空的，那么就调用 tcp_send_fin() 开始进行四次挥手过程的第一次挥手

```C
void tcp_close(struct sock *sk, long timeout)
{
  // 如果接收缓冲区有数据，那么清空数据
    while ((skb = __skb_dequeue(&sk->sk_receive_queue)) != NULL) {
        u32 len = TCP_SKB_CB(skb)->end_seq - TCP_SKB_CB(skb)->seq -
              tcp_hdr(skb)->fin;
        data_was_unread += len;
        __kfree_skb(skb);
    }

   if (data_was_unread) {
    // 如果接收缓冲区的数据被清空了，发 RST
        tcp_send_active_reset(sk, sk->sk_allocation);
     } else if (tcp_close_state(sk)) {
    // 正常四次挥手, 发 FIN
        tcp_send_fin(sk);
    }
    // 等待关闭
    sk_stream_wait_close(sk, timeout);
}
```
## 如果发送缓冲区有数据时，执行close了，会怎么样？
![](.socket_images/sendbuf_nonEmpty.png)
```c
void tcp_send_fin(struct sock *sk)
{
  // 获得发送缓冲区的最后一块数据
    struct sk_buff *skb, *tskb = tcp_write_queue_tail(sk);
    struct tcp_sock *tp = tcp_sk(sk);

  // 如果发送缓冲区还有数据
    if (tskb && (tcp_send_head(sk) || sk_under_memory_pressure(sk))) {
        TCP_SKB_CB(tskb)->tcp_flags |= TCPHDR_FIN; // 把最后一块数据值为 FIN 
        TCP_SKB_CB(tskb)->end_seq++;
        tp->write_seq++;
    }  else {
    // 发送缓冲区没有数据，就造一个FIN包
  }
  // 发送数据
    __tcp_push_pending_frames(sk, tcp_current_mss(sk), TCP_NAGLE_OFF);
}
```
此时，还有些数据没发出去，内核会把发送缓冲区最后一个数据块拿出来。然后置为 FIN

    socket 缓冲区是个先进先出的队列，这种情况是指内核会等待TCP层安静把发送缓冲区数据都发完，最后再执行四次挥手的第一次挥手（FIN包）。
    
    有一点需要注意的是，只有在接收缓冲区为空的前提下，我们才有可能走到 tcp_send_fin() 。而只有在进入了这个方法之后，我们才有可能考虑发送缓冲区是否为空的场景。

## UDP也有缓冲区吗
UDP socket 也是 socket，一个socket 就是会有收和发两个缓冲区，跟用什么协议关系不大
```shell
int udp_sendmsg()
{
    // corkreq 为 true 表示是 MSG_MORE 的方式，仅仅组织报文，不发送；
    int corkreq = up->corkflag || msg->msg_flags&MSG_MORE；

    //  将要发送的数据，按照MTU大小分割，每个片段一个skb；并且这些
    //  skb会放入到套接字的发送缓冲区中；该函数只是组织数据包，并不执行发送动作。
    err = ip_append_data(sk, fl4, getfrag, msg->msg_iov, ulen,
                 sizeof(struct udphdr), &ipc, &rt,
                 corkreq ? msg->msg_flags|MSG_MORE : msg->msg_flags);

    // 没有启用 MSG_MORE 特性，那么直接将发送队列中的数据发送给IP。 
    if (!corkreq)
        err = udp_push_pending_frames(sk);

}
```

而我们大部分情况下，都不会用  MSG_MORE，也就是来一个数据包就直接发一个数据包。从这个行为上来说，虽然UDP用上了发送缓冲区，但实际上并没有起到"缓冲"的作用