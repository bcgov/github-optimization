package utils

import (
	"errors"
	"fmt"
	"os"
)

func HandleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func WriteLineToFile(f *os.File, cells ...string) {
	var line string
	for i, b := range cells {
		if i == 0 {
			line += b
		} else {
			line += "," + b
		}
	}

	f.WriteString(line + "\n")
}

func CheckEnv() (string, string) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		HandleError(errors.New("GITHUB_TOKEN is missing"))
	}

	org, ok := os.LookupEnv("GITHUB_ORG")
	if !ok {
		HandleError(errors.New("GITHUB_ORG is missing"))
	}

	return token, org
}
