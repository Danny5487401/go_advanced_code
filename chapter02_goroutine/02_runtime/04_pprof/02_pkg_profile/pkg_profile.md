<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/pkg/profile](#githubcompkgprofile)
  - [使用](#%E4%BD%BF%E7%94%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/pkg/profile


依赖 github.com/felixge/fgprof

## 使用

```go
	// 开始性能分析, 返回一个停止接口, 只能是一个
	stopper := profile.Start(profile.CPUProfile, profile.ProfilePath("."))

	// 在main()结束时停止性能分析
	defer stopper.Stop()
```



选项
```go
const (
	cpuMode = iota
	memMode
	mutexMode
	blockMode
	traceMode
	threadCreateMode
	goroutineMode
	clockMode
)


// cpu 采样
func CPUProfile(p *Profile) { p.mode = cpuMode }

// DefaultMemProfileRate is the default memory profiling rate.
// See also http://golang.org/pkg/runtime/#pkg-variables
const DefaultMemProfileRate = 4096

// 内存采样
func MemProfile(p *Profile) {
	p.memProfileRate = DefaultMemProfileRate
	p.mode = memMode
}
```


启动

```go
func Start(options ...func(*Profile)) interface {
	Stop()
} {
	// 只能启动一次
	if !atomic.CompareAndSwapUint32(&started, 0, 1) {
		log.Fatal("profile: Start() already called")
	}

	var prof Profile
	for _, option := range options {
		option(&prof)
	}

	path, err := func() (string, error) {
		if p := prof.path; p != "" {
			return p, os.MkdirAll(p, 0777)
		}
		return ioutil.TempDir("", "profile")
	}()

	if err != nil {
		log.Fatalf("profile: could not create initial output directory: %v", err)
	}

	logf := func(format string, args ...interface{}) {
		if !prof.quiet {
			log.Printf(format, args...)
		}
	}

	if prof.memProfileType == "" {
		prof.memProfileType = "heap"
	}

	switch prof.mode {
	case cpuMode: // cpu 
		fn := filepath.Join(path, "cpu.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create cpu profile %q: %v", fn, err)
		}
		logf("profile: cpu profiling enabled, %s", fn)
		pprof.StartCPUProfile(f)
		prof.closer = func() {
			pprof.StopCPUProfile()
			f.Close()
			logf("profile: cpu profiling disabled, %s", fn)
		}

	case memMode: // 内存
		fn := filepath.Join(path, "mem.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create memory profile %q: %v", fn, err)
		}
		old := runtime.MemProfileRate
		runtime.MemProfileRate = prof.memProfileRate
		logf("profile: memory profiling enabled (rate %d), %s", runtime.MemProfileRate, fn)
		prof.closer = func() {
			pprof.Lookup(prof.memProfileType).WriteTo(f, 0)
			f.Close()
			runtime.MemProfileRate = old
			logf("profile: memory profiling disabled, %s", fn)
		}

	case mutexMode: // 锁
		fn := filepath.Join(path, "mutex.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create mutex profile %q: %v", fn, err)
		}
		runtime.SetMutexProfileFraction(1)
		logf("profile: mutex profiling enabled, %s", fn)
		prof.closer = func() {
			if mp := pprof.Lookup("mutex"); mp != nil {
				mp.WriteTo(f, 0)
			}
			f.Close()
			runtime.SetMutexProfileFraction(0)
			logf("profile: mutex profiling disabled, %s", fn)
		}

	case blockMode:
		fn := filepath.Join(path, "block.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create block profile %q: %v", fn, err)
		}
		runtime.SetBlockProfileRate(1)
		logf("profile: block profiling enabled, %s", fn)
		prof.closer = func() {
			pprof.Lookup("block").WriteTo(f, 0)
			f.Close()
			runtime.SetBlockProfileRate(0)
			logf("profile: block profiling disabled, %s", fn)
		}

	case threadCreateMode:
		fn := filepath.Join(path, "threadcreation.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create thread creation profile %q: %v", fn, err)
		}
		logf("profile: thread creation profiling enabled, %s", fn)
		prof.closer = func() {
			if mp := pprof.Lookup("threadcreate"); mp != nil {
				mp.WriteTo(f, 0)
			}
			f.Close()
			logf("profile: thread creation profiling disabled, %s", fn)
		}

	case traceMode:
		fn := filepath.Join(path, "trace.out")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create trace output file %q: %v", fn, err)
		}
		if err := trace.Start(f); err != nil {
			log.Fatalf("profile: could not start trace: %v", err)
		}
		logf("profile: trace enabled, %s", fn)
		prof.closer = func() {
			trace.Stop()
			logf("profile: trace disabled, %s", fn)
		}

	case goroutineMode:
		fn := filepath.Join(path, "goroutine.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create goroutine profile %q: %v", fn, err)
		}
		logf("profile: goroutine profiling enabled, %s", fn)
		prof.closer = func() {
			if mp := pprof.Lookup("goroutine"); mp != nil {
				mp.WriteTo(f, 0)
			}
			f.Close()
			logf("profile: goroutine profiling disabled, %s", fn)
		}

	case clockMode:
		fn := filepath.Join(path, "clock.pprof")
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create clock profile %q: %v", fn, err)
		}
		logf("profile: clock profiling enabled, %s", fn)
		stop := fgprof.Start(f, fgprof.FormatPprof)
		prof.closer = func() {
			stop()
			f.Close()
			logf("profile: clock profiling disabled, %s", fn)
		}
	}

	if !prof.noShutdownHook {
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c

			log.Println("profile: caught interrupt, stopping profiles")
			prof.Stop()

			os.Exit(0)
		}()
	}

	return &prof
}

```



