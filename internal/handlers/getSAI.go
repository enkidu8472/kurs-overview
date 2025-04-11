package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"kursoverview/internal/config"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DataPoint struct {
	Date   string  `json:"date"`
	Price  float64 `json:"price"`
	PPrice float64 `json:"percentPrice"` // takes the first price of the data set and set it as 100%, in relation to this value all other prices are transformed
}

type additionalCurveInfo struct {
	Title       string `json:"title"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	SymbolCode  string `json:"symbolCode"`
	Exchange    string `json:"exchange"`
	Type        string `json:"type"`
	Country     string `json:"country"`
	Currency    string `json:"currency"`
	Color       string `json:"color"`
	Duration    string `json:"duration"`
	Perf1m      string `json:"perf1m"`
	Perf6m      string `json:"perf6m"`
	Perf1y      string `json:"perf1y"`
	Perf5y      string `json:"perf5y"`
	Tags        string `json:"tags"`
}

type Curve struct {
	Name              string              `json:"name"`
	Values            []DataPoint         `json:"values"`
	Color             string              `json:"color"`
	TotalPercentage   float64             `json:"totalPercentage"`   // change during the current period in percentage
	PercentagePerYear float64             `json:"percentagePerYear"` // change in percentage over the current time, in relation to one year
	AdditionalInfo    additionalCurveInfo `json:"additionalInfo"`
}

type requestBodyGAV struct {
	Isins    []string `json:"isins"`
	FromDate string   `json:"fromDate"`
	ToDate   string   `json:"toDate"`
}

type Boundaries struct {
	Date   []string  `json:"date"`
	Price  []float64 `json:"price"`
	PPrice []float64 `json:"pprice"`
}

type Periode struct {
	Years string
	Days  int
}

type Assets struct {
	CurveBoundaries Boundaries `json:"curveBoundaries"`
	Curves          []Curve    `json:"curves"`
	PeriodeDuration Periode
}

// root function of this file
func GetSelectedAssetInformation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint called Get-Selected-Asset-Information")

	// Create a slice of Asset structs
	var inputBody requestBodyGAV
	var out []Curve

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error 146")
		return
	}

	// Parse the JSON body into the slice
	err = json.Unmarshal(body, &inputBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		fmt.Println("error 154")
		return
	}

	for _, isin := range inputBody.Isins {
		out = append(out, readDataFile(isin))
	}

	// final manuipulation: find boundaries
	out2 := treatFinishedResponse(out, inputBody.FromDate, inputBody.ToDate)

	// convert the struct into JSON
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(out2); err != nil {
		fmt.Println("error 172")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readDataFile(inpIsin string) Curve {

	var fileName string
	var out Curve
	var dataPoints []DataPoint

	for _, fn := range AllDataFileNames {
		if strings.Contains(fn, inpIsin) {
			fileName = fn
			break
		}
	}
	if fileName == "" {
		log.Fatal("Error no file name found! " + inpIsin)
		return out
	}

	// os.Open() opens specific file in. read-only mode and this return. a pointer of type os.File
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error while reading the file", err)
		return out
	}
	defer file.Close()

	// The csv.NewReader() function is called in  which the object os.File passed as its parameter
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file and Returns them as slice of slices of string  and an error if any
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records")
		return out
	}

	for _, eachRecord := range records {
		var dp DataPoint
		dp.Date = eachRecord[0]
		decimalNumber, _ := strconv.ParseFloat(eachRecord[4], 64)
		dp.Price = decimalNumber
		dataPoints = append(dataPoints, dp)
	}

	out.Values = dataPoints[1:] // removes the header from the csv file
	out.Name = inpIsin
	out.enrichWithAdditionalCurveData()

	// handle plot color of the currrent curve
	if out.AdditionalInfo.Color != "" {

		out.Color = out.AdditionalInfo.Color
	} else {

		// source := rand.NewSource(time.Now().UnixNano())
		// rng := rand.New(source)
		// randomNumber := rng.Intn(7)

		// colorPalett := []string{"#F94144", "#F3722C", "#F8961E", "#F9C74F", "#90BE6D", "#43AA8B", "#577590"}
		// out.Color = colorPalett[randomNumber]

		out.Color = "#e6cfcf"
		out.Color = "#ccb1b1"
		out.Color = "#ad8e8e"
	}

	return out
}

func treatFinishedResponse(curveData []Curve, fromDate string, toDate string) Assets {
	var out Assets
	out.Curves = curveData
	var minPrice float64 = 9999
	var maxPrice float64 = -9999
	var minPPrice float64 = 9999
	var maxPPrice float64 = -9999
	var boundarieDates []string

	// filter on date input (maybe need to merge this function with this funtion)
	out = selectCurveDataInDateRange(curveData, fromDate, toDate)

	// make percentage calculation of price
	for i := range out.Curves {
		curve := out.Curves[i]
		var newValues []DataPoint

		referencePrice := curve.Values[len(curve.Values)-1].Price
		for _, ele := range curve.Values {
			ele.PPrice = (ele.Price - referencePrice) * 100 / referencePrice
			newValues = append(newValues, ele)
		}

		out.Curves[i].Values = newValues
		out.Curves[i].TotalPercentage = newValues[0].PPrice

		// calculate pecentage change per year
		layout := "2006-01-02" // Define the date format (assuming YYYY-MM-DD format)
		startDate, err := time.Parse(layout, newValues[len(newValues)-1].Date)
		if err != nil {
			fmt.Println("Error parsing Date 141:", err)
			continue
		}

		endDate, err := time.Parse(layout, newValues[0].Date)
		if err != nil {
			fmt.Println("Error parsing Date 143:", err)
			continue
		}

		timeDiff := (endDate.Sub(startDate).Hours() / 24) / 365 // time diff in years
		out.Curves[i].PercentagePerYear = out.Curves[i].TotalPercentage / timeDiff
	}

	// find lowest and higest values, over all selected assets for plot boundaries
	for _, curve := range out.Curves {
		for _, ele := range curve.Values {
			if ele.Price > maxPrice {
				maxPrice = ele.Price
			}
			if ele.Price < minPrice {
				minPrice = ele.Price
			}

			if ele.PPrice > maxPPrice {
				maxPPrice = ele.PPrice
			}
			if ele.PPrice < minPPrice {
				minPPrice = ele.PPrice
			}

		}

		boundarieDates = append(boundarieDates, curve.Values[0].Date)
		boundarieDates = append(boundarieDates, curve.Values[len(curve.Values)-1].Date)
	}

	sort.Strings(boundarieDates)

	out.CurveBoundaries.Price = []float64{minPrice, maxPrice}
	out.CurveBoundaries.PPrice = []float64{minPPrice, maxPPrice}
	out.CurveBoundaries.Date = []string{boundarieDates[0], boundarieDates[len(boundarieDates)-1]}

	// calculate periode for current data
	layout := "2006-01-02" // Define the date format (assuming YYYY-MM-DD format)
	startDate, err := time.Parse(layout, out.CurveBoundaries.Date[0])
	if err != nil {
		fmt.Println("Error parsing Date 188:", err)
	}

	endDate, err := time.Parse(layout, out.CurveBoundaries.Date[1])
	if err != nil {
		fmt.Println("Error parsing Date 194:", err)
	}

	timeDiff := endDate.Sub(startDate).Hours() / 24
	out.PeriodeDuration.Days = int(timeDiff)
	years := timeDiff / 365
	out.PeriodeDuration.Years = strconv.FormatFloat(years, 'f', 1, 64) // foramt a int64 to a string, round one decimal position

	return out
}

func selectCurveDataInDateRange(curveData []Curve, fromDate string, toDate string) Assets {
	var out Assets
	out.Curves = []Curve{}

	// Parse the input date strings into time.Time objects
	layout := "2006-01-02" // Define the date format (assuming YYYY-MM-DD format)
	var fromTime, toTime time.Time
	var err error

	if fromDate != "" {
		fromTime, err = time.Parse(layout, fromDate)
		if err != nil {
			fmt.Println("Error parsing fromDate:", err)
			fromDate = ""
		}
	}
	if toDate != "" {
		toTime, err = time.Parse(layout, toDate)
		if err != nil {
			fmt.Println("Error parsing toDate:", err)
			toDate = ""
		}
	}

	// Iterate over each curve and its values
	for _, curve := range curveData {
		var filteredValues []DataPoint // Temporary slice to store the filtered values
		for _, ele := range curve.Values {
			// Parse the ele.Date string to time.Time
			eleDate, err := time.Parse(layout, ele.Date)
			if err != nil {
				fmt.Println("Error parsing ele.Date:", err)
				continue
			}

			// Apply filtering based on the available fromDate and toDate
			if (fromDate == "" || eleDate.After(fromTime) || eleDate.Equal(fromTime)) &&
				(toDate == "" || eleDate.Before(toTime) || eleDate.Equal(toTime)) {

				filteredValues = append(filteredValues, ele) // Add the value to the filtered result
			}
		}

		// If any values match, add the filtered curve to the result
		if len(filteredValues) > 0 {
			curve.Values = filteredValues          // Assign the filtered values back to the curve
			out.Curves = append(out.Curves, curve) // Add the curve with filtered values to the final result
		}
	}

	return out
}

func (obj *Curve) enrichWithAdditionalCurveData() {

	isin := obj.Name

	// get inforamtion from EOD metadata file
	var fileContent responseDID
	filename := config.DataIsinPath + isin + "/" + config.EodLastUsedSymbolFileName
	fileContent, _ = readJsonFile(filename, fileContent)

	obj.AdditionalInfo.Exchange = fileContent.ExchangeCode
	obj.AdditionalInfo.Title = fileContent.Name
	obj.AdditionalInfo.SymbolCode = fileContent.AssetCode
	obj.AdditionalInfo.Type = fileContent.Type
	obj.AdditionalInfo.Country = fileContent.Country
	obj.AdditionalInfo.Currency = fileContent.Currency

	diffYears, _ := calculateYearDifference(obj.Values[len(obj.Values)-1].Date, obj.Values[0].Date)
	obj.AdditionalInfo.Duration = diffYears
	obj.calculatePerformenceValues("1m")
	obj.calculatePerformenceValues("6m")
	obj.calculatePerformenceValues("1y")
	obj.calculatePerformenceValues("5y")

	// get info from Custom Additional information
	var fileStructureDefenition customIsinInfoFile
	filename = config.DataIsinPath + isin + "/" + config.CustomIsinInfoFileName
	fileData, err1 := readJsonFile(filename, fileStructureDefenition)
	if err1 == nil {
		obj.AdditionalInfo.Nickname = fileData.Nick
		obj.AdditionalInfo.Description = fileData.Description
		obj.AdditionalInfo.Color = fileData.Color

		var readableTags string
		for _, v := range fileData.Tags {
			readableTags = readableTags + v + ", "
		}

		// remove last comma when tag exist
		end := len(readableTags) - 2
		if end > 1 {
			readableTags = readableTags[0:end]
			obj.AdditionalInfo.Tags = readableTags
		}
	}

}

func (obj *Curve) calculatePerformenceValues(durationCode string) {

	var duration int
	if durationCode == "1m" {
		duration = 30
	} else if durationCode == "6m" {
		duration = 30.5 * 6
	} else if durationCode == "1y" {
		duration = 365
	} else if durationCode == "5y" {
		duration = 365 * 5
	} else {
		fmt.Println("Error, calculate-Performence-Values. not managed code:" + durationCode)
		return
	}

	timeDiffs := make(map[int]int)
	today := time.Now()

	// find out which value has the lowest distance to the time period
	for key, value := range obj.Values {
		diff, _ := calculateYearDifference2(today, value.Date)
		if (diff - duration) >= 0 {
			timeDiffs[key] = diff - duration
		} else {
			timeDiffs[key] = (diff - duration) * -1
		}
		if timeDiffs[key] == 0 {
			break
		}
	}

	lowest := [2]int{0, timeDiffs[0]}
	for key, value := range timeDiffs {
		if lowest[1] > value {
			lowest[0] = key
			lowest[1] = value
		}
	}

	baseValue := obj.Values[lowest[0]].Price
	todayValue := obj.Values[0].Price
	percentagePerf := (todayValue - baseValue) / baseValue * 100

	percentOnYear := percentagePerf * 365 / float64(duration)
	result := fmt.Sprintf("%.1f", percentagePerf) + " (" + fmt.Sprintf("%.1f", percentOnYear) + ")"

	if durationCode == "1m" {
		obj.AdditionalInfo.Perf1m = result
	} else if durationCode == "6m" {
		obj.AdditionalInfo.Perf6m = result
	} else if durationCode == "1y" {
		obj.AdditionalInfo.Perf1y = result
	} else if durationCode == "5y" {
		obj.AdditionalInfo.Perf5y = result
	}

}

func calculateYearDifference(date1, date2 string) (string, error) {
	// Define the date format
	const layout = "2006-01-02"

	// Parse the dates
	d1, err1 := time.Parse(layout, date1)
	d2, err2 := time.Parse(layout, date2)
	if err1 != nil || err2 != nil {
		return "err", fmt.Errorf("error parsing dates: %v, %v", err1, err2)
	}

	// Calculate the total difference in days
	totalDays := d2.Sub(d1).Hours() / 24

	// Convert the difference to years
	differenceInYears := totalDays / 365.25
	result := fmt.Sprintf("%.1f", differenceInYears)
	return result, nil
}

func calculateYearDifference2(date1 time.Time, date2 string) (int, error) {
	// Define the date format
	const layout = "2006-01-02"

	// Parse the dates
	d2, err := time.Parse(layout, date2)
	if err != nil {
		return 0, fmt.Errorf("error parsing dates: %v", err)
	}

	// Calculate the total difference in days
	totalDays := date1.Sub(d2).Hours() / 24

	return int(totalDays), nil
}
