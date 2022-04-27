# Time

## 三个比较常用的数据结

- time.Time
- time.Duration
- time.C

## 两个重要概念:单调时间和壁挂时间
- wall time壁挂时间:
  挂钟时间，实际上就是指的是现实的时间，这是由变量xtime来记录的。系统每次启动时将CMOS上的RTC时间读入xtime，
  这个值是"自1970-01-01起经历的秒数、本秒中经历的纳秒数"，每来一个timer interrupt，也需要去更新xtime。，常见的场景就是通过修改时间延长收费软件的试用期。
- monotonic time单调时间:
  是单调时间,实际它指的是系统启动以后流逝的时间,这是由变量jiffies记录系统每次启动时jiffies初始化为0，
  每来一个timer interrupt，jiffies加1，也就是说它代表系统启动后流逝的tick数。jiffies一定是单调递增的，因为时间不够逆

CLOCK_MONOTONIC是monotonic time;CLOCK_REALTIME是wall time。

### 应用
在操作系统中如果需要 显示 时间的时候，会使用 wall time ，而需要 测量 时间的时候，会使用 monotonic time 

### wall time 相关函数
- time.Since(start)
- time.Until(deadline)
- time.Now().Before(deadline)

### monotonic time 相关函数
time.Since(start)
time.Until(deadline)
time.Now().Before(deadline)

### CLOCK_MONOTONIC 和 CLOCK_REALTIME




### monotonic time Vs. wall time
wall time不一定单调递增的。wall time是指现实中的实际时间，如果系统要与网络中某个节点时间同步、或者由系统管理员觉得这个wall time与现实时间不一致，
有可能任意的改变这个wall time。最简单的例子是，我们用户可以去任意修改系统时间，这个被修改的时间应该就是wall time，即xtime，它甚至可以被写入RTC而永久保存。

### Time时注意事项
- Time能够代表纳秒精度的时间。
- 因为Time并非并发安全,所以在存储或者传递的时候,都应该使用值引用。
- 在Go中, == 运算符不仅仅比较时刻,还会比较Location以及单调时钟,因此在不保证所有时间设置为相同的位置的时候,不应该将time.Time作为map或者database的健。如果必须要使用,应该通过UTC或者Local方法将单调时间剥离。

## time源码结构
```go
// src/time/time.go

// 1. Time 能够代表纳秒精度的时间
// 2. 因为 Time 并非并发安全，所以在存储或者是传递的时候，都应该使用值引用。
// 3. 在 Go 中， == 运算符不仅仅会比较时刻，还会比较 Location 以及单调时钟，
//    因此在不保证所有时间设置为相同的位置的时候，不应将 time.Time 
//    作为 map 或者 database 的键。如果必须要使用，应先通过 UTC 或者 Local 
//    方法将单调时钟剥离。
//
type Time struct {
    // wall 和 ext 字段共同组成 wall time 秒级、纳秒级，monotonic time 纳秒级
    // 的时间精度，先看下 wall 这个 无符号64位整数 的结构。
    //
    //          +------------------+--------------+--------------------+
    // wall =>  | 1 (hasMonotonic) | 33 (second)  |  30 (nanosecond)   |
    //          +------------------+--------------+--------------------+
    // 
	// 
    // 所以 wall 字段会有两种情况，分别为
    // 1. 当 wall 字段的 hasMonotonic 为 0 时，second 位也全部为 0，ext 字段会存储
    //    自1年1月1日以来开始的秒级精度时间作为 wall time 。
    // 2. 当 wall 字段的 hasMonotonic 为 1 时，second 位会存储从 1885-1-1 开始的秒
    //    级精度时间作为 wall time，并且 ext 字段会存储从操作系统启动后的纳秒级精度时间
    //    作为 monotonic time ,这是大部分情况
    wall uint64
    ext  int64
    
    // Location 作为当前时间的时区，可用于确定时间是否处在正确的位置上。
    // 当 loc 为 nil 时，则表示为 UTC 时间。
    // 因为北京时区为东八区，比 UTC 时间要领先 8 个小时，
    // 所以我们获取到的时间默认会记为 +0800
    loc *Location
}
```

首先要wall提供一个简单的“挂钟”读数值，并以单调时钟的形式ext提供此扩展信息。

wall它: hasMonotonic的最高位包含一个1位标志。然后是33位用于跟踪秒数；最后是30位，用于跟踪纳秒，范围为[0，999999999]

对于Go> = 1.9，该hasMonotonic标志始终处于启用状态，其日期介于1885年至2157年之间，但是由于兼容性承诺以及极端情况，Go还要确保正确处理这些时间值


存储机制
```go
// src/time/time.go

// 定义的一些常量
const (
    // 代表无符号64位整数的首位
    hasMonotonic = 1 << 63
    
    // maxWall 和 minWall 是指 hasMonotonic 为 1 的情况下，
    // wall time 的最大以及最小的时间范围。
    maxWall      = wallToInternal + (1<<33 - 1) // year 2157
    minWall      = wallToInternal               // year 1885
    
    // 纳秒位位置的辅助常量
    nsecMask     = 1<<30 - 1
    nsecShift    = 30
)

// 注意：以下方法均为包内辅助方法，会通过指针接收器进行操作，减轻调用负担，
// 但我们在使用 time.Time 时应该尽量避免使用指针，以免出现竞态争用。

// sec 返回时间的秒数
func (t *Time) sec() int64 {
    if t.wall&hasMonotonic != 0 {
        // hasMonotonic 为 1，则 second 位记录从 1885-1-1 开始的秒数
        // 则返回值为以下两个值相加：
        // wallToInternal = 1885-1-1 00:00:00 的秒数 = 59453308800
        // t.wall<<1>>(nsecShift+1) = 1885-1-1 00:00:00 到现在的秒数
        return wallToInternal + int(t.wall<<1>>(nsecShift+1))
    }
    return int64(t.ext)  // hasMonotonic 为 0 ，返回 ext 为 wall time 秒数
}

// addSec 在当前时间基础上加上 d 秒
func (t *Time) addSec(d int64) {
    if t.wall&hasMonotonic != 0 {
        // 同上，获取当前秒数
        sec := int64(t.wall << 1 >> (nsecShift + 1))
        dsec := sec + d
        if 0 <= dsec && dsec <= 1<<33-1 { // 判断 wall 的 second 不会溢出
            t.wall = t.wall&nsecMask | uint64(dsec)<<nsecShift | hasMonotonic
            return
        }
        // second 位已经不足以存下了 wall time 的秒数，需要去掉单调时钟，并
        // 其移动到 ext 字段中，移动完成后，执行下面的 t.ext += d 语句即可
        t.stripMono()
    }
    
    // 如果 hasMonotonic 为 0，直接就在 ext 字段上面添加就好了
    t.ext += d
}

// stripMono 去除单调时钟
func (t *Time) stripMono() {
    if t.wall&hasMonotonic != 0 {
        t.ext = t.sec()
        t.wall &= nsecMask
    }
}
```
### 时间比较
```go
// src/time/time.go

// 判断 t 时间是否晚于 u 时间
func (t Time) After(u Time) bool {
    if t.wall&u.wall&hasMonotonic != 0 { // 判断 t 和 u 是否都有单调时钟
        return t.ext > u.ext             // 只需判断单调时钟即可
    }
    ts := t.sec()    // 否则需要从 wall 字段中获取秒数
    us := u.sec()
    // 判断 t 的秒数是否大于 u
    // 如果秒数相同，则比较纳秒数
    return ts > us || ts == us && t.nsec() > u.nsec() 
}

// 判断 t 时间是否早于 u 时间
func (t Time) Before(u Time) bool {
    if t.wall&u.wall&hasMonotonic != 0 { // 同上
        return t.ext < u.ext
    }
    // 同上反之
    return t.sec() < u.sec() || t.sec() == u.sec() && t.nsec() < u.nsec()
}

// 判断 t 与 u 时间是否相同，不判断 location
// 如 6:00 +0200 CEST  与  4:00 UTC  会返回 true
// 如需同时判断 location ，可以使用 == 操作符
func (t Time) Equal(u Time) bool {
    if t.wall&u.wall&hasMonotonic != 0 { // 同上
        return t.ext == u.ext
    }
    // 判断 t 与 u 的秒数以及纳秒数是否相同
    return t.sec() == u.sec() && t.nsec() == u.nsec()
}
```
Now函数：当前时间
```go
// src/time/time.go

const (
    // unix 时间戳为从 1970-01-01 00:00:00 开始到当前的秒数
    // time 包的 internal 时间为从 0000-00-00 00:00:00 开始的秒数
    // 以下两个常量用于在 internal 与 unix 时间戳之间转换的辅助常量
    unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
    internalToUnix int64 = -unixToInternal

    // minWall = wallToInternal
    wallToInternal int64 = (1884*365 + 1884/4 - 1884/100 + 1884/400) * secondsPerDay
    internalToWall int64 = -wallToInternal
)

// Provided by package runtime.  
// runtime包定义
func now() (sec int64, nsec int32, mono uint64)
// darwin,amd64 darwin,386 windows => src/runtime/timeasm.go
// other                           => src/runtime/timestub.go

//go:linkname time_now time.now
func time_now() (sec int64, nsec int32, mono int64) {
    sec, nsec = walltime()
    return sec, nsec, nanotime()
}


// runtime/time_nofake.go
//go:nosplit
func nanotime() int64 {
    return nanotime1()
}

func walltime() (sec int64, nsec int32) {
    return walltime1()
}

// 两者nanotime1和walltime1分别为几种 不同的 平台 和体系结构定义 


// 返回当前本机时间
func Now() Time {
    sec, nsec, mono := now()
    // 计算从 1885-1-1 开始到现在的秒数
    // unixToInternal = 1970-01-01 00:00:00
    // minWall        = 1885-01-01 00:00:00
    sec += unixToInternal - minWall
    if uint64(sec)>>33 != 0 { // 如果有溢出，则不能用 wall 的 second 保存完整的时间戳
        // 返回自 1970-01-01 00:00:00 开始的秒数
        return Time{uint64(nsec), sec + minWall, Local}
    }
    return Time{hasMonotonic | uint64(sec)<<nsecShift | uint64(nsec), mono, Local}
}
```

walltime是如何计算的: AMD64的Linux这里
```assembly
// func walltime1() (sec int64, nsec int32)
// non-zero frame-size means bp is saved and restored
TEXT runtime·walltime1(SB),NOSPLIT,$16-12
	// 由于我们不知道代码需要多少堆栈空间，因此我们切换到g0
	
	// In particular, a kernel configured with CONFIG_OPTIMIZE_INLINING=n
	// and hardening can use a full page of stack space in gettime_sym
	// due to stack probes inserted to avoid stack/heap collisions.
	// See issue #20427.

	MOVQ	SP, R12	// Save old SP; R12 unchanged by C code.

	get_tls(CX)
	MOVQ	g(CX), AX
	MOVQ	g_m(AX), BX // BX unchanged by C code.

	// Set vdsoPC and vdsoSP for SIGPROF traceback.
	// Save the old values on stack and restore them on exit,
	// so this function is reentrant.
	// 把 vdsoPC and vdsoSP（程序计数器和堆栈指针）压栈
	MOVQ	m_vdsoPC(BX), CX
	MOVQ	m_vdsoSP(BX), DX
	MOVQ	CX, 0(SP)
	MOVQ	DX, 8(SP)

	LEAQ	sec+0(FP), DX
	MOVQ	-8(DX), CX
	MOVQ	CX, m_vdsoPC(BX)
	MOVQ	DX, m_vdsoSP(BX)

    // 检查它是否已经打开
	CMPQ	AX, m_curg(BX)	// Only switch if on curg.
	JNE	noswitch

	MOVQ	m_g0(BX), DX
	MOVQ	(g_sched+gobuf_sp)(DX), SP	// Set SP to g0 stack

noswitch:
	SUBQ	$16, SP		// 让出空间
	ANDQ	$~15, SP	// Align for C code

	MOVL	$0, DI // CLOCK_REALTIME
	LEAQ	0(SP), SI
	MOVQ	runtime·vdsoClockgettimeSym(SB), AX
	CMPQ	AX, $0
	JEQ	fallback
	CALL	AX
ret:
	MOVQ	0(SP), AX	// sec
	MOVQ	8(SP), DX	// nsec
	MOVQ	R12, SP		// Restore real SP
	// 还原 vdsoPC, vdsoSP
	// We don't worry about being signaled between the two stores.
	// If we are not in a signal handler, we'll restore vdsoSP to 0,
	// and no one will care about vdsoPC. If we are in a signal handler,
	// we cannot receive another signal.
	MOVQ	8(SP), CX
	MOVQ	CX, m_vdsoSP(BX)
	MOVQ	0(SP), CX
	MOVQ	CX, m_vdsoPC(BX)
	MOVQ	AX, sec+0(FP)
	MOVL	DX, nsec+8(FP)
	RET
fallback:
	MOVQ	$SYS_clock_gettime, AX
	SYSCALL
	JMP ret

```

为什么__vdso_clock_gettime优先选择而不是__x64_sys_clock_gettime，它们之间有什么区别？

vDSO代表虚拟动态共享对象，是一种内核机制，用于将内核空间例程的子集导出到用户空间应用程序，以便可以在进程内调用这些内核空间例程，而不会造成从用户模式切换到内核模式的性能损失。


函数通过时间戳获取 time.Time对象 

time.Unix(sec, nsec) 函数通过传入 unix timestamp 获取 time.Time 结构，默认返回的是 UTC 时区。
```go
func Unix(sec int64, nsec int64) Time {
	if nsec < 0 || nsec >= 1e9 {
		n := nsec / 1e9
		sec += n
		nsec -= n * 1e9
		if nsec < 0 {
			nsec += 1e9
			sec--
		}
	}
	return unixTime(sec, int32(nsec))
}

func unixTime(sec int64, nsec int32) Time {
  return Time{uint64(nsec), sec + unixToInternal, Local}
}
```

Note:如果想获取带有 单调时钟 的时间只能通过 time.Now() 获取，而由于 wall 的 second 有 33 位，所以只要我们在 2157-01-01 00:00:00 UTC 前调用 time.Now() 获取到的时间都是带有 单调时钟 的

获取local

```go
var Local *Location = &localLoc

var localLoc Location
var localOnce sync.Once

func (l *Location) get() *Location {
    if l == nil {
        return &utcLoc
    }
    if l == &localLoc {
        localOnce.Do(initLocal)
    }
    return l
}
```
该initLocal()函数查找的内容$TZ以找到要使用的时区。
