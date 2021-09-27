package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	//处理路由为 /healthz 的健康检查方法
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Client ip is : %s , the httpcode is : %d", r.RemoteAddr, http.StatusOK)
		fmt.Fprintln(w, "200")
	})

	//处理路由为 / 的方法
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		version := os.Getenv("VERSION")
		if version == "" {
			version = "unknow"
		}
		w.Header().Set("VERSION", version)
		fmt.Printf("Client ip is : %s , the httpcode is : %d", r.RemoteAddr, http.StatusOK)
		fmt.Fprintln(w, "hello")
	})
	//监听8080端口
	http.ListenAndServe(":8080", nil)
}
