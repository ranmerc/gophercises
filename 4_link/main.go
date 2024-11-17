package main

import (
	"fmt"
	"os"

	"github.com/ranmerc/gophercises/link/parser"
)

func main() {
	fileNames := []string{
		"ex1.html",
		"ex2.html",
		"ex3.html",
		"ex3.html",
	}

	for _, name := range fileNames {
		parseAndPrintLinks(name)
	}
}

func parseAndPrintLinks(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("failed to read %s: %v", fileName, err)
		os.Exit(0)
	}

	links, err := parser.ParseLinks(file)
	if err != nil {
		fmt.Printf("failed to parse html in %s: %v", fileName, err)
		os.Exit(0)
	}

	fmt.Printf("\n\nLinks in %s:\n", fileName)
	fmt.Print(links)
}
