package main

import (
	"log"
	"net/http"

	"github.com/HelixSpiral/cmpscraper"
)

func main() {
	httpClient := &http.Client{}
	stats, _ := cmpscraper.GetStats(httpClient)
	log.Printf("%+v", stats)
}
