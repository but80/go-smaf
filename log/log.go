package log

import (
	"fmt"
	"os"
	"strings"
)

type LogLevel int

const (
	LogLevel_None = iota
	LogLevel_Warn
	LogLevel_Info
	LogLevel_Debug
)

var Level LogLevel = LogLevel_Info

func Warnf(f string, args ...interface{}) {
	if LogLevel_Warn <= Level {
		fmt.Fprintf(os.Stderr, "[WARNING] "+f+"\n", args...)
	}
}

func Infof(f string, args ...interface{}) {
	if LogLevel_Info <= Level {
		fmt.Fprintf(os.Stderr, f+"\n", args...)
	}
}

var indent = 0

func Debugf(f string, args ...interface{}) {
	if LogLevel_Debug <= Level {
		fmt.Fprintf(os.Stderr, strings.Repeat("  ", indent)+f+"\n", args...)
	}
}

func Enter() {
	indent++
}

func Leave() {
	indent--
}
