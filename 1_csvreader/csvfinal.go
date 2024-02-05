package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := ""
	var timeLimit time.Duration
	file, timeLimit := parseFlags(&csvFilename, &timeLimit)
	problems := readCSV(file)
	guess(problems, &timeLimit)
}

func parseFlags(csvFilename *string, timeLimit *time.Duration) (*os.File, time.Duration) {
	defaultFile := "problems.csv"
	defaultTime := time.Second * 30
	flag.StringVar(csvFilename, "csv", defaultFile, "Provide a comma separated CSV file.")
	flag.DurationVar(timeLimit, "limit", defaultTime, "Provide a time limit in seconds eg.: 30s (default)")
	flag.Parse()

	if *timeLimit <= time.Second*0 {
		fmt.Println("Timer can't be <= 0.")
		*timeLimit = defaultTime
	}

	fmt.Printf("Time limit is set to %v\n", timeLimit)
	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open the CSV file: %s. Using default.\n", *csvFilename)
		file, err = os.Open(defaultFile)
		if err != nil {
			fmt.Printf("Failed to open default CSV file: %s.\n", defaultFile)
			os.Exit(1)
		}
	}
	return file, *timeLimit
}

func guess(problems []problem, timeLimit *time.Duration) {
	for {
		var input string
		fmt.Print("press Enter to begin")
		_, err := fmt.Scanln(&input)
		if err != nil && strings.TrimSpace(input) != "\n" {
			break
		}
	}

	timer := time.NewTimer(*timeLimit)
	defer timer.Stop()

	correct := 0
	total := len(problems)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, total)
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("\nYou scored %d out of %d.\n", correct, total)
}

func readCSV(file *os.File) []problem {
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	return problems
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}
