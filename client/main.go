package main

import (
	"fmt"
	"license-tool/client/service"
)

func main() {

	fmt.Println(service.VerifyLicense("3209497350222647.license"))
}
