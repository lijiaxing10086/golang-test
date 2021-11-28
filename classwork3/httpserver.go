package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
)

//初始化配置
var myconf Myconf

//用于读取本地config文件的配置中的json，实现配置变更与代码逻辑变更分离
type Myconf struct {
	//目前仅有一个timeout的参数设置
	TimeoutTime string `json:"timeout"`
	MyHttpPort  string `json:"httpport"`
}

//初始化函数，读取配置文件，获取配置文件中的超时配置
func init() {
	//获取环境变量设置的config的路径，没有时设置默认值
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config"
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &myconf)
}

func main() {

	//解析flag参数，供glog使用
	flag.Parse()
	defer glog.Flush()

	//处理路由为 /healthz 的健康检查方法
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("Client ip is : %s , the httpcode is : %d\n", r.RemoteAddr, http.StatusOK)
		fmt.Fprintln(w, "200")
	})

	//带有超时设置
	//timeout
	http.HandleFunc("/timeout", TimeoutHandler)

	//处理路由为 /hellow 的健康检查方法
	http.HandleFunc("/hellow", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		version := os.Getenv("VERSION")
		if version == "" {
			version = "unknow"
		}
		w.Header().Set("VERSION", version)
		fmt.Printf("Client ip is : %s , the httpcode is : %d", r.RemoteAddr, http.StatusOK)
		fmt.Fprintln(w, "hello word")
	})

	path := ":" + myconf.MyHttpPort

	//监听端口
	http.ListenAndServe(path, nil)
}

//实现带有超时context的handler
func TimeoutHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("Request is coming")
	c := make(chan int)
	//timeoutContext
	var ctx context.Context
	var cancel context.CancelFunc

	//设置超时配置，已请求中的最优先，请求中没有则获取是否存在于config中
	if qp := r.URL.Query().Get("timeout"); len(qp) > 0 {
		timeout, err := time.ParseDuration(qp)
		if err != nil {
			fmt.Println("err is :", err)
			ctx, cancel = context.WithCancel(r.Context())
		} else {
			ctx, cancel = context.WithTimeout(r.Context(), timeout)
		}
	} else if myconf.TimeoutTime != "" {
		timeout, err := time.ParseDuration(myconf.TimeoutTime)
		if err != nil {
			fmt.Println("err is :", err)
			ctx, cancel = context.WithCancel(r.Context())
		} else {
			ctx, cancel = context.WithTimeout(r.Context(), timeout)
		}
	} else {
		ctx, cancel = context.WithCancel(r.Context())
	}

	defer cancel()

	//对请求进行5s的等待，模拟后端处理数据的延迟
	go func() {
		time.Sleep(5 * time.Second)
		c <- 1
	}()

	//区分出等待结束和超时时返回的信息
	select {
	case <-ctx.Done():
		glog.Errorf("Context interrupt or timeout: %v\n", ctx.Err())
		fmt.Fprintln(w, "timeout")
		return
	case value := <-c:
		fmt.Fprintln(w, "hello---this is timeout method")
		glog.Info("time is ready", value)
	}
}
