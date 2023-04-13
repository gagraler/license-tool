package main

import (
	"fmt"
	"net/http"
	"server/config"
	"server/router"
)

func init() {
	config.InitConfig()

}

func main() {

	host := config.GetConfig("server", "address")
	port := config.GetConfig("server", "port")

	r := router.SetupRouter()
	if r == nil {
		// 路由器配置失败，无法启动服务器
		return
	} else {
		fmt.Println("Starting Successful")
	}

	// 启动服务器
	if err := http.ListenAndServe(host+":"+port, r); err != nil {
		panic(err)
	}

}
