package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const licenseEndpoint = "http://example.com/license"

type License struct {
	Type         string    `json:"type"`
	Expiration   time.Time `json:"expiration"`
	AllowedUsers int       `json:"allowed_users"`
}

/*
 * @author: xuzongsheng
 * @date:2023/03/29
 * @return: error: 指向license对象的指针和在过程中遇到的任何错误
 * @parmas: none
 */
func ValidateLicense() (*License, error) {

	// 向licenseEndpoint发送GET请求
	resp, err := http.Get(licenseEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 将响应主体读入字节数组中
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	license := &License{}

	// 将许可证信息反序列化为license对象
	err = json.Unmarshal(body, license)
	if err != nil {
		return nil, err
	}

	// 检查许可证是否已过期
	if license.Expiration.Before(time.Now()) {
		return nil, fmt.Errorf("license expired")
	}

	// 返回经过验证的许可证及无误差
	return license, nil
}

func main() {
	license, err := ValidateLicense()
	if err != nil {
		fmt.Println("Invalid license:", err)
	} else {
		fmt.Printf("Valid %s license for %d users until %v\n", license.Type, license.AllowedUsers, license.Expiration)

		// Run the main application logic here...
	}
}
