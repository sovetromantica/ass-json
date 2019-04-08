package main

import (
	"bufio"
	"log"
	"os"

	"./ass2json"

	"github.com/pborman/getopt/v2"
)

func main() {
	optASS := getopt.StringLong("ass", 'a', "", "Open ASS/SSA File")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}
	if len(*optASS) > 1 {
		file, err := os.Open(*optASS)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		ass2json.Ass2json(scanner)
	}

}
