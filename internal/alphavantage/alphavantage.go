package alphavantage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const CACHE_VALIDITY_PERIOD = time.Millisecond * 10 * 60 * 1000

type cache struct {
	price   float32
	updated time.Time
}

type AlphaVantage struct {
	key   string
	m     *sync.RWMutex
	cache map[string]cache
}

type priceInfo struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Open             string `json:"02. open"`
		High             string `json:"03. high"`
		Low              string `json:"04. low"`
		Price            string `json:"05. price"`
		Volume           string `json:"06. volume"`
		LatestTradingDay string `json:"07. latest trading day"`
		PreviousClose    string `json:"08. previous close"`
		Change           string `json:"09. change"`
		ChangePercent    string `json:"10. change percent"`
	} `json:"Global Quote"`
}

func NewAlphaVantage(key string) *AlphaVantage {
	return &AlphaVantage{key: key, m: &sync.RWMutex{}, cache: map[string]cache{}}
}

func (a *AlphaVantage) Price(symbol string) (float32, error) {
	a.m.RLock()
	if v, ok := a.cache[symbol]; ok && time.Since(v.updated) < CACHE_VALIDITY_PERIOD {
		defer a.m.RUnlock()
		return v.price, nil
	}
	a.m.RUnlock()

	a.m.Lock()
	defer a.m.Unlock()

	resp, err := http.Get(fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, a.key))
	if err != nil {
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}
	// log.Println(string(bytes))

	priceInfo := new(priceInfo)
	if err := json.Unmarshal(bytes, priceInfo); err != nil {
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}

	price, err := strconv.ParseFloat(priceInfo.GlobalQuote.Price, 32)
	if err != nil {
		if v, ok := a.cache[symbol]; ok {
			log.Printf("return expired cache value for %s", symbol)
			return v.price, nil
		}
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}

	a.cache[symbol] = cache{price: float32(price), updated: time.Now()}

	return float32(price), nil
}
