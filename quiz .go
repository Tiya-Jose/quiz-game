package main
import ("flag"
"os"
"log"
"encoding/csv"
"io"
"fmt"
)

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
