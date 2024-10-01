package functions

import (
	"fmt"
)

func AsciiRep(fileLines []string) [][]string {
	var asciiRep [][]string
	var arr []string
	counter := 0
	// counter is incremented by the number of lines in filelines
	for _, line := range fileLines {
		if counter == 0 {
			counter++
			continue
		}
		arr = append(arr, line)
		counter++
		if counter == 9 {
			// append that character array to asciiRep
			asciiRep = append(asciiRep, arr)
			arr = nil
			counter = 0
		}
	}
	return asciiRep
}

func PrintStr(inputString string, asciiRep [][]string) [][]string {
	charHeight := 8

	output := make([][]string, charHeight)

	for i := range output {
		output[i] = make([]string, 0, len(inputString))
	}

	for _, char := range inputString {
		if char >= 32 && char <= 126 {
			index := int(char) - 32
			if index < 0 || index >= len(asciiRep) {
				fmt.Printf("Warning: ASCII representation for character '%c' not found.\n", char)
				continue
			}
			for i := 0; i < charHeight; i++ {
				output[i] = append(output[i], asciiRep[index][i])
			}
		} else {
			for i := 0; i < charHeight; i++ {
				output[i] = append(output[i], " ") // Handle non-printable characters
			}

		}
	}

	return output
}
