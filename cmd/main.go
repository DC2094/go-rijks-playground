package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rijks/internal/ingest"
)

const BaseURL = "https://www.rijksmuseum.nl/api/oai/"

func main() {
	fmt.Println("Hello World")

	apiKey := os.Getenv("RIJKS_API_KEY")

	client := &http.Client{}
	rh := ingest.NewRijksHandler(apiKey, BaseURL, client)
	record, err := rh.GetRecord("oai:rijksmuseum.nl:sk-c-5")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record: ", record.GetRecord)

	rh.ListRecords()
}
