<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [exec](#exec)
  - [Shell 脚本的 5 种执行方式](#shell-%E8%84%9A%E6%9C%AC%E7%9A%84-5-%E7%A7%8D%E6%89%A7%E8%A1%8C%E6%96%B9%E5%BC%8F)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第三方使用--istio pilot agent 调用 envoy](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8--istio-pilot-agent-%E8%B0%83%E7%94%A8-envoy)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# exec 


创建一个新进程分为两个步骤，一个是fork系统调用，一个是execve 系统调用，fork调用会复用父进程的堆栈，而execve直接覆盖当前进程的堆栈，并且将下一条执行指令指向新的可执行文件。

## Shell 脚本的 5 种执行方式


1 使用绝对路径执行：/root/sleep.sh

2 使用相对路径执行：./sleep.sh （需要 x 权限）

3 使用 sh 或 bash 命令来执行：bash /root/sleep.sh

4 使用 . (空格)脚本名称来执行：. /root/sleep.sh

5 使用 source 来执行(一般用于生效配置文件)：source /root/sleep.sh


## 源码分析

结构
```go
// go1.21.5/src/os/exec/exec.go

// Cmd结构体代表一个准备或正在执行的外部命令
// 一个Cmd的对象不能在Run，Output或者CombinedOutput方法调用之后重复使用。
type Cmd struct {
    // Path代表运行命令的路径
    // 这个字段是唯一一个需要被赋值的字段，不能是空字符串，
    // 并且如果Path是相对路径，那么参照的是Dir这个字段的所指向的目录
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env代表执行命令进程的环境变量
	// 每个Env数组中的条目都以“key=value”的形式存在
	// 如果Env是nil，那边运行命令所创建的进程将使用当前进程的环境变量
	// 如果Env中存在重复的key，那么会使用这个key中排在最后一个的值。
	// 在Windows中存在特殊的情况, 如果系统中缺失了SYSTEMROOT，或者这个环境变量没有被设置成空字符串，那么它操作都是追加操作。
	Env []string

    // Dir代表命令的运行路径
    // 如果Dir是空字符串，那么命令就会运行在当前进程的运行路径
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error), or because writing to the pipe returned an error,
	// or because a nonzero WaitDelay was set and expired.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error or a nonzero WaitDelay
	// expires.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// ExtraFiles specifies additional open files to be inherited by the
	// new process. It does not include standard input, standard output, or
	// standard error. If non-nil, entry i becomes file descriptor 3+i.
	//
	// ExtraFiles is not supported on Windows.
	ExtraFiles []*os.File

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr

	// 当命令运行之后，Process就是该命令运行所代表的进程
	Process *os.Process

	// ProcessState包含关于一个退出的进程的信息，在调用Wait或Run后可用。
	ProcessState *os.ProcessState

	// ctx is the context passed to CommandContext, if any.
	ctx context.Context

	Err error // LookPath error, if any.

	// If Cancel is non-nil, the command must have been created with
	// CommandContext and Cancel will be called when the command's
	// Context is done. By default, CommandContext sets Cancel to
	// call the Kill method on the command's Process.
	//
	// Typically a custom Cancel will send a signal to the command's
	// Process, but it may instead take other actions to initiate cancellation,
	// such as closing a stdin or stdout pipe or sending a shutdown request on a
	// network socket.
	//
	// If the command exits with a success status after Cancel is
	// called, and Cancel does not return an error equivalent to
	// os.ErrProcessDone, then Wait and similar methods will return a non-nil
	// error: either an error wrapping the one returned by Cancel,
	// or the error from the Context.
	// (If the command exits with a non-success status, or Cancel
	// returns an error that wraps os.ErrProcessDone, Wait and similar methods
	// continue to return the command's usual exit status.)
	//
	// If Cancel is set to nil, nothing will happen immediately when the command's
	// Context is done, but a nonzero WaitDelay will still take effect. That may
	// be useful, for example, to work around deadlocks in commands that do not
	// support shutdown signals but are expected to always finish quickly.
	//
	// Cancel will not be called if Start returns a non-nil error.
	Cancel func() error

	// If WaitDelay is non-zero, it bounds the time spent waiting on two sources
	// of unexpected delay in Wait: a child process that fails to exit after the
	// associated Context is canceled, and a child process that exits but leaves
	// its I/O pipes unclosed.
	//
	// The WaitDelay timer starts when either the associated Context is done or a
	// call to Wait observes that the child process has exited, whichever occurs
	// first. When the delay has elapsed, the command shuts down the child process
	// and/or its I/O pipes.
	//
	// If the child process has failed to exit — perhaps because it ignored or
	// failed to receive a shutdown signal from a Cancel function, or because no
	// Cancel function was set — then it will be terminated using os.Process.Kill.
	//
	// Then, if the I/O pipes communicating with the child process are still open,
	// those pipes are closed in order to unblock any goroutines currently blocked
	// on Read or Write calls.
	//
	// If pipes are closed due to WaitDelay, no Cancel call has occurred,
	// and the command has otherwise exited with a successful status, Wait and
	// similar methods will return ErrWaitDelay instead of nil.
	//
	// If WaitDelay is zero (the default), I/O pipes will be read until EOF,
	// which might not occur until orphaned subprocesses of the command have
	// also closed their descriptors for the pipes.
	WaitDelay time.Duration

	// childIOFiles holds closers for any of the child process's
	// stdin, stdout, and/or stderr files that were opened by the Cmd itself
	// (not supplied by the caller). These should be closed as soon as they
	// are inherited by the child process.
	childIOFiles []io.Closer

	// parentIOPipes holds closers for the parent's end of any pipes
	// connected to the child's stdin, stdout, and/or stderr streams
	// that were opened by the Cmd itself (not supplied by the caller).
	// These should be closed after Wait sees the command and copying
	// goroutines exit, or after WaitDelay has expired.
	parentIOPipes []io.Closer

	// goroutine holds a set of closures to execute to copy data
	// to and/or from the command's I/O pipes.
	goroutine []func() error

	// If goroutineErr is non-nil, it receives the first error from a copying
	// goroutine once all such goroutines have completed.
	// goroutineErr is set to nil once its error has been received.
	goroutineErr <-chan error

	// If ctxResult is non-nil, it receives the result of watchCtx exactly once.
	ctxResult <-chan ctxResult

	// The stack saved when the Command was created, if GODEBUG contains
	// execwait=2. Used for debugging leaks.
	createdByStack []byte

	// For a security release long ago, we created x/sys/execabs,
	// which manipulated the unexported lookPathErr error field
	// in this struct. For Go 1.19 we exported the field as Err error,
	// above, but we have to keep lookPathErr around for use by
	// old programs building against new toolchains.
	// The String and Start methods look for an error in lookPathErr
	// in preference to Err, to preserve the errors that execabs sets.
	//
	// In general we don't guarantee misuse of reflect like this,
	// but the misuse of reflect was by us, the best of various bad
	// options to fix the security problem, and people depend on
	// those old copies of execabs continuing to work.
	// The result is that we have to leave this variable around for the
	// rest of time, a compatibility scar.
	//
	// See https://go.dev/blog/path-security
	// and https://go.dev/issue/43724 for more context.
	lookPathErr error
}
```


```go
// go1.21.5/src/os/exec.go
func StartProcess(name string, argv []string, attr *ProcAttr) (*Process, error) {
	testlog.Open(name)
	return startProcess(name, argv, attr)
}
```

golang中 fork 和execve 创建子进程 的过程 被封装成了一个统一的方法forkExec，它能够控制子进程，只继承特定的文件描述符，而对其他文件描述符则进行关闭。
而内核fork系统调用则是会对父进程的所有文件描述符进行复制，那么golang又是如何做到只继承特定的文件描述符的呢？


os.StartProcess 底层会调用到 forkAndExecInChild1 方法
```go
// go1.21.5/src/syscall/exec_linux.go
func forkAndExecInChild1(argv0 *byte, argv, envv []*byte, chroot, dir *byte, attr *ProcAttr, sys *SysProcAttr, pipe int) (pid uintptr, err1 Errno, mapPipe [2]int, locked bool) {
    // ... 

	// Guard against side effects of shuffling fds below.
	// Make sure that nextfd is beyond any currently open files so
	// that we can't run the risk of overwriting any of them.
	fd := make([]int, len(attr.Files))
	nextfd = len(attr.Files)
	for i, ufd := range attr.Files {
		if nextfd < int(ufd) {
			nextfd = int(ufd)
		}
		fd[i] = int(ufd)
	}
	nextfd++

	// Allocate another pipe for parent to child communication for
	// synchronizing writing of User ID/Group ID mappings.
	if sys.UidMappings != nil || sys.GidMappings != nil {
		if err := forkExecPipe(mapPipe[:]); err != nil {
			err1 = err.(Errno)
			return
		}
	}

	flags = sys.Cloneflags
	if sys.Cloneflags&CLONE_NEWUSER == 0 && sys.Unshareflags&CLONE_NEWUSER == 0 {
		flags |= CLONE_VFORK | CLONE_VM
	}
	// Whether to use clone3.
	if sys.UseCgroupFD {
		clone3 = &cloneArgs{
			flags:      uint64(flags) | CLONE_INTO_CGROUP,
			exitSignal: uint64(SIGCHLD),
			cgroup:     uint64(sys.CgroupFD),
		}
	} else if flags&CLONE_NEWTIME != 0 {
		clone3 = &cloneArgs{
			flags:      uint64(flags),
			exitSignal: uint64(SIGCHLD),
		}
	}

	// About to call fork.
	// No more allocation or calls of non-assembly functions.
	// 这里便进行了fork调用创建新进程了，不过可以看到这里用的是clone系统调用，其实它和fork类似，不过区别在于clone系统调用可以通过flags指定新进程 对于 父进程的哪些属性需要继承，哪些属性不需要继承，比如子进程需要新的网络命名空间，则需要指定flags为syscall.CLONE_NEWNS
	runtime_BeforeFork()
	locked = true
	if clone3 != nil {
		pid, err1 = rawVforkSyscall(_SYS_clone3, uintptr(unsafe.Pointer(clone3)), unsafe.Sizeof(*clone3))
	} else {
		flags |= uintptr(SIGCHLD)
		if runtime.GOARCH == "s390x" {
			// On Linux/s390, the first two arguments of clone(2) are swapped.
			pid, err1 = rawVforkSyscall(SYS_CLONE, 0, flags)
		} else {
			pid, err1 = rawVforkSyscall(SYS_CLONE, flags, 0)
		}
	}
	if err1 != 0 || pid != 0 {
		// If we're in the parent, we must return immediately
		// so we're not in the same stack frame as the child.
		// This can at most use the return PC, which the child
		// will not modify, and the results of
		// rawVforkSyscall, which must have been written after
		// the child was replaced.
		return
	}

	// Fork succeeded, now in child.

	// Enable the "keep capabilities" flag to set ambient capabilities later.
	if len(sys.AmbientCaps) > 0 {
		_, _, err1 = RawSyscall6(SYS_PRCTL, PR_SET_KEEPCAPS, 1, 0, 0, 0, 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// Wait for User ID/Group ID mappings to be written.
	if sys.UidMappings != nil || sys.GidMappings != nil {
		if _, _, err1 = RawSyscall(SYS_CLOSE, uintptr(mapPipe[1]), 0, 0); err1 != 0 {
			goto childerror
		}
		pid, _, err1 = RawSyscall(SYS_READ, uintptr(mapPipe[0]), uintptr(unsafe.Pointer(&err2)), unsafe.Sizeof(err2))
		if err1 != 0 {
			goto childerror
		}
		if pid != unsafe.Sizeof(err2) {
			err1 = EINVAL
			goto childerror
		}
		if err2 != 0 {
			err1 = err2
			goto childerror
		}
	}

	// Session ID
	if sys.Setsid {
		_, _, err1 = RawSyscall(SYS_SETSID, 0, 0, 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// Set process group
	if sys.Setpgid || sys.Foreground {
		// Place child in process group.
		_, _, err1 = RawSyscall(SYS_SETPGID, 0, uintptr(sys.Pgid), 0)
		if err1 != 0 {
			goto childerror
		}
	}

	if sys.Foreground {
		pgrp = int32(sys.Pgid)
		if pgrp == 0 {
			pid, _ = rawSyscallNoError(SYS_GETPID, 0, 0, 0)

			pgrp = int32(pid)
		}

		// Place process group in foreground.
		_, _, err1 = RawSyscall(SYS_IOCTL, uintptr(sys.Ctty), uintptr(TIOCSPGRP), uintptr(unsafe.Pointer(&pgrp)))
		if err1 != 0 {
			goto childerror
		}
	}

	// Restore the signal mask. We do this after TIOCSPGRP to avoid
	// having the kernel send a SIGTTOU signal to the process group.
	runtime_AfterForkInChild()

	// Unshare
	if sys.Unshareflags != 0 {
		_, _, err1 = RawSyscall(SYS_UNSHARE, sys.Unshareflags, 0, 0)
		if err1 != 0 {
			goto childerror
		}

		if sys.Unshareflags&CLONE_NEWUSER != 0 && sys.GidMappings != nil {
			dirfd = int(_AT_FDCWD)
			if fd1, _, err1 = RawSyscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(&psetgroups[0])), uintptr(O_WRONLY), 0, 0, 0); err1 != 0 {
				goto childerror
			}
			pid, _, err1 = RawSyscall(SYS_WRITE, uintptr(fd1), uintptr(unsafe.Pointer(&setgroups[0])), uintptr(len(setgroups)))
			if err1 != 0 {
				goto childerror
			}
			if _, _, err1 = RawSyscall(SYS_CLOSE, uintptr(fd1), 0, 0); err1 != 0 {
				goto childerror
			}

			if fd1, _, err1 = RawSyscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(&pgid[0])), uintptr(O_WRONLY), 0, 0, 0); err1 != 0 {
				goto childerror
			}
			pid, _, err1 = RawSyscall(SYS_WRITE, uintptr(fd1), uintptr(unsafe.Pointer(&gidmap[0])), uintptr(len(gidmap)))
			if err1 != 0 {
				goto childerror
			}
			if _, _, err1 = RawSyscall(SYS_CLOSE, uintptr(fd1), 0, 0); err1 != 0 {
				goto childerror
			}
		}

		if sys.Unshareflags&CLONE_NEWUSER != 0 && sys.UidMappings != nil {
			dirfd = int(_AT_FDCWD)
			if fd1, _, err1 = RawSyscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(&puid[0])), uintptr(O_WRONLY), 0, 0, 0); err1 != 0 {
				goto childerror
			}
			pid, _, err1 = RawSyscall(SYS_WRITE, uintptr(fd1), uintptr(unsafe.Pointer(&uidmap[0])), uintptr(len(uidmap)))
			if err1 != 0 {
				goto childerror
			}
			if _, _, err1 = RawSyscall(SYS_CLOSE, uintptr(fd1), 0, 0); err1 != 0 {
				goto childerror
			}
		}

		// The unshare system call in Linux doesn't unshare mount points
		// mounted with --shared. Systemd mounts / with --shared. For a
		// long discussion of the pros and cons of this see debian bug 739593.
		// The Go model of unsharing is more like Plan 9, where you ask
		// to unshare and the namespaces are unconditionally unshared.
		// To make this model work we must further mark / as MS_PRIVATE.
		// This is what the standard unshare command does.
		if sys.Unshareflags&CLONE_NEWNS == CLONE_NEWNS {
			_, _, err1 = RawSyscall6(SYS_MOUNT, uintptr(unsafe.Pointer(&none[0])), uintptr(unsafe.Pointer(&slash[0])), 0, MS_REC|MS_PRIVATE, 0, 0)
			if err1 != 0 {
				goto childerror
			}
		}
	}

	// Chroot
	if chroot != nil {
		_, _, err1 = RawSyscall(SYS_CHROOT, uintptr(unsafe.Pointer(chroot)), 0, 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// User and groups
	if cred = sys.Credential; cred != nil {
		ngroups = uintptr(len(cred.Groups))
		groups = uintptr(0)
		if ngroups > 0 {
			groups = uintptr(unsafe.Pointer(&cred.Groups[0]))
		}
		if !(sys.GidMappings != nil && !sys.GidMappingsEnableSetgroups && ngroups == 0) && !cred.NoSetGroups {
			_, _, err1 = RawSyscall(_SYS_setgroups, ngroups, groups, 0)
			if err1 != 0 {
				goto childerror
			}
		}
		_, _, err1 = RawSyscall(sys_SETGID, uintptr(cred.Gid), 0, 0)
		if err1 != 0 {
			goto childerror
		}
		_, _, err1 = RawSyscall(sys_SETUID, uintptr(cred.Uid), 0, 0)
		if err1 != 0 {
			goto childerror
		}
	}

	if len(sys.AmbientCaps) != 0 {
		// Ambient capabilities were added in the 4.3 kernel,
		// so it is safe to always use _LINUX_CAPABILITY_VERSION_3.
		caps.hdr.version = _LINUX_CAPABILITY_VERSION_3

		if _, _, err1 = RawSyscall(SYS_CAPGET, uintptr(unsafe.Pointer(&caps.hdr)), uintptr(unsafe.Pointer(&caps.data[0])), 0); err1 != 0 {
			goto childerror
		}

		for _, c = range sys.AmbientCaps {
			// Add the c capability to the permitted and inheritable capability mask,
			// otherwise we will not be able to add it to the ambient capability mask.
			caps.data[capToIndex(c)].permitted |= capToMask(c)
			caps.data[capToIndex(c)].inheritable |= capToMask(c)
		}

		if _, _, err1 = RawSyscall(SYS_CAPSET, uintptr(unsafe.Pointer(&caps.hdr)), uintptr(unsafe.Pointer(&caps.data[0])), 0); err1 != 0 {
			goto childerror
		}

		for _, c = range sys.AmbientCaps {
			_, _, err1 = RawSyscall6(SYS_PRCTL, PR_CAP_AMBIENT, uintptr(PR_CAP_AMBIENT_RAISE), c, 0, 0, 0)
			if err1 != 0 {
				goto childerror
			}
		}
	}

	// Chdir
	if dir != nil {
		_, _, err1 = RawSyscall(SYS_CHDIR, uintptr(unsafe.Pointer(dir)), 0, 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// Parent death signal
	if sys.Pdeathsig != 0 {
		_, _, err1 = RawSyscall6(SYS_PRCTL, PR_SET_PDEATHSIG, uintptr(sys.Pdeathsig), 0, 0, 0, 0)
		if err1 != 0 {
			goto childerror
		}

		// Signal self if parent is already dead. This might cause a
		// duplicate signal in rare cases, but it won't matter when
		// using SIGKILL.
		pid, _ = rawSyscallNoError(SYS_GETPPID, 0, 0, 0)
		if pid != ppid {
			pid, _ = rawSyscallNoError(SYS_GETPID, 0, 0, 0)
			_, _, err1 = RawSyscall(SYS_KILL, pid, uintptr(sys.Pdeathsig), 0)
			if err1 != 0 {
				goto childerror
			}
		}
	}

	// Pass 1: look for fd[i] < i and move those up above len(fd)
	// so that pass 2 won't stomp on an fd it needs later.
	if pipe < nextfd {
		_, _, err1 = RawSyscall(SYS_DUP3, uintptr(pipe), uintptr(nextfd), O_CLOEXEC)
		if err1 != 0 {
			goto childerror
		}
		pipe = nextfd
		nextfd++
	}
	for i = 0; i < len(fd); i++ {
		if fd[i] >= 0 && fd[i] < i {
			if nextfd == pipe { // don't stomp on pipe
				nextfd++
			}
			_, _, err1 = RawSyscall(SYS_DUP3, uintptr(fd[i]), uintptr(nextfd), O_CLOEXEC)
			if err1 != 0 {
				goto childerror
			}
			fd[i] = nextfd
			nextfd++
		}
	}

	// Pass 2: dup fd[i] down onto i.
	for i = 0; i < len(fd); i++ {
		if fd[i] == -1 {
			RawSyscall(SYS_CLOSE, uintptr(i), 0, 0)
			continue
		}
		if fd[i] == i {
			// dup2(i, i) won't clear close-on-exec flag on Linux,
			// probably not elsewhere either.
			_, _, err1 = RawSyscall(fcntl64Syscall, uintptr(fd[i]), F_SETFD, 0)
			if err1 != 0 {
				goto childerror
			}
			continue
		}
		// The new fd is created NOT close-on-exec,
		// which is exactly what we want.
		_, _, err1 = RawSyscall(SYS_DUP3, uintptr(fd[i]), uintptr(i), 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// By convention, we don't close-on-exec the fds we are
	// started with, so if len(fd) < 3, close 0, 1, 2 as needed.
	// Programs that know they inherit fds >= 3 will need
	// to set them close-on-exec.
	for i = len(fd); i < 3; i++ {
		RawSyscall(SYS_CLOSE, uintptr(i), 0, 0)
	}

	// Detach fd 0 from tty
	if sys.Noctty {
		_, _, err1 = RawSyscall(SYS_IOCTL, 0, uintptr(TIOCNOTTY), 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// Set the controlling TTY to Ctty
	if sys.Setctty {
		_, _, err1 = RawSyscall(SYS_IOCTL, uintptr(sys.Ctty), uintptr(TIOCSCTTY), 1)
		if err1 != 0 {
			goto childerror
		}
	}

	// Restore original rlimit.
	if rlimOK && rlim.Cur != 0 {
		rawSetrlimit(RLIMIT_NOFILE, &rlim)
	}

	// Enable tracing if requested.
	// Do this right before exec so that we don't unnecessarily trace the runtime
	// setting up after the fork. See issue #21428.
	if sys.Ptrace {
		_, _, err1 = RawSyscall(SYS_PTRACE, uintptr(PTRACE_TRACEME), 0, 0)
		if err1 != 0 {
			goto childerror
		}
	}

	// 进行execve 系统调用
	// Time to exec.
	_, _, err1 = RawSyscall(SYS_EXECVE,
		uintptr(unsafe.Pointer(argv0)),
		uintptr(unsafe.Pointer(&argv[0])),
		uintptr(unsafe.Pointer(&envv[0])))

childerror:
	// send error code on pipe
	RawSyscall(SYS_WRITE, uintptr(pipe), uintptr(unsafe.Pointer(&err1)), unsafe.Sizeof(err1))
	for {
		RawSyscall(SYS_EXIT, 253, 0, 0)
	}
}

```


## 第三方使用--istio pilot agent 调用 envoy 
```go
// https://github.com/istio/istio/blob/b57e1af32437ce4a340442ce7578826c447a9666/pkg/envoy/proxy.go#L165
func (e *envoy) Run(abort <-chan error) error {
	// spin up a new Envoy process
	args := e.args(e.ConfigPath, istioBootstrapOverrideVar.Get())
	log.Infof("Envoy command: %v", args)

	/* #nosec */
	cmd := exec.Command(e.BinaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if e.AgentIsRoot {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: 1337,
			Gid: 1337,
		}
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-abort:
		log.Warnf("Aborting proxy")
		if errKill := cmd.Process.Kill(); errKill != nil {
			log.Warnf("killing proxy caused an error %v", errKill)
		}
		return err
	case err := <-done:
		return err
	}
}
```


## 参考


- [Linux下Fork与Exec使用](https://www.cnblogs.com/hicjiajia/archive/2011/01/20/1940154.html)
- context 超时后将这个信号发送给子进程以关闭所有子进程: https://github.com/golang/go/issues/53400
- [从源码角度剖析 golang 如何fork一个进程](https://www.cnblogs.com/hobbybear/p/17449737.html)
