package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"kursoverview/internal/config"
	"net/http"
	"strings"
)

type resquestBodyPAI struct {
	Isin     string `json:"isin"`
	Property string `json:"property"`
	Value    string `json:"value"`
}

type customIsinInfoFile struct {
	Nick        string   `json:"nick"`
	Color       string   `json:"color"`
	Tags        []string `json:"tags"`
	Description string   `json:"desc"`
}

func ApiPutAdditionalData(w http.ResponseWriter, r *http.Request) {

	var inputBody resquestBodyPAI
	out := "Went fine"

	var fileData customIsinInfoFile

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

	fmt.Printf("endpoint AID called v3 data :: %v\n", inputBody)

	// get current file situation
	filename := config.DataIsinPath + inputBody.Isin + "/" + config.CustomIsinInfoFileName
	fileData, _ = readJsonFile(filename, fileData)

	// update file info acording to API call
	switch inputBody.Property {
	case "Nickname(~):":
		fileData.Nick = inputBody.Value
	case "Color(~):":
		fileData.Color = inputBody.Value
	case "Tags(~):":
		tmp := strings.Split(inputBody.Value, ",")
		for k, v := range tmp {
			v = strings.ReplaceAll(v, " ", "")
			tmp[k] = strings.ToLower(v)
		}
		fileData.Tags = tmp
	case "Description(~):":
		fileData.Description = inputBody.Value
	default:
		fmt.Println("Error 67, ", inputBody.Property)
		http.Error(w, "Unknown key send:"+inputBody.Property, http.StatusInternalServerError)
		return
	}

	// save used info to file
	writeToJsonFile(filename, fileData)

	// convert the response into JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
