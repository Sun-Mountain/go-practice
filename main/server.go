package main

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
	"log"
)

func main() {
	// Use the thumbnailHandler function
	http.HandleFunc("/api/thumbnail", thumbnailHandler)

	// Serve static files from the frontend/dist directory
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	// Starts Server
	fmt.Println("Server listening on port 3000")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello world")
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type thumbnailRequest struct {
	Url string `json:"url"`
}

type screenshotAPIRequest struct {
	Token						string `json:"token"`
	Url							string `json:"url"`
	Output					string `json:"output"`
	Width						int 	 `json:"width"`
	Height					int		 `json:"height"`
	ThumbnailWidth	int 	 `json:"thumbnail_width"`
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	var decoded thumbnailRequest

	// Try to decode the request into the thumbnailRequest struct.
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a struct with the parameters needed to call the ScreenshotAPI.
	apiRequest := screenshotAPIRequest{
		Token:					"SWSZDNR414URKK5BXI7ES4TKGU8ESQC6",
		Url: 						decoded.Url,
		Output:					"json",
		Width:					1920,
		Height:					1080,
		ThumbnailWidth:	300,
	}

	// Convert the struct to a JSON string.
	jsonString, err := json.Marshal(apiRequest)
	checkError(err)

	// Create a HTTP request.
	req, err := http.NewRequest("POST", "https://screenshotapi.net/api/v1/screenshot", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	checkError(err)

	// Tell Go to close the response at the end of the function.
	defer response.Body.Close()

	// Read the raw response into a Go struct.
	type screenshotAPIResponse struct {
		Screenshot string `json:"screenshot"`
	}
	var apiResponse screenshotAPIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	checkError(err)

	// Pass back the screenshot URL to frontend
	_, err = fmt.Fprintf(w, `{ "screenshot": "%s" |`, apiResponse.Screenshot)
	checkError(err)
}

