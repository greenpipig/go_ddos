package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)
var THREAD_NUM int64= 60 // 并发线程总数
var	ONE_WORKER_NUM int64= 100  //每个线程的循环次数
var	LOOP_SLEEP int64 = 1 // 每次请求时间间隔(毫秒)
var ERROR_NUM int//出错数

var wg sync.WaitGroup

func webPush(url string){
	_, err :=http.Get(url)
	if err != nil {
		ERROR_NUM++
	}
}

func working(url string){
	for i:=int64(0);i<ONE_WORKER_NUM;i++{
		webPush(url)
		time.Sleep(time.Duration(LOOP_SLEEP)*time.Microsecond)
	}
	defer wg.Done()
}

func Run(url string){
	wg.Add(int(THREAD_NUM))
	startTime:=time.Now().UnixNano()
	for i:=int64(0);i<THREAD_NUM;i++{
		go working(url)
	}
	wg.Wait()
	endTime:=time.Now().UnixNano()
	fmt.Println("total num is ",THREAD_NUM*ONE_WORKER_NUM)
	fmt.Println("total time is ",float64(float64(endTime-startTime)/1000000000))
	fmt.Println("qps is ",float64(float64(THREAD_NUM)*float64(ONE_WORKER_NUM)/(float64(endTime-startTime)/1000000000)))
	fmt.Println("true persecond is ",float64(1.0/float64(float64(THREAD_NUM)*float64(ONE_WORKER_NUM)/(float64(endTime-startTime)/1000000000))))
	fmt.Println("error num is ",ERROR_NUM)
}
