package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"server/request"
)

func SetupRouter() http.Handler {

	// 创建新的路由器实例
	r := mux.NewRouter()

	// 为 "/generate_license" 路径注册生成许可证的处理函数
	r.HandleFunc("/generate_license", request.GetLicenseRequest).Methods("POST")
	return r
}
