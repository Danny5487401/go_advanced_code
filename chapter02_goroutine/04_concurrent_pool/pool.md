#协程池：
	能够达到协程资源复用。
#案例：
	有基于链表实现的Tidb，有基于环形队列实现的Jaeger，有基于数组栈实现的FastHTTP等
#分类：
	1。提前创建协程：Jaeger，Istio，Tars等。
	2。按需创建协程：Tidb，FastHTTP，Ants等。
