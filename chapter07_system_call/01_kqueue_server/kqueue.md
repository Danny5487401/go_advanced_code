<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [kqueue](#kqueue)
  - [源码](#%E6%BA%90%E7%A0%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# kqueue

kqueue是BSD中使用的内核事件通知机制,一个kqueue指的是一个描述符，这个描述符会阻塞等待直到一个特定的类型和种类的事件发生


Microsoft Windows 是一套操作系统，FreeBSD 也是一套操作系统。Mac OS X 就是使用 FreeBSD 做为系统核心


## 源码


```go
// go1.21.5/src/syscall/ztypes_darwin_arm64.go
type Kevent_t struct {
	Ident  uint64 //该事件关联的描述符，常见的有socket fd，file fd， signal fd等
	Filter int16 //事件的类型，比如读事件EVFILT_READ，写事件EVFILT_WRITE，信号事件EVFILT_SIGNAL
	//事件的行为，也就是对kqueue的操作，下面介绍几个常用的
	//如EV_ADD：添加到kqueue中，EV_DELETE从kqueue中删除
	//EV_ONESHOT：一次性或事件，kevent返回后从kqueue中删除
	//EV_CLEAR：事件通知给用户后，事件的状态会重置，
	Flags  uint16
	Fflags uint32
	Data   int64
	Udata  *byte //用户指定的数据
}

```


## 参考

- https://dev.to/frosnerd/writing-a-simple-tcp-server-using-kqueue-cah