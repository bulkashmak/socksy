package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("failed to open %s: %s", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	linesCh := getLinesChannel(file)
	for line := range linesCh {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		defer file.Close()
		defer close(linesCh)

		line := ""

		for {
			buffer := make([]byte, 8)
			n, err := file.Read(buffer)
			if err != nil {
				if line != "" {
					linesCh <- line
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("failed to read a file: %s\n", err.Error())
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				linesCh <- fmt.Sprintf("%s%s", line, parts[i])
				line = ""
			}
			line += parts[len(parts)-1]
		}
	}()

	return linesCh
}
