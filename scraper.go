package cmpscraper

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const URL = "https://ecmp.cmpco.com/OutageReports/CMP.html"

func GetStats(httpClient *http.Client) (CMP, error) {
	var stats CMP
	stats.Counties = make(map[string]Outage)

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return stats, fmt.Errorf("error loading time information: %w", err)
	}

	regTotals := regexp.MustCompile("Total</th><th>([0-9,]+)</th><th>([0-9,]+)</th>")
	counties := regexp.MustCompile(`([a-zA-Z]+\.html)'>([a-zA-Z]+)</a>.+?([0-9,]+)</t.+?([0-9,]+)</t`)
	updatedAt := regexp.MustCompile("Update: ([^<]+)")

	resp, err := httpClient.Get(URL)
	if err != nil {
		return stats, fmt.Errorf("error in http GET: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return stats, fmt.Errorf("error reading http response body: %w", err)
	}
	defer resp.Body.Close()

	if strings.Contains(string(body), "No reported electricity outages are in our system.") {
		stats.NoOutages = true

		return stats, nil
	}

	match := regTotals.FindStringSubmatch(string(body))
	stats.Total = match[1]
	stats.WithoutPower = match[2]

	match2 := updatedAt.FindStringSubmatch(string(body))

	stats.LastUpdate, err = time.ParseInLocation("Jan 02, 2006 03:04 PM", match2[1], loc)
	if err != nil {
		return stats, fmt.Errorf("erorr parsing time location: %w", err)
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
