package main

import "fmt"

func main ()  {
	fmt.Println("####################")
	fmt.Println("# 1.http_flood     #")
	fmt.Println("# 2.syn_flood      #")
	fmt.Println("# 3.ack_flood      #")
	fmt.Println("# 4.syn&ack_flood  #")
	fmt.Println("####################")
	fmt.Println("please choose the model")
	var model int
	fmt.Scan(&model)
	if model==1{
		Run("http://127.0.0.1:9090")
	}else if model==2{
		SynStart(1)
	}else if model==3{
		SynStart(2)
	}else if model==4{
		SynStart(3)
	}else{
		judge:=HostAddrCheck("127.0.0.1:3456")
		fmt.Println(judge)
		workDns("127.0.0.1","23.23.23.23",45)
	}
}
