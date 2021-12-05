# Time

## 三个比较常用的数据结

- time.Time
- time.Duration
- time.C

## 两个重要概念
- wall time:
  就是挂在墙上的时钟，我们在计算机中能看到的当前时间就是 wall time ，但是这个时间是可以通过 人为设置 或者 NTP服务同步 被修改，常见的场景就是通过修改时间延长收费软件的试用期。
- monotonic time:
  一个单调递增的时间，当操作系统被初始化时， jiffies 变量被初始化为 0 ，每当接收到一个 timer interrupt ，则 jiffies 自增 1 ，所以它必然是一个不可修改的单调递增时间。

应用：在操作系统中如果需要 显示 时间的时候，会使用 wall time ，而需要 测量 时间的时候，会使用 monotonic time 
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
    // 所以 wall 字段会有两种情况，分别为
    // 1. 当 wall 字段的 hasMonotonic 为 0 时，second 位也全部为 0，ext 字段会存储
    //    从 1-1-1 开始的秒级精度时间作为 wall time 。
    // 2. 当 wall 字段的 hasMonotonic 为 1 时，second 位会存储从 1885-1-1 开始的秒
    //    级精度时间作为 wall time，并且 ext 字段会存储从操作系统启动后的纳秒级精度时间
    //    作为 monotonic time 。
    wall uint64
    ext  int64
    
    // Location 作为当前时间的时区，可用于确定时间是否处在正确的位置上。
    // 当 loc 为 nil 时，则表示为 UTC 时间。
    // 因为北京时区为东八区，比 UTC 时间要领先 8 个小时，
    // 所以我们获取到的时间默认会记为 +0800
    loc *Location
}
```
存储机制
```go
// src/time/time.go

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
时间比较
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
当前时间
```go
// src/time/time.go

const (
    // unix 时间戳为从 1970-01-01 00:00:00 开始到当前的秒数
    // time 包的 internal 时间为从 0000-00-00 00:00:00 开始的秒数
    //
    // 以下两个常量用于在 internal 与 unix 时间戳之间转换的辅助常量
    unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
    internalToUnix int64 = -unixToInternal

    // minWall = wallToInternal
    wallToInternal int64 = (1884*365 + 1884/4 - 1884/100 + 1884/400) * secondsPerDay
    internalToWall int64 = -wallToInternal
)

// Provided by package runtime.
func now() (sec int64, nsec int32, mono uint64)
// darwin,amd64 darwin,386 windows => src/runtime/timeasm.go
// other                           => src/runtime/timestub.go

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

函数获取 time.Time 

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
