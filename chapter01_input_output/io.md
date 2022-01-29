# Linux的文件系统
![](.io_images/linux_file_system.png)
## 文件系统的特点
1. 文件系统要有严格的组织形式，使得文件能够以块为单位进行存储。

2. 文件系统中也要有索引区，用来方便查找一个文件分成的多个块都存放在了什么位置。

3. 如果文件系统中有的文件是热点文件，近期经常被读取和写入，文件系统应该有缓存层。

4. 文件应该用文件夹的形式组织起来，方便管理和查询。

5. Linux内核要在自己的内存里面维护一套数据结构，来保存哪些文件被哪些进程打开和使用。


## I/O输入输出分类
IO 分为网络和存储 IO 两种类型（其实网络 IO 和磁盘 IO 在 Go 里面有着根本性区别 ）
1. 网络 IO 对应的是网络数据传输过程，网络是分布式系统的基石，通过网络把离散的物理节点连接起来，形成一个有机的系统
2. 存储 IO 对应的就是数据存储到物理介质的过程，通常我们物理介质对应的是磁盘，磁盘上一般会分个区，然后在上面格式化个文件系统出来

其实适合 Go 的程序类型准确的来讲是：网络 IO 密集型的程序。

其中差异就在于：
- 网络 fd 可以用 epoll 池来管理事件，实现异步 IO；
- 文件 fd 不能用 epoll 池来管理事件，只能同步 IO；

题外话：文件要实现异步 IO 当前在 Linux 下有两种方式：

- Linux 系统提供的 AIO ：但 Go 没封装实现，因为这个有不少坑；
- Linux 系统提供的 io_uring ：但内核版本要求高，线上没普及；  
一句话，Go 语言级别把网络 IO 做了异步化，但是文件 IO 还是同步的调用。

在 Golang 里可以归类出两种读写文件的方式：

- 标准库封装：操作对象 File;
- 系统调用 ：操作对象 fd;

## 文件的I/O操作
根据是否使用内存做缓存，我们可以把文件的I/O操作分为两种类型。

第一种类型是缓存I/O。大多数文件系统的默认I/O操作都是缓存I/O。对于读操作来讲，操作系统会先检查，内核的缓冲区有没有需要的数据。如果已经缓存了，那就直接从缓存中返回；否则从磁盘中读取，然后缓存在操作系统的缓存中。对于写操作来讲，操作系统会先将数据从用户空间复制到内核空间的缓存中。这时对用户程序来说，写操作就已经完成。至于什么时候再写到磁盘中由操作系统决定，除非显式地调用了sync同步命令。

第二种类型是直接IO，就是应用程序直接访问磁盘数据，而不经过内核缓冲区，从而减少了在内核缓存和用户程序之间数据复制。

```C
ssize_t
generic_file_read_iter(struct kiocb *iocb, struct iov_iter *iter)
{
......
    if (iocb->ki_flags & IOCB_DIRECT) {
......
        struct address_space *mapping = file->f_mapping;
......
        retval = mapping->a_ops->direct_IO(iocb, iter);
    }
......
    retval = generic_file_buffered_read(iocb, iter, retval);
}

ssize_t __generic_file_write_iter(struct kiocb *iocb, struct iov_iter *from)
{
......
    if (iocb->ki_flags & IOCB_DIRECT) {
......
        written = generic_file_direct_write(iocb, from);
......
    } else {
......
        written = generic_perform_write(file, from, iocb->ki_pos);
......
    }
}
```
如果在写的逻辑__generic_file_write_iter里面，发现设置了IOCB_DIRECT，则调用generic_file_direct_write，
里面同样会调用address_space的direct_IO的函数， 将数据直接写入硬盘

## 整个 Unix 体系结构
![](.io_images/unix_structure.png)


* 内核是最核心的实现，包括了和 IO 设备，硬件交互等功能。与内核紧密的一层是内核提供给外部调用的系统调用，系统调用提供了用户态到内核态调用的一个通道；

* 对于系统调用，各个语言的标准库会有一些封装，比如 C 语言的 libc 库，Go 语言的 os ，syscall 库都是类似的地位，这个就是所谓的公共库。
这层封装的作用最主要是简化普通程序员使用效率，并且屏蔽系统细节，为跨平台提供基础（同样的，为了跨平台的特性，可能会阉割很多不兼容的功能，所以才会有直接调用系统掉调用的需求）；

* 当然，我们右上角还看到一个缺口，应用程序除了可以使用公共函数库，其实是可以直接调用系统调用的，但是由此带来的复杂性又应用自己承担。这种需求也是很常见的，
标准库封装了通用的东西，同样割舍了很多系统调用的功能，这种情况下，只能通过系统调用来获取;
## fd文件描述符
文件描述符File descriptor是一个非负整数，本质上是一个索引值（这句话非常重要）。

什么时候拿到的 fd ？

当打开一个文件时，内核向进程返回一个文件描述符（ open 系统调用得到 ），后续 read、write 这个文件时，则只需要用这个文件描述符来标识该文件，将其作为参数传入 read、write 。
fd 的值范围是什么？

在 POSIX 语义中，0，1，2 这三个 fd 值已经被赋予特殊含义，分别是标准输入（ STDIN_FILENO ），标准输出（ STDOUT_FILENO ），标准错误（ STDERR_FILENO ）。
ulimit 命令查看当前系统的配置文件描述符

```shell
ulimit -n
#我的mac os系统上进程默认最多打开 256 文件
```
## 标准库中对IO封装
![](.io_images/io_realized.png)

###  IO 接口描述（语义）
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

io 库的内容，如果按照接口的定义维度大致可以分为 3 大类：
1. 基础类型
   1. Reader、Writer、Closer、ReaderAt、WriterAt、Seeker、ByteReader、ByteWrieter、RuneReader、StringWriter
2. 组合类型
   1. ReaderCloser、WriteCloser、WriteSeeker
```go
type ReadWriter interface {
    Reader
    Writer
}
```
3. 进阶类型
   1. TeeReader: 这是一个分流器的实现，如果把数据比作一个水流，那么通过 TeeReader 之后，将会分叉出一个新的数据流
   2. LimitReader: 则是给 Reader 加了一个边界期限
   3. MultiReader:则是把多个数据流合成一股
   4. SectionReader、MultiWriter、PipeReader、PipeWriter 等