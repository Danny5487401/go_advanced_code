package main

//不是伪代码，可以直接go run

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math"
	"net/http"
	"time"
)

// 2. 内置collector
var (
	// myTestCounter Counter (累加指标)
	myTestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		//因为Name不可以重复，所以建议规则为："部门名_业务名_模块名_标量名_类型"
		Name: "my_test_counter", //唯一id，不可重复Register()，可以Unregister()
		Help: "my test counter", //对此Counter的描述
	})
	// 通过 NewCounterVec() 方法创建带"标签"的 CounterVec 对象
	hdFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hd_errors_total",
		Help: "Number of hard-disk errors.",
	},
		[]string{"device"},
	)

	// myTestGauge Gauge (测量指标)
	myTestGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "my_test_gauge",
		Help: "my test gauge",
	})

	// myTestHistogram Histogram (直方图)
	myTestHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "my_test_histogram",
		Help:    "my test histogram",
		Buckets: prometheus.LinearBuckets(20, 5, 5), //第一个桶20起，每个桶间隔5，共5个桶。 所以20, 25, 30, 35, 40
	})

	// myTestSummary Summary (概略图)
	myTestSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "my_test_summary",
		Help:       "my test summary",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //返回五分数， 九分数， 九九分数
	})
)

func init() {
	//不能注册多次Name相同的Metrics
	//MustRegister注册失败将直接panic()，如果想捕获error，建议使用Register()
	prometheus.MustRegister(myTestCounter)
	prometheus.MustRegister(hdFailures)
	prometheus.MustRegister(myTestGauge)
	prometheus.MustRegister(myTestHistogram)
	prometheus.MustRegister(myTestSummary)
}

func main() {
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	go func() {
		var i float64
		for {
			i++
			myTestCounter.Add(10000)                                         //每次加常量
			myTestGauge.Add(i)                                               //每次加增量
			myTestHistogram.Observe(30 + math.Floor(120*math.Sin(i*0.1))/10) //每次观察一个18 - 42的量
			myTestSummary.Observe(30 + math.Floor(120*math.Sin(i*0.1))/10)

			time.Sleep(time.Second * 5)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:8888", nil)) //多个进程不可监听同一个端口
}
