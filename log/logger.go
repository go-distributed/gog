package log

import (
	"flag"
	"fmt"
	"os"
)

const (
	verboseError = iota
	verboseWarning
	verboseInfo
	verboseDebug
)

var verbose int

func init() {
	flag.IntVar(&verbose, "v", verboseDebug, "The log veboseness")
}

func Errorf(format string, args ...interface{}) {
	if verbose < verboseError {
		return
	}
	fmt.Fprintf(os.Stderr, "[ERROR]:    ")
	fmt.Fprintf(os.Stderr, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	if verbose < verboseError {
		return
	}
	fmt.Fprintf(os.Stderr, "[FATAL]:    ")
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func Warningf(format string, args ...interface{}) {
	if verbose < verboseWarning {
		return
	}
	fmt.Fprintf(os.Stderr, "[WARNING]: ")
	fmt.Fprintf(os.Stderr, format, args...)
}

func Infof(format string, args ...interface{}) {
	if verbose < verboseInfo {
		return
	}
	fmt.Fprintf(os.Stderr, "[INFO]:     ")
	fmt.Fprintf(os.Stderr, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if verbose < verboseDebug {
		return
	}
	fmt.Fprintf(os.Stderr, "[DEBUG]:    ")
	fmt.Fprintf(os.Stderr, format, args...)
}
