package log

import (
	"fmt"
)

func Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR]: ")
	fmt.Printf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	fmt.Printf("[WARNING]: ")
	fmt.Printf(format, args...)
}

func Infoln(args ...interface{}) {
	fmt.Printf("[INFO]: ")
	fmt.Println(args...)
}
