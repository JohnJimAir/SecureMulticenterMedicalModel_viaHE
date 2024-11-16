package nread

func V_Comma(filename string) ([]float64) {
	return Read(filename, ",", false, "")[0]
}

func V_Comma_EmitFirstline(filename string) ([]float64) {
	return Read(filename, ",", true, "")[0]
}

func V_Comma_Trim(filename string, cutset string) ([]float64) {
	return Read(filename, ",", false, cutset)[0]
}

func V_Blank(filename string) ([]float64) {
	return Read(filename, " ", false, "")[0]
}

func V_Blank_Trim(filename string, cutset string) ([]float64) {
	return Read(filename, " ", false, cutset)[0]
}