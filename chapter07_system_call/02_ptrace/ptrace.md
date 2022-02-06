# strace 
strace 可用于追踪进程与内核的交互情况，包括系统调用，信号，进程状态等信息。

strace 常用于性能分析、问题定位等场景。

strace 基于 linux 内核特性 ptrace 开发，相当于直接 hook 了系统调用。这意味着 strace 可以追踪所有的用户进程，即使没有被追踪程序的源码，或者程序是非 debug 版本，或者程序日志不完善，都可以通过 strace 追踪到不少有用信息\

## strace 的两种启动方式
```shell
# 另一种是追踪已启动的进程，使用-p加上进程号
$ strace -p <pid> #查看某个进程的系统调用

# 一种是使用 strace 启动被追踪的程序，
$ strace <commond> #查看某条commond指令或进程的系统调用
```

```C
#include <unistd.h>
#include <stdio.h>
int main(){
   for(;;){
       printf("pid=%d\n", getpid());
       sleep(2);
  }
   return 0;
}
```
$ gcc -o print print.c

```shell
centos@xxxxxx:/app/gowork/stramgrpc/c$ strace ./print 
execve("./print", ["./print"], [/* 51 vars */]) = 0
.......
getpid() = 23419
fstat(1, {st_mode=S_IFCHR|0620, st_rdev=makedev(136, 0), ...}) = 0
brk(NULL) = 0x55e3343e2000
brk(0x55e334403000) = 0x55e334403000
write(1, "pid=23419\n", 10pid=23419
) = 10
nanosleep({tv_sec=2, tv_nsec=0}, 0x7ffd2a6d37f0) = 0
getpid() = 23419
write(1, "pid=23419\n", 10pid=23419
) = 10
nanosleep({tv_sec=2, tv_nsec=0}, ^Cstrace: Process 23419 detached
<detached ...>
```
strace命令是c语言实现的，基于Ptrace系统调用。

看一下c标准库的定义规则
```shell
long ptrace(int request, pid_t pid, void *addr, void *data);
```
ptrace在需要传入四个参数：
1. pid用于传入目标进程，也就是要跟性进程的pid；

2. addr和data用于传入内存地址和附加地址，通常会在系统调用结束后读取传入的参数获取系统调用结果，会因操作的不同而不同。

3. request用于选择一个符号标志，内核会根据这个标志决定要选用那个内核函数来执行，接下来介绍一下重点要使用的几个符号标志。

    - 其他的用法不过多展开，感兴趣的同学可以自己探索一下。

    - PTRACE_ATTACH发出一个请求，连接到一个进程并开始跟踪，相反，PTRACE_DETACH从该进程断开并结束跟踪。在调用该指令后，被跟踪进程会发送信号给跟踪者进程，跟踪者进程需要使用waitpid获取该信号，并进行后续的系统调用跟踪。

    - PTRACE_SYSCALL发出系统调用追踪的指令，当使用了该选项时候，被追踪的进程就会在进入系统调用之前或者结束后停下来，这时候追踪者进程可以使用waitpid系统调用时候收到被追踪者发来的通知，从而分析此时的地址空间以及系统调用相关等信息；

    - PTRACE_GETREGS和PTRACE_SETREGS用来设置和读取CPU寄存器，在x86_64的Linux上，系统调用的编号存储在orig_rax寄存器，其他参数是在rdi、rsi、rdx等寄存器，在返回时，返回值存储在rax寄存器；

    - PTRACE_TRACEME：此进程允许被其父进程跟踪(用于strace+命令形式)