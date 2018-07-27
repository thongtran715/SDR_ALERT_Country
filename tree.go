package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Building two trees -> outer tree will consist of Country and inner tree will consist of Operator

// Operators -> Store each of the OPerator
type Operators struct {
	operatorName                   string
	totalMessagesDelivered         int
	totalMessagesRejected          int
	totalMessagesInitial           int
	totalMessagesDeliveredDirect   int
	totalMessagesExpired           int
	totalMessagesDeleted           int
	totalMessagesUndeliverable     int
	totalMessagesReceived          int
	totalMessagesLessThan10mins    int
	totalMessagesLessThan1min      int
	totalMessagesLessThan1hour     int
	totalMessagesLessThan2hour     int
	totalMessagesDLRS              int
	totalMessagesLessThan10seconds int
}

// OperatorNode -> Consist of left and right node
type OperatorNode struct {
	leftOperatorNode  *OperatorNode
	rightOperatorNode *OperatorNode
	operas            Operators
}

// CountryNode -> Store each of the CountryNode
type CountryNode struct {
	countryName                    string
	operatorRoot                   *OperatorNode
	leftCountryNode                *CountryNode
	rightCountryNode               *CountryNode
	totalMessagesReceived          int
	totalMessagesLessThan10mins    int
	totalMessagesLessThan1min      int
	totalMessagesLessThan1hour     int
	totalMessagesLessThan2hour     int
	totalMessagesDLRS              int
	totalMessagesLessThan10seconds int
}

// CountryTree -> CountryNode root
type CountryTree struct {
	CountryNodeRoot *CountryNode
}

// ****************************************************************************************************************
// 										FUNCTIONS

// ****************************************************************************************************************
// 											ADD FUNCTION

func (node *OperatorNode) addOperator(ops Operators) {
	if node == nil {
		return
	}
	compare := strings.Compare(node.operas.operatorName, ops.operatorName)
	if compare == 0 {
		return
	} else if compare == -1 {
		if node.leftOperatorNode == nil {
			node.leftOperatorNode = &OperatorNode{
				operas: ops,
			}
			return
		}
		node.leftOperatorNode.addOperator(ops)
	} else if compare == 1 {
		if node.rightOperatorNode == nil {
			node.rightOperatorNode = &OperatorNode{
				operas: ops,
			}
			return
		}
		node.rightOperatorNode.addOperator(ops)
	}
}

func (node *CountryNode) addCountry(country string, ops Operators) {
	if node == nil {
		return
	}
	compare := strings.Compare(node.countryName, country)
	if compare == 0 {
		// Add the new Operator
		node.operatorRoot.addOperator(ops)
		return
	} else if compare == -1 {
		if node.leftCountryNode == nil {
			node.leftCountryNode = &CountryNode{
				countryName: country,
				operatorRoot: &OperatorNode{
					operas: ops,
				},
			}
			return
		}
		node.leftCountryNode.addCountry(country, ops)
	} else if compare == 1 {
		if node.rightCountryNode == nil {
			node.rightCountryNode = &CountryNode{
				countryName: country,
				operatorRoot: &OperatorNode{
					operas: ops,
				},
			}
			return
		}
		node.rightCountryNode.addCountry(country, ops)
	}
}

func (tree *CountryTree) addCountry(country string, ops Operators) {
	if tree.CountryNodeRoot == nil {
		tree.CountryNodeRoot = &CountryNode{
			countryName: country,
			operatorRoot: &OperatorNode{
				operas: ops,
			},
		}
		return
	}
	tree.CountryNodeRoot.addCountry(country, ops)
}

// ****************************************************************************************************************
// 												DISPLAY

func (node *OperatorNode) displaySingleOperator() {
	fmt.Println("Name of Operator: ", node.operas.operatorName)
	fmt.Println("Total Message Received: ", node.operas.totalMessagesReceived)
	fmt.Println("Total Message Less Than 10 seconds: ", node.operas.totalMessagesLessThan10seconds)
	fmt.Println("Total Message Less Than 1 min: ", node.operas.totalMessagesLessThan1min)
	fmt.Println("Total Message Less Than 10 mins: ", node.operas.totalMessagesLessThan10mins)
	fmt.Println("Total Message Less Than 1 hour: ", node.operas.totalMessagesLessThan1hour)
	fmt.Println("Total Message Less Than 2 hours: ", node.operas.totalMessagesLessThan2hour)
	fmt.Println("Total Message Received status 4: ", node.operas.totalMessagesDLRS)
	fmt.Println("Total Messages that are Initial: ", node.operas.totalMessagesInitial)
	fmt.Println("Total Messages that are Rejected: ", node.operas.totalMessagesRejected)
	fmt.Println("Total Messages that are Delivered: ", node.operas.totalMessagesDelivered)
	fmt.Println("Total Messages that are Expired: ", node.operas.totalMessagesExpired)
	fmt.Println("Total Messages that are Undeliverable: ", node.operas.totalMessagesUndeliverable)
	fmt.Println("Total Messages that are Deleted: ", node.operas.totalMessagesDeleted)
	fmt.Println("Total Messages that are Delivered Direct: ", node.operas.totalMessagesDeliveredDirect)

}

func (node *OperatorNode) display() {
	if node == nil {
		return
	}
	fmt.Println("Name of Operator: ", node.operas.operatorName)
	fmt.Println("Total Message Received: ", node.operas.totalMessagesReceived)
	fmt.Println("Total Message Less Than 10 seconds: ", node.operas.totalMessagesLessThan10seconds)
	fmt.Println("Total Message Less Than 1 min: ", node.operas.totalMessagesLessThan1min)
	fmt.Println("Total Message Less Than 10 mins: ", node.operas.totalMessagesLessThan10mins)
	fmt.Println("Total Message Less Than 1 hour: ", node.operas.totalMessagesLessThan1hour)
	fmt.Println("Total Message Less Than 2 hours: ", node.operas.totalMessagesLessThan2hour)
	fmt.Println("Total Message Received status 4: ", node.operas.totalMessagesDLRS)
	fmt.Println("Total Messages that are Initial: ", node.operas.totalMessagesInitial)
	fmt.Println("Total Messages that are Rejected: ", node.operas.totalMessagesRejected)
	fmt.Println("Total Messages that are Delivered: ", node.operas.totalMessagesDelivered)
	fmt.Println("Total Messages that are Expired: ", node.operas.totalMessagesExpired)
	fmt.Println("Total Messages that are Undeliverable: ", node.operas.totalMessagesUndeliverable)
	fmt.Println("Total Messages that are Deleted: ", node.operas.totalMessagesDeleted)
	fmt.Println("Total Messages that are Delivered Direct: ", node.operas.totalMessagesDeliveredDirect)
	node.leftOperatorNode.display()
	node.rightOperatorNode.display()
}

func (node *CountryNode) display() {
	if node == nil {
		return
	}
	fmt.Println("Country Name: ", node.countryName)
	node.operatorRoot.display()
	node.leftCountryNode.display()
	node.rightCountryNode.display()
}

func (tree *CountryTree) display() {
	if tree.CountryNodeRoot == nil {
		fmt.Println("Emtpy Tree")
		return
	}
	tree.CountryNodeRoot.display()
}

func (node *CountryNode) displaySouce() {
	if node == nil {
		return
	}
	fmt.Println("Country Name: ", node.countryName)
	fmt.Println("Total Messages Less Than 1 minute: ", node.totalMessagesLessThan1min)
	fmt.Println("Total Messages Less Than 10 minutes: ", node.totalMessagesLessThan10mins)
	fmt.Println("Total Messages Less Than 1 hour: ", node.totalMessagesLessThan1hour)
	fmt.Println("Total Messages Less Than 2 hour: ", node.totalMessagesLessThan2hour)
	fmt.Println("Total Messages Less Than 10 seconds: ", node.totalMessagesLessThan10seconds)
	fmt.Println("Total Messages have status 4(DL): ", node.totalMessagesDLRS)
	fmt.Println("Total Messages Receved", node.totalMessagesReceived)
	fmt.Println("**************************************************")
	node.leftCountryNode.displaySouce()
	node.rightCountryNode.displaySouce()
}
func (tree *CountryTree) displaySource() {
	if tree.CountryNodeRoot == nil {
		return
	}
	tree.CountryNodeRoot.displaySouce()
}

func (node *OperatorNode) findCountryAndOperatorDisplay(name string) {
	if node == nil {
		return
	}
	compare := strings.Compare(node.operas.operatorName, name)
	if compare < 0 {
		node.leftOperatorNode.findCountryAndOperatorDisplay(name)
	} else if compare > 0 {
		node.rightOperatorNode.findCountryAndOperatorDisplay(name)
	} else {
		node.displaySingleOperator()
		return
	}
}
func (node *CountryNode) findCountryAndOperatorDisplay(countryName, opsName string) {
	if node == nil {
		return
	}
	compare := strings.Compare(node.countryName, countryName)
	if compare < 0 {
		node.leftCountryNode.findCountryAndOperatorDisplay(countryName, opsName)
	} else if compare > 0 {
		node.rightCountryNode.findCountryAndOperatorDisplay(countryName, opsName)
	} else {
		node.operatorRoot.findCountryAndOperatorDisplay(opsName)
		return
	}
}

func (tree *CountryTree) findCountryAndOperatorDisplay(countryName, opsName string) {
	if tree.CountryNodeRoot == nil {
		return
	}
	tree.CountryNodeRoot.findCountryAndOperatorDisplay(countryName, opsName)
}

// ****************************************************************************************************************
// 												FIND AND INCREMENT

func (node *OperatorNode) findAndIncrementOPerator(typeIncrement, opsName string) bool {
	if node == nil {
		return false
	}
	compare := strings.Compare(node.operas.operatorName, opsName)
	if compare == 0 {
		switch typeIncrement {
		case "totalMessagesLessThan10mins":
			node.operas.totalMessagesLessThan10mins++
		case "totalMessagesLessThan1min":
			node.operas.totalMessagesLessThan1min++
		case "totalMessagesLessThan1hour":
			node.operas.totalMessagesLessThan1hour++
		case "totalMessagesLessThan2hour":
			node.operas.totalMessagesLessThan2hour++
		case "totalMessagesDLRS":
			node.operas.totalMessagesDLRS++
		case "totalMessagesLessThan10seconds":
			node.operas.totalMessagesLessThan10seconds++
		case "totalMessagesReceived":
			node.operas.totalMessagesReceived++
		case "totalMessagesDelivered":
			node.operas.totalMessagesDelivered++
		case "totalMessagesRejected":
			node.operas.totalMessagesRejected++
		case "totalMessagesInitial":
			node.operas.totalMessagesInitial++
		case "totalMessagesExpired":
			node.operas.totalMessagesExpired++
		case "totalMessagesDeleted":
			node.operas.totalMessagesDeleted++
		case "totalMessagesUndeliverable":
			node.operas.totalMessagesUndeliverable++
		case "totalMessagesDeliveredDirect":
			node.operas.totalMessagesDeliveredDirect++
		default:
		}
		return true
	} else if compare == -1 {
		return node.leftOperatorNode.findAndIncrementOPerator(typeIncrement, opsName)
	} else {
		return node.rightOperatorNode.findAndIncrementOPerator(typeIncrement, opsName)
	}
}

func (node *CountryNode) findAndIncrementOPerator(typeIncrement, opsName, countryName string) bool {
	if node == nil {
		return false
	}
	compare := strings.Compare(node.countryName, countryName)
	if compare == 0 {
		return node.operatorRoot.findAndIncrementOPerator(typeIncrement, opsName)
	} else if compare == -1 {
		return node.leftCountryNode.findAndIncrementOPerator(typeIncrement, opsName, countryName)
	} else {
		return node.rightCountryNode.findAndIncrementOPerator(typeIncrement, opsName, countryName)
	}
}

func (tree *CountryTree) findAndIncrement(typeName, opsName, countryName string) bool {
	if tree.CountryNodeRoot == nil {
		return false
	}
	return tree.CountryNodeRoot.findAndIncrementOPerator(typeName, opsName, countryName)
}

// ****************************************************************************************************************
// 												TOTAL UP

func (node *OperatorNode) totalUpOperators(a, b, c, d, e, f, g *int) {
	if node == nil {
		return
	}
	*a += node.operas.totalMessagesDLRS
	*b += node.operas.totalMessagesLessThan10mins
	*c += node.operas.totalMessagesLessThan10seconds
	*d += node.operas.totalMessagesLessThan1hour
	*e += node.operas.totalMessagesLessThan1min
	*f += node.operas.totalMessagesLessThan2hour
	*g += node.operas.totalMessagesReceived
	node.leftOperatorNode.totalUpOperators(a, b, c, d, e, f, g)
	node.rightOperatorNode.totalUpOperators(a, b, c, d, e, f, g)
}

func (node *CountryNode) totalUpCountry() {
	if node == nil {
		return
	}
	node.operatorRoot.totalUpOperators(&node.totalMessagesDLRS, &node.totalMessagesLessThan10mins, &node.totalMessagesLessThan10seconds, &node.totalMessagesLessThan1hour, &node.totalMessagesLessThan1min, &node.totalMessagesLessThan2hour, &node.totalMessagesReceived)
	node.leftCountryNode.totalUpCountry()
	node.rightCountryNode.totalUpCountry()

}
func (tree *CountryTree) totalTree() {
	if tree.CountryNodeRoot == nil {
		return
	}
	tree.CountryNodeRoot.totalUpCountry()
}

// ****************************************************************************************************************
// 												Number of Country in the Tree

func (node *CountryNode) numberOfCountryNodes() int {
	if node == nil {
		return 0
	}
	return node.leftCountryNode.numberOfCountryNodes() + node.rightCountryNode.numberOfCountryNodes() + 1
}

func (tree *CountryTree) numberOfCountry() int {
	if tree.CountryNodeRoot == nil {
		return 0
	}
	return tree.CountryNodeRoot.numberOfCountryNodes()
}

// ****************************************************************************************************************
// 												Number of OPerators in One Country

func (node *OperatorNode) numberOfOperatorsNode() int {
	if node == nil {
		return 0
	}
	return node.leftOperatorNode.numberOfOperatorsNode() + node.rightOperatorNode.numberOfOperatorsNode() + 1
}

func (node *CountryNode) numberOfOperatorsPerCountry(countryName string) int {
	if node == nil {
		return 0
	}
	compare := strings.Compare(node.countryName, countryName)
	if compare == 0 {
		return node.operatorRoot.numberOfOperatorsNode()
	} else if compare < 0 {
		return node.leftCountryNode.numberOfOperatorsPerCountry(countryName)
	} else {
		return node.rightCountryNode.numberOfOperatorsPerCountry(countryName)
	}
}
func (tree *CountryTree) numberOfOperatorsPerCountry(name string) int {
	if tree.CountryNodeRoot == nil {
		fmt.Println("Tree is empty")
		return 0
	}
	return tree.CountryNodeRoot.numberOfOperatorsPerCountry(name)
}

// ****************************************************************************************************************
// 									SDR ALERT BASED ON DATA

func sendAlertToSlack(payload map[string]interface{}) {
	webHookURL := "https://hooks.slack.com/services/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"
	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
}

func (node *OperatorNode) calculateSDRPercentage(sourceName string) {
	if node == nil {
		return
	}
	averageOneMinute := (float64(node.operas.totalMessagesLessThan1min) / float64(node.operas.totalMessagesReceived)) * 100
	averageTenMinutes := (float64(node.operas.totalMessagesLessThan10mins) / float64(node.operas.totalMessagesReceived)) * 100
	averageOneHour := (float64(node.operas.totalMessagesLessThan1hour) / float64(node.operas.totalMessagesReceived)) * 100
	averageTwoHours := (float64(node.operas.totalMessagesLessThan2hour) / float64(node.operas.totalMessagesReceived)) * 100
	averageTenSeconds := (float64(node.operas.totalMessagesLessThan10seconds) / float64(node.operas.totalMessagesReceived)) * 100
	averageStatus := (float64(node.operas.totalMessagesDLRS) / float64(node.operas.totalMessagesReceived)) * 100
	payload := make(map[string]interface{})
	payload["username"] = "SDR BASED ALERT"
	var str bytes.Buffer
	str.WriteString("Country Name: " + sourceName + "\nOperator Name: " + node.operas.operatorName)
	if averageOneMinute > 60 {
		str.WriteString("\nPercentage Under 1 minute: " + strconv.FormatFloat(averageOneMinute, 'f', 6, 64))
	}
	if averageTenMinutes > 60 {
		str.WriteString("\nPercentage Under Ten Minutes: " + strconv.FormatFloat(averageTenMinutes, 'f', 6, 64))
	}
	if averageOneHour > 60 {
		str.WriteString("\nPercentage Under 1 Hour: " + strconv.FormatFloat(averageOneHour, 'f', 6, 64))
	}
	if averageTwoHours > 60 {
		str.WriteString("\nPercentage Under 2 Hours: " + strconv.FormatFloat(averageTwoHours, 'f', 6, 64))
	}
	if averageTenSeconds > 60 {
		str.WriteString("\nPercentage Under 10 seconds: " + strconv.FormatFloat(averageTenSeconds, 'f', 6, 64))
	}
	if averageStatus > 60 {
		str.WriteString("\nPercentage Status (4): " + strconv.FormatFloat(averageStatus, 'f', 6, 64))
	}
	payload["text"] = str.String()
	sendAlertToSlack(payload)
	node.leftOperatorNode.calculateSDRPercentage(sourceName)
	node.rightOperatorNode.calculateSDRPercentage(sourceName)
}

func (node *CountryNode) calculateSDRPercentage() {
	if node == nil {
		return
	}
	node.operatorRoot.calculateSDRPercentage(node.countryName)
	node.leftCountryNode.calculateSDRPercentage()
	node.rightCountryNode.calculateSDRPercentage()
}
func (tree *CountryTree) calculateSDRPercentage() {
	if tree.CountryNodeRoot == nil {
		return
	}
	tree.CountryNodeRoot.calculateSDRPercentage()
}

// ****************************************************************************************************************
// 												SDR in one country

func (node *CountryNode) sdrAlertOneCountry(name string) {
	if node == nil {
		return
	}
	compare := strings.Compare(node.countryName, name)
	if compare == 0 {
		node.operatorRoot.calculateSDRPercentage(node.countryName)
		return
	} else if compare < 0 {
		node.leftCountryNode.sdrAlertOneCountry(name)
	} else {
		node.rightCountryNode.sdrAlertOneCountry(name)
	}
}

func (tree *CountryTree) sdrAlertOneCountry(name string) {
	if tree.CountryNodeRoot == nil {
		return
	}
	tree.CountryNodeRoot.sdrAlertOneCountry(name)
}
