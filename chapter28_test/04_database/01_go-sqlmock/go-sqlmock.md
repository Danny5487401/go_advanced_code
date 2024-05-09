<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [sql mock](#sql-mock)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [go-sqlmock 源码分析](#go-sqlmock-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# sql mock

这里使用 github.com/DATA-DOG/go-sqlmock


## 使用

```go
// 使用相等匹配器
dbEqual, mockEqual, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
logx.Must(err)
// 默认的 sql 匹配器是正则
dbRegexp, mockRegexp, err := sqlmock.New()
logx.Must(err)
```

## go-sqlmock 源码分析

注册 driver
```go
// go-sqlmock@v1.5.2/driver.go
func init() {
	pool = &mockDriver{
		conns: make(map[string]*sqlmock),
	}
	sql.Register("sqlmock", pool)
}

func (d *mockDriver) Open(dsn string) (driver.Conn, error) {
	d.Lock()
	defer d.Unlock()

	c, ok := d.conns[dsn]
	if !ok {
		return c, fmt.Errorf("expected a connection to be available, but it is not")
	}

	c.opened++
	return c, nil
}
```
这里 mockDriver 实现了  Driver 接口：Open(name string) (Conn, error) ， 返回的是 *sqlmock(实现 driver.Conn 接口 )
```go
// go1.20/src/database/sql/driver/driver.go
type Conn interface {
	// Prepare returns a prepared statement, bound to this connection.
	Prepare(query string) (Stmt, error)

	// Close invalidates and potentially stops any current
	// prepared statements and transactions, marking this
	// connection as no longer in use.
	//
	// Because the sql package maintains a free pool of
	// connections and only calls Close when there's a surplus of
	// idle connections, it shouldn't be necessary for drivers to
	// do their own connection caching.
	//
	// Drivers must ensure all network calls made by Close
	// do not block indefinitely (e.g. apply a timeout).
	Close() error

	// Begin starts and returns a new transaction.
	//
	// Deprecated: Drivers should implement ConnBeginTx instead (or additionally).
	Begin() (Tx, error)
}
```



sqlmock 结构体介绍
```go
type sqlmock struct {
	ordered      bool
	dsn          string
	opened       int
	drv          *mockDriver
	converter    driver.ValueConverter
	queryMatcher QueryMatcher
	monitorPings bool

	expected []expectation // 匹配的信息
}
```


初始化
```go
func New(options ...func(*sqlmock) error) (*sql.DB, Sqlmock, error) {
	pool.Lock()
	// 这里的 dsn 同步加1 表示惟一 
	dsn := fmt.Sprintf("sqlmock_db_%d", pool.counter)
	pool.counter++

	smock := &sqlmock{dsn: dsn, drv: pool, ordered: true}
	pool.conns[dsn] = smock
	pool.Unlock()

	return smock.open(options)
}
```

```go
func (c *sqlmock) open(options []func(*sqlmock) error) (*sql.DB, Sqlmock, error) {
	// 指定 driver 是sqlmock 
	db, err := sql.Open("sqlmock", c.dsn)
	if err != nil {
		return db, c, err
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return db, c, err
		}
	}
	if c.converter == nil {
		// DefaultParameterConverter是ValueConverter接口的默认实现
		c.converter = driver.DefaultParameterConverter
	}
	if c.queryMatcher == nil {
		// 这里默认用 正则匹配器
		c.queryMatcher = QueryMatcherRegexp
	}

	if c.monitorPings {
		// We call Ping on the driver shortly to verify startup assertions by
		// driving internal behaviour of the sql standard library. We don't
		// want this call to ping to be monitored for expectation purposes so
		// temporarily disable.
		c.monitorPings = false
		defer func() { c.monitorPings = true }()
	}
	return db, c, db.Ping()
}
```




写匹配的sql : 正则作为案例开始匹配
```go
rows := sqlmock.NewRows([]string{"id", "title", "body"}).
    AddRow(1, "post 1", "hello").
    AddRow(2, "post 2", "world")
mock.ExpectQuery("^SELECT (.+) FROM posts$").WillReturnRows(rows)
```

匹配条件信息
```go
// go-sqlmock@v1.5.2/sqlmock.go
func (c *sqlmock) ExpectQuery(expectedSQL string) *ExpectedQuery {
	e := &ExpectedQuery{}
	e.expectSQL = expectedSQL
	e.converter = c.converter
	c.expected = append(c.expected, e)
	return e
}
```

匹配成功返回信息
```go
// 表明返回的列标题
func NewRows(columns []string) *Rows {
	return &Rows{
		cols:      columns,
		nextErr:   make(map[int]error),
		converter: driver.DefaultParameterConverter, 
	}
}

// 表明返回的列具体值
func (r *Rows) AddRow(values ...driver.Value) *Rows {
	// 判断列数等于数值数量
	if len(values) != len(r.cols) {
		panic(fmt.Sprintf("Expected number of values to match number of columns: expected %d, actual %d", len(values), len(r.cols)))
	}

	row := make([]driver.Value, len(r.cols))
	for i, v := range values {
		// Convert user-friendly values (such as int or driver.Valuer)
		// to database/sql native value (driver.Value such as int64)
		var err error
		v, err = r.converter.ConvertValue(v)
		if err != nil {
			panic(fmt.Errorf(
				"row #%d, column #%d (%q) type %T: %s",
				len(r.rows)+1, i, r.cols[i], values[i], err,
			))
		}

		row[i] = v
	}

	r.rows = append(r.rows, row)
	return r
}
```


```go
func (e *ExpectedQuery) WillReturnRows(rows ...*Rows) *ExpectedQuery {
	defs := 0
	sets := make([]*Rows, len(rows))
	for i, r := range rows {
		sets[i] = r
		if r.def != nil {
			defs++
		}
	}
	if defs > 0 && defs == len(sets) {
		e.rows = &rowSetsWithDefinition{&rowSets{sets: sets, ex: e}}
	} else {
		e.rows = &rowSets{sets: sets, ex: e}
	}
	return e
}
```




开启执行 sql 




