package services

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"kursoverview/internal/config"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Results struct {
	Datum_s     string
	Amount_euro float64
	Shares      float64 // for interest: this field is used as the avarage amount of euro
	Share_price float64 // for interest: this field is used for the interesst rate
	Fee         float64 // for interest: this field is used for the tax
	Type        string
	Asset       string
}

var TargetOutput = config.DataPath + config.TrLogFile

func LogTrTransaction(input string) string {

	// input: is the clipboard content we get from the webUI

	responseText := ""

	// how to use: go to trade republic -> Profile -> open one transaction
	// -> copy all content of the browser window to the clipboard (CTRL+C)
	// execute this in a terminal "go run main.go"

	fmt.Println("TR transaction parser v0.3")

	overContent := getOverviewContent(input)
	rContent := takeRelevantPart(input)

	var transactionDetails Results

	transactionDetails.addOverviewValues(rContent)
	if transactionDetails.Type == "interest+" {
		transactionDetails.addInterest(rContent)
	} else if transactionDetails.Type == "cash-dividend" {
		transactionDetails.addDividend(rContent)
	} else if transactionDetails.Type == "banktransfer+" || transactionDetails.Type == "banktransfer-" {
		transactionDetails.addBankAmount(rContent)
	} else {
		// details for sell or buy of shares/etfs
		transactionDetails.addTransactionValues(rContent)
	}

	csvLine, _ := transactionDetails.ToCsvFormat()

	responseText = saveToFile(csvLine, overContent)

	return responseText
}

func saveToFile(newCsvString string, overviewContent []string) string {
	currentStorageContent := readCurrentStorageFile()
	outLoglines := ""
	fmt.Println("sdfsdfpenis ss")

	notPresent := true
	for _, existingLine := range currentStorageContent {
		if existingLine == newCsvString {
			notPresent = false
		}
	}

	if notPresent {
		err := writeToFile(append(currentStorageContent, newCsvString))
		if err != nil {
			log.Fatal("Write file:", err)
		}
		outLoglines = outLoglines + "\n-> successful added to file."
	} else {
		outLoglines = outLoglines + "\n-> not added to file because already there."
	}

	// check if all records that are shown on the website are downloaded
	outLoglines = outLoglines + "\n   " + newCsvString
	outLoglines = outLoglines + "\n"

	fmt.Println("\nin storage - from clipboard overview")
	for i, overviewValue := range overviewContent {
		if len(currentStorageContent) == i {
			outLoglines = outLoglines + "\n-> looks like all transaction are here."
			break
		}

		storedValues := strings.Split(currentStorageContent[i], ",")
		if len(storedValues) == 0 {
			fmt.Print("corupted line 102:" + currentStorageContent[i])
			break
		}

		fmt.Printf("compare: %v :: %v - %v \n", i, storedValues[1], overviewValue)

		if storedValues[1] != overviewValue {
			outLoglines = outLoglines + "\n-> A transaction is missing in the storage."
			outLoglines = outLoglines + "\n   overview record number (from the top): " + fmt.Sprintf("%d / %d", i+1, len(overviewContent))
			outLoglines = outLoglines + "\n   with the value[€]: " + overviewValue
			break
		}
	}

	return outLoglines
}

func takeRelevantPart(input string) []string {
	var output []string

	contentLines := strings.Split(input, "\n")
	counter := 0
	storeLine := false
	findBegin := true
	stopLineKeyWords := map[string]bool{
		"Documents": true,
		"word2":     true,
	}

	for _, line := range contentLines {
		if line == "" {
			continue
		}

		if findBegin {
			if line[0] != ' ' {
				// add to the counter when the first character is not empty (header)
				counter++
			}

			if counter == 6 {
				// after six sections we have the beginning of the transaction. no overview data anymore.
				findBegin = false
				storeLine = true
			}
		}

		if stopLineKeyWords[line] {
			storeLine = false
		}

		if storeLine {
			output = append(output, line)
		}
	}

	return output
}

func (r Results) ToCsvFormat() (string, error) {
	// converts the Results struct to a CSV string

	var buffer strings.Builder

	// Create a CSV writer
	writer := csv.NewWriter(&buffer)

	// Convert the struct fields to a slice of strings
	record := []string{
		r.Datum_s,
		fmt.Sprintf("%.2f", r.Amount_euro),
		fmt.Sprintf("%.6f", r.Shares),
		fmt.Sprintf("%.2f", r.Share_price),
		fmt.Sprintf("%.2f", r.Fee),
		r.Type,
		r.Asset,
	}

	// Write the record to the CSV
	if err := writer.Write(record); err != nil {
		return "", err
	}

	// Flush the writer to ensure all data is written
	writer.Flush()

	// Check for errors during flushing
	if err := writer.Error(); err != nil {
		return "", err
	}

	// Return the CSV string and not error
	resultCsvString := strings.Replace(buffer.String(), "\n", "", 1)
	return resultCsvString, nil
}

func writeToFile(completeText []string) error {

	// sort descending
	sort.Slice(completeText, func(i, j int) bool {
		return completeText[i] > completeText[j]
	})

	// make a string with header
	newContent := strings.Join(completeText, "\n")
	newContent = "date,ammount,shares,price,tax/fee,type,partner\n" + newContent

	// Write the string to the file
	err := os.WriteFile(TargetOutput, []byte(newContent), 0644)
	if err != nil {
		log.Fatal("Error writing to file:", err)
		return nil
	}

	return nil
}

func secondCleaning(input []string) []string {
	output := input

	if input[9] == "Performance" {
		output = input[:9]
		output = append(output, input[14:]...)
	}

	return output
}

func readCurrentStorageFile() []string {
	var output []string

	// Open the file
	file, err := os.Open(TargetOutput)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close() // Ensure the file is closed when done

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read each line
	for scanner.Scan() {
		line := scanner.Text() // Get the current line as a string
		output = append(output, line)
	}

	// skip header
	return output[1:]
}

func transformDate(input string) string {

	input = addCurrentYearIfYearIsMissing(input)

	// Define layout
	inputLayout := "2 January 2006 at 15:04"
	outputLayout := "2006-01-02 15:04"

	// Parse the input string into a time.Time object
	parsedTime, err := time.Parse(inputLayout, input)
	if err != nil {
		fmt.Println("Error 262 parsing date :: ", err)
		return ""
	}

	// Format the time into the desired numeric format
	output := parsedTime.Format(outputLayout)

	return output
}

func addCurrentYearIfYearIsMissing(input string) string {

	output := input

	// Compile the regular expression
	pattern := "(^\\d+\\s+[A-z]+\\s+)(at.*)"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Error compiling regex:", err)
	}

	// execute regex and store match groups in variable if pattern was found
	matches := regex.FindStringSubmatch(input)
	if matches != nil {

		currentYearStr := fmt.Sprintf(" %d", time.Now().Year())
		output = matches[1] + currentYearStr + " " + matches[2]
	}

	return output
}

func (r *Results) addOverviewValues(input []string) {

	var err error
	r.Datum_s = transformDate(input[1])

	inOverview := false
	for i, ele := range input {
		if ele == "Transaction" {
			break
		}
		if ele == "Overview" {
			inOverview = true
		}
		if inOverview {
			// used for share buy, sell and devidend actions
			if strings.ToLower(ele) == "order type" {
				r.Type = strings.ToLower(strings.ReplaceAll(input[i+1][4:], " ", "-"))
			}
			if ele == "Asset" {
				r.Asset = strings.ReplaceAll(input[i+1][4:], ",", ";")
			}

			// used for banktransfer
			if ele == "From" {
				r.Type = "banktransfer+"
			}
			if ele == "To" {
				r.Type = "banktransfer-"
			}
			if ele == "IBAN" {
				r.Asset = strings.ReplaceAll(input[i+1][4:], " ", "")
			}

			// used for interest
			if ele == "Average balance" {
				r.Type = "interest+"
				r.Shares, err = strconv.ParseFloat(strings.ReplaceAll(input[i+1][7:], ",", ""), 64)
				if err != nil {
					log.Fatal("Error 326:", err)
				}
			}
			if r.Type == "interest+" && strings.ToLower(ele) == "annual rate" {
				r.Share_price, err = strconv.ParseFloat(strings.ReplaceAll(input[i+1][4:], " %", ""), 64)
				if err != nil {
					log.Fatal("Error 329:", err)
				}
			}

			// used for dividend
			if ele == "Event" {
				if input[i+1][4:] == "Cash dividend" {
					r.Type = "cash-dividend"
				}
			}
		}
	}

	// check if treatment was success
	if r.Type == "" {
		for _, ele := range input {
			fmt.Println(">>> " + ele)
		}
		log.Fatal("Transaction Type empty")
	}
	if r.Datum_s == "" {
		log.Fatal("Transaction Date empty")
	}
	if r.Asset == "" {
		log.Fatal("Transaction Asset empty")
	}
}

func (r *Results) addTransactionValues(input []string) {

	inTransaction := false
	for i, ele := range input {

		var err error

		if ele == "Transaction" {
			inTransaction = true
		}
		if inTransaction {
			if ele == "Shares" {
				r.Shares, err = strconv.ParseFloat(input[i+1][4:], 64)
				if err != nil {
					log.Fatal("Error 76:", err)
				}
			}

			if ele == "Share price" {
				r.Share_price, err = strconv.ParseFloat(strings.ReplaceAll(input[i+1][7:], ",", ""), 64)
				if err != nil {
					log.Fatal("Error 81:", err)
				}
			}

			if ele == "Fee" {
				r.Fee, err = strconv.ParseFloat(input[i+1][7:], 64)
				if err != nil {
					log.Fatal("Error 86:", err)
				}
			}

			if ele == "Total" {
				tmp := strings.NewReplacer(",", "", "€", "", "+", "", "-", "", " ", "").Replace(input[i+1])
				r.Amount_euro, err = strconv.ParseFloat(tmp, 64)
				if err != nil {
					log.Fatal("Error 91:", err)
				}
			}
		}

	}

	// check if treatment was success
	if r.Shares == 0 {
		log.Fatal("Transaction Shares empty")
	}
	if r.Share_price == 0 {
		log.Fatal("Transaction Share-Price empty")
	}
	if r.Fee == 0 {
		log.Fatal("Transaction Fee empty")
	}
	if r.Amount_euro == 0 {
		log.Fatal("Transaction Amount empty")
	}
}

func (r *Results) addBankAmount(input []string) {

	firstLine := strings.ReplaceAll(input[0], ",", "")

	// Compile the regular expression
	pattern := "[0-9.]+"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Error compiling regex 382:", err)
	}

	// execute regex and store match groups in variable if pattern was found
	matches := regex.FindStringSubmatch(firstLine)
	if matches != nil {
		r.Amount_euro, err = strconv.ParseFloat(matches[0], 64)
		if err != nil {
			log.Fatal("Error 391:", err)
		}
	}

	if r.Amount_euro == 0 {
		log.Fatal("Transaction-BT Amount empty")
	}
}

func (r *Results) addInterest(input []string) {

	r.Fee = -1
	r.Amount_euro = -1

	inTransaction := false
	for i, ele := range input {

		var err error

		if ele == "Transaction" {
			inTransaction = true
		}
		if inTransaction {

			if ele == "Accrued" {
				r.Amount_euro, err = strconv.ParseFloat(strings.ReplaceAll(input[i+1][7:], ",", ""), 64)
				if err != nil {
					log.Fatal("Error 454:", err)
				}
			}

			if ele == "Tax" {
				r.Fee, err = strconv.ParseFloat(input[i+1][7:], 64)
				if err != nil {
					log.Fatal("Error 461:", err)
				}
				fmt.Println(r.Fee)
			}
		}

	}

	// sometimes transaction data are not present. then take it from the title
	if !inTransaction && strings.Contains(input[0], "EUR") {

		// Compile the regular expression
		pattern := "[0-9.,]+"
		regex, err := regexp.Compile(pattern)
		if err != nil {
			log.Fatal("Error compiling regex:", err)
		}

		// execute regex and store match groups in variable if pattern was found
		matches := regex.FindStringSubmatch(input[0])
		if matches != nil {
			r.Amount_euro, err = strconv.ParseFloat(strings.ReplaceAll(matches[0], ",", ""), 64)
			if err != nil {
				log.Fatal("Error 492:", err)
			}
		}
	}

	// check if treatment was success
	if r.Shares == 0 {
		log.Fatal("Interest Shares empty")
	}
	if r.Share_price == 0 {
		log.Fatal("Interest Share-Price empty")
	}
	if r.Fee == -1 && inTransaction {
		log.Fatal("Interest Fee empty")
	}
	if r.Amount_euro == -1 {
		log.Fatal("Interest Amount empty")
	}
}

func (r *Results) addDividend(input []string) {

	r.Fee = -1
	r.Amount_euro = -1

	inTransaction := false
	for i, ele := range input {

		var err error

		if ele == "Transaction" {
			inTransaction = true
		}
		if inTransaction {

			if ele == "Shares" {
				r.Shares, err = strconv.ParseFloat(input[i+1][4:], 64)
				if err != nil {
					log.Fatal("Error 514:", err)
				}
			}

			if ele == "Dividend per share" {
				r.Share_price = getFloatNumber(input[i+1])
			}

			if ele == "Tax" {
				r.Fee, err = strconv.ParseFloat(input[i+1][7:], 64)
				if err != nil {
					log.Fatal("Error 528:", err)
				}
			}

			if ele == "Total" {
				tmp := strings.NewReplacer(",", "", "€", "", "+", "", "-", "", " ", "").Replace(input[i+1])
				r.Amount_euro, err = strconv.ParseFloat(tmp, 64)
				if err != nil {
					log.Fatal("Error 536:", err)
				}
			}
		}

	}

	// check if treatment was success
	if r.Shares == 0 {
		log.Fatal("Dividend Shares empty")
	}
	if r.Share_price == 0 {
		log.Fatal("Dividend Share-Price empty")
	}
	if r.Fee == -1 {
		log.Fatal("Dividend Fee empty")
	}
	if r.Amount_euro == -1 {
		log.Fatal("Dividend Amount empty")
	}
}

func getOverviewContent(input string) []string {
	var OverviewTransactionValues []string

	contentLines := strings.Split(input, "\n")
	// Compile the regular expression
	pattern := "[+]*€([0-9.,]+)"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Error compiling regex:", err)
	}

	for _, line := range contentLines {

		// execute regex and store match groups in variable if pattern was found
		matches := regex.FindStringSubmatch(line)
		if matches != nil {
			OverviewTransactionValues = append(OverviewTransactionValues, strings.Replace(matches[1], ",", "", 1))
		}
	}

	return OverviewTransactionValues
}

func getFloatNumber(input string) float64 {

	var output float64

	// remove unwanted character
	input = strings.NewReplacer(" ", "", ",", "").Replace(input)

	// Compile the regular expression
	pattern := "[0-9.]+"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Error 613 compiling regex:", err)
	}

	// execute regex
	match := regex.FindStringSubmatch(input)
	if match == nil {
		log.Fatal("Error 621 no match found regex:", err, input)
	}

	// transform to float
	output, err = strconv.ParseFloat(match[0], 64)
	if err != nil {
		log.Fatal("Error 628 floatParsing:", err, input)
	}

	return output
}
