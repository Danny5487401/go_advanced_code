<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [testing](#testing)
  - [内置 testing 库中的包类型](#%E5%86%85%E7%BD%AE-testing-%E5%BA%93%E4%B8%AD%E7%9A%84%E5%8C%85%E7%B1%BB%E5%9E%8B)
    - [公共类 testing.common](#%E5%85%AC%E5%85%B1%E7%B1%BB-testingcommon)
    - [T 平常使用的单元测试](#t-%E5%B9%B3%E5%B8%B8%E4%BD%BF%E7%94%A8%E7%9A%84%E5%8D%95%E5%85%83%E6%B5%8B%E8%AF%95)
    - [B 基准测试](#b-%E5%9F%BA%E5%87%86%E6%B5%8B%E8%AF%95)
    - [M 可预置测试前后的操作](#m-%E5%8F%AF%E9%A2%84%E7%BD%AE%E6%B5%8B%E8%AF%95%E5%89%8D%E5%90%8E%E7%9A%84%E6%93%8D%E4%BD%9C)
    - [F 模糊测试](#f-%E6%A8%A1%E7%B3%8A%E6%B5%8B%E8%AF%95)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# testing


## 内置 testing 库中的包类型


在testing包中包含一下结构体:

- testing.T: 这就是我们平常使用的单元测试
- testing.F: 模糊测试, 可以自动生成测试用例
- testing.B: 基准测试. 对函数的运行时间进行统计.
- testing.M: 测试的钩子函数, 可预置测试前后的操作.
- testing.PB: 测试时并行执行.


### 公共类 testing.common


testing.T和testing.B 属于testing.common的扩展。
```go
// common holds the elements common between T and B and
// captures common methods such as Errorf.
type common struct {
	mu          sync.RWMutex         // 读写锁，仅用于控制本数据内的成员访问
	output      []byte               // 存储当前测试产生的日志，每产生一条日志则追加到该切片中，待测试结束后再一并输出.
	w           io.Writer            // 子测试执行结束需要把产生的日志输送到父测试中的output切片中，传递时需要考虑缩进等格式调整，通过w把日志传递到父测试
	ran         bool                 // 仅表示是否已执行过。比如，跟据某个规范筛选测试，如果没有测试被匹配到的话，则common.ran为false，表示没有测试运行过。.
	failed      bool                 // 如果当前测试执行失败，则置为true
	skipped     bool                 // 标记当前测试是否已跳过.
	done        bool                 // 表示当前测试及其子测试已结束，此状态下再执行Fail()之类的方法标记测试状态会产生panic。.
	helperPCs   map[uintptr]struct{} // functions to be skipped when writing file/line info
	helperNames map[string]struct{}  // helperPCs converted to function names
	cleanups    []func()             // optional functions to be called at the end of the test
	cleanupName string               // Name of the cleanup function.
	cleanupPc   []uintptr            // The stack trace at the point where Cleanup was called.
	finished    bool                 // 如果当前测试结束，则置为true.
	inFuzzFn    bool                 // Whether the fuzz target, if this is one, is running.

	chatty         *chattyPrinter // 对应命令行中的-v参数，默认为false，true则打印更多详细日志.
	bench          bool           // Whether the current test is a benchmark.
	hasSub         atomic.Bool    // 标记当前测试是否包含子测试，当测试使用t.Run()方法启动子测试时，t.hasSub则置为1.
	cleanupStarted atomic.Bool    // Registered cleanup callbacks have started to execute
	raceErrors     int            // 竞态检测错误数。.
	runner         string         // 执行当前测试的函数名.
	isParallel     bool           // Whether the test is parallel.

	parent   *common  // 如果当前测试为子测试，则置为父测试的指针
	level    int       // 测试嵌套层数，比如创建子测试时，子测试嵌套层数就会加1.
	creator  []uintptr // 测试函数调用栈
	name     string    // 记录每个测试函数名，比如测试函数TestAdd(t *testing.T), 其中t.name即“TestAdd”。 测试结束，打印测试结果会用到该成员.
	start    time.Time // 记录测试开始的时间
	duration time.Duration // 记录测试所花费的时间。
	barrier  chan bool // 用于控制父测试和子测试执行的channel，如果测试为Parallel，则会阻塞等待父测试结束后再继续。
	signal   chan bool // 通知当前测试结束。
	sub      []*T      // 子测试列表。

	tempDirMu  sync.Mutex
	tempDir    string
	tempDirErr error
	tempDirSeq int32
}
```


### T 平常使用的单元测试

```go
// go1.20/src/testing/testing.go
type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type testContext struct {
	match    *matcher // 匹配器，用于管理测试名称匹配、过滤等
	deadline time.Time

	// isFuzzing is true in the context used when generating random inputs
	// for fuzz targets. isFuzzing is false when running normal tests and
	// when running fuzz tests as unit tests (without -fuzz or when -fuzz
	// does not match).
	isFuzzing bool

	mu sync.Mutex

	// 用于通知测试可以并发执行的控制管道，测试并发达到最大限制时，需要阻塞等待该管道的通知事件
	startParallel chan bool

	// 当前并发执行的测试个数
	running int

	// 等待并发执行的测试个数，所有等待执行的测试都阻塞在startParallel管道处；
	numWaiting int

	// 最大并发数，默认为系统CPU数
	maxParallel int
}
```


### B 基准测试

```go
type B struct {
	common
	importPath       string // import path of the package containing the benchmark
	context          *benchContext
	N                int
	previousN        int           // number of iterations in the previous run
	previousDuration time.Duration // total duration of the previous run
	benchFunc        func(b *B)
	benchTime        durationOrCountFlag
	bytes            int64
	missingBytes     bool // one of the subbenchmarks does not have bytes set.
	timerOn          bool
	showAllocResult  bool
	result           BenchmarkResult
	parallelism      int // RunParallel creates parallelism*GOMAXPROCS goroutines
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64
	startBytes  uint64
	// The net total of this test after being run.
	netAllocs uint64
	netBytes  uint64
	// Extra metrics collected by ReportMetric.
	extra map[string]float64
}
```


### M 可预置测试前后的操作

Go 1.14 版本引进了 TestMain 的能力，当需要针对一个单测或benchmark执行一些额外的操作（比如资源的创建/回收），可以使用 testing 支持的 Main 函数机制。

```go
// M is a type passed to a TestMain function to run the actual tests.
type M struct {
	deps        testDeps
	tests       []InternalTest // 单元测试
	benchmarks  []InternalBenchmark // 性能测试
	fuzzTargets []InternalFuzzTarget
	examples    []InternalExample

	timer     *time.Timer
	afterOnce sync.Once

	numRun int

	// value to pass to os.Exit, the outer test func main
	// harness calls os.Exit with this code. See #34129.
	exitCode int
}
```







### F 模糊测试

```go
type F struct {
	common
	fuzzContext *fuzzContext
	testContext *testContext

	// inFuzzFn is true when the fuzz function is running. Most F methods cannot
	// be called when inFuzzFn is true.
	inFuzzFn bool

	// corpus is a set of seed corpus entries, added with F.Add and loaded
	// from testdata.
	corpus []corpusEntry

	result     fuzzResult
	fuzzCalled bool
}
```

开启条件 ： go test with the -fuzz flag

## 参考

- [Jay Conrod 关于 Internals of Go's new fuzzing system](https://jayconrod.com/posts/123/internals-of-go-s-new-fuzzing-system)