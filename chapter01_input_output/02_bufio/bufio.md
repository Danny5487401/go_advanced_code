# bufio
    bufio实现了带缓冲的IO功能，它是在io.Reader和io.Writer接口对象上提供了进一步的封装

设计原因：

    io操作本身的效率并不低，低的是频繁的访问本地磁盘的文件。
解决：

    所以bufio就提供了缓冲区(分配一块内存)，读和写都先在缓冲区中，最后再读写文件，来降低访问本地磁盘的次数，从而提高效率。

原理：

    把文件读取进缓冲（内存）之后再读取的时候就可以避免文件系统的io 从而提高速度。
    同理，在进行写操作时，先把文件写入缓冲（内存），然后由缓冲写入文件系统。
    看完以上解释有人可能会表示困惑了，直接把 内容->文件 和 内容->缓冲->文件相比， 缓冲区好像没有起到作用嘛。
    其实缓冲区的设计是为了存储多次的写入，最后一口气把缓冲区内容写入文件
分类：

    主要分三部分Reader、Writer、Scanner,分别是读数据、写数据和扫描器三种数据类型
主要读取方式

    ReadLine和ReadString方法：buf.ReadLine()，buf.ReadString("\n")都是按行读，只不过ReadLine读出来的是[]byte，后者直接读出了string，
    最终他们底层调用的都是ReadSlice方法

## 源码分析

```go
//结构体
type Reader struct {
    buf          []byte    //缓冲区的数据
    rd           io.Reader // 底层的io.Reader
    r, w         int       //  r ,w分别表示 buf中读和写的指针位置
    err          error    //记录本次读取的error，后续操作中调用readErr函数后会重置err
    lastByte     int      //记录读取的最后一个字节（用于撤销）
    lastRuneSize int      //记录读取的最后一个字符(Rune)的长度（用于撤销）
}

// 初始化
func NewReaderSize(rd io.Reader, size int) *Reader {
    // Is it already a Reader?
    b, ok := rd.(*Reader)
    if ok && len(b.buf) >= size {
    return b
    }
    if size < minReadBufferSize {
    size = minReadBufferSize
    }
    r := new(Reader)
    r.reset(make([]byte, size), rd)
    return r
}
//size用于指定缓冲区的大小，如果size小于minReadBufferSize，则重置size的值为minReadBufferSize（16）。
//如果该rd是*bufio.Reader对象，并且rd的缓冲区大于size，则不会创建Reader对象，而是直接返回原来的rd对象。
//否则会创建一个*bufio.Reader对象，并指定buf的大小为size

```
常用方法
```go
 //Size方法返回底层缓冲区的大小
func (b *Reader) Size() int { return len(r.buf) }


//   Reset方法丢弃所有缓冲数据，重置所有状态，并切换reader到r，从而从r中读取数据。
func (b *Reader) Reset(r io.Reader) {
    b.reset(b.buf, r)
}

func (b *Reader) reset(buf []byte, r io.Reader) {
    *b = Reader{
    buf:          buf,
    rd:           r,
    lastByte:     -1,
    lastRuneSize: -1,
}

// fill方法用于将缓冲区读满，可以读入的最大长度是：len(buf)-未读的字节数，如果尝试读取了maxConsecutiveEmptyReads(100)次都没有读取到数据，则会返回    
func (b *Reader) fill() {
// Slide existing data to beginning.
    if b.r > 0 {
        /*
            将buf中未读的数据copy到buf中首部位置
            重置r和w 位置
        */
        copy(b.buf, b.buf[b.r:b.w])
        b.w -= b.r
        b.r = 0
    }
    
    if b.w >= len(b.buf) {
    panic("bufio: tried to fill full buffer")
    }
    
    /*
        maxConsecutiveEmptyReads是最多尝试次数
        从rd Reader中读取数据到缓冲区buf中，并重置w的位置索引
    */
    for i := maxConsecutiveEmptyReads; i > 0; i-- {
        n, err := b.rd.Read(b.buf[b.w:])
        if n < 0 {
            panic(errNegativeRead)
            }
            b.w += n
            if err != nil {
            b.err = err
            return
            }
            if n > 0 {
            return
            }
            }
            b.err = io.ErrNoProgress
}


//Peek返回输入流的前n个字节，而不会移动读取位置。该操作不会将数据读出，只是引用，引用的数据在下一次读取操作之前是有效的，
//如果Peek返回的切片长度比n小，它也会返会一个错误说明原因。如果n比缓冲尺寸还大，返回的错误将是ErrBufferFull
func (b *Reader) Peek(n int) ([]byte, error) {
    if n < 0 {
    return nil, ErrNegativeCount
    }

    /*
        如果buf中未读的字节数小于n
        并且buf中未读的字节数小于buf的总大小
        并且err为nil
        满足以上3个条件，则调用fill方法尝试从Reader中读取部分数据块
    */
    for b.w-b.r < n && b.w-b.r < len(b.buf) && b.err == nil {
        b.fill() // b.w-b.r < len(b.buf) => buffer is not full
    }

    /*
        如果n大于buf的长度，则返回所有未读的内容，和ErrBufferFull的错误信息
    */
    if n > len(b.buf) {
        return b.buf[b.r:b.w], ErrBufferFull
    }

    // 0 <= n <= len(b.buf)
    var err error
    /*
        如果n大于可读的长度，则返回error信息
    */
    if avail := b.w - b.r; avail < n {
        // not enough data in buffer
        n = avail
        err = b.readErr()
        if err == nil {
        err = ErrBufferFull
        }
    }
    return b.buf[b.r : b.r+n], err
}
    
```
```go
// Discard 方法跳过后续的 n 个字节的数据，返回跳过的字节数。如果结果小于 n，将返回错误信息。如果 n 小于缓存中的数据长度，则不会从底层提取数据
func (b *Reader) Discard(n int) (discarded int, err error) {
	if n < 0 {
		return 0, ErrNegativeCount
	}
	if n == 0 {
		return
	}
	remain := n  //remain记录剩余要跳过的字节长度
	for {
		skip := b.Buffered() //获取buf中可读的数据长度
		if skip == 0 {//如果可读的数据长度为0，则尝试从底层rd  Reader中读取数据到buf中，然后再获取buf中可读的数据长度
			b.fill()
			skip = b.Buffered()
		}
		if skip > remain {//设置skip
			skip = remain
		}
		b.r += skip  //设置r的位置，改变读的索引值r
		remain -= skip
		if remain == 0 { //如果已经跳过n个字节，则返回
			return n, nil
		}
		if b.err != nil {
			return n - remain, b.readErr()
		}
	}
}

```

```go
/*
Read从Reader对象b中读出数据到p中，n是返回读取的字节数。如果Reader对象b中的缓冲区buf不为空，则只能读取缓冲中的数据，不会从底层的io.Reader中读取数据。

如果b中的缓冲buf为空，则：

1. len(p) >= 缓冲大小，则跳过缓存，直接从底层io.Reader中读出到p总

2.len(p)<缓存大小，则先将数据从底层的io.Reader中读取到缓存中，再从缓存buf读取到p中。

 */
func (b *Reader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {  //如果p的长度为0，则直接返回
		return 0, b.readErr()
	}
	if b.r == b.w {  //缓冲区没有可读的数据
		if b.err != nil { //如果上次读取有error，则直接返回
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) { //如果 len(p) >=  缓冲区大小
			// Large read, empty buffer.
			// Read directly into p to avoid copy.
			n, b.err = b.rd.Read(p)  //直接读取rd中的数据到p中
			if n < 0 {
				panic(errNegativeRead)
			}
			if n > 0 {  //有读取的数据
				b.lastByte = int(p[n-1])
				b.lastRuneSize = -1
			}
			return n, b.readErr()
		}
		// One read.
		// Do not use b.fill, which will loop.
		/*
		如果len(p) < 缓冲区大小，将读(r)的位置和写(w)的位置设为0，并从rd中读取数据到缓冲区buf中
		*/
		b.r = 0
		b.w = 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n  //更新写的位置
	}
 
	// copy as much as we can
	n = copy(p, b.buf[b.r:b.w])  //将数据copy到p中
	b.r += n  //并更新读的位置
	b.lastByte = int(b.buf[b.r-1])  //记录读取的最后一个字节
	b.lastRuneSize = -1
	return n, nil
}
```

```go
//fun    ReadRune  读取一个UTF-8字符，并且返回该字符r，该字符的所占的字节数size，以及可能出现的error，如果UTF-8编码的字符无效，则消耗一个字节并返回大小为1的unicode.ReplacementChar（U + FFFDc (b *Reader) ReadRune() (r rune, size int, err error) {
 
	/*
	当满足：
	1.可读的字节数小于UTFMax（UTFMax是一个UTF字符占用的最大字节数）
	2.可读的数据不够填充一个UTF-8字符
	3.没有error
	4.可读的字节数小于len(b.buf)
	当满足以上几个条件会调用fill来从io.Reader中读出内容到buf中
	*/
	for b.r+utf8.UTFMax > b.w && !utf8.FullRune(b.buf[b.r:b.w]) && b.err == nil && b.w-b.r < len(b.buf) {
		b.fill() // b.w-b.r < len(buf) => buffer is not full
	}
	b.lastRuneSize = -1  //重置lastRuneSize
	if b.r == b.w {  //当没有可读数据
		return 0, 0, b.readErr()
	}
	r, size = rune(b.buf[b.r]), 1
	if r >= utf8.RuneSelf { //不是ascii码
		r, size = utf8.DecodeRune(b.buf[b.r:b.w]) //读取一个字符r和长度size
	}
	b.r += size
	b.lastByte = int(b.buf[b.r-1])  //设置lastByte
	b.lastRuneSize = size  //设置lastRuneSize
	return r, size, nil
}

//  Buffered方法返回buf中可读的数据长度
func (b *Reader) Buffered() int { return b.w - b.r }
```

```go

/*
 ReadSlice方法在b中查找delim字节并返回delim及之前的所有数据，返回的是buf的切片

如果找到delim，则返回查找结果，err为nil

如果未找到delim，则：

1.缓存不满，则将缓存填满后再次查找

2.缓存是满的，则返回整个缓存，err返回ErrBufferFull

3.如果未找到delim且遇到错误(通常是io.EOF)，则返回缓存中的所有数据和遇到的错误。
因为返回的数据有可能被下一次的读写操作修改，所以大多数操作应该使用ReadBytes或ReadString，它们返回的数据的拷贝

*/
func (b *Reader) ReadSlice(delim byte) (line []byte, err error) {
	for {
		// Search buffer.
		//从buf中未读的数据中查找delim字节
		//如果可以找到则返回buf中的切片，并设置读的索引位置r，并跳出循环
		if i := bytes.IndexByte(b.buf[b.r:b.w], delim); i >= 0 {
			line = b.buf[b.r : b.r+i+1]
			b.r += i + 1
			break
		}
 
		// Pending error?
		if b.err != nil {  //如果上次读操作出现error，则返回
			line = b.buf[b.r:b.w]
			b.r = b.w
			err = b.readErr()
			break
		}
 
		// Buffer full?
		if b.Buffered() >= len(b.buf) { //如果buf是满的，则返回整个buf，并设置错误ErrBufferFull
			b.r = b.w
			line = b.buf
			err = ErrBufferFull
			break
		}
 
		b.fill() // buffer is not full 如果buf没有满，则向buf中读入数据
	}
 
	// Handle last byte, if any.
	if i := len(line) - 1; i >= 0 {
		b.lastByte = int(line[i])  //设置最后一个字节
		b.lastRuneSize = -1
	}
 
	return
}
```

```go
/*
ReadLine方法是一个低水平的行读取操作，大多数情况下，应该使用

ReadBytes('\n')或ReadString('\n')，或者使用Scanner

ReadLine 通过调用 ReadSlice 方法实现，返回的也是缓存的切片。用于读取一行数据，不包括行尾标记（\n 或 \r\n）

 */
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error) {
	//通过ReadSlice读取一行内容，返回的是buf的切片，在下次的读写操作时内容可能会改变
	line, err = b.ReadSlice('\n')
	if err == ErrBufferFull {
		// Handle the case where "\r\n" straddles the buffer.
		if len(line) > 0 && line[len(line)-1] == '\r' {
			// Put the '\r' back on buf and drop it from line.
			// Let the next call to ReadLine check for "\r\n".
			if b.r == 0 {
				// should be unreachable
				panic("bufio: tried to rewind past start of buffer")
			}
			b.r--
			line = line[:len(line)-1]
		}
		return line, true, nil
	}
 
	if len(line) == 0 {
		if err != nil {
			line = nil
		}
		return
	}
	err = nil
 
	if line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		line = line[:len(line)-drop]
	}
	return
}
```


```go
//ReadBytes 方法从b中读取数据 直到第一次出现delim字节

//功能类似ReadSlice
/*
不同点：

1.只不过返回的是缓存buf的copy，这样当下次读写操作也不会影响返回的内容

2.ReadSlice最多在整个len(buf)中查找，如果未找到返回ErrBufferFull。但是ReadBytes还会继续从io.Reader中读取数据，直到找到该delim字节，或者io.EOF，或者出现unexpected error

也就是ReadSlice最多返回的切片长度为len(buf),而ReadBytes返回的切片长度可能会大于len(buf)

 */
func (b *Reader) ReadBytes(delim byte) ([]byte, error) {
	// Use ReadSlice to look for array,
	// accumulating full buffers.
	var frag []byte
	var full [][]byte  //可能多次从底层io.Reader中读取数据
	var err error
	for {//可能多次调用ReadSlice去查找delim字节，直到找到或unexpected error(包括io.EOF)
		var e error
		frag, e = b.ReadSlice(delim)
		if e == nil { // got final fragment
			break
		}
		if e != ErrBufferFull { // unexpected error
			err = e
			break
		}
 
		// Make a copy of the buffer.
		buf := make([]byte, len(frag))
		copy(buf, frag)
		full = append(full, buf)
	}
 
	// Allocate new buffer to hold the full pieces and the fragment.
	n := 0
	for i := range full {
		n += len(full[i])
	}
	n += len(frag)
 
	// Copy full pieces and fragment in.
	buf := make([]byte, n)
	n = 0
	for i := range full {
		n += copy(buf[n:], full[i])
	}
	copy(buf[n:], frag)
	return buf, err
}


func (b *Reader) ReadString(delim byte) (string, error) {
    bytes, err := b.ReadBytes(delim)
    return string(bytes), err
}
//ReadString方法功能同ReadBytes，只不过把ReadBytes返回的切片转成字符串
```