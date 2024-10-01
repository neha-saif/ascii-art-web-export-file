package functions

import (
	"bufio"
	"fmt"
	"os"
)

func Read(fileName string) []string {
	file, err := os.Open(fileName + ".txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil 
	}
	defer file.Close() 
	Scanner := bufio.NewScanner(file)
	var fileLines []string
	for Scanner.Scan() {
		fileLines = append(fileLines, Scanner.Text())
	}
	if err := Scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return fileLines
}

