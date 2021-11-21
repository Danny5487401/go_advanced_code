#优雅退出
在服务端程序更新或重启时，如果我们直接 kill -9 杀掉旧进程并启动新进程，会有以下几个问题：

    * 旧的请求未处理完，如果服务端进程直接退出，会造成客户端链接中断（收到 RST）
    * 新请求打过来，服务还没重启完毕，造成 connection refused
    * 即使是要退出程序，直接 kill -9 仍然会让正在处理的请求中断
很直接的感受就是：在重启过程中，会有一段时间不能给用户提供正常服务；同时粗鲁关闭服务，也可能会对业务依赖的数据库等状态服务造成污染。

所以我们服务重启或者是重新发布过程中，要做到新旧服务无缝切换，同时可以保障变更服务 零宕机时间
##思路
对 http 服务来说，一般的思路就是关闭对 fd 的 listen , 确保不会有新的请求进来的情况下处理完已经进入的请求, 然后退出。

###连接的状态
```go
type ConnState int

const (
	// 新的连接，并且马上准备发送请求了. Connections begin at this
	// state and then transition to either StateActive or
	// StateClosed.
	StateNew ConnState = iota

	// 表明一个connection已经接收到一个或者多个字节的请求数据. The Server.ConnState hook for
	// StateActive fires before the request has entered a handler
	// and doesn't fire again until the request has been
	// handled. After the request is handled, the state
	// transitions to StateClosed, StateHijacked, or StateIdle.
	// For HTTP/2, StateActive fires on the transition from zero
	// to one active request, and only transitions away once all
	// active requests are complete. That means that ConnState
	// cannot be used to do per-request work; ConnState only notes
	// the overall state of the connection.
	StateActive

	// 表明一个connection已经处理完成一次请求，但因为是keepalived的，所以不会close，继续等待下一次请求.
	// Connections transition from StateIdle
	// to either StateActive or StateClosed.
	StateIdle

	// 表明外部调用了hijack，最终状态.
	// This is a terminal state. It does not transition to StateClosed.
	StateHijacked

	// 表明connection已经结束掉了.
	// This is a terminal state. Hijacked connections do not
	// transition to StateClosed.
	StateClosed
)

var stateName = map[ConnState]string{
	StateNew:      "new",
	StateActive:   "active",
	StateIdle:     "idle",
	StateHijacked: "hijacked",
	StateClosed:   "closed",
}
```


###源码分析:http 中提供了 server.ShutDown()

启动
```go
// go1.16: net/http/server.go
func (srv *Server) ListenAndServe() error {
    // 判断 Server 是否被关闭了
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}
func (s *Server) shuttingDown() bool {
    // 非0表示被关闭
	return s.inShutdown.isSet()
}
```

监听
```go
func (srv *Server) Serve(l net.Listener) error {
    ...
    // 将注入的 listener 加入内部的 map 中
    // 方便后续控制从该 listener 链接到的请求
    if !srv.trackListener(&l, true) {
        return ErrServerClosed
    }
    defer srv.trackListener(&l, false)
   ...
}

// Serve 中注册到内部 listeners map 中 listener，在 ShutDown 中就可以直接从 listeners 中获取到，然后执行 listener.Close()，TCP四次挥手后，新的请求就不会进入了。
func (s *Server) trackListener(ln *net.Listener, add bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listeners == nil {
		s.listeners = make(map[*net.Listener]struct{})
	}
	if add {
		if s.shuttingDown() {
			return false
		}
		s.listeners[ln] = struct{}{}
	} else {
		delete(s.listeners, ln)
	}
	return true
}
```

closeIdleConns：将目前 Server 中记录的活跃链接变成变成空闲状态
```go
func (srv *Server) Serve(l net.Listener) error {
  ...
  for {
    rw, err := l.Accept()
    // 此时 accept 会发生错误，因为前面已经将 listener close了
    if err != nil {
      select {
      // 又是一个标志：doneChan
      case <-srv.getDoneChan():
        return ErrServerClosed
      default:
      }
    }
  }
}
```
###go-zero流程
gracefulStop 的流程如下

	* 取消监听信号，毕竟要退出了，不需要重复监听了
	* wrap up，关闭目前服务请求，以及资源
	* time.Sleep() ，等待资源处理完成，以后关闭完成
	* shutdown ，通知退出
	* 如果主goroutine还没有退出，则主动发送 SIGKILL 退出进程

源码分析:go-zero/core/proc/signals.go
```

//go:build linux || darwin
// +build linux darwin

package proc

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

const timeFormat = "0102150405"

var done = make(chan struct{})

func init() {
	go func() {
		var profiler Stopper

		// https://golang.org/pkg/os/signal/#Notify
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				dumpGoroutines()
			case syscall.SIGUSR2:
				if profiler == nil {
					profiler = StartProfile()
				} else {
					profiler.Stop()
					profiler = nil
				}
			case syscall.SIGTERM:
				select {
				case <-done:
					// already closed
				default:
					close(done)
				}

				gracefulStop(signals)
			default:
				logx.Error("Got unregistered signal:", v)
			}
		}
	}()
}

// Done returns the channel that notifies the process quitting.
func Done() <-chan struct{} {
	return done
}

var noopStopper nilStopper

type (
	// Stopper interface wraps the method Stop.
	Stopper interface {
		Stop()
	}

	nilStopper struct{}
)

func (ns nilStopper) Stop() {
}


var (
	wrapUpListeners          = new(listenerManager)
	shutdownListeners        = new(listenerManager)
	delayTimeBeforeForceQuit = waitTime
)

func gracefulStop(signals chan os.Signal) {
	signal.Stop(signals)

	logx.Info("Got signal SIGTERM, shutting down...")
	wrapUpListeners.notifyListeners()

	time.Sleep(wrapUpTime)
	shutdownListeners.notifyListeners()

	time.Sleep(delayTimeBeforeForceQuit - wrapUpTime)
	logx.Infof("Still alive after %v, going to force kill the process...", delayTimeBeforeForceQuit)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}


```

##流程
![](.grateful_stop_images/graceful_stop.png)

我们目前 go 程序都是在 docker 容器中运行，所以在服务发布过程中，k8s 会向容器发送一个 SIGTERM 信号，然后容器中程序接收到信号，开始执行 ShutDown

但是还有平滑重启，这个就依赖 k8s 了，基本流程如下：

    * old pod 未退出之前，先启动 new pod
    * old pod 继续处理完已经接受的请求，并且不再接受新请求
    * new pod接受并处理新请求的方式
    * old pod 退出
这样整个服务重启就算成功了，如果 new pod 没有启动成功，old pod 也可以提供服务，不会对目前线上的服务造成影响。