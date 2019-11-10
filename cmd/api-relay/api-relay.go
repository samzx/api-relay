package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// GOOS=linux GOARCH=amd64 go build
func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "2001"
	}
	fmt.Println("Running on port " + port)
	http.HandleFunc("/", RelayServer)
	http.ListenAndServe(":"+port, nil)
}

func RelayServer(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodOptions {
		return
	}

	target := r.URL

	relay, err := fetch(target)

	if err == nil {
		fmt.Fprintf(w, "%s", relay)
	} else {
		fmt.Fprintf(w, "{error: %s}", err)
	}
}

func fetch(url *url.URL) (string, error) {
	// Form proper https URL
	urlString := "https://" + url.String()[1:]

	// Fetch
	resp, err := http.Get(urlString)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Return body
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), err
}
