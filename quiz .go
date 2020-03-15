package main
import ("flag"
"os"
"log"
"encoding/csv"
"io"
"fmt"
"strconv"

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
	var total=0
	var right=0
	quizfile:=flag.String("quizFile","prob.csv","csv quiz filename ")
	flag.Parse()
	file, err := os.Open(*quizfile) // For read access.TODO: remove this opening
     if err != nil {
	log.Fatal(err)
	}
	var ans int
	for rec := range processCSV(file) {
		total++
		fmt.Println(rec[0])
		fmt.Scan(&ans)
		anss:=strconv.Itoa(ans)
		if anss==rec[1]{
right++
		}
		}
		fmt.Printf("Mark: %d/%d\n",right,total)
}
