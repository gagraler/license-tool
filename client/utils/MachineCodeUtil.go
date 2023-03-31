package utils

import (
	"crypto/sha256"
	"fmt"
	"net"
)

/*
 * getMacAddr 用于获取本地网络接口的MAC地址
 * @return: success 返回MAC地址及nil
 *			failed  返回空字符串及错误信息
 * @params: null
 */
func getMacAddr() (string, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifi := range ifs {
		if ifi.Flags&net.FlagUp != 0 && ifi.Flags&net.FlagLoopback == 0 {
			address, err := ifi.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range address {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.To4() != nil {
						return ifi.HardwareAddr.String(), nil
					}
				case *net.IPAddr:
					if v.IP.To4() != nil {
						return ifi.HardwareAddr.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("failed to get mac address")
}

/*
 * MachineCode 生成机器码
 * @return: success 返回32位的机器码字符串
 *			failed  失败时会抛出panic异常
 * @params: null
 */
func MachineCode() string {
	macAddr, err := getMacAddr()
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256([]byte(macAddr))
	machineCode := fmt.Sprintf("%x", hash)[:32]
	return machineCode
}
