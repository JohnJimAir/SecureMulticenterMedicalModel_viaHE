package nread

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Read(filename string, delimiter string, emit bool, cutset string) ([][]float64) {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data [][]float64
	scanner := bufio.NewScanner(file)

	if emit {
		if scanner.Scan() {}
	}
	for scanner.Scan() {

		line := scanner.Text()
		line = strings.Trim(line, cutset)
		elements := strings.Split(line, delimiter)

		var floatElements []float64
		for _, element := range elements {

			element = strings.TrimSpace(element)
			if element != "" {
				floatValue, err := strconv.ParseFloat(element, 64)
				if err != nil {
					panic(err)
				}
				floatElements = append(floatElements, floatValue)
			}
		}
		
		if len(floatElements) > 0 {
			data = append(data, floatElements)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return data
}