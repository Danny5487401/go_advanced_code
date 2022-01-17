# Prometheus
Prometheus将所有数据存储为时间序列

## 基本概念

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


### 指标类型
Prometheus client库提供四种核心度量标准类型。注意是客户端。Prometheus服务端没有区分类型，将所有数据展平为无类型时间序列.

1. Counter就是一个计数器，表示一种累积型指标，该指标只能单调递增或在重新启动时重置为零，例如，您可以使用计数器来表示所服务的请求数，已完成的任务或错误.

2. Gauge是最简单的度量类型，只有一个简单的返回值，可增可减，也可以set为指定的值。所以Gauge通常用于反映当前状态，比如当前温度或当前内存使用情况；当然也可以用于“可增加可减少”的计数指标。

3. Histogram主要用于在设定的分布范围内(Buckets)记录大小或者次数。
   例如http请求响应时间：0-100ms、100-200ms、200-300ms、>300ms 的分布情况，Histogram会自动创建3个指标，分别为：

- 事件发送的总次数<basename>_count：比如当前一共发生了2次http请求
- 所有事件产生值的大小的总和<basename>_sum：比如发生的2次http请求总的响应时间为150ms
- 事件产生的值分布在bucket中的次数<basename>_bucket{le="上限"}：比如响应时间0-100ms的请求1次，100-200ms的请求1次，其他的0次

4. Summary：数据分布统计图
   ummary和Histogram类似，都可以统计事件发生的次数或者大小，以及其分布情况。

Summary和Histogram都提供了对于事件的计数_count以及值的汇总_sum，因此使用_count,和_sum时间序列可以计算出相同的内容。

同时Summary和Histogram都可以计算和统计样本的分布情况，比如中位数，n分位数等等。不同在于Histogram可以通过histogram_quantile函数在服务器端计算分位数。
而Sumamry的分位数则是直接在客户端进行定义。因此对于分位数的计算。 Summary在通过PromQL进行查询时有更好的性能表现，而Histogram则会消耗更多的资源。相对的对于客户端而言Histogram消耗的资源更少。


### 作业和实例
在Prometheus中，一个可以拉取数据的端点IP:Port叫做一个实例（instance），而具有多个相同类型实例的集合称作一个作业（job）
```yaml
 - job: api-server
      - instance 1: 1.2.3.4:5670
      - instance 2: 1.2.3.4:5671
      - instance 3: 5.6.7.8:5670
      - instance 4: 5.6.7.8:5671

```

当Prometheus拉取指标数据时，会自动生成一些标签（label）用于区别抓取的来源：

- job：配置的作业名；
- instance：配置的实例名，若没有实例名，则是抓取的IP:Port

对于每一个实例（instance）的抓取，Prometheus会默认保存以下数据：

- up{job="<job>", instance="<instance>"}：如果实例是健康的，即可达，值为1，否则为0；
- scrape_duration_seconds{job="<job>", instance="<instance>"}：抓取耗时；
- scrape_samples_post_metric_relabeling{job="<job>", instance="<instance>"}：指标重新标记后剩余的样本数。
- scrape_samples_scraped{job="<job>", instance="<instance>"}：实例暴露的样本数
该up指标对于监控实例健康状态很有用。

