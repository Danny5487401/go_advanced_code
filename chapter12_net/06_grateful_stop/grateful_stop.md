#优雅退出
在服务端程序更新或重启时，如果我们直接 kill -9 杀掉旧进程并启动新进程，会有以下几个问题：

    * 旧的请求未处理完，如果服务端进程直接退出，会造成客户端链接中断（收到 RST）
    * 新请求打过来，服务还没重启完毕，造成 connection refused
    * 即使是要退出程序，直接 kill -9 仍然会让正在处理的请求中断
很直接的感受就是：在重启过程中，会有一段时间不能给用户提供正常服务；同时粗鲁关闭服务，也可能会对业务依赖的数据库等状态服务造成污染。

所以我们服务重启或者是重新发布过程中，要做到新旧服务无缝切换，同时可以保障变更服务 零宕机时间
##思路
对 http 服务来说，一般的思路就是关闭对 fd 的 listen , 确保不会有新的请求进来的情况下处理完已经进入的请求, 然后退出。

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
// 参考Go-zero关闭流程




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