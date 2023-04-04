/*
 * package request http请求处理层，处理http请求
 */
package request

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"server/service"
	"server/utils"
	"strconv"
	"time"
)

// ResponseData 响应数据结构体
type ResponseData struct {
	ID           string `json:"id"`
	License      string `json:"license"`
	Date         string `json:"date"`
	Signature    string `json:"signatureCode"`
	Type         string `json:"type"`
	Expiration   string `json:"expiration"`
	AllowedUsers uint   `json:"usersNum,string"`
	Project      string `json:"project"`
	Module       string `json:"module"`
}

// ResponseMsg 响应信息结构体
type ResponseMsg struct {
	Status string       `json:"status"`
	Data   ResponseData `json:"data"`
	Code   int          `json:"code"`
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
	msgData := buildResponseData(license)

	// 将响应结构体转换为JSON格式
	response, err := json.Marshal(msgData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建一个以json数据的id字段作为文件名的文件
	licenseFileName := msgData.ID + ".license"
	file, err := createFile(licenseFileName)
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

	utils.RespondWithJSON(w, http.StatusOK, msgData)
}

// buildResponseData 构建响应结构体
func buildResponseData(license *service.License) ResponseData {
	return ResponseData{
		ID:           utils.GenerateUniqueID(),
		Signature:    license.SignatureCode,
		License:      license.LicenseID,
		Date:         license.Date.Format("2006-01-02 15:04:05"),
		Type:         license.Type,
		Expiration:   license.ExpirationDate.Format("2006-01-02"),
		AllowedUsers: license.AllowedUsers,
		Project:      license.Project,
		Module:       license.Module,
	}
}

// createFile 创建一个以json数据的id字段作为文件名的文件
func createFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file:", err)
		return nil, err
	}
	return file, nil
}
