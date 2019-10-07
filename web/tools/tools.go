package tools

import (
	"io/ioutil"
	"runtime"
	"strings"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func WriteToFile(fileName string, lines []string) {
	data := []byte(strings.Join(lines, "\n"))
	ioutil.WriteFile(fileName, data, 777)
}

func ReadFromFile(fileName string) []string {
	data, error := ioutil.ReadFile(fileName)
	if error != nil {
		return nil
	}
	return strings.Split(string(data), "\n")
}
