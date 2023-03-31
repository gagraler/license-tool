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
	"license-tool/server/utils"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Authorized 授权详细信息
type Authorized struct {
	Id            string `json:"id"`
	License       string `json:"license"`
	Date          string `json:"date"`
	SignatureCode string `json:"signatureCode"`
	Type          string `json:"type"`
	Expiration    string `json:"expiration"`
	AllowedUsers  string `json:"usersNum"`
	Project       string `json:"project"`
	Module        string `json:"module"`
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
	expirationString := r.URL.Query().Get("expiration")
	usersNumberString := r.URL.Query().Get("usersNum")
	obj := r.URL.Query().Get("object")
	module := r.URL.Query().Get("module")

	// 校验 signatureCode 是否为 32 位纯数字
	signatureCode := r.URL.Query().Get("signatureCode")
	matched, err := regexp.MatchString(`^.{1,32}$`, signatureCode)
	if err != nil {
		http.Error(w, "Invalid signature code format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !matched {
		http.Error(w, "Invalid signature code format....", http.StatusBadRequest)
		return
	}

	// 校验输入参数并检查是否正确
	expiration, err := time.Parse("2006-01-02", expirationString)
	if err != nil {
		http.Error(w, "Invalid expiration date format: "+err.Error(), http.StatusBadRequest)
		return
	}

	usersNumber, err := strconv.ParseUint(usersNumberString, 10, 32)
	if err != nil {
		http.Error(w, "Invalid allowed users value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 使用输入参数调用GenerateLicense函数生成许可证
	license, err := service.GenerateLicense(signatureCode, licenseType, expiration, uint(usersNumber), obj, module)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建响应结构体
	msg := Msg{
		Authorized: Authorized{
			Id:            utils.GenerateUniqueID(),
			SignatureCode: license.SignatureCode,
			License:       license.LicenseID,
			Date:          license.Date.Format("2006-01-02 15:04:05"),
			Type:          license.Type,
			Expiration:    license.ExpirationDate.Format("2006-01-02"),
			AllowedUsers:  strconv.FormatUint(uint64(license.AllowedUsers), 10),
			Project:       license.Project,
			Module:        license.Module,
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

	// 以json数据的id字段作为文件名存储到文件中
	licenseFileName := msg.Authorized.Id + ".license"
	file, err := os.Create(licenseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// 对数据进行加密操作
	encrypted, err := utils.ObfuscationUtil(response, signatureCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将加密后的数据写入文件
	_, err = file.WriteString(encrypted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头并将JSON格式的响应写入HTTP响应写入器中
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(response))

}
