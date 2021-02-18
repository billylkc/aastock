package aastock

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type company struct {
	Code string
	Name string
}

// GetCompanyName looks up company name from HKEX
func GetCompanyName(c int) (company, error) {
	var result company

	// Handle input, e.g. code = 00005, date 2021-02-01
	targetCode := fmt.Sprintf("%05d", c) // zfill to 5 digit
	currentTime := time.Now()
	d := currentTime.Format("2006-01-02")
	d = strings.ReplaceAll(d, "-", "") // date in string format

	url := fmt.Sprintf("https://www.hkexnews.hk/sdw/search/stocklist_c.aspx?sortby=stockcode&shareholdingdate=%s", d)
	res, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return result, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return result, err
	}

	// Find the review items
	doc.Find("table.table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title

		content := s.Find("td").Text()
		regex := *regexp.MustCompile(`\s*(\d{5})\s*(.*)`)
		matched := regex.FindAllStringSubmatch(content, -1)
		for i := range matched {
			codeStr := matched[i][1]
			companyStr := matched[i][2]

			if codeStr == targetCode {
				result = company{
					Code: targetCode,
					Name: companyStr,
				}
				break // find then break
			}
		}
	})
	return result, nil
}
