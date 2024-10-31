package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

var fileName string

func init() {
	flag.StringVar(&fileName, "file", "problems.csv", "Specify the file to be used for load questions from")
}

func main() {
	fmt.Println("Welcome to the quiz!")

	flag.Parse()

	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("CSV file %s not found\n", fileName)
		os.Exit(1)
	}

	problems := make([]problem, 0)

	csvReader := csv.NewReader(csvFile)
	for {
		line, err := csvReader.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Println("There are problems in your CSV file, please check")
			}

			break
		}

		problems = append(problems, problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		})
	}

	var score uint8
	for i, problem := range problems {
		fmt.Printf("Question %d: %s\n", i+1, problem.question)
		fmt.Print("Your answer: ")

		consoleReader := bufio.NewReader(os.Stdin)

		answer, err := consoleReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Unable to understand your answer")
			continue
		}

		answer = strings.TrimSpace(answer)
		if strings.EqualFold(answer, problem.answer) {
			score++
		}

		fmt.Println()
	}

	fmt.Printf("Your score is %d/%d!\n", score, len(problems))
}
