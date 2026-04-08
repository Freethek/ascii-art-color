package banner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Load(fileName string) map[rune][]string {
	fileData, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Readng File: %v", err)
	}

	//close the opened file
	defer fileData.Close()

	//declaring variables that wil hold the all lines of text in the file
	var lines []string

	//going into the file the was opened
	scanner := bufio.NewScanner(fileData)

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimRight(line, "\r")

		lines = append(lines, line)
	}

	bannerMap := make(map[rune][]string)
	for i := 32; i <= 126; i++ {
		start := (i-32)*9 + 1
		grab8lines := lines[start : start+8]
		bannerMap[rune(32)] = grab8lines
	}

	return bannerMap
}
