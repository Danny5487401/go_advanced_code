# 映射Map、归约reduce与过滤filter
![](.introduction_images/.map-reduce.png)
## 介绍
- Map 是一对一的场景，是 循环中对数据加工处理
- Reduce 是多对一，是 数据聚合处理
- Filter是过滤的处理，是 数据有效性

## 场景
- 统计消费总额 - Reduce
- 统计用户A - Filter
- 统计本月 - Filter
- 费用转化为美金 - Map

1. 调用第三方接口的时候， 一个需求你需要调用不同的接口，做数据组装。
2. 一个应用首页可能依托于很多服务。那就涉及到在加载页面时需要同时请求多个服务的接口。这一步往往是由后端统一调用组装数据再返回给前端，也就是所谓的 BFF(Backend For Frontend) 层。

## [Go-zero优秀源码--mapReduceWithPanicChan](chapter10_function/02_advanced_function/04_mapReduce/mapReduce/mapReduce.go)
![](.introduction_images/mapReduceWithPanicChan.png)
官网推荐链接： https://go-zero.dev/cn/mapreduce.html

