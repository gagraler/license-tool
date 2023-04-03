package main

import (
	"net/http"
	"server/router"
)

func main() {

	r := router.SetupRouter()
	if r == nil {
		// 路由器配置失败，无法启动服务器
		return
	}

	// 启动服务器
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
