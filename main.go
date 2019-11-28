package main

func main() {
	file := `tmp/resource.xlsx`
	data := readXls(file)
	GetResult(data)
}
