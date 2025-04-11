package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Base    string `json:"base"`
	Rates   map[string]struct {
		EURXAU float64 `json:"EURXAU"`
	} `json:"rates"`
}

type RelevantData struct {
	Date         string
	GoldPriceEur string
}

func main() {

	fmt.Println("-> start metal-price thing v3.1")

	// basic config
	const apiKey = "ae945d3518dbfb8dc564ee795a8e9251"

	const baseURL = "https://api.metalpriceapi.com/v1/timeframe"
	const periodeDuration = 6
	const numberOfIterations = 100
	const storageFile = "output.csv"

	// Start and end dates

	endDate, err := getNewBeginDateFromFile(storageFile)
	if err != nil {
		fmt.Printf("Dont get good begin time %v\n", err)
		//endDate = time.Date(2020, 11, 20, 0, 0, 0, 0, time.UTC)
		return
	}
	startDate := endDate.AddDate(0, 0, periodeDuration*-1+1)

	// Open CSV file for writing
	file, err := os.OpenFile(storageFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	//writer.Write([]string{"-", "-", "-"})

	// Loop to make 100 requests
	for i := 0; i < numberOfIterations; i++ {
		var allRelDat []RelevantData

		// Format the request URL
		url := fmt.Sprintf("%s?api_key=%s&start_date=%s&end_date=%s&base=EUR&currencies=XAU",
			baseURL, apiKey, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

		// Make the API request
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		defer resp.Body.Close()

		// Parse the response JSON
		var response Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println("Error decoding JSON response:", err)
			continue
		}

		// Write data to slice
		for date, rate := range response.Rates {
			allRelDat = append(allRelDat, RelevantData{
				Date:         date,
				GoldPriceEur: fmt.Sprintf("%.4f", rate.EURXAU),
			})
		}

		// Sort the data by date and write to CSV
		sort.Slice(allRelDat, func(i, j int) bool {
			return allRelDat[i].Date > allRelDat[j].Date
		})

		for _, v := range allRelDat {
			// i add 3 empty fields that it is compadible with the csv format of alphaV
			writer.Write([]string{v.Date, "", "", "", v.GoldPriceEur})
		}

		// Move the date range by periode
		startDate = startDate.AddDate(0, 0, periodeDuration*-1)
		endDate = endDate.AddDate(0, 0, periodeDuration*-1)

		// fmt.Println(url)
		fmt.Print("-> summary :: i:")
		fmt.Print(i)
		fmt.Printf(" / first element: %v\n", allRelDat[0])

		time.Sleep(1 * time.Second)
	}

	fmt.Println("Data written to output.csv")
}

// ExtractDateFromFile extracts the date from the last line of the file and returns it as a time.Time object
func getNewBeginDateFromFile(fileName string) (time.Time, error) {
	// Open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Use a scanner to read the file line by line
	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text() // Keep updating the last line
	}

	// Handle potential errors from scanning
	if err := scanner.Err(); err != nil {
		return time.Time{}, fmt.Errorf("failed to read the file: %v", err)
	}

	// Extract the date part (before the first comma)
	parts := strings.Split(lastLine, ",")
	dateStr := parts[0] // Assuming the date is the first part in YYYY-MM-DD format

	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %v", err)
	}

	// remove one day
	date = date.AddDate(0, 0, -1)

	return date, nil
}
