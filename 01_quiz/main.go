package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const defaultFile = "resources/problems.csv"
const sleepTime = 20 * time.Millisecond

type Problem struct {
	question string
	answer   string
}

type ProblemsLoader struct {
	reader *csv.Reader
}

func NewFileProblemsLoader(path string) *ProblemsLoader {
	// potential issue with opened file
	// (cannot defer close because panic)
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open %s", path)
	}

	return NewReaderProblemsLoader(f)
}

func NewReaderProblemsLoader(reader io.Reader) *ProblemsLoader {
	csvReader := csv.NewReader(reader)
	return &ProblemsLoader{csvReader}
}

func (p *ProblemsLoader) Load() []Problem {
	problems := []Problem{}

	for {
		row, err := p.reader.Read()

		if err == io.EOF {
			break
		}

		problems = append(problems, Problem{row[0], row[1]})
	}

	return problems
}

func main() {
	fmt.Printf("quiz game\n")

	csvFilePtr := flag.String("file", defaultFile, "path to file with problems")
	shufflePtr := flag.Bool("shuffle", false, "set this flag to shuffle the questions")
	timerPtr := flag.Int("timer", 30, "set max number of seconds for each question")

	flag.Parse()

	fmt.Printf("reading file: %s, shuffle: %v, timer: %ds\n", *csvFilePtr, *shufflePtr, *timerPtr)
	pl := NewFileProblemsLoader(*csvFilePtr)
	problems := pl.Load()
	score := 0
	scoreMax := len(problems)

	if *shufflePtr == true {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

	stdInReader := bufio.NewReader(os.Stdin)
	userInputChannel := make(chan string, 1)

	fmt.Printf("Starting the quiz.\n\nPlease type in your answer, and press <enter>. Press Q to quit.\n\n")
	for i, problem := range problems {
		fmt.Printf("Question %d: %s =  \n", i, problem.question)

		timeoutChannel := make(chan struct{}, 1)

		go func(waitSeconds int) {
			time.Sleep(time.Duration(waitSeconds) * time.Second)
			close(timeoutChannel)
		}(*timerPtr)

		go func() {
			text, err := stdInReader.ReadString('\n')

			if err != nil {
				log.Fatalf("error reading answer")
			}

			userInputChannel <- text
		}()

		var resultFlag string

		select {
		case text := <-userInputChannel:
			answer := strings.TrimSpace(strings.ToLower(text))

			if answer == "q" {
				fmt.Printf("Exiting")
				os.Exit(0)
			}

			resultFlag = "WRONG"
			if answer == problem.answer {
				resultFlag = "OK"
				score++
			}

		case <-timeoutChannel:
			fmt.Printf("Destroyed coroutine")
			resultFlag = "TIMEOUT"
		}

		fmt.Printf(resultFlag + "\n")
	}

	scorePercent := 100 * float64(score) / float64(scoreMax)
	fmt.Printf("\nRESULT: %d/%d (%.2f%%)", score, scoreMax, scorePercent)
}
