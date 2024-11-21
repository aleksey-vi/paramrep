package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	payload := flag.String("p", "", "Payload to replace parameter values (mutually exclusive with -pl)")
	payloadFile := flag.String("pl", "", "File with multiple payloads (mutually exclusive with -p)")
	inputFile := flag.String("i", "", "Input file with URLs (optional, stdin is used if not provided)")
	outputFile := flag.String("o", "", "Output file to write results (optional, stdout is used if not provided)")
	injectPath := flag.Bool("path", false, "Inject payloads into URL path segments")

	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	if *payload == "" && *payloadFile == "" {
		fmt.Println("Error: -p <payload> or -pl <payload file> is required.")
		flag.Usage()
		os.Exit(1)
	}

	if *payload != "" && *payloadFile != "" {
		fmt.Println("Error: -p and -pl cannot be used together.")
		flag.Usage()
		os.Exit(1)
	}

	// Load payloads
	payloads := []string{}
	if *payload != "" {
		payloads = append(payloads, *payload)
	} else {
		file, err := os.Open(*payloadFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open payload file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			payloads = append(payloads, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading payload file: %s\n", err)
			os.Exit(1)
		}
	}

	// Open input and output
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

	// Process each URL
	for inputScanner.Scan() {
		rawURL := inputScanner.Text()
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid URL: %s [%s]\n", rawURL, err)
			continue
		}

		// Inject payloads into query parameters
		queryParams := parsedURL.Query()
		for key := range queryParams {
			originalValue := queryParams.Get(key)
			for _, payload := range payloads {
				queryParams.Set(key, payload)
				parsedURL.RawQuery = queryParams.Encode()
				fmt.Fprintln(outputWriter, parsedURL.String())
			}
			queryParams.Set(key, originalValue)
		}

		// Inject payloads into path if -path flag is enabled
		if *injectPath {
			pathSegments := strings.Split(parsedURL.Path, "/")

			// Проверка: если количество сегментов больше одного, корректируем инъекцию
			if len(pathSegments) > 2 {
				// 1. Инъекция пейлоада в начало пути (замена каждого сегмента по очереди)
				for i := 1; i < len(pathSegments); i++ {
					if pathSegments[i] != "" {
						for _, payload := range payloads {
							newSegments := append([]string{}, pathSegments[:i]...)
							newSegments = append(newSegments, payload)
							newPath := strings.Join(newSegments, "/")
							fmt.Fprintln(outputWriter, parsedURL.Scheme+"://"+parsedURL.Host+newPath)
						}
					}
				}

				// 2. Добавляем payload в конец каждого сегмента без вставки между ними
				for i := 1; i < len(pathSegments); i++ {
					if pathSegments[i] != "" {
						for _, payload := range payloads {
							newSegments := append([]string{}, pathSegments...)
							newSegments[i] += payload
							newPath := strings.Join(newSegments, "/")
							fmt.Fprintln(outputWriter, parsedURL.Scheme+"://"+parsedURL.Host+newPath)
						}
					}
				}

				// 3. Добавляем payload в конец всего пути через слэш
				for _, payload := range payloads {
					newPath := strings.Join(pathSegments, "/") + "/" + payload
					// Убираем лишние слэши, чтобы не было конструкций типа "//"
					newPath = strings.ReplaceAll(newPath, "//", "/")
					fmt.Fprintln(outputWriter, parsedURL.Scheme+"://"+parsedURL.Host+newPath)
				}
			} else {
				// Если сегментов меньше или равно одному, просто выполняем стандартные операции
				for i := 1; i < len(pathSegments); i++ {
					if pathSegments[i] != "" {
						for _, payload := range payloads {
							newSegments := append([]string{}, pathSegments...)
							newSegments[i] = payload
							newPath := strings.Join(newSegments, "/")
							fmt.Fprintln(outputWriter, parsedURL.Scheme+"://"+parsedURL.Host+newPath)
						}
					}
				}

				for i := 1; i < len(pathSegments); i++ {
					if pathSegments[i] != "" {
						for _, payload := range payloads {
							newSegments := append([]string{}, pathSegments...)
							newSegments[i] += payload
							newPath := strings.Join(newSegments, "/")
							fmt.Fprintln(outputWriter, parsedURL.Scheme+"://"+parsedURL.Host+newPath)
						}
					}
				}

				for i := 1; i < len(pathSegments); i++ {
					if pathSegments[i] != "" {
						for _, payload := range payloads {
							newSegments := append([]string{}, pathSegments...)
							newSegments[i] = newSegments[i] + "/" + payload
							newPath := strings.Join(newSegments, "/")
							newPath = strings.ReplaceAll(newPath, "//", "/")
							fmt.Fprintln(outputWriter, parsedURL.Scheme+"://"+parsedURL.Host+newPath)
						}
					}
				}
			}
		}
	}

	if err := inputScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
	}
}
