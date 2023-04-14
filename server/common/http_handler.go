package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

/*
 * RespondWithError 通过 HTTP 响应返回 JSON 格式的错误信息
 * @params: w http.ResponseWriter：表示 HTTP 响应的输出流
 * 			code int：表示 HTTP 状态码，例如 404、500 等
 * 			message string：表示错误信息的字符串
 * @return: null
 */
func RespondWithError(w http.ResponseWriter, code int, message string) {

	// 设置响应头部信息
	w.Header().Set("Content-Type", "application/json")
	// 构造响应体
	response := map[string]interface{}{
		"status": http.StatusText(code),
		"code":   strconv.Itoa(code),
		"msg":    message,
	}
	// 将响应体编码为 JSON 格式，并写入 ResponseWriter
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}

	// 设置 HTTP 状态码
	w.WriteHeader(code)
}

// ResponseBody 响应消息结构体
type ResponseBody struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Code   int         `json:"code"`
}

/*
 * RespondWithProper 封装HTTP响应写入器
 * @params: w http.ResponseWriter: HTTP响应写入器，用于将JSON格式的响应写入到HTTP响应中
 * 			status int: HTTP状态码，响应的状态
 * 			data interface{}: 要序列化为JSON格式响应数据。可以是任何Go类型
 * @return: 返回值 void 或 error 类型，如果在生成JSON响应时发生错误，则返回一个非空的 error 对象
 */
func RespondWithProper(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	resp := ResponseBody{
		Status: http.StatusText(status),
		Data:   data,
		Code:   status,
	}

	// 序列化ResponseBody结构体，并检查是否存在错误
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	_, _ = w.Write(jsonResp)
}
