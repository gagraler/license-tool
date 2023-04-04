/*
 * package request http请求处理层，处理http请求
 */
package request

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"server/service"
	"server/utils"
	"time"
)

// ResponseData 响应数据结构体
type ResponseData struct {
	ID         string `json:"id"`            // 许可证唯一ID
	License    string `json:"license"`       // 许可证字符串
	Date       string `json:"date"`          // 许可证生成日期UTC
	Signature  string `json:"signatureCode"` // 被授权系统所在机器特征码
	Object     string `json:"object"`        // 许可对象
	Type       string `json:"type"`          // 许可类型
	Expiration string `json:"expiration"`    // 到期时间
	Project    string `json:"project"`       // 许可项目
	Module     string `json:"module"`        // 许可模块
}

// ResponseMsg 响应信息结构体
type ResponseMsg struct {
	Status string       `json:"status"`
	Data   ResponseData `json:"data"`
	Code   int          `json:"code"`
}

/*
 * GetLicenseRequest 处理获取许可证请求
 * @params: w http.ResponseWriter - 用于返回http响应的对象
 * 			r *http.Request - 包含http请求信息的指针
 * @returns: null
 */
func GetLicenseRequest(w http.ResponseWriter, r *http.Request) {
	// 从查询参数中解析出许可证所需参数
	params, err := parseLicenseParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 根据许可证所需参数生成许可证并封装响应信息
	license, err := generateLicenseAndResponse(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将生成的许可证写入到文件中
	err = writeLicenseToFile(license)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将响应信息返回给客户端
	utils.RespondWithJSON(w, http.StatusOK, license)
}

/*
 * parseLicenseParams 从http请求的URL中解析并校验获取许可证所需要的参数
 * @params: queryParams url.Values - 包含http请求中查询参数的对象
 * @returns: *service.LicenseParams - 许可证参数结构体指针
 * 			 error - 任何可能发生的错误
 */
func parseLicenseParams(queryParams url.Values) (*service.LicenseParams, error) {

	// 对获取的特征码进行格式检查
	signatureCode := queryParams.Get("signatureCode")
	matched, err := regexp.MatchString(`^.{1,32}$`, signatureCode)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("invalid signature code format")
	}

	// 将获取的过期时间字符串转换为时间对象
	expirationString := queryParams.Get("expiration")
	expiration, err := time.Parse("2006-01-02", expirationString)
	if err != nil {
		return nil, err
	}

	// 判断 type 参数是否为空
	typeParam := queryParams.Get("type")
	if typeParam == "" {
		return nil, fmt.Errorf("type parameter cannot be empty")
	}
	// 判断 type 参数是否为有效值
	// temporary - 临时
	// permanent - 永久
	if typeParam != "temporary" && typeParam != "permanent" {
		return nil, fmt.Errorf("invalid license type: %s", typeParam)
	}

	return &service.LicenseParams{
		SignatureCode: signatureCode,
		Object:        queryParams.Get("object"),
		Type:          typeParam,
		Expiration:    expiration,
		Project:       queryParams.Get("project"),
		Module:        queryParams.Get("module"),
	}, nil
}

/*
 * generateLicenseAndResponse 生成许可证并封装响应信息
 * @params: params *service.LicenseParams - 许可证参数结构体指针
 * @returns: *ResponseData - 响应信息结构体指针
 *       	 error - 任何可能发生的错误
 */
func generateLicenseAndResponse(params *service.LicenseParams) (*ResponseData, error) {

	// 生成许可证
	license, err := service.GenerateLicense(
		params.SignatureCode,
		params.Type,
		params.Expiration,
		params.Object,
		params.Project,
		params.Module,
	)
	if err != nil {
		return nil, err
	}

	// 构建返回的消息数据体
	responseData := &ResponseData{
		ID:         utils.GenerateUniqueID(),
		Signature:  license.SignatureCode,
		License:    license.LicenseID,
		Date:       license.Date.Format("2006-01-02 15:04:05 UTC"),
		Type:       license.Type,
		Expiration: license.Expiration.Format("2006-01-02"),
		Object:     license.Object,
		Project:    license.Project,
		Module:     license.Module,
	}

	return responseData, nil
}

/*
 * writeLicenseToFile 将生成的许可证写入到文件中
 * @params: license *ResponseData - 响应信息结构体指针
 * @returns: error - 任何可能发生的错误
 */
func writeLicenseToFile(license *ResponseData) error {

	// 将许可证数据转化为JSON格式
	response, err := json.Marshal(license)
	if err != nil {
		return err
	}

	// 根据许可证ID创建文件名
	fileName := license.ID + ".license"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("error creating file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("error closing file: ", err)
		}
	}(file)

	// 对许可证数据进行混淆处理
	encrypted, err := utils.ObfuscationUtil(response, license.Signature)
	if err != nil {
		return err
	}

	// 将混淆后的许可证数据写入文件
	_, err = file.WriteString(encrypted)
	if err != nil {
		log.Println("error writing to file:", err)
		return err
	}

	return nil
}
