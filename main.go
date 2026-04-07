package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	//declaration of args
	var substring string
	var input string
	var colorName string
	bannerName := "standard"

	//assinng values to the variables declared based on the argument
	if len(os.Args) == 1 {
		Error()
	} else if len(os.Args) == 2 {
		if checkForFlagPrefix(os.Args[1]) {
			Error()
		} else if !checkForFlagPrefix(os.Args[1]) {
			substring = ""
			input = os.Args[1]
		}

	} else if len(os.Args) == 3 {
		if checkForFlagPrefix(os.Args[1]) {
			substring = ""
			input = os.Args[2]
		} else {
			Error()
		}
	} else if len(os.Args) == 4 {
		if checkForFlagPrefix(os.Args[1]) {
			substring = os.Args[2]
			input = os.Args[3]
		} else {
			Error()
		}
	} else if len(os.Args) == 5 {
		if checkForFlagPrefix(os.Args[1]) {
			if isValidBanner(os.Args[4]) {
				substring = os.Args[2]
				input = os.Args[3]
				bannerName = os.Args[4]
			} else {
				Error()
			}
		} else {
			Error()
		}
	}

	
}

// checking for the right flag func
func checkForFlagPrefix(flag string) bool {
	if strings.HasPrefix("flag", "--color=") {
		return true
	} else {
		return false
	}
}

// checking if the banner is a valid banner file
func isValidBanner(bannerName string) bool {
	switch bannerName {
	case "thinkertoy", "shadow", "standard":
		return true
	default:
		return false
	}
}

// printing the error message
func Error() {
	//direct the user on the expected format
	fmt.Fprintf(os.Stderr, "Usage: go run . [OPTION] [STRING].\nEX: go run . --color=<color> <substring to be colored> \"something\"")
	//exiting the program
	os.Exit(1)
}
