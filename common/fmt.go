package common

import "fmt"

var FPrintln = fmt.Fprintln
var FPrintf = fmt.Fprintf
var Print = fmt.Print
var Printf = fmt.Printf
var Println = fmt.Println
var Sprintf = fmt.Sprintf

func PrintJSON(v interface{}) {
	json := MustMarshalIndentJSON(v)
	Println(string(json))
}

func PrintNewlines(n int) {
	Print(RepeatStr("\n", n))
}
