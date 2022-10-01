package cmpscraper

import (
	"io"
	"net/http"
	"regexp"
	"time"
)

const URL = "https://ecmp.cmpco.com/OutageReports/CMP.html"

func GetStats() (CMP, error) {
	var stats CMP
	stats.Counties = make(map[string]Outage)

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return stats, err
	}

	regTotals := regexp.MustCompile("Total</th><th>([0-9,]+)</th><th>([0-9,]+)</th>")
	counties := regexp.MustCompile(`([a-zA-Z]+\.html)'>([a-zA-Z]+)</a>.+?([0-9,]+)</t.+?([0-9,]+)</t`)
	updatedAt := regexp.MustCompile("Update: ([^<]+)")

	resp, err := http.Get(URL)
	if err != nil {
		return stats, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()

	match := regTotals.FindStringSubmatch(string(body))
	stats.Total = match[1]
	stats.WithoutPower = match[2]

	match2 := updatedAt.FindStringSubmatch(string(body))

	stats.LastUpdate, err = time.ParseInLocation("Jan 02, 2006 03:04 PM", match2[1], loc)
	if err != nil {
		return stats, err
	}

	match3 := counties.FindAllStringSubmatch(string(body), -1)

	for _, y := range match3 {
		stats.Counties[y[2]] = Outage{
			Total:        y[3],
			WithoutPower: y[4],
		}
	}

	return stats, nil
}
