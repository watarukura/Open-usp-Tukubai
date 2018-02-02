package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
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
	var reader *bufio.Reader
	if !terminal.IsTerminal(syscall.Stdin) {
		reader = bufio.NewReader(os.Stdin)
	} else {
		fileName := param[len(param)-1]
		// debug: fmt.Println(fileName)
		file, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
		if err != nil {
			fmt.Println("failed to read file")
			fmt.Println(err)
			os.Exit(1)
		}

		reader = bufio.NewReader(file)
		// defer file.Close()

		// パラメータの末尾を削除
		param = param[:len(param)-1]
	}
	// debug: fmt.Println(param)

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
				field = string(r[startNum-1 : startNum-1+lenNum-1])
				record = append(record, field)
			case strings.Contains(p, "/"):
				// fromTo := strings.Split(p, "/")
				// from, to := fromTo[0], fromTo[1]

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
