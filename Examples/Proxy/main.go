package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/HelixSpiral/cmpscraper"
	"golang.org/x/net/proxy"
)

func main() {
	// This needs to be changed out for a real SOCKS5 proxy in a region that CMP accepts
	// connections from.
	proxyDial, err := proxy.SOCKS5("tcp", "0.0.0.0", nil, proxy.Direct)
	if err != nil {
		log.Fatalln("Cannot connect to proxy:", err)
	}

	httpTransport := &http.Transport{
		Dial: proxyDial.Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	cmp, err := cmpscraper.New(&cmpscraper.CMP{
		Client: httpClient,
	})
	if err != nil {
		log.Fatal(err)
	}

	load, err := cmp.GetCurrentLoad()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(load)

	stats, _ := cmp.GetOutageStats()
	log.Printf("%+v", stats)
}
