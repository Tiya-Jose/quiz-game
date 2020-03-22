package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var total = 10
var right = 0
var ch = make(chan []string, 10)

func processCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)

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
func checkTime(ticker *time.Ticker) {
	select {
	case <-ticker.C:
		fmt.Println("Your timer expired!!!")
		fmt.Printf("Score: %d/%d\n", right, total)
		close(ch)
		ticker.Stop()
		os.Exit(1)
		return
	}
}
func startTimer(start string, timer *int) *time.Ticker {
	if start == "y" {
		duration := int(time.Second) * *timer
		newDuration := time.Duration(duration)
		return time.NewTicker(newDuration)
	}
	return time.NewTicker(0)
}
func main() {
	var start string
	quizfile := flag.String("quizFile", "problems.csv", "csv quiz filename ")
	timer := flag.Int("timer", 5, "time to complete the quiz")

	flag.Parse()
	defer close(ch)

	file, err := os.Open(*quizfile) // For read access.TODO: remove this opening
	if err != nil {
		log.Fatal(err)
	}
	intro := fmt.Sprintf("Press y to start your timer for %d seconds", *timer)
	fmt.Println(intro)
	fmt.Scan(&start)
	ticker := startTimer(start, timer)
	defer ticker.Stop()
	go process(file)

	for {
		select {
		case rec := <-ch:
			go checkTime(ticker)
			var ans int
			fmt.Println(rec[0])
			fmt.Scan(&ans)
			anss := strconv.Itoa(ans)
			if anss == rec[1] {
				right++
			}

		}
	}

}
