package helpers

import (
	"os"
	"strconv"
	"time"
)

type PrintHelper struct {
	withTimer   bool
	currentTime time.Time
}

func CreatePrintHelper(withTimer bool) *PrintHelper {
	var ph PrintHelper
	ph.withTimer = withTimer
	ph.DoDefaultStuff()
	return &ph
}
func (pHelper *PrintHelper) DoDefaultStuff() {
	if pHelper.withTimer {
		pHelper.currentTime = time.Now()
	}
}
func (pHelper *PrintHelper) P(message string) {
	var timeString string = ""
	if pHelper.withTimer {
		var curTime = time.Now()
		timeString = strconv.Itoa(int(curTime.Sub(pHelper.currentTime).Seconds()))
		pHelper.currentTime = curTime
	}
	print("\n" + timeString + "\t" + message)
}
func PadString(message string, size int) string {
	for len(message) < size {
		message += " "
	}
	return message
}
func PrintRow(values []string, length int) {
	var str string
	for i := 0; i < len(values); i++ {
		str += PadString(values[i], length)
	}
	println(str)
}

/*
Write message to file
*/
func (pHelper *PrintHelper) WriteToFile(message string, path string) {
	file, _ := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(message + "\n")
	file.Close()
}
func (pHealer *PrintHelper) Line(str string) {

	println(str)
}
