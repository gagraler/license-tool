package utils

import (
	"encoding/json"
	"net/http"
)

// ResponseBody 响应消息结构体
type ResponseBody struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Code   int         `json:"code"`
}

/*
 * RespondWithJSON 封装HTTP响应写入器
 * @params: w http.ResponseWriter：HTTP响应写入器，用于将JSON格式的响应写入到HTTP响应中
 * 			status int：HTTP状态码，响应的状态
 * 			data interface{}：要序列化为JSON格式响应数据。可以是任何Go类型
 * @return: 返回值 void 或 error 类型，如果在生成JSON响应时发生错误，则返回一个非空的 error 对象
 */
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(jsonResp)
}
