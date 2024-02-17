package main

import (
	"log"

	"github.com/HelixSpiral/cmpscraper"
)

func main() {
	cmp, err := cmpscraper.New(&cmpscraper.CMP{})
	if err != nil {
		log.Fatal(err)
	}

	load, err := cmp.GetCurrentLoad()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(load)

	stats, err := cmp.GetOutageStats()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", stats)
}
