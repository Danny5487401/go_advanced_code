<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Fileinfo 文件属性](#fileinfo-%E6%96%87%E4%BB%B6%E5%B1%9E%E6%80%A7)
  - [接口](#%E6%8E%A5%E5%8F%A3)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Fileinfo 文件属性

文件属性，也即文件元数据。在 Go 中，文件属性具体信息通过 os.FileInfo 接口获取。
函数 Stat、Lstat 和 File.Stat 可以得到该接口的实例。这三个函数对应三个系统调用：stat、lstat 和 fstat


这三个函数的区别：

- stat 会返回所命名文件的相关信息。
- lstat 与 stat 类似，区别在于如果文件是符号链接，那么所返回的信息针对的是符号链接自身（而非符号链接所指向的文件）。
- fstat 则会返回由某个打开文件描述符（Go 中则是当前打开文件 File）所指代文件的相关信息

## 接口
```go
type FileInfo interface {
    Name() string       // 文件的名字（不含扩展名）
    Size() int64        // 普通文件返回值表示其大小,单位byte字节；其他文件的返回值含义各系统不同
    Mode() FileMode     // 文件的模式位
    ModTime() time.Time // 文件的修改时间
    IsDir() bool        // 等价于Mode().IsDir()
    Sys() interface{}   // 底层数据来源（可以返回nil）
}
```

FileMode 文件的模式位
```go
type FileMode uint32

// The defined file mode bits are the most significant bits of the FileMode.
// The nine least-significant bits are the standard Unix rwxrwxrwx permissions.
// The values of these bits should be considered part of the public API and
// may be used in wire protocols or disk representations: they must not be
// changed, although new bits might be added.
const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
	ModeAppend                                     // a: append-only
	ModeExclusive                                  // l: exclusive user
	ModeTemporary                                  // T: temporary file; Plan 9 only
	ModeSymlink                                    // L: symbolic link
	ModeDevice                                     // D: device file
	ModeNamedPipe                                  // p: named pipe (FIFO)
	ModeSocket                                     // S: Unix domain socket
	ModeSetuid                                     // u: setuid
	ModeSetgid                                     // g: setgid
	ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
	ModeSticky                                     // t: sticky
	ModeIrregular                                  // ?: non-regular file; nothing else is known about this file

	// Mask for the type bits. For regular files, none will be set.
	ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeCharDevice | ModeIrregular

	ModePerm FileMode = 0777 // Unix permission bits
)
```

Sys() 底层数据的 C 语言 结构 statbuf 格式如下：
```cgo
struct stat {
    dev_t    st_dev;    // 设备 ID
    ino_t    st_ino;    // 文件 i 节点号
    mode_t    st_mode;    // 位掩码，文件类型和文件权限
    nlink_t    st_nlink;    // 硬链接数
    uid_t    st_uid;    // 文件属主，用户 ID
    gid_t    st_gid;    // 文件属组，组 ID
    dev_t    st_rdev;    // 如果针对设备 i 节点，则此字段包含主、辅 ID
    off_t    st_size;    // 常规文件，则是文件字节数；符号链接，则是链接所指路径名的长度，字节为单位；对于共享内存对象，则是对象大小
    blksize_t    st_blsize;    // 分配给文件的总块数，块大小为 512 字节
    blkcnt_t    st_blocks;    // 实际分配给文件的磁盘块数量
    time_t    st_atime;        // 对文件上次访问时间
    time_t    st_mtime;        // 对文件上次修改时间
    time_t    st_ctime;        // 文件状态发生改变的上次时间
}
```
对应 Go 中 syscal.Stat_t 与该结构对应
