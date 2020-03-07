package main

import (
	"fmt"
	"time"
)

func Test(){
	start:=time.Now().Unix()
	fmt.Println(time.Now().Unix())
	time.Sleep(3*time.Second)
	fmt.Println(time.Now().Unix()-start)
}
