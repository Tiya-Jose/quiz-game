package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	total    = 10
	rightAns = 0
	ch       = make(chan []string)
)

func startTimer(start string, myTime *int) *time.Timer {
	if start == "y" {
		return time.NewTimer(time.Duration(*myTime) * time.Second)
	}
	return time.NewTimer(0)
}

func processCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string)
	go func() {
		r := csv.NewReader(rc)
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(err)
				return

			}
			ch <- rec
		}
	}()
	return ch
}

func process(file *os.File) {
	for rec := range processCSV(file) {
		ch <- rec
	}
}

func checkTime(timer *time.Timer) {
	select {
	case <-timer.C:
		fmt.Println("\nYour timer expired!!!")
		fmt.Printf("Score: %d/%d\n", rightAns, total)
		close(ch)
		os.Exit(1)
	}
}

func main() {

	quizfile := flag.String("quizFile", "problems.csv", "csv quizFile file in the format of 'question,answer' ")
	myTime := flag.Int("timer", 30, "time to complete the quiz")
	flag.Parse()

	defer close(ch)

	file, err := os.Open(*quizfile)
	if err != nil {
		log.Println(err)
		return
	}

	intro := fmt.Sprintf("Press y to start your timer for %d seconds", *myTime)
	fmt.Println(intro)
	var start string
	fmt.Scan(&start)
	timer := startTimer(start, myTime)

	go process(file)
	problemCount := 0
	for {

		select {
		case rec := <-ch:
			go checkTime(timer)
			var ans int
			problemCount++
			fmt.Printf("Problem #%d: %s\n", problemCount, rec[0])
			fmt.Scan(&ans)
			anss := strconv.Itoa(ans)
			if anss == strings.TrimSpace(rec[1]) {
				rightAns++
			}

		}
	}

}
