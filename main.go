package main

import (
	"fmt"
	"kursoverview/internal/handlers"
	"log"
	"net/http"
)

func main() {
	fmt.Println("start-kurs-overview. v0.51")

	// Serve static files (e.g., HTML, CSS, JS) from the "frontend" directory
	http.Handle("/", http.FileServer(http.Dir("./views/html")))

	// all endpoints of the application
	http.HandleFunc("/api/getPossibleAssets", handlers.GetPossibleAssets)
	http.HandleFunc("/api/getAssetData", handlers.GetSelectedAssetInformation)
	http.HandleFunc("/api/downloadFullIsinData", handlers.DownloadIsinData)
	http.HandleFunc("/api/putAdditionalIsinData", handlers.ApiPutAdditionalData)
	http.HandleFunc("/api/putTRtransaction", handlers.ApiPutTRtransaction)

	// Start the webserver
	log.Println("Starting server on http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
