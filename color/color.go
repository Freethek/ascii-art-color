package color

import "strings"

// variable map that hold all the ansicolor and their corresponding ansicode
var ansicodes = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"cyan":    "\033[34m",
	"magenta": "\033[35m",
	"white":   "\033[36m",
	"orange":  "\033[38;2;255;165;0m",
}

// functions to get ansi code from the ansicodes map
func GetAnsiCode(colorName string) string {
	ansicode, exist := ansicodes[colorName]

	if exist {
		return ansicode
	} else {
		return ""
	}
}

// finding the color indexes
func FinfColorIndex(substring, input string) map[int]bool {
	indexes := make(map[int]bool)

	//if the input is empty
	if input == "" {
		return indexes
	}

	if substring == "" {
		//mark all true if the substring is empty
		for i := range []rune(input) {
			indexes[i] = true
		}
	} else {
		for {
			start := 0
			//checking for the first index of the subtring in the input
			ind := strings.Index(input[start:], substring)
			if ind == -1 {
				break
			}
			//marking all subtring index in the input indec and stored for rendering
			for j := 0; j < len([]rune(substring)); j++ {
				indexes[start+j+ind] = true
			}

			start = ind + len([]rune(substring))
		}

	}

	return indexes
}
