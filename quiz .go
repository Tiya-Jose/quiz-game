package main
import ("flag"
"os"
"log"
"encoding/csv"
"io"
"fmt"
)
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
	return
}
func main(){
	quizfile:=flag.String("quizFile","problems.csv","csv quiz filename ")
	flag.Parse()
	file, err := os.Open(*quizfile) // For read access.
     if err != nil {
	log.Fatal(err)
    }
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
	}
}
