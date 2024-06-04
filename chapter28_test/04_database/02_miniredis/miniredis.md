<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [miniredis: 纯Go实现的Redis测试服务器](#miniredis-%E7%BA%AFgo%E5%AE%9E%E7%8E%B0%E7%9A%84redis%E6%B5%8B%E8%AF%95%E6%9C%8D%E5%8A%A1%E5%99%A8)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [启动服务](#%E5%90%AF%E5%8A%A8%E6%9C%8D%E5%8A%A1)
    - [设置值](#%E8%AE%BE%E7%BD%AE%E5%80%BC)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# miniredis: 纯Go实现的Redis测试服务器

Miniredis是一个专为Go语言单元测试设计的纯Go实现的轻量级Redis模拟服务器。在开发过程中，有时我们希望在不进行完整的集成测试的情况下测试与Redis相关的代码，这就是Miniredis的作用。
它提供了一个真实的TCP接口，使你能快速地在内存中搭建一个Redis替代品，方便你在测试时直接查询和操作数据。


## 特点

Miniredis实现了Redis命令的大部分功能，包括连接管理、键值操作、事务处理、字符串、哈希表、列表、发布订阅、集合、有序集合、脚本、时间戳等。此外，
它还支持了部分Redis 7.2.0的新特性，如RESP3协议。通过引入这个库，你可以轻松地在你的测试代码中模拟复杂的Redis交互场景，无需依赖外部Redis实例



## 源码分析

结构体
```go
// /Users/python/go/pkg/mod/github.com/alicebob/miniredis/v2@v2.22.0/miniredis.go
type Miniredis struct {
	sync.Mutex
	srv         *server.Server
	port        int
	passwords   map[string]string // username password
	dbs         map[int]*RedisDB
	selectedDB  int               // DB id used in the direct Get(), Set() &c.
	scripts     map[string]string // sha1 -> lua src
	signal      *sync.Cond
	now         time.Time // time.Now() if not set.
	subscribers map[*Subscriber]struct{}
	rand        *rand.Rand
	Ctx         context.Context
	CtxCancel   context.CancelFunc
}

// RedisDB holds a single (numbered) Redis database.
type RedisDB struct {
	master        *Miniredis               // pointer to the lock in Miniredis
	id            int                      // db id
	keys          map[string]string        // Master map of keys with their type
	stringKeys    map[string]string        // GET/SET &c. keys
	hashKeys      map[string]hashKey       // MGET/MSET &c. keys
	listKeys      map[string]listKey       // LPUSH &c. keys
	setKeys       map[string]setKey        // SADD &c. keys
	hllKeys       map[string]*hll          // PFADD &c. keys
	sortedsetKeys map[string]sortedSet     // ZADD &c. keys
	streamKeys    map[string]*streamKey    // XADD &c. keys
	ttl           map[string]time.Duration // effective TTL values
	keyVersion    map[string]uint          // used to watch values
}
```

### 启动服务
```go
func RunT(t Tester) *Miniredis {
	m := NewMiniRedis()
	if err := m.Start(); err != nil {
		t.Fatalf("could not start miniredis: %s", err)
		// not reached
	}
	t.Cleanup(m.Close)
	return m
}

func NewMiniRedis() *Miniredis {
	m := Miniredis{
		dbs:         map[int]*RedisDB{},
		scripts:     map[string]string{},
		subscribers: map[*Subscriber]struct{}{},
	}
	m.Ctx, m.CtxCancel = context.WithCancel(context.Background())
	m.signal = sync.NewCond(&m)
	return &m
}
```
启动了一个tcp服务器
```go
// Start starts a server. It listens on a random port on localhost. See also
// Addr().
func (m *Miniredis) Start() error {
	s, err := server.NewServer(fmt.Sprintf("127.0.0.1:%d", m.port))
	if err != nil {
		return err
	}
	return m.start(s)
}
```



```go
func newServer(l net.Listener) *Server {
	s := Server{
		cmds:  map[string]Cmd{},
		peers: map[net.Conn]struct{}{},
		l:     l,
	}
    // 每接受一个请求，起一个协程来处理请求
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.serve(l)

		s.mu.Lock()
		for c := range s.peers {
			c.Close()
		}
		s.mu.Unlock()
	}()
	return &s
}
```

命令注册
```go
func (m *Miniredis) start(s *server.Server) error {
	m.Lock()
	defer m.Unlock()
	m.srv = s
	m.port = s.Addr().Port

	commandsConnection(m)
	commandsGeneric(m)
	commandsServer(m)
	commandsString(m)
	commandsHash(m)
	commandsList(m)
	commandsPubsub(m)
	commandsSet(m)
	commandsSortedSet(m)
	commandsStream(m)
	commandsTransaction(m)
	commandsScripting(m)
	commandsGeo(m)
	commandsCluster(m)
	commandsCommand(m)
	commandsHll(m)

	return nil
}
```
举例
```go
// commandsString handles all string value operations.
func commandsString(m *Miniredis) {
	m.srv.Register("APPEND", m.cmdAppend)
	m.srv.Register("BITCOUNT", m.cmdBitcount)
	m.srv.Register("BITOP", m.cmdBitop)
	m.srv.Register("BITPOS", m.cmdBitpos)
	m.srv.Register("DECRBY", m.cmdDecrby)
	m.srv.Register("DECR", m.cmdDecr)
	m.srv.Register("GETBIT", m.cmdGetbit)
	m.srv.Register("GET", m.cmdGet)
	m.srv.Register("GETEX", m.cmdGetex)
	m.srv.Register("GETRANGE", m.cmdGetrange)
	m.srv.Register("GETSET", m.cmdGetset)
	m.srv.Register("GETDEL", m.cmdGetdel)
	m.srv.Register("INCRBYFLOAT", m.cmdIncrbyfloat)
	m.srv.Register("INCRBY", m.cmdIncrby)
	m.srv.Register("INCR", m.cmdIncr)
	m.srv.Register("MGET", m.cmdMget)
	m.srv.Register("MSET", m.cmdMset)
	m.srv.Register("MSETNX", m.cmdMsetnx)
	m.srv.Register("PSETEX", m.cmdPsetex)
	m.srv.Register("SETBIT", m.cmdSetbit)
	m.srv.Register("SETEX", m.cmdSetex)
	m.srv.Register("SET", m.cmdSet)
	m.srv.Register("SETNX", m.cmdSetnx)
	m.srv.Register("SETRANGE", m.cmdSetrange)
	m.srv.Register("STRLEN", m.cmdStrlen)
}
```
```go
func (s *Server) Register(cmd string, f Cmd) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cmd = strings.ToUpper(cmd)
	if _, ok := s.cmds[cmd]; ok {
		return fmt.Errorf("command already registered: %s", cmd)
	}
	s.cmds[cmd] = f
	return nil
}
```

### 设置值
```go
// /Users/python/go/pkg/mod/github.com/alicebob/miniredis/v2@v2.22.0/server/server.go
// ServeConn handles a net.Conn. Nice with net.Pipe()
func (s *Server) ServeConn(conn net.Conn) {
	s.wg.Add(1)
	s.mu.Lock()
	s.peers[conn] = struct{}{}
	s.infoConns++
	s.mu.Unlock()

	go func() {
		defer s.wg.Done()
		defer conn.Close()

		s.servePeer(conn)

		s.mu.Lock()
		delete(s.peers, conn)
		s.mu.Unlock()
	}()
}

func (s *Server) servePeer(c net.Conn) {
	r := bufio.NewReader(c)
	peer := &Peer{
		w: bufio.NewWriter(c),
	}
	defer func() {
		for _, f := range peer.onDisconnect {
			f()
		}
	}()

	for {
		args, err := readArray(r)
		if err != nil {
			return
		}
		s.Dispatch(peer, args)
		peer.Flush()

		s.mu.Lock()
		closed := peer.closed
		s.mu.Unlock()
		if closed {
			c.Close()
		}
	}
}
```

取出对应的命令处理函数：
```go
func (s *Server) Dispatch(c *Peer, args []string) {
	cmd, args := args[0], args[1:]
	cmdUp := strings.ToUpper(cmd)
	s.mu.Lock()
	h := s.preHook
	s.mu.Unlock()
	if h != nil {
		if h(c, cmdUp, args...) {
			return
		}
	}

	s.mu.Lock()
	cb, ok := s.cmds[cmdUp]
	s.mu.Unlock()
	if !ok {
		c.WriteError(errUnknownCommand(cmd, args))
		return
	}

	s.mu.Lock()
	s.infoCmds++
	s.mu.Unlock()
	cb(c, cmdUp, args)
}
```