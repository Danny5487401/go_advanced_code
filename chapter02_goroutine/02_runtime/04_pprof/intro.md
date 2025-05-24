<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [pprof](#pprof)
  - [内置库:两个包](#%E5%86%85%E7%BD%AE%E5%BA%93%E4%B8%A4%E4%B8%AA%E5%8C%85)
  - [使用](#%E4%BD%BF%E7%94%A8)
    - [Go test 使用 pprof](#go-test-%E4%BD%BF%E7%94%A8-pprof)
    - [pprof 文件分析](#pprof-%E6%96%87%E4%BB%B6%E5%88%86%E6%9E%90)
  - [net/http/pprof 源码分析](#nethttppprof-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [/debug/pprof/profile 路由](#debugpprofprofile-%E8%B7%AF%E7%94%B1)
    - [/debug/pprof/ 路由](#debugpprof-%E8%B7%AF%E7%94%B1)
  - [runtime/pprof](#runtimepprof)
    - [heap](#heap)
    - [cpu](#cpu)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# pprof

pprof 是用于可视化和分析性能分析数据的工具.

pprof 以 profile.proto 读取分析样本的集合，并生成报告以可视化并帮助分析数据（支持文本和图形报告）

profile.proto 是一个 Protocol Buffer v3 的描述文件，它描述了一组 callstack 和 symbolization 信息， 作用是表示统计分析的一组采样的调用栈，是很常见的 stacktrace 配置文件格式



## 内置库:两个包
1. net/http/pprof
   使用场景：在线服务（一直运行着的程序）

2. runtime/pprof
   使用场景：工具型应用（比如说定制化的分析小工具、集成到公司监控系统）

这两个包都是可以监控代码性能的， 只不过net/http/pprof 是通过http端口方式暴露出来的，内部封装的仍然是runtime/pprof。

## 使用

![use-pprof.png](use-pprof.png)

```shell
(⎈|kubeasz-test:descheduler)➜  ~ go tool pprof --help
usage:

Produce output in the specified format.

   pprof <format> [options] [binary] <source> ...

Omit the format to get an interactive shell whose commands can be used
to generate various views of a profile

   pprof [options] [binary] <source> ...

Omit the format and provide the "-http" flag to get an interactive web
interface at the specified host:port that can be used to navigate through
various views of a profile.

   pprof -http [host]:[port] [options] [binary] <source> ...

Details:
  Output formats (select at most one):
    -callgrind       Outputs a graph in callgrind format
    -comments        Output all profile comments
    -disasm          Output assembly listings annotated with samples
    -dot             Outputs a graph in DOT format
    -eog             Visualize graph through eog
    -evince          Visualize graph through evince
    -gif             Outputs a graph image in GIF format
    -gv              Visualize graph through gv
    -kcachegrind     Visualize report in KCachegrind
    -list            Output annotated source for functions matching regexp
    -pdf             Outputs a graph in PDF format
    -peek            Output callers/callees of functions matching regexp
    -png             Outputs a graph image in PNG format
    -proto           Outputs the profile in compressed protobuf format
    -ps              Outputs a graph in PS format
    -raw             Outputs a text representation of the raw profile
    -svg             Outputs a graph in SVG format
    -tags            Outputs all tags in the profile
    -text            Outputs top entries in text form
    -top             Outputs top entries in text form
    -topproto        Outputs top entries in compressed protobuf format
    -traces          Outputs all profile samples in text form
    -tree            Outputs a text rendering of call graph
    -web             Visualize graph through web browser
    -weblist         Display annotated source in a web browser

  Options:
    -call_tree       Create a context-sensitive call tree
    -compact_labels  Show minimal headers
    -divide_by       Ratio to divide all samples before visualization
    -drop_negative   Ignore negative differences
    -edgefraction    Hide edges below <f>*total
    -focus           Restricts to samples going through a node matching regexp
    -hide            Skips nodes matching regexp
    -ignore          Skips paths going through any nodes matching regexp
    -intel_syntax    Show assembly in Intel syntax
    -mean            Average sample value over first value (count)
    -nodecount       Max number of nodes to show
    -nodefraction    Hide nodes below <f>*total
    -noinlines       Ignore inlines.
    -normalize       Scales profile based on the base profile.
    -output          Output filename for file-based outputs
    -prune_from      Drops any functions below the matched frame.
    -relative_percentages Show percentages relative to focused subgraph
    -sample_index    Sample value to report (0-based index or name)
    -show            Only show nodes matching regexp
    -show_from       Drops functions above the highest matched frame.
    -source_path     Search path for source files
    -tagfocus        Restricts to samples with tags in range or matched by regexp
    -taghide         Skip tags matching this regexp
    -tagignore       Discard samples with tags in range or matched by regexp
    -tagleaf         Adds pseudo stack frames for labels key/value pairs at the callstack leaf.
    -tagroot         Adds pseudo stack frames for labels key/value pairs at the callstack root.
    -tagshow         Only consider tags matching this regexp
    -trim            Honor nodefraction/edgefraction/nodecount defaults
    -trim_path       Path to trim from source paths before search
    -unit            Measurement units to display

  Option groups (only set one per group):
    granularity
      -functions       Aggregate at the function level.
      -filefunctions   Aggregate at the function level.
      -files           Aggregate at the file level.
      -lines           Aggregate at the source code line level.
      -addresses       Aggregate at the address level.
    sort
      -cum             Sort entries based on cumulative weight
      -flat            Sort entries based on own weight

  Source options:
    -seconds              Duration for time-based profile collection
    -timeout              Timeout in seconds for profile collection
    -buildid              Override build id for main binary
    -add_comment          Free-form annotation to add to the profile
                          Displayed on some reports or with pprof -comments
    -diff_base source     Source of base profile for comparison
    -base source          Source of base profile for profile subtraction
    profile.pb.gz         Profile in compressed protobuf format
    legacy_profile        Profile in legacy pprof format
    http://host/profile   URL for profile handler to retrieve
    -symbolize=           Controls source of symbol information
      none                  Do not attempt symbolization
      local                 Examine only local binaries
      fastlocal             Only get function names from local binaries
      remote                Do not examine local binaries
      force                 Force re-symbolization
    Binary                  Local path or build id of binary for symbolization
    -tls_cert             TLS client certificate file for fetching profile and symbols
    -tls_key              TLS private key file for fetching profile and symbols
    -tls_ca               TLS CA certs file for fetching profile and symbols

  Misc options:
   -http              Provide web interface at host:port.
                      Host is optional and 'localhost' by default.
                      Port is optional and a randomly available port by default.
   -no_browser        Skip opening a browser for the interactive web UI.
   -tools             Search path for object tools

  Legacy convenience options:
   -inuse_space           Same as -sample_index=inuse_space
   -inuse_objects         Same as -sample_index=inuse_objects
   -alloc_space           Same as -sample_index=alloc_space
   -alloc_objects         Same as -sample_index=alloc_objects
   -total_delay           Same as -sample_index=delay
   -contentions           Same as -sample_index=contentions
   -mean_delay            Same as -mean -sample_index=delay

  Environment Variables:
   PPROF_TMPDIR       Location for saved profiles (default $HOME/pprof)
   PPROF_TOOLS        Search path for object-level tools
   PPROF_BINARY_PATH  Search path for local binary files
                      default: $HOME/pprof/binaries
                      searches $buildid/$name, $buildid/*, $path/$buildid,
                      ${buildid:0:2}/${buildid:2}.debug, $name, $path
   * On Windows, %USERPROFILE% is used instead of $HOME
```

- inuse_space：分析应用程序的常驻内存占用情况。
- inuse_object: 正在使用，尚未释放的对象
- alloc_space：分析应用程序的内存临时分配情况,所有分配的空间，包含已释放的。
- alloc_objects: 所有分配的对象，包含已释放的

- http: 开启web 服务


```shell
// 分析
go tool pprof -http=:8082 http://172.16.7.33:30901/debug/pprof/allocs\?debug\=1
```


### Go test 使用 pprof

Golang在运行测试用例或压测时也可以通过添加参加输出测试过程中的CPU、内存和trace情况


### pprof 文件分析

pprof 文件是二进制的，不是给人读的，需要翻译一下，而 golang 原生就给我们提供了分析工具，直接执行下面命令即可，会生成一张很直观的 svg 图片，
直接用 chrome 就可以打开，当然也可以生成别的格式（pdf，png 都可以），可以用 go tool pprof -h 命令查看支持的输出类型





## net/http/pprof 源码分析
```go
// net/http/pprof/pprof.go
func init() {
	http.HandleFunc("/debug/pprof/", Index)
	http.HandleFunc("/debug/pprof/cmdline", Cmdline)
	http.HandleFunc("/debug/pprof/profile", Profile)
	http.HandleFunc("/debug/pprof/symbol", Symbol)
	http.HandleFunc("/debug/pprof/trace", Trace)
}
```
1. 
直接使用如下命令，则不需要通过点击浏览器上的链接就能进入命令行交互模式：
```go
go tool pprof http://x.x.x.x:8080/debug/pprof/profile

```

当然也是需要先后台采集一段时间的数据，再将数据文件下载到本地，最后进行分析。上述的 Url 后面还可以带上时间参数：?seconds=60，自定义 CPU Profiling 的时长。


类似的命令还有：
```shell
# 下载 cpu profile，默认从当前开始收集 30s 的 cpu 使用情况，需要等待 30s
go tool pprof http://127.0.0.1:8080/debug/pprof/profile
# wait 120s
go tool pprof http://127.0.0.1:8080/debug/pprof/profile?seconds=120     

# 下载 heap profile
go tool pprof http://127.0.0.1:8080/debug/pprof/heap

# 下载 goroutine profile
go tool pprof http://127.0.0.1:8080/debug/pprof/goroutine

# 下载 block profile
go tool pprof http://127.0.0.1:8080/debug/pprof/block

# 下载 mutex profile
go tool pprof http://127.0.0.1:8080/debug/pprof/mutex
```

### /debug/pprof/profile 路由
```go
func Profile(w http.ResponseWriter, r *http.Request) {
    // ...
	if err := pprof.StartCPUProfile(w); err != nil {
		// StartCPUProfile failed, so no writes yet.
		serveError(w, http.StatusInternalServerError,
			fmt.Sprintf("Could not enable CPU profiling: %s", err))
		return
	}
	sleep(w, time.Duration(sec)*time.Second)
	pprof.StopCPUProfile()
}
```
这个函数也是调用runtime/pprof的StartCPUProfile(w)方法开始 CPU profiling，然后睡眠一段时间（这个时间就是采样间隔），最后调用pprof.StopCPUProfile()停止采用.

StartCPUProfile()方法传入的是http.ResponseWriter类型变量，所以采样结果直接写回到 HTTP 的客户端

### /debug/pprof/ 路由


```go
func Index(w http.ResponseWriter, r *http.Request) {
	if name, found := strings.CutPrefix(r.URL.Path, "/debug/pprof/"); found {
		if name != "" {
			handler(name).ServeHTTP(w, r)
			return
		}
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var profiles []profileEntry
	for _, p := range pprof.Profiles() {
		profiles = append(profiles, profileEntry{
			Name:  p.Name(),
			Href:  p.Name(),
			Desc:  profileDescriptions[p.Name()],
			Count: p.Count(),
		})
	}

	// Adding other profiles exposed from within this package
	for _, p := range []string{"cmdline", "profile", "trace"} {
		profiles = append(profiles, profileEntry{
			Name: p,
			Href: p,
			Desc: profileDescriptions[p],
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	if err := indexTmplExecute(w, profiles); err != nil {
		log.Print(err)
	}
}

```


![](.intro_images/pprof_args.png)

```go
// go1.23.0/src/net/http/pprof/pprof.go

// 展示
var profileDescriptions = map[string]string{
	"allocs":       "A sampling of all past memory allocations", // 采样某段时间内所有分配过的内存
	"block":        "Stack traces that led to blocking on synchronization primitives", // 追踪阻塞在同步原语(Mutex、RWMutex、WaitGroup、Cond、chan等)处的堆栈
	"cmdline":      "The command line invocation of the current program", // 当前程序运行时执行的命令行
	"goroutine":    "Stack traces of all current goroutines. Use debug=2 as a query parameter to export in the same format as an unrecovered panic.", // 追踪当前运行的所有协程堆栈。
	"heap":         "A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.", // 采样当前存活对象的内存分配信息
	"mutex":        "Stack traces of holders of contended mutexes",
	"profile":      "CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.",
	"threadcreate": "Stack traces that led to the creation of new OS threads",
	"trace":        "A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.",
}
```

allocs 和 heap 采样的信息一致，不过前者是所有对象的内存分配，而 heap 则是活跃对象的内存分配


- CPU profiling（CPU 性能分析）：这是最常使用的一种类型。用于分析函数或方法的执行耗时；
- Memory profiling：这种类型也常使用。用于分析程序的内存占用情况；
- Block profiling：这是 Go 独有的，用于记录 goroutine 在等待共享资源花费的时间；
- Goroutine Profiling: 报告goroutines的使用情况，有哪些 goroutines，它们的调用关系是怎样的。
- Mutex profiling：与 Block profiling 类似，但是只记录因为锁竞争导致的等待或延迟。

## runtime/pprof

```go
// go1.23.0/src/runtime/pprof/pprof.go

func lockProfiles() {
	profiles.mu.Lock()
	if profiles.m == nil {
		// Initial built-in profiles.
		profiles.m = map[string]*Profile{
          "goroutine":    goroutineProfile,  //显示当前所有协程的堆栈信息
          "threadcreate": threadcreateProfile, // 系统线程创建情况的采样信息
          "heap":         heapProfile,  // 堆上的内存分配情况的采样信息
          "allocs":       allocsProfile,  //内存分配情况的采样信息
          "block":        blockProfile,  //阻塞操作情况的采样信息
          "mutex":        mutexProfile,  // 锁竞争情况的采样信息
        }
	}
}


// profiles records all registered profiles.
var profiles struct {
	mu sync.Mutex
	m  map[string]*Profile
}

var goroutineProfile = &Profile{
	name:  "goroutine",
	count: countGoroutine,
	write: writeGoroutine,
}

var threadcreateProfile = &Profile{
	name:  "threadcreate",
	count: countThreadCreate,
	write: writeThreadCreate,
}

var heapProfile = &Profile{
	name:  "heap",
	count: countHeap,
	write: writeHeap,
}

var allocsProfile = &Profile{
	name:  "allocs",
	count: countHeap, // identical to heap profile
	write: writeAlloc,
}

var blockProfile = &Profile{
	name:  "block",
	count: countBlock,
	write: writeBlock,
}

var mutexProfile = &Profile{
	name:  "mutex",
	count: countMutex,
	write: writeMutex,
}

```

默认情况下是不追踪block和mutex的信息的，如果想要看这两个信息，需要在代码中加上两行
```go
runtime.SetBlockProfileRate(1) // 开启对阻塞操作的跟踪，block  
runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪，mutex
```



### heap 
内存概要文件用于保存在用户程序执行期间的内存使用情况。这里所说的内存使用情况，其实就是程序运行过程中堆内存的分配情况。


主要结构体
```go
type MemStats struct {
    // 常规统计信息（General statistics）
    Alloc      uint64 // 已申请且仍在使用的字节数
    TotalAlloc uint64 // 已申请的总字节数（已释放的部分也算在内）
    Sys        uint64 // 从系统中获取的字节数（下面XxxSys之和）
    Lookups    uint64 // 指针查找的次数
    Mallocs    uint64 // 申请内存的次数
    Frees      uint64 // 释放内存的次数
	
    // 主分配堆统计
    HeapAlloc    uint64 // 已申请且仍在使用的字节数
    HeapSys      uint64 // 从系统中获取的字节数
    HeapIdle     uint64 // 闲置span中的字节数
    HeapInuse    uint64 // 非闲置span中的字节数
    HeapReleased uint64 // 释放到系统的字节数
    HeapObjects  uint64 // 已分配对象的总个数
	
	// 栈内存统计（Stack memory statistics）
    // L低层次、大小固定的结构体分配器统计，Inuse为正在使用的字节数，Sys为从系统获取的字节数
    StackInuse  uint64 // 引导程序的堆栈
    StackSys    uint64
    MSpanInuse  uint64 // mspan结构体
    MSpanSys    uint64
	
	// 堆外内存统计信息（Off-heap memory statistics）
    MCacheInuse uint64 // mcache结构体
    MCacheSys   uint64
    BuckHashSys uint64 // profile桶散列表
    GCSys       uint64 // GC元数据
    OtherSys    uint64 // 其他系统申请
	
    // 垃圾收集器统计
    NextGC       uint64 // 会在HeapAlloc字段到达该值（字节数）时运行下次GC
    LastGC       uint64 // 上次运行的绝对时间（纳秒）
    PauseTotalNs uint64
    PauseNs      [256]uint64 // 近期GC暂停时间的循环缓冲，最近一次在[(NumGC+255)%256]
    NumGC        uint32
    EnableGC     bool
    DebugGC      bool
	
	// 按 per-size class 大小分配统计
    // 每次申请的字节数的统计，61是C代码中的尺寸分级数
    BySize [61]struct {
        Size    uint32
        Mallocs uint64
        Frees   uint64
    }
}
```


```go
// go1.23.0/src/runtime/pprof/pprof.go

func writeHeap(w io.Writer, debug int) error {
	return writeHeapInternal(w, debug, "")
}


func writeHeapInternal(w io.Writer, debug int, defaultSampleType string) error {
	var memStats *runtime.MemStats
	if debug != 0 {
		// Read mem stats first, so that our other allocations
		// do not appear in the statistics.
		memStats = new(runtime.MemStats)
		runtime.ReadMemStats(memStats)
	}

	// Find out how many records there are (the call
	// pprof_memProfileInternal(nil, true) below),
	// allocate that many records, and get the data.
	// There's a race—more records might be added between
	// the two calls—so allocate a few extra records for safety
	// and also try again if we're very unlucky.
	// The loop should only execute one iteration in the common case.
	var p []profilerecord.MemProfileRecord
	n, ok := pprof_memProfileInternal(nil, true)
	for {
		// Allocate room for a slightly bigger profile,
		// in case a few more entries have been added
		// since the call to MemProfile.
		p = make([]profilerecord.MemProfileRecord, n+50)
		n, ok = pprof_memProfileInternal(p, true)
		if ok {
			p = p[0:n]
			break
		}
		// Profile grew; try again.
	}

	if debug == 0 {
		return writeHeapProto(w, p, int64(runtime.MemProfileRate), defaultSampleType)
	}

	slices.SortFunc(p, func(a, b profilerecord.MemProfileRecord) int {
		return cmp.Compare(a.InUseBytes(), b.InUseBytes())
	})

	b := bufio.NewWriter(w)
	tw := tabwriter.NewWriter(b, 1, 8, 1, '\t', 0)
	w = tw

	var total runtime.MemProfileRecord
	for i := range p {
		r := &p[i]
		total.AllocBytes += r.AllocBytes
		total.AllocObjects += r.AllocObjects
		total.FreeBytes += r.FreeBytes
		total.FreeObjects += r.FreeObjects
	}

	// Technically the rate is MemProfileRate not 2*MemProfileRate,
	// but early versions of the C++ heap profiler reported 2*MemProfileRate,
	// so that's what pprof has come to expect.
	rate := 2 * runtime.MemProfileRate

	// pprof reads a profile with alloc == inuse as being a "2-column" profile
	// (objects and bytes, not distinguishing alloc from inuse),
	// but then such a profile can't be merged using pprof *.prof with
	// other 4-column profiles where alloc != inuse.
	// The easiest way to avoid this bug is to adjust allocBytes so it's never == inuseBytes.
	// pprof doesn't use these header values anymore except for checking equality.
	inUseBytes := total.InUseBytes()
	allocBytes := total.AllocBytes
	if inUseBytes == allocBytes {
		allocBytes++
	}

	fmt.Fprintf(w, "heap profile: %d: %d [%d: %d] @ heap/%d\n",
		total.InUseObjects(), inUseBytes,
		total.AllocObjects, allocBytes,
		rate)

	for i := range p {
		r := &p[i]
		fmt.Fprintf(w, "%d: %d [%d: %d] @",
			r.InUseObjects(), r.InUseBytes(),
			r.AllocObjects, r.AllocBytes)
		for _, pc := range r.Stack {
			fmt.Fprintf(w, " %#x", pc)
		}
		fmt.Fprintf(w, "\n")
		printStackRecord(w, r.Stack, false)
	}

	// Print memstats information too.
	// Pprof will ignore, but useful for people
	s := memStats
    
	// 打印...

	// Also flush out MaxRSS on supported platforms.
	addMaxRSS(w)

	tw.Flush()
	return b.Flush()
}

```

读取内存信息 

runtime.ReadMemStats方法是需要stw的，所以尽量不要在线上调用
```go
// go1.23.0/src/runtime/mstats.go

func ReadMemStats(m *MemStats) {
	_ = m.Alloc // nil check test before we switch stacks, see issue 61158
	stw := stopTheWorld(stwReadMemStats)

	systemstack(func() {
		readmemstats_m(m)
	})

	startTheWorld(stw)
}
```




### cpu 

在默认情况下，Go语言的运行时系统会以100 Hz的的频率对CPU使用情况进行取样。

```go
func StartCPUProfile(w io.Writer) error {
	// The runtime routines allow a variable profiling rate,
	// but in practice operating systems cannot trigger signals
	// at more than about 500 Hz, and our processing of the
	// signal is not cheap (mostly getting the stack trace).
	// 100 Hz is a reasonable choice: it is frequent enough to
	// produce useful data, rare enough not to bog down the
	// system, and a nice round number to make it easy to
	// convert sample counts to seconds. Instead of requiring
	// each client to specify the frequency, we hard code it.
	const hz = 100

	cpu.Lock()
	defer cpu.Unlock()
	if cpu.done == nil {
		cpu.done = make(chan bool)
	}
	// Double-check.
	if cpu.profiling {
		return fmt.Errorf("cpu profiling already in use")
	}
	cpu.profiling = true
	runtime.SetCPUProfileRate(hz)
	go profileWriter(w)
	return nil
}

```









## 参考
- [pprof使用详解和源码分析](https://zhuanlan.zhihu.com/p/666945970)
- [万字长文讲解Golang pprof 的使用](https://juejin.cn/post/7343428554686611495)
- [7.8 内存统计](https://golang.design/under-the-hood/zh-cn/part2runtime/ch07alloc/mstats/)




