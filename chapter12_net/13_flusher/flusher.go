package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`
		   <html>
			 <body>
		  `))

		w.(http.Flusher).Flush()

		// 这里对每次循环输出都进行Flush刷新输出
		for i := 0; i < 10; i++ {
			w.Write([]byte(fmt.Sprintf(`
				<h3>%d</h3>
		   `, i)))
			//w.Flush()
			w.(http.Flusher).Flush()
			time.Sleep(time.Duration(1) * time.Second)
		}

		w.Write([]byte(`
			 </body>
		   </html>
		  `))
		w.(http.Flusher).Flush()

	})
	http.ListenAndServe(":8080", nil)
}
