package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

	limitPtr := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	seconds := *limitPtr
	// open file
	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer file.Close()

	// read csv values using csv.Reader
	r := csv.NewReader(file)

	correctAnswers := 0
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	timer := time.NewTimer(time.Second * time.Duration(seconds))
	
	for i, record := range records {
		fmt.Printf("Problem #%d: %s = ", i + 1, record[0])
		answerCh := make(chan string)
		go func() {
			scanner.Scan()
			err = scanner.Err()
			if err != nil {
				log.Fatal(err)
			}
			answerCh <- scanner.Text()
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correctAnswers, len(records))
			return
		case answer := <-answerCh:
			if answer == record[1] {
				correctAnswers++
			}
		}
	}
	
}

