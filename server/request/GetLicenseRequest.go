/*
 * GetLicenseRequest 处理获取许可证请求的HTTP处理程序
 * @params: w http.ResponseWriter - HTTP响应写入器
 *			r *http.Request - HTTP请求指针
 * @returns: null
 */

package request

import (
	"encoding/json"
	"fmt"
	"license-tool/server/service"
	"net/http"
	"strconv"
	"time"
)

// Authorized 授权详细信息
type Authorized struct {
	License      string `json:"license"`
	Date         string `json:"date"`
	Type         string `json:"type"`
	Expiration   string `json:"expiration"`
	AllowedUsers string `json:"allowed_users"`
	Days         string `json:"days"`
	Project      string `json:"project"`
	Module       string `json:"module"`
}

// Msg 授权信息、状态和代码
type Msg struct {
	Authorized Authorized `json:"authorized"`
	Status     string     `json:"status"`
	Code       int        `json:"code"`
}

/*
 * GetLicenseRequest 处理获取许可证请求的HTTP处理程序
 * @params:  w http.ResponseWriter - HTTP响应写入器
 * 			 r *http.Request - HTTP请求指针
 * @returns: null
 */
func GetLicenseRequest(w http.ResponseWriter, r *http.Request) {

	licenseType := r.URL.Query().Get("type")
	expirationString := r.URL.Query().Get("authorization_deadline")
	allowedUsersString := r.URL.Query().Get("allowed_users")
	daysString := r.URL.Query().Get("days")
	obj := r.URL.Query().Get("object")
	module := r.URL.Query().Get("module")

	// 校验输入参数并检查是否正确
	expiration, err := time.Parse("2006-01-02", expirationString)
	if err != nil {
		http.Error(w, "Invalid expiration date format: "+err.Error(), http.StatusBadRequest)
		return
	}
	allowedUsers, err := strconv.ParseUint(allowedUsersString, 10, 32)
	if err != nil {
		http.Error(w, "Invalid allowed users value: "+err.Error(), http.StatusBadRequest)
		return
	}
	days, err := strconv.ParseUint(daysString, 10, 32)
	if err != nil {
		http.Error(w, "Invalid days value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 使用输入参数调用GenerateLicense函数生成许可证
	license, err := service.GenerateLicense(licenseType, expiration, uint(allowedUsers), uint(days), obj, module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建响应结构体
	msg := Msg{
		Authorized: Authorized{
			License:      license.LicenseID,
			Date:         license.Date.Format("2006-01-02 15:04:05"),
			Type:         license.Type,
			Expiration:   license.ExpirationDate.Format("2006-01-02"),
			AllowedUsers: strconv.FormatUint(uint64(license.AllowedUsers), 10),
			Days:         strconv.FormatUint(uint64(license.Days), 10),
			Project:      license.Project,
			Module:       license.Module,
		},
		Status: http.StatusText(200),
		Code:   http.StatusOK,
	}

	// 将响应结构体转换为JSON格式
	response, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头并将JSON格式的响应写入HTTP响应写入器中
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(response))
}
