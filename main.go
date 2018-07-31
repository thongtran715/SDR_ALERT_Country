package main

import (
	"bufio"
	"flag"
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

	fileName := flag.String("fileName", "", "a string")

	flag.Parse()

	if *fileName == "" {
		log.Fatal("Missing flag for file name. Put the flag -fileName='name of your CDR file'")
	}
	//	"../CDR_20180702_2.txt"
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	// Closing file later
	defer file.Close()

	// Making buffer Read
	scanner := bufio.NewScanner(file)

	// Initiate tree
	tree := CountryTree{
		CountryNodeRoot: nil,
	}

	// Network storing country and operator
	countryOperatorMCCMNC := network()

	// Scanner for scan
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), ";")

		if arr[TYPE] == "M" {
			operator := strings.Split(arr[ROUTING_DEST], ":")
			if operator[0] != "null" {
				aCDR := CDR{
					timestamp: arr[ENTRY_TS],
					state:     arr[STATE],
					id:        arr[ID],
					operator:  operator[2],
				}
				id := aCDR.id
				_, ok := cdrMTable[id]
				if !ok {
					cdrMTable[id] = aCDR
				}
				imsi := aCDR.operator
				if imsi != "000000" {
					str := strings.Split(countryOperatorMCCMNC[imsi], ";")
					ops := Operators{
						operatorName: str[1],
					}
					tree.addCountry(str[0], ops)
					switch aCDR.state {
					case "Initial":
						tree.findAndIncrement("totalMessagesInitial", str[1], str[0])
					case "Rejected":
						tree.findAndIncrement("totalMessagesRejected", str[1], str[0])
					case "Delivered":
						tree.findAndIncrement("totalMessagesDelivered", str[1], str[0])
					case "Expired":
						tree.findAndIncrement("totalMessagesExpired", str[1], str[0])
					case "Undeliverable":
						tree.findAndIncrement("totalMessagesUndeliverable", str[1], str[0])
					case "Deleted":
						tree.findAndIncrement("totalMessagesDeleted", str[1], str[0])
					case "Delivered direct":
						tree.findAndIncrement("totalMessagesDeliveredDirect", str[1], str[0])
					}
					tree.findAndIncrement("totalMessagesReceived", str[1], str[0])
				}
			}
		} else {
			operator := strings.Split(arr[ROUTING_DEST], ":")
			if operator[0] != "null" {
				aCDR := CDR{
					timestamp: arr[STATE_TS],
					state:     arr[STATE],
					id:        arr[REF_ID],
					operator:  operator[2],
					statusIND: arr[ATTEMPTS],
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

	fmt.Println("Processing .......")
	// Value of M and N to verify
	for keyM, valueM := range cdrMTable {
		valueN, existed := cdrNTable[keyM]
		imsi := valueM.operator
		operator, country := "", ""
		if imsi != "000000" {
			str := strings.Split(countryOperatorMCCMNC[imsi], ";")
			operator = str[1]
			country = str[0]
		}
		if existed {
			status := valueN.statusIND

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
		} else {
			// If the M does not the N respond than it will be pending
			tree.findAndIncrement("totalMessagesPending", operator, country)
		}
	}

	tree.totalTree()
	//	tree.displaySource()

	// Uncomment this if we want to get all Country
	//tree.calculateSDRPercentage()

	// Uncomment this if we want to get one pariticular country
	//tree.sdrAlertOneCountry("Japan")

	// Uncomment this if you want to see the details of the country and its operator
	tree.findCountryAndOperatorDisplay("Japan", "au")
}
