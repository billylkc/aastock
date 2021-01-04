package aastock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// GetList gets a list of stocks that we are interested in
func GetList() []string {
	var result []string
	c := []string{"1", "5", "9988", "700"}

	for _, cc := range c {
		result = append(result, fmt.Sprintf("%05s", cc))
	}
	return result
}

// StockPrice as the response body from yahoo finance chart
type StockPrice struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PreviousClose        float64 `json:"previousClose"`
				Scale                int     `json:"scale"`
				PriceHint            int     `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"pre"`
					Regular struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"regular"`
					Post struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				TradingPeriods [][]struct {
					Timezone  string `json:"timezone"`
					Start     int    `json:"start"`
					End       int    `json:"end"`
					Gmtoffset int    `json:"gmtoffset"`
				} `json:"tradingPeriods"`
				DataGranularity string   `json:"dataGranularity"`
				Range           string   `json:"range"`
				ValidRanges     []string `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Low    []float64     `json:"low"`
					High   []float64     `json:"high"`
					Close  []float64     `json:"close"`
					Volume []interface{} `json:"volume"`
					Open   []float64     `json:"open"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

// GetCurrentPrice gets the current price of the stock from Yahoo finance
func GetCurrentPrice(c int) (StockPrice, error) {
	var stock StockPrice

	code := fmt.Sprintf("%04d", c)
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s.HK?region=US&lang=en-US&includePrePost=false&interval=5m&range=1d&corsDomain=finance.yahoo.com&.tsrc=finance", code)
	payload := strings.NewReader("")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, payload)

	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Add("Accept", "")
	req.Header.Add("Referer", fmt.Sprintf("https://finance.yahoo.com/quote/%s.HK/", code))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/87.0.4280.88 Safari/537.36")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return stock, errors.Wrap(err, "something's wrong with ReadAll")
	}

	err = json.Unmarshal(body, &stock)
	if err != nil {
		return stock, errors.Wrap(err, "can not marshell response")
	}

	// check last price - close
	prices := stock.Chart.Result[0].Indicators.Quote[0].Close
	fmt.Println(getLastValue(prices))

	return stock, nil
}

// getLastValue gets the last non zero value from a list
func getLastValue(values []float64) float64 {
	var last float64
	for _, v := range values {
		if v > 0 {
			last = v
		}
	}
	return last
}
