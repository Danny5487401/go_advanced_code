package _4_database_sql

//自动注册
import (
	 "database/sql"
	_ "github.com/go-sql-driver/mysql"//自动执行init()函数
）
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
	   connector driver.Connector// 用于获取driver.Conn 可以由驱动层实现，否则用sql.dsnConnector
	
	   numClosed uint64 // 是一个原子计数器，代表总的关闭连接数量
	
	   mu           sync.Mutex
	   freeConn     []*driverConn //空闲连接池
	   connRequests map[uint64]chan connRequest // 无可用连接时，处于 Pending 状态的连接请求
	   nextRequest  uint64
	   numOpen      int    // 打开和准备打开的连接总数
	
	   openerCh    chan struct{} // 用来传信号的管道 表示需要多少新连接
	   resetterCh  chan *driverConn // 用来传需要重置 Session 的 driverConn
	   closed      bool
	   dep         map[finalCloser]depSet // 依赖记录
	   lastPut     map[*driverConn]string
	   maxIdle     int
	   maxOpen     int
	   maxLifetime time.Duration // 连接的生命后期
	   cleanerCh   chan struct{} // 传信号 表示需要清理freeConn空闲池中已经关掉的driverConn
	
	   stop func()
	}
	func Open(driverName, dataSourceName string) (*DB, error) {
		driversMu.RLock()
		driveri, ok := drivers[driverName]//获得已经注册过的driver
		driversMu.RUnlock()
		.......
		//先判断是否实现了driver.DriverContext的接口
		if driverCtx, ok := driveri.(driver.DriverContext); ok {
			connector, err := driverCtx.OpenConnector(dataSourceName)//得到mysql的connecter
			if err != nil {
				return nil, err
			}
			return OpenDB(connector), nil//最后通过connector参数调用OpenDB
		}

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
		go db.connectionResetter(ctx)//connResetter单独运行在一个goroutine中

		return db
	}


*/
func main() {
sql.Open("mysql", )
}
