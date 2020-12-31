package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var p string

func init() {
	flag.StringVar(&p, "path", path.Join("csv", "problems.csv"), "Define the path to the csv file")
	flag.Parse()
}

func main() {
	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("No csv file found at path path %s\n", p)
	}
	r := csv.NewReader(f)
	var countLine, countOk int
L:
	for {
		rec, err := r.Read()
		switch err {
		case nil:
			countLine++
			fmt.Printf("%d. %s\n", countLine, rec[0])
			scanner := bufio.NewScanner(os.Stdin)
			_ = scanner.Scan()
			a := strings.TrimSpace(scanner.Text())
			if a == rec[1] {
				countOk++
			}
		case io.EOF:
			fmt.Printf("Finished!\nResults : %d correct on %d questions\n", countOk, countLine)
			break L

		default:
			countLine++
			log.Printf("Error reading record: %s.\n Trying to read the next one. \n", err)
		}
	}
}
