package log

import (
	"fmt"
	"os"
	"strings"
)

// LogLevel は、ログレベルを表す列挙子です。
type LogLevel int

const (
	// LogLevel_None は、ログを記録しないログレベルです。
	LogLevel_None = iota
	// LogLevel_Warn は、警告のみを記録するログレベルです。
	LogLevel_Warn
	// LogLevel_Info は、参考情報までを記録するログレベルです。
	LogLevel_Info
	// LogLevel_Debug は、デバッグ情報までを記録するログレベルです。
	LogLevel_Debug
)

// Level は、現在のログレベルです。
var Level LogLevel = LogLevel_Info

// Warnf は、Warnレベルのログをフォーマットして記録します。
func Warnf(f string, args ...interface{}) {
	if LogLevel_Warn <= Level {
		fmt.Fprintf(os.Stderr, "[WARNING] "+f+"\n", args...)
	}
}

// Infof は、Infoレベルのログをフォーマットして記録します。
func Infof(f string, args ...interface{}) {
	if LogLevel_Info <= Level {
		fmt.Fprintf(os.Stderr, f+"\n", args...)
	}
}

var indent = 0

// Debugf は、Debugレベルのログをフォーマットして記録します。
func Debugf(f string, args ...interface{}) {
	if LogLevel_Debug <= Level {
		fmt.Fprintf(os.Stderr, strings.Repeat("  ", indent)+f+"\n", args...)
	}
}

// Enter は、この後に記録するログのインデントレベルを一段深くします。
func Enter() {
	indent++
}

// Leave は、この後に記録するログのインデントレベルを一段浅くします。
func Leave() {
	indent--
}
