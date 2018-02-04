package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	flag.Parse()
	param := flag.Args()
	// debug: fmt.Println(param)

	validParam, file := validateParam(param)

	fields := selectField(validParam, file)
	// debug: fmt.Println(fields)

	writeFields(fields)
}

func validateParam(param []string) ([]string, *bufio.Reader) {
	if len(param) == 0 {
		fmt.Println("failed to read param")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	return param, reader
}

func selectField(param []string, file *bufio.Reader) [][]string {
	csvr := csv.NewReader(file)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm

	orgRecord, err := csvr.ReadAll()
	if err != nil {
		fmt.Println("failed to read file as csv")
		fmt.Println(err)
		os.Exit(1)
	}

	var result [][]string
	var field string
	var record []string
	var toNum int
	for _, line := range orgRecord {
		for _, p := range param {
			switch {
			case p == "NF":
				field = line[len(line)-1]
				record = append(record, field)
			case strings.Contains(p, "."):
				nfStartLength := strings.Split(p, ".")
				nf, start, length := nfStartLength[0], nfStartLength[1], nfStartLength[2]
				num, _ := strconv.Atoi(nf)
				startNum, _ := strconv.Atoi(start)
				lenNum, _ := strconv.Atoi(length)
				str := line[num-1]
				r := []rune(str)
				field = string(r[startNum-1 : startNum-1+lenNum])
				record = append(record, field)
			case strings.Contains(p, "/"):
				fromTo := strings.Split(p, "/")
				from, to := fromTo[0], fromTo[1]
				fromNum, _ := strconv.Atoi(from)
				if to == "NF" {
					toNum = len(line)
				} else {
					toNum, _ = strconv.Atoi(to)
				}
				fields := make([]string, len(line[fromNum-1:toNum]))
				copy(fields, line[fromNum-1:toNum])
				record = append(record, fields...)
			default:
				num, _ := strconv.Atoi(p)
				field = line[num-1]
				record = append(record, field)
			}
		}
		// debug: fmt.Println(record)
		result = append(result, record)
		record = []string{}
	}

	// debug: fmt.Println(result)
	return result
}

func writeFields(fields [][]string) {
	csvw := csv.NewWriter(os.Stdout)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvw.Comma = delm

	for _, line := range fields {
		csvw.Write(line)
	}
	csvw.Flush()
}
