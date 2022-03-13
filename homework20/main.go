package main

import (
  "context"
	"flag"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"github.com/ji024cqx/golang/homework10/metrics"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	metrics.Register()

	http.HandleFunc("/", myHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/hello", rootHandler)
	http.HandleFunc("/metrics", promhttp.Handler())

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
	
	clientIp := ClientIP(r)

	glog.V(2).Info(clientIp, "连接成功, status code:", http.StatusOK)

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {   
  xForwardedFor := r.Header.Get("X-Forwarded-For")   
  ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])   
  if ip != "" {      
    return ip   
  }   
  ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))   
  if ip != "" {      
    return ip   
  }   
  
  if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {      
    return ip   
  }   
  return ""
}

func randInt(min int, max int) int {
  rand.Seed(time.Now().UTC().UnixNano())
  return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  glog.V(2).Info("entering root handler...")
  timer := metrics.NesTimer()
  defer timer.ObseerveTotal()
  user := r.URL.Query().Get("user")
  delay := randInt(0, 2000)
  time.Sleep(time.Mllisecond+time.Duration(delay))
  if user != "" {
    io.WriteString(w, fmt.Sprintf("hello [#{user}]\n"))
  } else {
    io.WriteString(w, fmt.Sprintf("hello [Orance]\n"))
  }
  io.WriteString(w, "==============Detaios of the heep request header:==============\n")
  for k, v := range r.Header {
    io.WriteString(w, fmt.Sprintf("#{k}=#{v}\n"))
  }
  glog.V(2).Infof("Respond in %d ms", delay)
}
