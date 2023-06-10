package guilog

import (
	"log"

	"github.com/lxn/walk"
)

var output *walk.TextEdit

var logString string

func SetControl(outTE *walk.TextEdit) {
	output = outTE
}
func Reset() {
	logString = ""
	output.SetText("")
}
func PrintDivideLine() {
	Println("===================================================")
}
func Println(strings ...string) {
	line := ""
	for _, str := range strings {
		line += str
	}
	if line == "\r\n" || line == "\n" || line == "\r" || line == "\n\r" {
		log.Println("===================================================")
	} else {
		if line != "" {
			log.Println(line)
		} else {
			log.Println("===================================================")
		}
	}

	logString += line
	logString += "\r\n"
	output.SetText(logString)
}
