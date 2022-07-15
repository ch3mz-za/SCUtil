package common

import (
	"fmt"
	"runtime"
	"strings"
)

func CleanInput(input string) string {
	fmt.Println("INPUT: ", input)
	os := runtime.GOOS
	switch os {
	case "windows":
		return strings.Replace(input, "\r\n", "", -1)
	case "darwin":
		return strings.Replace(input, "\n", "", -1)
	default:
		return input
	}
}
