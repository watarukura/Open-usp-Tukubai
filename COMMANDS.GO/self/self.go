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

	for _, p := range param {
		switch {
		case strings.Contains(p, "."):
			sp := strings.Split(p, ".")
			if len(sp) != 3 {
				if len(sp) != 2 {
					fmt.Println("invalid param: " + p)
					os.Exit(1)
				}
			}
			for _, spp := range sp {
				if spp != "NF" {
					_, err := strconv.Atoi(spp)
					if err != nil {
						fmt.Println("invalid param: " + p)
						os.Exit(1)
					}
				}
			}
		case strings.Contains(p, "/"):
			sp := strings.Split(p, "/")
			from, to := sp[0], sp[1]
			if len(sp) != 2 {
				fmt.Println("invalid param" + p)
				os.Exit(1)
			}
			if strings.HasPrefix(from, "NF") {
				if len(from) > 2 {
					sign := from[2:3]
					if sign != "-" {
						fmt.Println("invalid param: " + p)
						os.Exit(1)
					}
					_, err := strconv.Atoi(from[3:])
					if err != nil {
						fmt.Println("invalid param: " + p)
						os.Exit(1)
					}
				}
			} else {
				_, err := strconv.Atoi(from)
				if err != nil {
					fmt.Println("invalid param: " + p)
					os.Exit(1)
				}

			}

			if strings.HasPrefix(to, "NF") {
				if len(to) > 2 {
					sign := to[2:3]
					if sign != "-" {
						fmt.Println("invalid param: " + p)
						os.Exit(1)
					}
					_, err := strconv.Atoi(to[3:])
					if err != nil {
						fmt.Println("invalid param: " + p)
						os.Exit(1)
					}
				}
			} else {
				_, err := strconv.Atoi(to)
				if err != nil {
					fmt.Println("invalid param: " + p)
					os.Exit(1)
				}
			}
		case strings.HasPrefix(p, "NF"):
			if len(p) > 2 {
				sign := p[2:3]
				if sign != "-" {
					fmt.Println("invalid param: " + p)
					os.Exit(1)
				}
				_, err := strconv.Atoi(p[3:])
				if err != nil {
					fmt.Println("invalid param: " + p)
					os.Exit(1)
				}
			}
		default:
			_, err := strconv.Atoi(p)
			if err != nil {
				fmt.Println("invalid param: " + p)
				os.Exit(1)
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)

	return param, reader
}

func selectField(param []string, file *bufio.Reader) [][]string {
	csvr := csv.NewReader(file)
	delm, _ := utf8.DecodeLastRuneInString(" ")
	csvr.Comma = delm
	csvr.TrimLeadingSpace = true

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
			case p == "0":
				fields := make([]string, len(line))
				copy(fields, line)
				record = append(record, fields...)
			case strings.Contains(p, "."):
				nfStartLength := strings.Split(p, ".")
				var nf string
				var start string
				var length string
				var num int
				var startNum int
				var lenNum int
				var str string
				if len(nfStartLength) == 2 {
					nf, length = nfStartLength[0], nfStartLength[1]
					if nf == "NF" {
						num = len(line)
					} else {
						num, _ = strconv.Atoi(nf)
					}
					lenNum, _ = strconv.Atoi(length)
					str := line[num-1]
					startNum = utf8.RuneCountInString(str) - lenNum
					r := []rune(str)
					field = string(r[startNum:])
					record = append(record, field)
				} else {
					nf, start, length = nfStartLength[0], nfStartLength[1], nfStartLength[2]
					if nf == "NF" {
						num = len(line)
					} else {
						num, _ = strconv.Atoi(nf)
					}
					startNum, _ = strconv.Atoi(start)
					lenNum, _ = strconv.Atoi(length)
					str = line[num-1]
					r := []rune(str)
					field = string(r[startNum-1 : startNum-1+lenNum])
					record = append(record, field)
				}
			case strings.Contains(p, "/"):
				fromTo := strings.Split(p, "/")
				from, to := fromTo[0], fromTo[1]
				var fromNum int
				var toNum int
				if strings.HasPrefix(from, "NF") {
					if len(from) > 2 {
						nfMinus, _ := strconv.Atoi(from[3:])
						fromNum = len(line) - nfMinus
					} else {
						fromNum = len(line)
					}
				} else {
					fromNum, _ = strconv.Atoi(from)
				}

				if strings.HasPrefix(to, "NF") {
					if len(to) > 2 {
						nfMinus, _ := strconv.Atoi(to[3:])
						toNum = len(line) - nfMinus
					} else {
						toNum = len(line)
					}
				} else {
					toNum, _ = strconv.Atoi(to)
				}
				fields := make([]string, len(line[fromNum-1:toNum]))
				copy(fields, line[fromNum-1:toNum])
				record = append(record, fields...)
			case strings.HasPrefix(p, "NF"):
				var num int
				if len(p) > 2 {
					nfMinus, _ := strconv.Atoi(p[3:])
					num = len(line) - 1 - nfMinus
					field = line[num]
				} else {
					field = line[len(line)-1]
				}
				record = append(record, field)
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
