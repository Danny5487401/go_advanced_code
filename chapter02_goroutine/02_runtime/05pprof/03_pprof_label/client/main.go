package main

import (
	"github.com/dlsniper/debugger"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	fakeTraffic()
}

func fakeTraffic() {
	// Wait for the server to start
	time.Sleep(1 * time.Second)

	pages := []string{"/", "/login", "/logout", "/products", "/product/{productID}", "/basket", "/about"}

	activeConns := make(chan struct{}, 10)

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	i := int64(0)

	for {
		activeConns <- struct{}{}
		i++

		page := pages[rand.Intn(len(pages))]

		// We need to launch this using a closure function to
		// ensure that we capture the correct value for the
		// two parameters we need: page and i
		go func(p string, rid int64) {
			// 方式一：不推荐，生产总影响少量性能
			//labels := pprof.Labels("request", "automated", "page", p, "rid", strconv.Itoa(int(rid)))
			//pprof.Do(context.Background(), labels, func(_ context.Context) {
			//	makeRequest(activeConns, c, p, rid)
			//})

			makeRequest(activeConns, c, p, rid)
		}(page, i)
	}
}

func makeRequest(done chan struct{}, c *http.Client, page string, i int64) {
	defer func() {
		// Unblock the next request from the queue
		<-done
	}()

	// 方式二：测试中go tools参数添加 -tags debugger。否则，该库将加载生产代码，标签将不起作用
	debugger.SetLabels(func() []string {
		return []string{
			"request", "automated",
			"page", page,
			"rid", strconv.Itoa(int(i)),
		}
	})

	page = strings.Replace(page, "{productID}", "abc-"+strconv.Itoa(int(i)), -1)
	r, err := http.NewRequest(http.MethodGet, "http://localhost:8080"+page, nil)
	if err != nil {
		return
	}

	resp, err := c.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(ioutil.Discard, resp.Body)

	time.Sleep(time.Duration(10+rand.Intn(40)) + time.Millisecond)
}
