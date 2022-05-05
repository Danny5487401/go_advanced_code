# Prometheus
![](.prometheus_images/prometheus_structure.png)
- Exporter :收集系统或进程信息，转换为 Prometheus 可以识别的数据指标，以 http 或 https 服务的方式暴露给 Prometheus
- AlertManager ：接收 Prometheus 推送过来的告警信息，通过告警路由，向集成的组件 / 工具发送告警信息.
- Prometheus Server: 负责实现对监控数据的获取，存储以及查询。Prometheus Server可以通过静态配置管理监控目标，也可以配合使用Service Discovery的方式动态管理监控目标，并从这些监控目标中获取数据。
  其次Prometheus Server需要对采集到的监控数据进行存储，Prometheus Server本身就是一个时序数据库，将采集到的监控数据按照时间序列的方式存储在本地磁盘当中。
  最后Prometheus Server对外提供了自定义的PromQL语言，实现对数据的查询以及分析。
- PushGateway: 用于支持临时任务的推送网关

## 基本概念

### 样本
```css
<--------------- metric ---------------------><-timestamp -><-value->
http_request_total{status="200", method="GET"}@1434417560938 => 94355
http_request_total{status="200", method="GET"}@1434417561287 => 94334

http_request_total{status="404", method="GET"}@1434417560938 => 38473
http_request_total{status="404", method="GET"}@1434417561287 => 38544

http_request_total{status="200", method="POST"}@1434417560938 => 4748
http_request_total{status="200", method="POST"}@1434417561287 => 4785
```
在时间序列中的每一个点称为一个样本（sample），样本由以下三部分组成：

- 指标（metric）：指标名称和描述当前样本特征的 labelsets；
- 时间戳（timestamp）：一个精确到毫秒的时间戳；
- 样本值（value）： 一个 float64 的浮点型数据表示当前样本的值

### 指标名和标签
每个时间序列都由指标名和一组键值对（也称为标签）唯一标识。
metric的格式如下：
```css
<metric name>{<label name>=<label value>, ...}
```
例如：
```css
http_requests_total{host="192.10.0.1", method="POST", handler="/messages"}
```
- http_requests_total是指标名；
- host、method、handler是三个标签(label)，也就是三个维度；
- 查询语句可以基于这些标签or维度进行过滤和聚合

在Prometheus的底层实现中指标名称实际上是以__name__=<metric name>的形式保存在数据库中的，因此上面的一条time-series等同于：
```css
{__name__="api_http_requests_total"，host="192.10.0.1",method="POST", handler="/messages"}
```


### 指标类型
Prometheus client库提供四种核心度量标准类型。注意是客户端。Prometheus服务端没有区分类型，将所有数据展平为无类型时间序列.

#### 1. Counter计数器
表示一种累积型指标，该指标只能单调递增或在重新启动时重置为零，例如，您可以使用计数器来表示所服务的请求数，已完成的任务或错误.

Counter 类型数据可以让用户方便的了解事件产生的速率的变化，在 PromQL 内置的相关操作函数可以提供相应的分析，比如以 HTTP 应用请求量来进行说明：
```shell
//通过rate()函数获取HTTP请求量的增长率
rate(http_requests_total[5m])
//查询当前系统中，访问量前10的HTTP地址
topk(10, http_requests_total)
```

#### 2. Gauge仪表盘
是最简单的度量类型，只有一个简单的返回值，可增可减，也可以set为指定的值。所以Gauge通常用于反映当前状态，比如当前温度或当前内存使用情况；当然也可以用于“可增加可减少”的计数指标。

对于 Gauge 类型的监控指标，通过 PromQL 内置函数 delta() 可以获取样本在一段时间内的变化情况，例如，计算 CPU 温度在两小时内的差异：
```shell
dalta(cpu_temp_celsius{host="zeus"}[2h])
```

#### 3. Histogram直方图

如果大多数 API 请求都维持在 100ms 的响应时间范围内，而个别请求的响应时间需要 5s，那么就会导致某些 WEB 页面的响应时间落到中位数的情况，而这种现象被称为长尾问题。

为了区分是平均的慢还是长尾的慢，最简单的方式就是按照请求延迟的范围进行分组

主要用于在设定的分布范围内(Buckets)记录大小或者次数。

例如http请求响应时间：0-100ms、100-200ms、200-300ms、>300ms 的分布情况，Histogram会自动创建3个指标，分别为：

- 事件发送的总次数<basename>_count：比如当前一共发生了2次http请求
- 所有事件产生值的大小的总和<basename>_sum：比如发生的2次http请求总的响应时间为150ms
- 事件产生的值分布在bucket中的次数<basename>_bucket{le="上限"}：比如响应时间0-100ms的请求1次，100-200ms的请求1次，其他的0次

#### 4. Summary摘要
与 Histogram 类型类似，用于表示一段时间内的数据采样结果（通常是请求持续时间或响应大小等），但它直接存储了分位数（通过客户端计算，然后展示出来），而不是通过区间来计算

Summary和Histogram都提供了对于事件的计数_count以及值的汇总_sum，因此使用_count,和_sum时间序列可以计算出相同的内容。

现在可以总结一下 Histogram 与 Summary 的异同：
- 它们都包含了 <basename>_sum 和 <basename>_count 指标
- Histogram 需要通过 <basename>_bucket 来计算分位数，而 Summary 则直接存储了分位数的值。

除了四种基本的数据指标类型外，Prometheus 数据模型的一个非常重要的部分是沿着称为 “标签” 的维度对数据指标样本进行划分，这就产生了数据指标向量 (metric vectors).
prometheus 分为为四种基本数据指标类型提供了相应的数据指标向量，分别是 CounterVec,GaugeVec,HistogramVec 和 SummaryVec.
```go
Collect(ch chan<- Metric)  // 实现 Collector 接口的 Collect() 方法
Describe(ch chan<- *Desc)  // 实现 Collector 接口的 Describe() 方法

CurryWith(labels Labels)  // 返回带有指定标签的向量指标及可能发生的错误.多用于 promhttp 包中的中间件.
Delete(labels Labels)  // 删除带有指定标签的向量指标.如果删除了指标,返回 true
DeleteLabelValues(lvs ...string)  // 删除带有指定标签和标签值的向量指标.如果删除了指标,返回 true
GetMetricWith(labels Labels)  // 返回带有指定标签的数据指标及可能发生的错误
GetMetricWithLabelValues(lvs ...string)  // 返回带有指定标签和标签值的数据指标及可能发生的错误
MustCurryWith(labels Labels)  // 与 CurryWith 相同,但如果出现错误,则引发 panics
Reset()  // 删除此指标向量中的所有数据指标
With(labels Labels)  // 与 GetMetricWithLabels 相同,但如果出现错误,则引发 panics
WithLabelValues(lvs ...string)  // 与 GetMetricWithLabelValues 相同,但如果出现错误,则引发 panics
```

### 作业和实例
在prometheus.yml配置文件中，添加如下配置
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
```

当前在每一个Job中主要使用了静态配置(static_configs)的方式定义监控目标。除了静态配置每一个Job的采集Instance地址以外，Prometheus还支持与DNS、Consul、E2C、Kubernetes等进行集成实现自动发现Instance实例，并从这些Instance上获取监控数据。

在Prometheus中，一个可以拉取数据的端点IP:Port叫做一个实例（instance），而具有多个相同类型实例的集合称作一个作业（job）
```yaml
- job: api-server
  - instance 1: 1.2.3.4:5670
  - instance 2: 1.2.3.4:5671
  - instance 3: 5.6.7.8:5670
  - instance 4: 5.6.7.8:5671

```
在Prometheus中，每一个暴露监控样本数据的HTTP服务称为一个实例。例如在当前主机上运行的node exporter可以被称为一个实例(Instance)


当Prometheus拉取指标数据时，会自动生成一些标签（label）用于区别抓取的来源：
![](.prometheus_images/target_in_ui.png)
- job：配置的作业名；
- instance：配置的实例名，若没有实例名，则是抓取的IP:Port

对于每一个实例（instance）的抓取，Prometheus会默认保存以下数据：

- up{job="<job>", instance="<instance>"}：如果实例是健康的，即可达，值为1，否则为0；
- scrape_duration_seconds{job="<job>", instance="<instance>"}：抓取耗时；
- scrape_samples_post_metric_relabeling{job="<job>", instance="<instance>"}：指标重新标记后剩余的样本数。
- scrape_samples_scraped{job="<job>", instance="<instance>"}：实例暴露的样本数
该up指标对于监控实例健康状态很有用。

## 数据模型
Prometheus 所有采集的监控数据均以指标（metric）的形式保存在内置的时间序列数据库当中（TSDB）：属于同一指标名称，同一标签集合的、有时间戳标记的数据流。
除了存储的时间序列，Prometheus 还可以根据查询请求产生临时的、衍生的时间序列作为返回结果。


## 部署
使用prometheus的docker环境
```yaml
# prometheus.yml
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
- job_name: "go-test"
  scrape_interval: 60s
  scrape_timeout: 60s
  metrics_path: "/metrics"

  static_configs:
  - targets: ["localhost:8888"]

```

可以看到配置文件中指定了一个job_name，所要监控的任务即视为一个job, scrape_interval和scrape_timeout是pro进行数据采集的时间间隔和频率，metrics_path指定了访问数据的http路径，target是目标的ip:port,这里使用的是同一台主机上的8888端口

```shell
docker run -p 9090:9090 -v /Users/python/Desktop/go_advanced_code/chapter02_goroutine/02_runtime/07prometheus/client/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```
![](.prometheus_images/prometheus_panel.png)
启动之后可以访问web页面http://localhost:9090/graph,在status下拉菜单中可以看到配置文件和目标的状态，此时目标状态为DOWN，因为我们所需要监控的服务还没有启动起来，那就赶紧步入正文，用pro golang client来实现程序吧。

![](.prometheus_images/server_state.png)

启动后状态
![](.prometheus_images/server_state2.png)

##  prometheus/client_golang 源码分析
![](.prometheus_images/client_golang_structure.png)
### collector
接口定义
```go
type Collector interface {
    // Describe 暴露全部可能的 Metric 描述列表
	Describe(chan<- *Desc)

	// 获取采样数据，然后通过 HTTP 接口暴露给 Prom Server
	Collect(chan<- Metric)
}
```
### counter相关函数
```go
func (c *counter) Add(v float64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in value"))
	}

	ival := uint64(v)
	if float64(ival) == v {
		atomic.AddUint64(&c.valInt, ival)
		return
	}

	for {
		oldBits := atomic.LoadUint64(&c.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&c.valBits, oldBits, newBits) {
			return
		}
	}
}
```
Add 中修改共享数据时采用了“无锁”实现，相比“有锁 (Mutex)”实现可以更充分利用多核处理器的并行计算能力，性能相比加 Mutex 的实现会有很大提升


### opts
```go
// 其中 GaugeOpts, CounterOpts 实际上均为 Opts 的别名
type CounterOpts Opts
type GaugeOpts Opts
type Opts struct {
    // Namespace, Subsystem, and Name  是 Metric 名称的组成部分(通过 "_" 将这些组成部分连接起来),只有 Name 是必需的.
	// strings.Join([]string{namespace, subsystem, name}, "_")
    Namespace string
    Subsystem string
    Name      string

    // Help 提供 Metric 的信息.具有相同名称的 Metric 必须具有相同的 Help 信息
    Help string

    // ConstLabels 用于将固定标签附加到该指标.很少使用.
    ConstLabels Labels
}

type HistogramOpts struct {
    // Namespace, Subsystem, and Name  是 Metric 名称的组成部分(通过 "_" 将这些组成部分连接起来),只有 Name 是必需的.
    Namespace string
    Subsystem string
    Name      string

    //  Help 提供 Metric 的信息.具有相同名称的 Metric 必须具有相同的 Help 信息
    Help string

    // ConstLabels 用于将固定标签附加到该指标.很少使用.
    ConstLabels Labels

    // Buckets 定义了观察值的取值区间.切片中的每个元素值都是区间的上限,元素值必须按升序排序.
    // Buckets 会隐式添加 `+Inf` 值作为取值区间的最大值
    // 默认值是 DefBuckets =  []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
    Buckets []float64
}

type SummaryOpts struct {
    // Namespace, Subsystem, and Name  是 Metric 名称的组成部分(通过 "_" 将这些组成部分连接起来),只有 Name 是必需的.
    Namespace string
    Subsystem string
    Name      string

    //  Help 提供 Metric 的信息.具有相同名称的 Metric 必须具有相同的 Help 信息
    Help string

    // ConstLabels 用于将固定标签附加到该指标.很少使用.
    ConstLabels Labels

    // Objectives 定义了分位数等级估计及其各自的绝对误差.如果 Objectives[q] = e，则 q 报告的值将是 [q-e, q + e]之间某个 φ 的 φ 分位数
    // 默认值为空 map,表示没有分位数的摘要
    Objectives map[float64]float64

    // MaxAge 定义观察值与摘要保持相关的持续时间.必须是正数.默认值为 DefMaxAge = 10 * time.Minute
    MaxAge time.Duration

    // AgeBuckets 用于从摘要中排除早于 MaxAge 的观察值的取值区间.默认值为 DefAgeBuckets = 5
    AgeBuckets uint32

    // BufCap 定义默认样本流缓冲区大小.默认值为 DefBufCap = 500.
    BufCap uint32
}
```

### WithLabelValues方法
1个指标由Metric name + Labels共同确定。

若Metric name相同，但Label的值不同，则是不同的Metric。

例如：http_request_count(endpoint="hello")，http_request_count(endpoint="world")是两个不同的指标

```go
// @Param lvs 表示label values
func (v *CounterVec) WithLabelValues(lvs ...string) Counter {
	c, err := v.GetMetricWithLabelValues(lvs...) // 根据label的值来找对应的Metric
	if err != nil {
		panic(err)
	}
	return c
}

// 根据label的值来找对应的Metric
// @Param lvs表示label value
func (v *CounterVec) GetMetricWithLabelValues(lvs ...string) (Counter, error) {
	metric, err := v.MetricVec.GetMetricWithLabelValues(lvs...)
	if metric != nil {
		return metric.(Counter), err
	}
	return nil, err
}

// 根据label值找对应的metric
func (m *MetricVec) GetMetricWithLabelValues(lvs ...string) (Metric, error) {
	h, err := m.hashLabelValues(lvs) // 获取label对应的hash值，非重点不展开讲，这块的核心是，若hash值一样，则对应的Metric是同一个
	if err != nil {
		return nil, err
	}
    // 根据hash值从metricMap里get对应的metric
    // 若不存在则新创建一个metric并放入到metricMap里
	return m.metricMap.getOrCreateMetricWithLabelValues(h, lvs, m.curry), nil
}
```
metricMap的结构，Metric最终存到一个map里，key=根据label值计算出的hash值，value=Metric元信息
```go
// metricMap定义，Exporter的Metric都存在这个结构中
type metricMap struct {
	mtx       sync.RWMutex // Protects metrics.
	metrics   map[uint64][]metricWithLabelValues // Metric最终存到一个map里，key=根据label值计算出的hash值，value=Metric元信息
	desc      *Desc
	newMetric func(labelValues ...string) Metric
}

type metricWithLabelValues struct {
	values []string // label的值
	metric Metric // Metric的meta信息
}
```


### Registry 的高级用法
prometheus 包提供了 MustRegister() 函数用于注册 Collector, 但如果注册过程中发生错误，程序会引发 panics. 而使用 Register() 函数可以实现注册 Collector 的同时处理可能发生的错误.


prometheus 通过 NewGoCollector() 和 NewProcessCollector() 函数创建 Go 运行时数据指标的 Collector 和进程数据指标的 Collector.

## Prometheus拉取Exporter的哪些数据
### promhttp 包
promhttp 包允许创建 http.Handler 实例通过 HTTP 公开 Prometheus 数据指标

```go
// Prometheus拉取的入口
http.Handle("/metrics", promhttp.Handler())

// http.go promhttp.Handler()
func Handler() http.Handler {
	return InstrumentMetricHandler(
		prometheus.DefaultRegisterer, HandlerFor(prometheus.DefaultGatherer, HandlerOpts{}),
	)
}


// http.go HandlerFor
func HandlerFor(reg prometheus.Gatherer, opts HandlerOpts) http.Handler {
	// 省略部分代码
	mfs, err := reg.Gather() // 收集Metric信息
	// 省略部分代码
}

// prometheus.DefaultGatherer
// registry.go 
var (
	defaultRegistry              = NewRegistry() // DefaultGatherer就是defaultRegistry
	DefaultRegisterer Registerer = defaultRegistry
	DefaultGatherer   Gatherer   = defaultRegistry
)

// registry.go 
// Gather implements Gatherer. 负责收集metrics信息
func (r *Registry) Gather() ([]*dto.MetricFamily, error) {
	// 省略部分代码
    // 声明Counter类型的Metric后，需要MustRegist注册到Registry，最终就是保存在collectorsByID里
    // Counter类型本身就是一个collector
	for _, collector := range r.collectorsByID {
		checkedCollectors <- collector
	}
	// 省略部分代码
	collectWorker := func() {
		for {
			select {
			case collector := <-checkedCollectors:
				collector.Collect(checkedMetricChan) // 执行Counter的Collect，见下文
			case collector := <-uncheckedCollectors:
				collector.Collect(uncheckedMetricChan)
			default:
				return
			}
			wg.Done()
		}
	}

	// 省略部分代码
}


// vec.go
// Collect implements Collector.
// Counter类型的Collect方法
func (m *MetricVec) Collect(ch chan<- Metric) { m.metricMap.Collect(ch) }

// vec.go
// Collect implements Collector.
// 返回metricMap里所有的Metric
func (m *metricMap) Collect(ch chan<- Metric) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	for _, metrics := range m.metrics {
		for _, metric := range metrics {
			ch <- metric.metric
		}
	}
}
```


## PromQL
Prometheus通过指标名称（metrics name）以及对应的一组标签（labelset）唯一定义一条时间序列。
指标名称反映了监控样本的基本标识，而label则在这个基本特征上为采集到的数据提供了多种特征维度。用户可以基于这些特征维度过滤，聚合，统计从而产生新的计算后的一条时间序列。


PromQL是Prometheus内置的数据查询语言，其提供对时间序列数据丰富的查询，聚合以及逻辑运算能力的支持。并且被广泛应用在Prometheus的日常应用当中，包括对数据查询、可视化、告警处理当中


### 表达式类型
- 瞬时向量表达式。
- 区间向量表达式
- 标量(Scalar)
- 字符串(String)

### 范围查询
直接通过类似于PromQL表达式http_requests_total查询时间序列时，返回值中只会包含该时间序列中的最新的一个样本值，这样的返回结果我们称之为瞬时向量。而相应的这样的表达式称之为瞬时向量表达式。

而如果我们想过去一段时间范围内的样本数据时，我们则需要使用区间向量表达式。
区间向量表达式和瞬时向量表达式之间的差异在于在区间向量表达式中我们需要定义时间选择的范围，时间范围通过时间范围选择器[]进行定义。例如，通过以下表达式可以选择最近5分钟内的所有样本数据：

```css
http_requests_total{}[5m]
```
除了使用m表示分钟以外，PromQL的时间范围选择器支持其它时间单位：
* s - 秒
* m - 分钟
* h - 小时
* d - 天
* w - 周
* y - 年

### Offset modifier时间位移操作
在瞬时向量表达式或者区间向量表达式中，都是以当前时间为基准
```css
http_request_total{} # 瞬时向量表达式，选择当前最新的数据
http_request_total{}[5m] # 区间向量表达式，选择以当前时间为基准，5分钟内的数据
```

如果我们想查询，5分钟前的瞬时样本数据，或昨天一天的区间内的样本数据呢? 这个时候我们就可以使用位移操作，位移操作的关键字为offset。

```promql
http_request_total{} offset 5m
http_request_total{}[1d] offset 1d
```

### 聚合操作
一般来说，如果描述样本特征的标签(label)在并非唯一的情况下，通过PromQL查询数据，会返回多条满足这些特征维度的时间序列。而PromQL提供的聚合操作可以用来对这些时间序列进行处理，形成一条新的时间序列
```css
# 查询系统所有http请求的总量
sum(http_request_total)

# 按照mode计算主机CPU的平均使用时间
avg(node_cpu) by (mode)

# 按照主机查询各个主机的CPU使用率
sum(sum(irate(node_cpu{mode!='idle'}[5m]))  / sum(irate(node_cpu[5m]))) by (instance)
```

### 标量(Scalar)和字符串(String)

1. 标量只有一个数字，没有时序，例如 10。

Note:当使用表达式count(http_requests_total)，返回的数据类型，依然是瞬时向量。用户可以通过内置函数scalar()将单个瞬时向量转换为标量

2. 直接使用字符串，作为PromQL表达式，则会直接返回字符串
```css
"this is a string"
'these are unescaped: \n \\ \t'
`these are not unescaped: \n ' " \t`
```


### 操作符
#### Binary operator precedence操作符优先级

查询主机的CPU使用率，可以使用表达式：
```css
100 * (1 - avg (irate(node_cpu{mode='idle'}[5m])) by(job) )
```
在PromQL操作符中优先级由高到低依次为：
- ^
- *, /, %
- +, -
- ==, !=, <=, <, >=, >
- and, unless
- or

#### 集合运算法
使用瞬时向量表达式能够获取到一个包含多个时间序列的集合，我们称为瞬时向量。 通过集合运算，可以在两个瞬时向量与瞬时向量之间进行相应的集合操作

- and (并且)
- or (或者)
- unless (排除)

使用解释
* vector1 and vector2 会产生一个由vector1的元素组成的新的向量。该向量包含vector1中完全匹配vector2中的元素组成。
* vector1 or vector2 会产生一个新的向量，该向量包含vector1中所有的样本数据，以及vector2中没有与vector1匹配到的样本数据。
* vector1 unless vector2 会产生一个新的向量，新向量中的元素由vector1中没有与vector2匹配的元素组成

#### 匹配模式
##### 1. 一对一 匹配模式会从操作符两边表达式获取的瞬时向量依次比较并找到唯一匹配(标签完全一致)的样本值
   
在操作符两边表达式标签不一致的情况下，可以使用on(label list)或者ignoring(label list）来修改便签的匹配行为。使用ignoreing可以在匹配时忽略某些便签。而on则用于将匹配行为限定在某些便签之内。
```css
<vector expr> <bin-op> ignoring(<label list>) <vector expr>
<vector expr> <bin-op> on(<label list>) <vector expr>
```   

当存在样本：
```css
method_code:http_errors:rate5m{method="get", code="500"}  24
method_code:http_errors:rate5m{method="get", code="404"}  30
method_code:http_errors:rate5m{method="put", code="501"}  3
method_code:http_errors:rate5m{method="post", code="500"} 6
method_code:http_errors:rate5m{method="post", code="404"} 21

method:http_requests:rate5m{method="get"}  600
method:http_requests:rate5m{method="del"}  34
method:http_requests:rate5m{method="post"} 120
```
使用PromQL表达式：返回在过去5分钟内，HTTP请求状态码为500的在所有请求中的比例
```css
method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m
```
如果没有使用ignoring(code)，操作符两边表达式返回的瞬时向量中将找不到任何一个标签完全相同的匹配项。
```css
# 结果
{method="get"}  0.04            //  24 / 600
{method="post"} 0.05            //   6 / 120
```
同时由于method为put和del的 样本 找不到匹配项，因此不会出现在结果当中。

##### 2. 多对一和一对多
多对一和一对多两种匹配模式指的是“一”侧的每一个向量元素可以与"多"侧的多个元素匹配的情况。
在这种情况下，必须使用group修饰符：group_left或者group_right来确定哪一个向量具有更高的基数（充当“多”的角色)

多对一和一对多两种模式一定是出现在操作符两侧表达式返回的向量标签不一致的情况。因此需要使用ignoring和on修饰符来排除或者限定匹配的标签列表。


```css
method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m
```
该表达式中，左向量method_code:http_errors:rate5m包含两个标签method和code。
而右向量method:http_requests:rate5m中只包含一个标签method，因此匹配时需要使用ignoring限定匹配的标签为code。 
在限定匹配标签后，右向量中的元素可能匹配到多个左向量中的元素 因此该表达式的匹配模式为多对一，需要使用group修饰符group_left指定左向量具有更好的基数。

```css
# 结果
{method="get", code="500"}  0.04            //  24 / 600
{method="get", code="404"}  0.05            //  30 / 600
{method="post", code="500"} 0.05            //   6 / 120
{method="post", code="404"} 0.175           //  21 / 120
```

### 聚合操作
内置的聚合操作符，这些操作符作用域瞬时向量。可以将瞬时表达式返回的样本数据进行聚合，形成一个新的时间序列。

* sum (求和)
* min (最小值)
* max (最大值)
* avg (平均值)
* stddev (标准差)
* stdvar (标准方差)
* count (计数)
* count_values (对value进行计数)
* bottomk (后n条时序)
* topk (前n条时序)
* quantile (分位数

使用聚合操作的语法如下：
```css
<aggr-op>([parameter,] <vector expression>) [without|by (<label list>)]
```
其中只有count_values, quantile, topk, bottomk支持参数(parameter)。

without用于从计算结果中移除列举的标签，而保留其它标签。by则正好相反，结果向量中只保留列出的标签，其余标签则移除。通过without和by可以按照样本的问题对数据进行聚合。

```css
sum(http_requests_total) without (instance)
```
等价于

```css
sum(http_requests_total) by (code,handler,job,method)
```


### 函数

#### 计算counter指标增长率
样本增长率反映出了样本变化的剧烈程度：
```css
increase(node_cpu[2m]) / 120
```
这里通过node_cpu[2m]获取时间序列最近两分钟的所有样本，increase计算出最近两分钟的增长量，最后除以时间120秒得到node_cpu样本在最近两分钟的平均增长率。并且这个值也近似于主机节点最近两分钟内的平均CPU使用率。

rate函数可以直接计算区间向量v在时间窗口内平均增长速率。因此，通过以下表达式可以得到与increase函数相同的结果：
```css
rate(node_cpu[2m])
```
需要注意的是使用rate或者increase函数去计算样本的平均增长速率，容易陷入“长尾问题”当中，其无法反应在时间窗口内样本数据的突发变化。
例如，对于主机而言在2分钟的时间窗口内，可能在某一个由于访问量或者其它问题导致CPU占用100%的情况，但是通过计算在时间窗口内的平均增长率却无法反应出该问题。


为了解决该问题，PromQL提供了另外一个灵敏度更高的函数irate(v range-vector)。irate同样用于计算区间向量的计算率，但是其反应出的是瞬时增长率。
irate函数是通过区间向量中最后两个样本数据来计算区间向量的增长速率。这种方式可以避免在时间窗口范围内的“长尾问题”，并且体现出更好的灵敏度，通过irate函数绘制的图标能够更好的反应样本数据的瞬时变化状态。
```css
irate(node_cpu[2m])
```

#### 动态标签替换
一般来说来说，使用PromQL查询到时间序列后，可视化工具会根据时间序列的标签来渲染图表。例如通过up指标可以获取到当前所有运行的Exporter实例以及其状态：
```css
up{instance="localhost:8080",job="cadvisor"}    1
up{instance="localhost:9090",job="prometheus"}    1
up{instance="localhost:9100",job="node"}    1
```

```go
label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string)
```
该函数会依次对v中的每一条时间序列进行处理，通过regex匹配src_label的值，并将匹配部分replacement写入到dst_label标签中。如下所示：
```css
label_replace(up, "host", "$1", "instance",  "(.*):.*")
```
结果
```css
up{host="localhost",instance="localhost:8080",job="cadvisor"}    1
up{host="localhost",instance="localhost:9090",job="prometheus"}    1
up{host="localhost",instance="localhost:9100",job="node"} 1
```



### 案例
合法表达式：所有的PromQL表达式都必须至少包含一个指标名称(例如http_request_total)，或者一个不会匹配到空字符串的标签过滤器(例如{code="200"})。
```css
http_request_total # 合法
http_request_total{} # 合法
{method="get"} # 合法
```
而如下表达式，则不合法：
```css
{job=~".*"} # 不合法
```
除了使用<metric name>{label=value}的形式以外，我们还可以使用内置的__name__标签来指定监控指标名称：
```css
{__name__=~"http_request_total"} # 合法
{__name__=~"node_disk_bytes_read|node_disk_bytes_written"} # 合法
```
> A match of env=~"foo" is treated as env=~"^foo$"
> 使用label=~regx表示选择那些标签符合正则表达式定义的时间序列；
> 反之使用label!~regx进行排除；

选用一个相对比较复杂的表达式为例：
```sql
sum(avg_over_time(go_goroutines{job=“prometheus”}[5m])) by (instance)
```
- sum(…) by (instance)：序列纵向分组合并序列（包含相同的 instance 会分配到一组）
- avg_over_time(…)
- go_goroutines{job=“prometheus”}[5m] 

## 参考链接
1. https://prometheus.fuckcloudnative.io/di-yi-zhang-jie-shao/overview#:~:text=%E2%80%8BPrometheus%20%E6%98%AF%E7%94%B1%E5%89%8D,%E4%BA%8E%E4%BB%BB%E4%BD%95%E5%85%AC%E5%8F%B8%E8%BF%9B%E8%A1%8C%E7%BB%B4%E6%8A%A4%E3%80%82
2. https://www.infoq.cn/article/prometheus-theory-source-code
3. https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/promql/prometheus-promql-functions
4. 官网https://prometheus.io/docs/prometheus/latest/querying/functions/
