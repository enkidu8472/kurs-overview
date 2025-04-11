package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"kursoverview/internal/services"
	"net/http"
	"os"
)

var AllDataFileNames []string

func ApiPutTRtransaction(w http.ResponseWriter, r *http.Request) {

	var inputBody resquestBodyTransactionLogger
	out := "Transaction log from Trade Republic went fine"

	// Read the request body
	rawRequestbody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error 1")
		return
	}

	// Parse the JSON request into the body
	err = json.Unmarshal(rawRequestbody, &inputBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		fmt.Println("error 2")
		return
	}

	// call the main function of the Service (for now no return, only log to console)
	out = services.LogTrTransaction(inputBody.Clipboard)

	// convert the response into JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readJsonFile[T any](filePath string, data T) (T, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// Read the contents of the file
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return data, fmt.Errorf("failed to read file: %v", err)
	}

	// Unmarshal the JSON into the output struct
	err = json.Unmarshal(fileContents, &data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return data, nil
}

func writeToJsonFile[T any](filePath string, content T) error {

	// Open or create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the struct to JSON and write it to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty formatting (optional)
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
