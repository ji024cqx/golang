package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")

	http.HandleFunc("/", myHandler)
	http.HandleFunc("/healthz", healthzHandler)

	err := http.ListenAndServe(":18023", nil)
	if err != nil {
		glog.V(2).Info("http.ListenAndServe err:", err)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method = ", r.Method) 	//请求方法
	//fmt.Println("URL = ", r.URL)		// 浏览器发送请求文件路径
	//fmt.Println("header = ", r.Header) // 请求头
	//fmt.Println("body = ", r.Body)		// 请求包体
	//fmt.Println(r.RemoteAddr, "连接成功")  	//客户端网络地址
	//fmt.Println(time.Now())
	for k, v := range r.Header {
		w.Header().Add(k, strings.Join(v, ""))
	}

	w.Header().Add("VERSION", os.Getenv("VERSION"))

	glog.V(2).Info(r.RemoteAddr, "连接成功, status code:", http.StatusOK)

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}
