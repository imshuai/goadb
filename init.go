package goadb

import (
	"fmt"
)

var (
	devices []string
)

func init() {
	err := refreshDevices()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(devices) == 0 {
		fmt.Println("No devices found")
	}
}
