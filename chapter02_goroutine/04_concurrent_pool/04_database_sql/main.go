package _4_database_sql

//自动注册
import (
	_ "github.com/go-sql-driver/mysql" //自动执行init()函数
)

/*
1.驱动注册
	drivers   = make(map[string]driver.Driver)

	func Register(name string, driver driver.Driver) {
		driversMu.Lock()
		defer driversMu.Unlock()
		if driver == nil {
			panic("sql: Register driver is nil")
		}
		if _, dup := drivers[name]; dup {
			panic("sql: Register called twice for driver " + name)
		}
		drivers[name] = driver
对外提供的注册函数，只要底层的驱动
	type Driver interface {
		Open(name string) (Conn, error)
	}
mysql的init函数
	func init() {
		sql.Register("mysql", &MySQLDriver{})
	}

2.打开DB句柄
	type DB struct {
		// 是一个包含Connect(context.Context) (Conn, error)和Driver() Driver方法的接口。
		// 这两个接口是需要实际的数据库驱动来实现，再使用
	   connector driver.Connector// 用于获取driver.Conn 可以由驱动层实现，否则用sql.dsnConnector

	   numClosed uint64 // 是一个原子计数器，是自sql.DB创建后关闭过的连接数；

	   mu           sync.Mutex
	   freeConn     []*driverConn //空闲连接池

		// 存储pending连接的map，当numOpen大于maxOpen时，连接会被暂存到该map中；
	   connRequests map[uint64]chan connRequest

		// 即connRequests的key，记录当前可用的最新的connRequest的key
	   nextRequest  uint64

		//  当前活跃连接数和当前pending的连接数的总和
	   numOpen      int

	   openerCh    chan struct{} // 用来传信号的管道 表示需要多少新连接
	   resetterCh  chan *driverConn // 用来传需要重置 Session 的 driverConn
	   closed      bool

      // 记录db与conn之间的依赖关系，维持连接池以及关闭时使用；
	   dep         map[finalCloser]depSet

	  // 最新入栈的连接，debug时使用，string中存的是栈buf相关信息；
	   lastPut     map[*driverConn]string

		// 最大空闲连接数；
	   maxIdle     int

		// 最大连接数
	   maxOpen     int

	  // 控线连接的最大存活时间
	   maxLifetime time.Duration

		// 标记超时连接的通知channel
	   cleanerCh   chan struct{} // 传信号 表示需要清理freeConn空闲池中已经关掉的driverConn

		// 通过context通知关闭connection opener和session resetter；
	   stop func()
	}
	注册驱动时不会真正连接数据库，只有在调用Open()方法的时候才会连接
	func Open(driverName, dataSourceName string) (*DB, error) {
		driversMu.RLock()
		driveri, ok := drivers[driverName]//获得已经注册过的driver
		driversMu.RUnlock()
		.......
		//先判断是否实现了driver.DriverContext的接口
		if driverCtx, ok := driveri.(driver.DriverContext); ok {
			connector, err := driverCtx.OpenConnector(dataSourceName)//得到mysql的connector
			// 调用驱动中的OpenConnector
			// func (d MySQLDriver) OpenConnector(dsn string) (driver.Connector, error) {
			//    cfg, err := ParseDSN(dsn)
			//    if err != nil {
			//        return nil, err
			//    }
			//    return &connector{
			//        cfg: cfg,
			//    }, nil
			// }
			// 实际是返回了dsn对应Config结构体指针

			if err != nil {
				return nil, err
			}
			//初始化DB结构体
			return OpenDB(connector), nil//最后通过connector参数调用OpenDB
		}
		// 这一句是在转换driver.DriverContext接口失败时执行
		// 因为go1.10之前默认是一下方法，1.10后使用driver.DriverContext接口
		// 下面是为了做兼容
		return OpenDB(dsnConnector{dsn: dataSourceName, driver: driveri}), nil
	}
	这里值得注意的是，先判断底层的驱动是否实现了driver.DriverContext的接口，如果没有实现，会默认调用sql自己实现的dsnConnector，两者是有区别的，
	前者有Context的使用权，后者没有使用

	func OpenDB(c driver.Connector) *DB {
		ctx, cancel := context.WithCancel(context.Background())
		db := &DB{
			connector:    c,
			openerCh:     make(chan struct{}, connectionRequestQueueSize),
			resetterCh:   make(chan *driverConn, 50),
			lastPut:      make(map[*driverConn]string),
			connRequests: make(map[uint64]chan connRequest),
			stop:         cancel,
		}

		go db.connectionOpener(ctx) //connOpener 运行在一个单独的goroutine中
		// 启动清理session的goroutine
		go db.connectionResetter(ctx)//connResetter单独运行在一个goroutine中

		return db
	}
3. 获取连接
	go从连接池中获取连接时有两个策略
	// a.请求新的连接
	alwaysNewConn connReuseStrategy = iota
	// b.从连接池中获取连接
	cachedOrNewConn
	// 一般查询代码
	rows, err := dbt.db.Query("SELECT * FROM test")
	// Query如下，包了QueryContext方法
	func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
		// context.Background()是为了创建一个根上下文，以便close的时候能够将资源彻底的释放
		return db.QueryContext(context.Background(), query, args...)
	}
	//QueryContext方法
	func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
		var rows *Rows
		var err error
		// maxBadConnRetries是个静态变量为2，这里最多会执行两次从连接池中获取连接，如果在两次获取
		// 过程中获取到可用连接则直接返回
		for i := 0; i < maxBadConnRetries; i++ {
			rows, err = db.query(ctx, query, args, cachedOrNewConn)
			if err != driver.ErrBadConn {
				break
			}
		}
		// 如果两次都获取不到可用连接，则以请求获取一个新连接的方式获取并返回
		if err == driver.ErrBadConn {
			return db.query(ctx, query, args, alwaysNewConn)
		}
		return rows, err
	}


	// query方法如下
	func (db *DB) query(ctx context.Context, query string, args []interface{}, strategy connReuseStrategy) (*Rows, error) {
		// 这里是重点，这是真正申请连接的过程
		dc, err := db.conn(ctx, strategy)
		if err != nil {
			return nil, err
		}
		// 这里是实际的查询过程，不过多介绍
		return db.queryDC(ctx, nil, dc, dc.releaseConn, query, args)
	}

	上面的OpenDB步骤并不会打开连接，真正打开连接是在执行query()前，分析一下DB.conn方法
	func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error)
	conn := db.freeConn[0]
	copy(db.freeConn, db.freeConn[1:])
	db.freeConn = db.freeConn[:numFree-1]
	相比
	conn := db.freeConn[0]
	db.freeConn = db.freeConn[1:]
	第一种方式：不会破坏数组的容量，把数组元素相当于向前移动一个位置，避免了append时再次分配内存。
	第二种方式：会破坏数组的容量，即容量减少一个
4，连接回收到连接池 func (db *DB) putConn(dc *driverConn, err error, resetSession bool)
	1)首先遍历dc.onPut，执行fn()

	2)如果发现该连接不可用，则调用maybeOpenNewConnections() 异步创建一个连接，并且关闭不可用的连接。

	3)如果连接成功被连接池回收，但db.resetterCh 阻塞了，则先标记连接为ErrBadConn,所以前面从连接池获取连接时每一次都会判断连接是否可用。

	4)如果连接池满了，没回放成功，则会关闭该连接
5.处理过期的连接 func (db *DB) connectionCleaner(d time.Duration)
	1）开定时器，每隔一段时间检测空闲连接池中的连接是否过期

	2）如果接收到db.cleanerCh的信号，也会遍历处理超时，db.cleanerCh的buffer只有1，一般在SetConnMaxLifetime检测生命周期配置变短时发送。

	3）为了遍历空闲队列里面连接的公平性，做了一个巧妙的处理，一旦发现队列前面的连接过期，则会把最后一个连接放到最前面，然后从当前开始遍历。

	4）遍历空闲队列发现超时的连接，把超时连接一个一个追加到关闭队列中append(closing, c)，然后遍历关闭的队列，一个一个关闭

*/

func main() {}
