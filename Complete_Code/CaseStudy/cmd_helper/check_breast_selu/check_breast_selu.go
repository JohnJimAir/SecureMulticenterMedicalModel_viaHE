package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	result_cipher, err := ReadFloat64MatrixFromFile_square_brackets("./result_breast_Selu_K4_d15.txt")
	if err != nil {
		panic(err)
	}	
	label_cipher := Label(result_cipher)

	result_plain, err := ReadCSVToFloat64Matrix("./plain_breast_selu_K4_d15.txt")
	if err != nil {
		panic(err)
	}
	label_plain := Label(result_plain)

	ratio_diff := Ratio_diff(label_cipher, label_plain)
	fmt.Println("diff between cipher and plain is ",ratio_diff)

	// label_std, err := ReadFileAs2DSlice_int("./data/breast-test_data_x_y_array2.txt")
	// if len(label_cipher) != len(label_std) {
	// 	fmt.Println("len of std not equal !!!")
	// }
	
	// label_std_T := Transpose2DSlice(label_std)
	// if Ratio_diff(label_std_T[0], label_cipher) > Ratio_diff(label_std_T[1], label_cipher) {
	// 	fmt.Println("same between std and cipher is ", 1.0 - Ratio_diff(label_std_T[1], label_cipher))
	// } else {
	// 	fmt.Println("same between std and cipher is ", 1.0 - Ratio_diff(label_std_T[0], label_cipher))
	// }
}


func Ratio_diff(label_1, label_2 []int) (ratio float64) {
	
	if len(label_1) != len(label_2) {
		fmt.Println("len of label different!")
	}
	
	count := 0.0
	for i:=0;i<len(label_1);i++ {
		if label_1[i] != label_2[i] {
			count++
			fmt.Println("diff at index: ", i)
		}
	}

	return count/float64(len(label_1))
}

func Label(input [][]float64) (label []int) {
	
	if len(input[0]) != 2 || len(input[1]) != 2 {
		fmt.Println("len != 2 !!!!")
	}

	label = make([]int, len(input))
	for i:=0;i<len(input);i++ {
		if input[i][0] < input[i][1] {
			label[i] = 1
		} else {
			label[i] = 0
		}
	}
	return label
}

func ReadCSVToFloat64Matrix(filePath string) ([][]float64, error) {
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("无法打开文件: %v", err)
    }
    defer file.Close()

    // 创建 CSV 读取器
    reader := csv.NewReader(file)

    // 读取所有行
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("读取文件失败: %v", err)
    }

    // 创建存储 float64 的矩阵
    matrix := make([][]float64, len(records))

    // 遍历每一行
    for i, row := range records {
        matrix[i] = make([]float64, len(row))
        for j, val := range row {
            // 将字符串转换为 float64
            floatVal, err := strconv.ParseFloat(val, 64)
            if err != nil {
                return nil, fmt.Errorf("转换 float64 失败: %v (row: %d, col: %d)", err, i+1, j+1)
            }
            matrix[i][j] = floatVal
        }
    }

    return matrix, nil
}


func ReadFloat64MatrixFromFile_square_brackets(filePath string) ([][]float64, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 存储二维float64切片
	var data [][]float64

	// 使用bufio.Scanner按行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 获取当前行
		line := scanner.Text()

		// 去掉左右的括号
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "[")
		line = strings.TrimSuffix(line, "]")

		// 用空格分割每一行的两个float64数值
		nums := strings.Fields(line)
		if len(nums) != 2 {
			return nil, fmt.Errorf("invalid format in line: %s", line)
		}

		// 将字符串转换为float64
		a, err1 := strconv.ParseFloat(nums[0], 64)
		b, err2 := strconv.ParseFloat(nums[1], 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("failed to parse float64 values in line: %s, errors: %v, %v", line, err1, err2)
		}

		// 将结果添加到二维切片中
		data = append(data, []float64{a, b})
	}

	// 检查扫描过程是否出错
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return data, nil
}

func ReadFileAs2DSlice_int(filename string) ([][]int, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 创建一个Scanner读取文件
	scanner := bufio.NewScanner(file)

	// 定义二维int切片来存储数据
	var data [][]int

	// 逐行读取文件
	for scanner.Scan() {
		// 获取当前行
		line := scanner.Text()

		// 将行按空格分割
		parts := strings.Fields(line)

		// 创建一个切片用于存储当前行的整数
		row := make([]int, len(parts))

		// 将每个字符串转为整数并存储到当前行切片中
		for i, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("failed to convert string to int: %v", err)
			}
			row[i] = num
		}

		// 将当前行切片添加到二维切片中
		data = append(data, row)
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return data, nil
}

func Transpose2DSlice(slice [][]int) [][]int {
	// 如果切片为空，直接返回空切片
	if len(slice) == 0 {
		return [][]int{}
	}

	// 获取原始行数和列数
	rows := len(slice)
	cols := len(slice[0])

	// 创建一个转置后的切片，行数和列数互换
	transposed := make([][]int, cols)
	for i := range transposed {
		transposed[i] = make([]int, rows)
	}

	// 进行转置操作
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = slice[i][j]
		}
	}

	return transposed
}