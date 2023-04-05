/*
 * package request http请求处理层，处理http请求
 */
package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"server/dao"
	"server/service"
	"server/utils"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

/*
 * GetLicenseRequest 处理获取许可证请求
 * @params: w http.ResponseWriter - 用于返回http响应的对象
 * 			r *http.Request - 包含http请求信息的指针
 * @returns: null
 */
func GetLicenseRequest(w http.ResponseWriter, r *http.Request) {

	// 从http请求参数中解析出许可证所需参数
	params, err := parseLicenseParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 根据许可证所需参数license字段值并封装响应信息
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

	//err = writeLicenseToMinio(license, license-file)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	// 将响应信息返回给客户端
	utils.RespondWithJSON(w, http.StatusOK, license)
}

/*
 * parseLicenseParams 从http请求的URL中解析并校验获取许可证所需要的参数
 * @params: queryParams url.Values - 包含http请求中查询参数的对象
 * @returns: *service.LicenseParams - 许可证参数结构体指针
 * 			 error - 任何可能发生的错误
 */
func parseLicenseParams(licenseData url.Values) (*dao.License, error) {

	// 对获取的特征码进行格式检查
	signature := licenseData.Get("signature")
	matched, err := regexp.MatchString(`^.{1,32}$`, signature)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("invalid signature code format")
	}

	// 将获取的过期时间字符串转换为时间对象
	expirationString := licenseData.Get("expiration")
	expiration, err := time.Parse("2006-01-02", expirationString)
	if err != nil {
		return nil, err
	}

	// 判断 type 参数是否为空
	typeParam := licenseData.Get("type")
	if typeParam == "" {
		return nil, fmt.Errorf("type parameter cannot be empty")
	}
	// 判断 type 参数是否为有效值
	// 0 - 临时
	// 1 - 永久
	// TODO 由前端来返回值 后端不做具体值
	if typeParam != "0" && typeParam != "1" {
		return nil, fmt.Errorf("invalid license type: %s", typeParam)
	}

	return &dao.License{
		Signature:  signature,
		Object:     licenseData.Get("object"),
		Type:       typeParam,
		Expiration: expiration,
		Project:    licenseData.Get("project"),
		Module:     licenseData.Get("module"),
	}, nil
}

/*
 * generateLicenseAndResponse 生成许可证并封装响应信息
 * @params: params *service.LicenseParams - 许可证参数结构体指针
 * @return: *ResponseData - 响应信息结构体指针
 *       	 error - 任何可能发生的错误
 */
func generateLicenseAndResponse(params *dao.License) (*dao.License, error) {

	// 生成许可证
	license, err := service.GenerateLicenseService(params)
	if err != nil {
		return nil, err
	}

	// 将许可证的日期和过期日期格式化为指定的格式
	dateStr := license.CreatedAt.Format("2006-01-02 15:04:05")
	expirationStr := license.Expiration.Format("2006-01-02")

	// 解析日期和过期日期
	// TODO UTC时间转为所在时区时间
	createDate, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return nil, err
	}

	expiration, err := time.Parse("2006-01-02", expirationStr)
	if err != nil {
		return nil, err
	}

	// 构建http返回的消息数据体
	responseData := &dao.License{
		ID:         license.ID,
		Signature:  license.Signature,
		License:    license.License,
		Type:       license.Type,
		Expiration: expiration,
		Object:     license.Object,
		Project:    license.Project,
		Module:     license.Module,
		CreatedAt:  createDate,
		UpdatedAt:  createDate,
	}

	return responseData, nil
}

/*
 * writeLicenseToFile 将生成的许可证写入到文件中
 * @params: license *ResponseData - 响应信息结构体指针
 * @returns: error - 任何可能发生的错误
 */
func writeLicenseToFile(licenseData *dao.License) error {

	// 将许可证数据转化为JSON格式
	// TODO 撤掉这层序列化
	licenseJSON, err := json.Marshal(licenseData)
	if err != nil {
		return err
	}

	// 根据许可证ID创建文件名
	fileName := strconv.Itoa(licenseData.ID) + ".license"
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
	encrypted, err := utils.ObfuscationUtil(licenseJSON, licenseData.Signature)
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

/*
 * writeLicenseToMinio 将生成的许可证写入到minio中
 * @params: license *dao.License - 许可证结构体指针
 * 		   bucketName string - 存储桶名称
 * @returns: error - 任何可能发生的错误
 */
func writeLicenseToMinio(license *dao.License, bucketName string) error {

	// 将许可证数据转化为JSON格式
	licenseJSON, err := json.Marshal(license)
	if err != nil {
		return err
	}

	// 获取minio的连接信息
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("1kXtZ2OvK65Frrw2")
	minioSecretKey := os.Getenv("6UyZ6gusbAO2PiPGb3ZN4dBKI2LcNm3G")

	// 初始化minio客户端
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	// 检查存储桶是否存在，如果不存在则创建
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			return err
		}
	}

	// 将许可证写入到minio中
	objectName := fmt.Sprintf("%s/%d.json", license.Project, license.ID)
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(licenseJSON), int64(len(licenseJSON)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
