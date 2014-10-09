package log

import (
	"fmt"
)

func Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Infoln(args ...interface{}) {
	fmt.Println(args...)
}
