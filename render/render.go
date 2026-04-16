package render

import (
	"ascii-art-color/color"
	"strings"
)

func Render(ansicode, input, substring string, bannerMap map[rune][]string) string {
	splitedForNewLine := strings.Split(input, "\\n")

	if len(splitedForNewLine) > 0 && splitedForNewLine[len(splitedForNewLine)-1] == "" {
		splitedForNewLine = splitedForNewLine[:len(splitedForNewLine)-1]
	}

	result := ""
	for _, segment := range splitedForNewLine {
		if segment == "" {
			result += "\n"
		} else {

			toBeColored := color.FindColorIndex(substring, segment)
			converted := []rune(segment)

			for row := 0; row <= 7; row++ {
				line := ""

				for i, ch := range converted {
					if ansicode != "" && toBeColored[i] && (i == 0 || toBeColored[i-1] == false) {
						line += ansicode
					}

					line += bannerMap[ch][row]

					if ansicode != "" && toBeColored[i] && (i == len(converted)-1 || !toBeColored[i+1]) {
						line += "\033[0m"
					}
				}
				result += line + "\n"
			}
		}

	}

	return result
}
