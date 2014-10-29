package log

import (
	"flag"
	"fmt"
)

const (
	verboseError = iota
	verboseWarning
	verboseInfo
	verboseDebug
)

var verbose int

func init() {
	flag.IntVar(&verbose, "v", verboseInfo, "The log veboseness")
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

func Infof(format string, args ...interface{}) {
	if verbose < verboseInfo {
		return
	}
	fmt.Printf("[INFO]: ")
	fmt.Printf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	if verbose < verboseDebug {
		return
	}
	fmt.Printf("[DEBUG]: ")
	fmt.Printf(format, args...)
}
