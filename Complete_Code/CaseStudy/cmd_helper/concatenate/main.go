package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func main() {

	combinedMatrix, err := ConcatMatrices("/home/chenjingwei/CaseStudy/result/real/result_breast_Selu_New.txt", "/home/chenjingwei/CaseStudy/data/3/breast-test_data_x_y_array2.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(combinedMatrix)
}


func ConcatMatrices(fileAPath, fileDPath string) (string, error) {
	// 打开文件
	fileA, err := os.Open(fileAPath)
	if err != nil {
		return "", err
	}
	defer fileA.Close()

	fileD, err := os.Open(fileDPath)
	if err != nil {
		return "", err
	}
	defer fileD.Close()

	// 创建 scanner 读取文件内容
	scannerA := bufio.NewScanner(fileA)
	scannerD := bufio.NewScanner(fileD)

	var result strings.Builder

	// 按行读取文件并拼接
	for scannerA.Scan() && scannerD.Scan() {
		lineA := scannerA.Text()
		lineD := scannerD.Text()

		// 拼接两行，用空格分隔
		combinedLine := lineA + " " + lineD
		result.WriteString(combinedLine + "\n")
	}

	// 检查扫描错误
	if err := scannerA.Err(); err != nil {
		return "", err
	}
	if err := scannerD.Err(); err != nil {
		return "", err
	}

	return result.String(), nil
}