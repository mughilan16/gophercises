package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	fileName, timeLimit := getArguments()
	problems := getProblems(fileName)
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	quiz(problems, timer)
}

func quiz(problems []problem, timer *time.Timer) {
	correct := 0
problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := getAnswerChannel()
		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d of %d\n", correct, len(problems))
}

func getAnswerChannel() chan string {
	answerCh := make(chan string)
	go func() {
		var answer string
		fmt.Scanf("%s\n", &answer)
		answerCh <- answer
	}()
	return answerCh
}

func getProblems(fileName string) []problem {
	data := readData(openCSV(fileName))
	return parseQuestion(data)
}

func getArguments() (string, int) {
	fileName := flag.String("csv", "problems.csv", "a csv file in the format of question,answer")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz")

	flag.Parse()
	return *fileName, *timeLimit
}

func openCSV(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		exit(fmt.Sprintf("Can't Open File : %s", fileName))
	}
	return file
}

func readData(file *os.File) [][]string {
	r := csv.NewReader(file)
	data, err := r.ReadAll()
	if err != nil {
		exit("Can't open CSV file")
	}
	return data
}

func parseQuestion(data [][]string) []problem {
	problems := make([]problem, len(data))
	for i, p := range data {
		problems[i] = problem{
			q: p[0],
			a: strings.TrimSpace(p[1]),
		}
	}
	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
