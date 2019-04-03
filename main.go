package main

import (
	"flag"
	"fmt"
	"os"
	"encoding/csv"
	"strings"
	"time"
)

type problem struct {
	ques string
	ans string
}

func parseLines (lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem {
			ques : line[0],
			ans : strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func exit (msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func main () {

	csvFile := flag.String("file", "problems.csv", "A csv file in the format of 'question,answer'")
	timeout := flag.Int("timeout", 30, "Time limit to attempt a question")
	flag.Parse()

	file, err := os.Open(*csvFile)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file : %s\n", *csvFile))
	}

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()

	if err != nil {
		exit("Failed to read the CSV file.\n")
	}

	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeout) * time.Second)
	
	correct := 0 

	problemLoop:
		for i, problem := range problems {
			fmt.Printf("Question #%d: %s = \n", i+1, problem.ques)

			answerChan := make(chan string)

			go func() {

				var answer string

				fmt.Scanf("%s\n", &answer)

				answerChan <- answer
			}()


			select {
			case <-timer.C:
				fmt.Println()
				break problemLoop

			case answer := <- answerChan:
				if answer == problem.ans {
					correct++
				}

			}

		}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}