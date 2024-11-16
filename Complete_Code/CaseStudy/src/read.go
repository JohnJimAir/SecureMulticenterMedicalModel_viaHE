package src

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// readFileToFloat64Matrix
func Read_matrix(filename string) ([][]float64, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var data [][]float64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		line = strings.Trim(line, " []")
		elements := strings.Split(line, ",")

		var floatElements []float64
		for _, element := range elements {

			element = strings.TrimSpace(element)
			element = strings.Trim(element, "[]")

			if element != "" {
				floatValue, err := strconv.ParseFloat(element, 64)
				if err != nil {
					return nil, fmt.Errorf("error converting to float: %v", err)
				}
				floatElements = append(floatElements, floatValue)
			}
		}
		if len(floatElements) > 0 {
			data = append(data, floatElements)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return data, nil
}

// readFloat64Vector
func Read_vector(filename string) ([]float64, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var data []float64
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {

		line := scanner.Text()
		line = strings.Trim(line, " []")
		elements := strings.Split(line, ",")

		for _, element := range elements {
			element = strings.TrimSpace(element)
			element = strings.Trim(element, "[]")

			if element != "" {
				floatValue, err := strconv.ParseFloat(element, 64)
				if err != nil {
					return nil, fmt.Errorf("error converting to float: %v", err)
				}
				data = append(data, floatValue)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return data, nil
}

func Readmatrix_nobracket(filePath string) ([][]float64, error) {
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    var data [][]float64

    // 创建一个Scanner读取文件
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // 读取每一行
        line := scanner.Text()
        // 按空格分割
        fields := strings.Fields(line)
        var row []float64
        for _, field := range fields {
            // 将字符串转换为float64
            num, err := strconv.ParseFloat(field, 64)
            if err != nil {
                return nil, fmt.Errorf("error parsing float: %v", err)
            }
            row = append(row, num)
        }
        data = append(data, row)
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading file: %v", err)
    }

    return data, nil
}