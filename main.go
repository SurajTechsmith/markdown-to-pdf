package main

import (
	"flag"
	"fmt"
	"md_to_pdf/converter"
	"os"
	"time"
)

func main() {
	start := time.Now()
	input := flag.String("in", "", "Path to input file")
	output := flag.String("out", "", "Path to output pdf file (optional)")

	flag.Parse()

	if *input == "" {
		fmt.Println("Please provide input file using -in flag")
		os.Exit(1)
	}

	err := converter.Convert(*input, *output)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	end := time.Since(start)
	fmt.Println("PDF created successfully hurray : ", end)
}
