package cmpscraper

import "net/http"

func New(newCMP *CMP) (CMP, error) {
	client := &http.Client{}
	reqHeaders := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.9999.99 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "keep-alive",
	}

	mwStatsUrl := "https://ecmp.cmpco.com/omni/content/cmpload.txt"
	powerStatsUrl := "https://ecmp.cmpco.com/OutageReports/CMP.html"

	if newCMP.Client != nil {
		client = newCMP.Client
	}

	if newCMP.ReqHeaders != nil {
		for k, v := range newCMP.ReqHeaders {
			reqHeaders[k] = v
		}
	}

	if newCMP.MWStatsUrl != "" {
		mwStatsUrl = newCMP.MWStatsUrl
	}

	if newCMP.PowerStatsUrl != "" {
		powerStatsUrl = newCMP.PowerStatsUrl
	}

	return CMP{
		Client:     client,
		ReqHeaders: reqHeaders,

		MWStatsUrl:    mwStatsUrl,
		PowerStatsUrl: powerStatsUrl,
	}, nil
}
