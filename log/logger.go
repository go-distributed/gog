package log

import (
	"flag"
	"fmt"
)

const (
	verboseError = iota
	verboseWarning
	verboseInfo
)

var verbose int

func init() {
	flag.IntVar(&verbose, "v", verboseWarning, "The log veboseness")
}

func Errorf(format string, args ...interface{}) {
	if verbose < verboseError {
		return
	}
	fmt.Printf("[ERROR]: ")
	fmt.Printf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	if verbose < verboseWarning {
		return
	}
	fmt.Printf("[WARNING]: ")
	fmt.Printf(format, args...)
}

func Infoln(args ...interface{}) {
	if verbose < verboseInfo {
		return
	}
	fmt.Printf("[INFO]: ")
	fmt.Println(args...)
}
