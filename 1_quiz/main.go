package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

const DEFAULT_TIME_LIMIT = 30 * time.Second

var fileName string
var timeLimit time.Duration

func init() {
	flag.StringVar(&fileName, "file", "problems.csv", "Specify the file to be used for load questions from")
	flag.DurationVar(&timeLimit, "limit", DEFAULT_TIME_LIMIT, "Specify the time limit for the quiz (in seconds)")
}

func main() {
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

	fmt.Println("Welcome to the quiz!")
	fmt.Println("You have " + timeLimit.String() + " to answer all the questions.")
	fmt.Println("Press enter/return key to begin")
	consoleReader := bufio.NewReader(os.Stdin)
	consoleReader.ReadLine()

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	done := make(chan struct{})
	var score uint8
	var timeoutOccurred bool

	go func() {
		go func() {
			select {
			case <-ctx.Done():
				timeoutOccurred = true
				os.Stdin.Close()
				close(done)
			case <-done:
				return
			}
		}()

		defer close(done)
		for i, problem := range problems {
			fmt.Printf("Question %d: %s\n", i+1, problem.question)
			fmt.Print("Your answer: ")

			consoleReader := bufio.NewReader(os.Stdin)

			answer, err := consoleReader.ReadString('\n')
			if err != nil {
				if timeoutOccurred {
					return
				}

				fmt.Printf("Unable to understand your answer")
				continue
			}

			answer = strings.TrimSpace(answer)
			if strings.EqualFold(answer, problem.answer) {
				score++
			}

			fmt.Println()
		}
	}()

	<-done

	if timeoutOccurred {
		fmt.Print("\nTime's up! ")
	}
	fmt.Printf("Your score is %d/%d!\n", score, len(problems))
}
