CMP Scraper
---

This is a simple package used to scrape the Central Maine Power website for outage information.

Usage
---

```go
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
```