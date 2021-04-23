package message

import (
	"fmt"

	"github.com/fatih/color"
)

var Cyan = color.New(color.FgCyan).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()

func Error(text string) {
	color.Red(text)
}

func Info(text string) {
	color.Cyan(text)
}

func Success(text string) {
	color.Green(text)
}

func Text(text string) {
	fmt.Println(text)
}

func Warn(text string) {
	color.Yellow(text)
}
