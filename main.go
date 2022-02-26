package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	problems := parseCsvToProblems(csvFilename)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correctCount := 0

problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)

			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nYou run out of time!")
			break problemloop
		case answer := <-answerCh:
			checkAnswer(&correctCount, answer, problem.answer)
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correctCount, len(problems))
}

func parseCsvToProblems(csvFilename *string) []Problem {
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	return problems
}

func parseLines(lines [][]string) []Problem {
	var problems []Problem

	for _, line := range lines {
		problem := Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
		problems = append(problems, problem)
	}

	return problems
}

func checkAnswer(score *int, answerFromInput string, correctAnswer string) {
	if answerFromInput == correctAnswer {
		*score++
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
