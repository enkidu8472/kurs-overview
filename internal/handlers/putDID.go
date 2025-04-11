package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"kursoverview/internal/config"
	"net/http"
	"os"
	"strings"
	"time"
)

type AssetInfoEOD struct {
	Code              string  `json:"Code"`              // Example: "VUSA"
	Exchange          string  `json:"Exchange"`          // Example: "LSE"
	Name              string  `json:"Name"`              // Example: "Vanguard S&P 500 UCITS ETF"
	Type              string  `json:"Type"`              // Example: "ETF"
	Country           string  `json:"Country"`           // Example: "UK"
	Currency          string  `json:"Currency"`          // Example: "GBP"
	ISIN              string  `json:"ISIN"`              // Example: "IE00B3XXRP09"
	PreviousClose     float64 `json:"previousClose"`     // Example: 90.8225
	PreviousCloseDate string  `json:"previousCloseDate"` // Example: "2024-12-13"
}

type requestBodyDID struct {
	IsinRaw string `json:"isin"`
}

type responseDID struct {
	ErrorMessage string
	Isin         string
	Name         string
	LastDate     string
	AssetCode    string
	ExchangeCode string
	NumberOfDays int
	FirstDate    string
	CurrentPrice float64
	Currency     string
	Type         string
	Country      string
}

// root function of this file: download info from a thrird party and store them localy
func DownloadIsinData(w http.ResponseWriter, r *http.Request) {
	// Create a slice of Asset structs
	var inputBody requestBodyDID

	// Read the request body
	rawRequestbody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error 53")
		return
	}
	err = json.Unmarshal(rawRequestbody, &inputBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		fmt.Println("error 59")
		return
	}

	fmt.Println("endpoint DID called:" + inputBody.IsinRaw)

	err = inputBody.sainityCheck()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, _ := inputBody.mainTreatment()

	// convert the response into JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getEodData(inputIsin string) ([]AssetInfoEOD, bool, error) {
	var output []AssetInfoEOD

	filePath := config.DataIsinPath + inputIsin + "/" + config.EodSearchFileName

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return output, true, nil
	}
	defer file.Close()

	// Read the contents of the file
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return output, false, fmt.Errorf("failed to read file: %v", err)
	}

	// Unmarshal the JSON into the output struct

	err = json.Unmarshal(fileContents, &output)
	if err != nil {
		return output, false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return output, false, nil
}

func downloadDailyStockData(possibleExchanges []AssetInfoEOD, isin string) AssetInfoEOD {

	var usedItem AssetInfoEOD
	var itemsEur []AssetInfoEOD
	var itemsOther []AssetInfoEOD
	var itemsGermany []AssetInfoEOD
	var itemsToCheck []AssetInfoEOD
	var itemLastSucces AssetInfoEOD
	var err error

	// read last successfull symbol
	// get inforamtion from EOD metadata file
	var fileContent responseDID
	filename := config.DataIsinPath + isin + "/" + config.EodLastUsedSymbolFileName
	fileContent, err = readJsonFile(filename, fileContent)
	if err == nil {
		itemLastSucces.Exchange = fileContent.ExchangeCode
		itemLastSucces.Code = fileContent.AssetCode
		itemLastSucces.ISIN = fileContent.Isin
	}

	// check from the data from https://eodhistoricaldata.com/api/search/ which one has valid information
	// > either from Germany or with EUR... stocks from Germany has priority
	for _, pe := range possibleExchanges {

		if itemLastSucces.ISIN != "" && itemLastSucces.Exchange == pe.Exchange && itemLastSucces.Code == pe.Code {
			itemLastSucces = pe
			continue
		}

		if pe.Currency == "EUR" && pe.Country == "Germany" {
			itemsGermany = append(itemsGermany, pe)
			continue
		}

		if pe.Currency == "EUR" {
			itemsEur = append(itemsEur, pe)
			continue
		}

		itemsOther = append(itemsOther, pe)
	}

	// merge info from three sources. 1. check should be last succeded item. 2. from symbols from Germany. 3. items in Euro
	if itemLastSucces.ISIN != "" {
		itemsToCheck = append(itemsToCheck, itemLastSucces)
	}
	itemsToCheck = append(itemsToCheck, itemsGermany...)
	itemsToCheck = append(itemsToCheck, itemsEur...)
	itemsToCheck = append(itemsToCheck, itemsOther...)

	// loop through possible symbols until download was successfull
	for _, itm := range itemsToCheck {
		symbol := itm.Code + "." + itm.Exchange
		succees := callAlphavantageEndpoint(symbol, itm.ISIN)
		if succees {
			fmt.Println("Successful download with: " + symbol + " for:" + itm.ISIN)
			usedItem = itm
			break
		}
	}

	return usedItem
}

func callEodEndpoint(isin string) {

	var responseData []AssetInfoEOD
	var storagePath = config.DataIsinPath + isin + "/"

	// make a new directory if not existing yet
	err := os.MkdirAll(storagePath, os.ModePerm) // os.ModePerm gives default permissions (0777)
	if err != nil {
		fmt.Printf("Error making directory: %v\n", err)
		return
	}

	// Construct the API URL
	url := fmt.Sprintf("https://eodhistoricaldata.com/api/search/%s?api_token=%s", isin, config.ApiTokenEodhd)

	// Make the HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", response.StatusCode)
		return
	}

	// Print the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Parse JSON into a slice of ETF structs
	errJson := json.Unmarshal(body, &responseData)
	if errJson != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	filePath := storagePath + config.EodSearchFileName
	err = writeToJsonFile(filePath, responseData)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
}

func callAlphavantageEndpoint(symbol string, isin string) bool {

	successfulTreatment := false

	var storagePath = config.DataIsinPath + isin + "/" + config.AlphaTimeSerieFileName

	// Construct the API URL
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s&datatype=csv&outputsize=full", symbol, config.ApiTokenAlpha)
	fmt.Println(url)

	// Make the HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return successfulTreatment
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", response.StatusCode)
		return successfulTreatment
	}

	// Print the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return successfulTreatment
	}

	if strings.Contains(string(body), "Information") {
		fmt.Println("Respons body contains informtation: " + string(body))
		return successfulTreatment
	}

	if strings.Contains(string(body), "Error Message") {
		fmt.Println("Respons body contains error message: " + string(body))
		return successfulTreatment
	}

	// Create the file or open it for writing
	file, err := os.Create(storagePath)
	if err != nil {
		fmt.Printf("failed to create file: %v", err)
		return successfulTreatment
	}
	defer file.Close()

	// Write to the file
	_, err = file.Write(body)
	if err != nil {
		fmt.Printf("failed to write to file: %v", err)
		return successfulTreatment
	}

	successfulTreatment = true
	return successfulTreatment
}

func (input requestBodyDID) sainityCheck() error {

	// not to long. iso norm say that all isins are exactly 12 character long
	if len(input.IsinRaw) > 12 {
		fmt.Println("Stop: ISIN to long")
		return fmt.Errorf("error: ISIN to long")
	}

	// isin is not upper case
	if strings.ToUpper(input.IsinRaw) != input.IsinRaw {
		fmt.Println("Stop: ISIN not only uppercase characters")
		return fmt.Errorf("error, ISIN not only uppercase characters")
	}

	// Check if the file was modified today
	filename := config.DataIsinPath + input.IsinRaw + "/" + config.AlphaTimeSerieFileName
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("No file found to check modification date")
		fmt.Printf("Error getting file info: %v\n", err)
		return nil
	}

	modTime := fileInfo.ModTime()

	today := time.Now().Truncate(24 * time.Hour) // Truncate to the start of the day
	lastModDate := modTime.Truncate(24 * time.Hour)
	fmt.Printf("Data file TS: %v\n", lastModDate)

	if lastModDate.Equal(today) {
		fmt.Println("Stop: file was already modfified today")
		return fmt.Errorf("file was already modified today. Skipping treatment")
	}

	return nil
}

func (input requestBodyDID) mainTreatment() (responseDID, error) {
	inputIsin := strings.ToUpper(input.IsinRaw)
	var out responseDID

	// call endpoint to get symbols for current isin from different exchanges
	dataEod, nofilefound, _ := getEodData(inputIsin)
	if nofilefound {
		callEodEndpoint(inputIsin)
		dataEod, _, _ = getEodData(inputIsin)
	}

	successItem := downloadDailyStockData(dataEod, inputIsin)

	if successItem.Code == "" {
		out.ErrorMessage = "No data downloaded"
	} else {
		out.Isin = successItem.ISIN
		out.AssetCode = successItem.Code
		out.ExchangeCode = successItem.Exchange
		out.CurrentPrice = successItem.PreviousClose
		out.Name = successItem.Name
		out.LastDate = successItem.PreviousCloseDate
		out.Currency = successItem.Currency
		out.Type = successItem.Type
		out.Country = successItem.Country
	}

	// save used symbol-info to "last used file"
	filename := config.DataIsinPath + out.Isin + "/" + config.EodLastUsedSymbolFileName
	writeToJsonFile(filename, out)

	return out, nil
}
