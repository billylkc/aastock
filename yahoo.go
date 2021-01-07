package aastock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Stock returns the price and timestamp of a stock
type Stock struct {
	Code        string
	MarketTime  string
	MarketPrice float64
}

// YHStockPrice as the response body from yahoo finance chart
type YHStockPrice struct {
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

// GetCurrentPrices gets the current prices of a provided list of stocks
func GetCurrentPrices(data ...int) []Stock {
	numWorkers := 5

	done := make(chan bool)
	defer close(done)

	// Send data
	in := send(data...)

	// Start workers to process the data
	workers := make([]<-chan Stock, numWorkers)
	for i := 0; i < len(workers); i++ {
		workers[i] = process(done, in)
	}

	// Merge all channels, and sort
	var result []Stock

	for n := range merge(done, workers...) {
		if n.MarketPrice != 0 {
			result = append(result, n)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Code < result[j].Code
	})

	return result
}

// GetPrice is a wrapper of getCurrentPrices to be used in concurrent call
// which would skip any error response
func getPrice(c int) Stock {
	result, _ := getCurrentPrices(c)
	return result
}

// getCurrentPrices gets the delayed price of the stock from Yahoo finance (15mins delay)
// Return only price within 30 mins
func getCurrentPrices(c int) (Stock, error) {
	var (
		stock   Stock        // Market price struct
		yhStock YHStockPrice // Original Yahoo Response
	)

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

	err = json.Unmarshal(body, &yhStock)
	if err != nil {
		return stock, errors.Wrap(err, "can not marshell response")
	}

	// Return stock object if market price is greater than 0
	r := yhStock.Chart.Result
	if len(r) != 0 {
		price := r[0].Meta.RegularMarketPrice
		rt := int64(r[0].Meta.RegularMarketTime) // Market time

		// Compare market time and current time
		tc := time.Now()
		tm := time.Unix(rt, 0)
		diff := tc.Sub(tm)
		threshold, _ := time.ParseDuration("30m")

		// Only return records within 30 mins of the call
		if diff <= threshold {
			if price > 0 {
				loc, _ := time.LoadLocation("Local")
				marketTime := tm.In(loc).Format("2006-01-02 15:04:05")
				stock = Stock{Code: code, MarketTime: marketTime, MarketPrice: price}
			}
		}
	}

	return stock, nil
}

func send(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()

	return out
}

func process(done <-chan bool, in <-chan int) <-chan Stock {
	out := make(chan Stock)
	go func() {
		defer close(out)

		for n := range in {
			select {
			case out <- getPrice(n):
			case <-done:
				return
			}
		}

	}()
	return out
}

func merge(done <-chan bool, cs ...<-chan Stock) <-chan Stock {
	var wg sync.WaitGroup
	out := make(chan Stock)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan Stock) {
		for n := range c {
			select {
			case out <- n:
			case <-done:
			}
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
