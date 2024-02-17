package cmpscraper

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func (cmp *CMP) rawReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("erorr creating request: %w", err)
	}

	// Set any headers provided
	for k, v := range cmp.ReqHeaders {
		req.Header.Set(k, v)
	}

	resp, err := cmp.Client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error in http GET: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading http response body: %w", err)
	}
	defer resp.Body.Close()

	return body, nil
}

func (cmp *CMP) GetCurrentLoad() (string, error) {
	body, err := cmp.rawReq(cmp.MWStatsUrl)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (cmp *CMP) GetOutageStats() (CMPPowerStats, error) {
	var stats CMPPowerStats
	stats.Counties = make(map[string]Outage)

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return stats, fmt.Errorf("error loading time information: %w", err)
	}

	regTotals := regexp.MustCompile("Total</th><th>([0-9,]+)</th><th>([0-9,]+)</th>")
	counties := regexp.MustCompile(`([a-zA-Z]+\.html)'>([a-zA-Z]+)</a>.+?([0-9,]+)</t.+?([0-9,]+)</t`)
	updatedAt := regexp.MustCompile("Update: ([^<]+)")

	body, err := cmp.rawReq(cmp.PowerStatsUrl)
	if err != nil {
		return stats, err
	}

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
