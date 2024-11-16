package nread

func M_Comma(filename string) ([][]float64) {
	return Read(filename, ",", false, "")
}

func M_Comma_EmitFirstline(filename string) ([][]float64) {
	return Read(filename, ",", true, "")
}

func M_Comma_Trim(filename string, cutset string) ([][]float64) {
	return Read(filename, ",", false, cutset)
}

func M_Blank(filename string) ([][]float64) {
	return Read(filename, " ", false, "")
}

func M_Blank_Trim(filename string, cutset string) ([][]float64) {
	return Read(filename, " ", false, cutset)
}