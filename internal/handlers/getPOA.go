package handlers

import (
	"encoding/json"
	"fmt"
	"kursoverview/internal/config"
	"net/http"
	"os"
)

// root function of this file: for init or when the user want to select items. give back
// a list of items that are in the "DB"
func GetPossibleAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set the content type to JSON
	fmt.Println("start endpoint get-possible-assets")

	fileIDs, out2 := getAvailableIsins()
	AllDataFileNames = out2

	// Encode the slice as JSON and send it to the frontend
	err := json.NewEncoder(w).Encode(fileIDs)
	if err != nil {
		http.Error(w, "Unable to encode slice", http.StatusInternalServerError)
	}
}

func getAvailableIsins() ([]string, []string) {

	var outputIsin []string
	var outputfileNames []string

	allIsinDirectories, _ := os.ReadDir(config.DataIsinPath)

	for _, isinFoler := range allIsinDirectories {
		isin := isinFoler.Name()
		allFilesInCurrentIsinDirectory, _ := os.ReadDir(config.DataIsinPath + isin)

		for _, file := range allFilesInCurrentIsinDirectory {

			fn := file.Name()
			if fn == config.AlphaTimeSerieFileName {
				filePath := config.DataIsinPath + isin + "/" + fn
				outputfileNames = append(outputfileNames, filePath)
				outputIsin = append(outputIsin, isin)
				break
			}
		}
	}

	return outputIsin, outputfileNames
}
