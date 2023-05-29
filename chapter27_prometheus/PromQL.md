<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [PromQL](#promql)
  - [表达式类型](#%E8%A1%A8%E8%BE%BE%E5%BC%8F%E7%B1%BB%E5%9E%8B)
  - [范围查询](#%E8%8C%83%E5%9B%B4%E6%9F%A5%E8%AF%A2)
  - [Offset modifier时间位移操作](#offset-modifier%E6%97%B6%E9%97%B4%E4%BD%8D%E7%A7%BB%E6%93%8D%E4%BD%9C)
  - [聚合操作](#%E8%81%9A%E5%90%88%E6%93%8D%E4%BD%9C)
  - [标量(Scalar)和字符串(String)](#%E6%A0%87%E9%87%8Fscalar%E5%92%8C%E5%AD%97%E7%AC%A6%E4%B8%B2string)
  - [操作符](#%E6%93%8D%E4%BD%9C%E7%AC%A6)
    - [Binary operator precedence操作符优先级](#binary-operator-precedence%E6%93%8D%E4%BD%9C%E7%AC%A6%E4%BC%98%E5%85%88%E7%BA%A7)
    - [集合运算法](#%E9%9B%86%E5%90%88%E8%BF%90%E7%AE%97%E6%B3%95)
    - [匹配模式](#%E5%8C%B9%E9%85%8D%E6%A8%A1%E5%BC%8F)
      - [1. 一对一 匹配模式会从操作符两边表达式获取的瞬时向量依次比较并找到唯一匹配(标签完全一致)的样本值](#1-%E4%B8%80%E5%AF%B9%E4%B8%80-%E5%8C%B9%E9%85%8D%E6%A8%A1%E5%BC%8F%E4%BC%9A%E4%BB%8E%E6%93%8D%E4%BD%9C%E7%AC%A6%E4%B8%A4%E8%BE%B9%E8%A1%A8%E8%BE%BE%E5%BC%8F%E8%8E%B7%E5%8F%96%E7%9A%84%E7%9E%AC%E6%97%B6%E5%90%91%E9%87%8F%E4%BE%9D%E6%AC%A1%E6%AF%94%E8%BE%83%E5%B9%B6%E6%89%BE%E5%88%B0%E5%94%AF%E4%B8%80%E5%8C%B9%E9%85%8D%E6%A0%87%E7%AD%BE%E5%AE%8C%E5%85%A8%E4%B8%80%E8%87%B4%E7%9A%84%E6%A0%B7%E6%9C%AC%E5%80%BC)
      - [2. 多对一和一对多](#2-%E5%A4%9A%E5%AF%B9%E4%B8%80%E5%92%8C%E4%B8%80%E5%AF%B9%E5%A4%9A)
  - [聚合操作](#%E8%81%9A%E5%90%88%E6%93%8D%E4%BD%9C-1)
  - [函数](#%E5%87%BD%E6%95%B0)
    - [计算counter指标增长率](#%E8%AE%A1%E7%AE%97counter%E6%8C%87%E6%A0%87%E5%A2%9E%E9%95%BF%E7%8E%87)
    - [动态标签替换](#%E5%8A%A8%E6%80%81%E6%A0%87%E7%AD%BE%E6%9B%BF%E6%8D%A2)
  - [案例](#%E6%A1%88%E4%BE%8B)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


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