package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	payload := flag.String("p", "", "Payload to replace parameter values (required)")
	inputFile := flag.String("i", "", "Input file with URLs (optional, stdin is used if not provided)")
	outputFile := flag.String("o", "", "Output file to write results (optional, stdout is used if not provided)")

	flag.Usage = func() {
		fmt.Println("Usage of paramrep:")
		fmt.Println("  -i string")
		fmt.Println("      Input file with URLs (optional, stdin is used if not provided)")
		fmt.Println("  -o string")
		fmt.Println("      Output file to write results (optional, stdout is used if not provided)")
		fmt.Println("  -p string")
		fmt.Println("      Payload to replace parameter values (required)")
		fmt.Println("\nExamples:")
		fmt.Println("  echo \"https://example.com?param1=1&param2=2\" | paramrep -p PAYLOAD")
		fmt.Println("  Output:")
		fmt.Println("    https://example.com?param1=PAYLOAD&param2=2")
		fmt.Println("    https://example.com?param1=1&param2=PAYLOAD")
		fmt.Println("\n  paramrep -p PAYLOAD -i input.txt -o output.txt")
		fmt.Println("  Processes URLs from input.txt and writes results to output.txt.")
	}

	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	if *payload == "" {
		fmt.Println("Error: -p <payload> is required.")
		flag.Usage()
		os.Exit(1)
	}

	var inputScanner *bufio.Scanner
	if *inputFile != "" {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open input file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		inputScanner = bufio.NewScanner(file)
	} else {
		inputScanner = bufio.NewScanner(os.Stdin)
	}

	var outputWriter *os.File
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create output file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		outputWriter = file
	} else {
		outputWriter = os.Stdout
	}

	for inputScanner.Scan() {
		rawURL := inputScanner.Text()
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid URL: %s [%s]\n", rawURL, err)
			continue
		}

		queryParams := parsedURL.Query()

		for key := range queryParams {
			originalValue := queryParams.Get(key)

			queryParams.Set(key, *payload)
			parsedURL.RawQuery = queryParams.Encode()
			fmt.Fprintln(outputWriter, parsedURL.String())

			queryParams.Set(key, originalValue)
		}
	}

	if err := inputScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
	}
}
