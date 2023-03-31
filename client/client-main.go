package main

import (
	"fmt"
	"license-tool/client/service"
)

func main() {

	fmt.Println(service.VerifyLicense("3477628109429210.license"))
}
