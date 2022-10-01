CMP Scraper
---

This is a simple package used to scrape the Central Maine Power website for outage information.

Usage
---

```go
package main

import (
	"log"

	"github.com/HelixSpiral/cmpscraper"
)

func main() {
	stats, _ := cmpscraper.GetStats()
	log.Printf("%+v", stats)
}
```