package main

import (
	"fmt"
	"io/ioutil"
	"license-tool/client/utils"
	"os"
)

func main() {
	fmt.Println(utils.MachineCode())

	// 打开许可文件
	ciphertext, err := os.Open("3624763275116859.license")
	if err != nil {
		panic(err)
	}
	defer func(ciphertext *os.File) {
		err := ciphertext.Close()
		if err != nil {
		}
	}(ciphertext)

	// 读取许可文件内容
	ciphertextBytes, err := ioutil.ReadAll(ciphertext)
	if err != nil {
		panic(err)
	}

	// 解密文件内容
	plaintextBytes, err := utils.DeobfuscateUtil(string(ciphertextBytes), utils.MachineCode())
	if err != nil {
		panic(err)
	}

	// 进行一些修改，并将新内容写回文件
	err = ioutil.WriteFile("./3624763275116859.license", plaintextBytes, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}
