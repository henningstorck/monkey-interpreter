package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"
const traceEnabled bool = false

func indentLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(msg string) {
	if traceEnabled {
		fmt.Printf("%s%s\n", indentLevel(), msg)
	}
}

func incIndent() { traceLevel = traceLevel + 1 }
func decIndent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	incIndent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIndent()
}
