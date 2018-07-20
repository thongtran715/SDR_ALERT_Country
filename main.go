package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CDR struct {
	timestamp string
	state     string
	id        string
	operator  string
	statusIND string
}

func network() map[string]string {
	//dic := make(map[string]string)
	csvFile, error := os.Open("operatorDetailReport.csv")
	if error != nil {
		log.Fatal(error)
	}
	dic := make(map[string]string)
	scanner := bufio.NewScanner(csvFile)
	for scanner.Scan() {
		text := scanner.Text()
		arr := strings.Split(text, "\t")
		imsi := fmt.Sprintf("%s%s", arr[4], arr[5])
		_, ok := dic[imsi]
		if !ok {
			str := fmt.Sprintf("%s%s%s", arr[0], ";", arr[1])
			dic[imsi] = str
		}
	}
	defer csvFile.Close()
	return dic
}

func main() {
	cdrMTable := map[string]CDR{}
	cdrNTable := map[string]CDR{}

	const longForm = "20060102150405"
	file, err := os.Open("../CDR_20180702_2.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), ";")

		if arr[4] == "M" {
			operator := strings.Split(arr[9], ":")
			if operator[0] != "null" {
				aCDR := CDR{
					timestamp: arr[2],
					state:     arr[5],
					id:        arr[19],
					operator:  operator[2],
				}
				id := aCDR.id
				_, ok := cdrMTable[id]
				if !ok {
					cdrMTable[id] = aCDR
				}
			}
		} else {
			operator := strings.Split(arr[9], ":")
			if operator[0] != "null" {
				aCDR := CDR{
					timestamp: arr[3],
					state:     arr[5],
					id:        arr[20],
					operator:  operator[2],
					statusIND: arr[27],
				}
				id := aCDR.id
				cdr, ok := cdrNTable[id]
				if !ok {
					cdrNTable[id] = aCDR
				} else {
					first, _ := strconv.Atoi(cdr.timestamp)
					second, _ := strconv.Atoi(aCDR.timestamp)
					if second > first {
						delete(cdrNTable, id)
						cdrNTable[id] = aCDR
					}
				}
			}

		}
	}
	fmt.Println("Processing.......")
	countryOperatorMCCMNC := network()

	tree := CountryTree{
		CountryNodeRoot: nil,
	}
	for keyM := range cdrMTable {
		valueN, existed := cdrNTable[keyM]
		if existed {
			imsi := valueN.operator
			if imsi != "000000" {
				str := strings.Split(countryOperatorMCCMNC[imsi], ";")
				ops := Operators{
					operatorName: str[1],
				}
				tree.addCountry(str[0], ops)
			}
		}
	}

	fmt.Println("Processing Increment.......")
	// Value of M and N to verify
	for keyM, valueM := range cdrMTable {
		valueN, existed := cdrNTable[keyM]
		if existed {
			status := valueN.statusIND
			imsi := valueN.operator
			operator, country := "", ""
			if imsi != "000000" {
				str := strings.Split(countryOperatorMCCMNC[imsi], ";")
				operator = str[1]
				country = str[0]
			}
			timeM, _ := time.Parse(longForm, valueM.timestamp)
			timeN, _ := time.Parse(longForm, valueN.timestamp)
			delay := (timeN.Day()-timeM.Day())*84000 + (timeN.Hour()-timeM.Hour())*3600 + (timeN.Minute()-timeM.Minute())*60 + (timeN.Second() - timeM.Second())
			if status == "4" {
				tree.findAndIncrement("totalMessagesDLRS", operator, country)
			}
			if delay < 10 {
				tree.findAndIncrement("totalMessagesLessThan10seconds", operator, country)
			} else if delay < 60 {
				tree.findAndIncrement("totalMessagesLessThan1min", operator, country)
			} else if delay < 600 {
				tree.findAndIncrement("totalMessagesLessThan10mins", operator, country)
			} else if delay < 3600 {
				tree.findAndIncrement("totalMessagesLessThan1hour", operator, country)
			} else if delay < 7200 {
				tree.findAndIncrement("totalMessagesLessThan2hour", operator, country)
			}
			tree.findAndIncrement("totalMessagesReceived", operator, country)
		}
	}

	tree.totalTree()
	tree.displaySource()

	// Uncomment this if we want to get all Country
	//tree.calculateSDRPercentage()

	// Uncomment this if we want to get one pariticular country
	tree.sdrAlertOneCountry("Japan")
}
