package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

var p string
var s bool
var t int

func init() {
	flag.StringVar(&p, "path", path.Join("csv", "problems.csv"), "Define the path to the csv file")
	flag.BoolVar(&s, "shuffle", false, "Determine if the questions follow the csv file order")
	flag.IntVar(&t, "timer", 30, "Define the time in seconds until the end of the quiz. Defaults to 30.")
	flag.Parse()
}

func main() {
	var countOk int
	c := make(chan int)
	//opening the file
	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("Error opening the csv file found at path %s. Error : %v\n", p, err)
	}
	//creating the reader from file and parsing it to records
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error parsing the csv file. Error : %v\n", err)
	}
	recordsNumber := len(records)
	// if the shuffle flag is true, shuffling the records
	if s {
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}
	//starting
	fmt.Printf("Ready when you are. Press enter.")
	scanner := bufio.NewScanner(os.Stdin)
	_ = scanner.Scan()
	ti := time.NewTimer(time.Duration(t) * time.Second)
	for _, record := range records {
		go askQ(record, scanner, c)
		select {
		case toAdd := <-c:
			countOk += toAdd
		case _ = <-ti.C:
			fmt.Printf("Time expired!\nResults : %d correct on %d questions\n", countOk, recordsNumber)
			os.Exit(0)
		}
	}
	fmt.Printf("Finished!\nResults : %d correct on %d questions\n", countOk, recordsNumber)
}

func askQ(record []string, scanner *bufio.Scanner, c chan int) {
	fmt.Printf("Q : %s\n", record[0])
	_ = scanner.Scan()
	a := strings.TrimSpace(scanner.Text())
	if a == record[1] {
		c <- 1
	} else {
		c <- 0
	}
}
