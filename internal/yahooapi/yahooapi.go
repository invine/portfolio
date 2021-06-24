package yahooapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const CACHE_VALIDITY_PERIOD = time.Millisecond * 10 * 60 * 1000

type cache struct {
	price   float32
	updated time.Time
}

type YahooAPI struct {
	m     *sync.RWMutex
	cache map[string]cache
}

type yahooQuote struct {
	QuoteSummary struct {
		Result []struct {
			Price struct {
				MaxAge          int `json:"maxAge"`
				PreMarketChange struct {
				} `json:"preMarketChange"`
				PreMarketPrice struct {
				} `json:"preMarketPrice"`
				PreMarketSource         string `json:"preMarketSource"`
				PostMarketChangePercent struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"postMarketChangePercent"`
				PostMarketChange struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"postMarketChange"`
				PostMarketTime  int `json:"postMarketTime"`
				PostMarketPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"postMarketPrice"`
				PostMarketSource           string `json:"postMarketSource"`
				RegularMarketChangePercent struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketChangePercent"`
				RegularMarketChange struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketChange"`
				RegularMarketTime int `json:"regularMarketTime"`
				PriceHint         struct {
					Raw     int    `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"priceHint"`
				RegularMarketPrice struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketPrice"`
				RegularMarketDayHigh struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketDayHigh"`
				RegularMarketDayLow struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketDayLow"`
				RegularMarketVolume struct {
					Raw     int    `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"regularMarketVolume"`
				AverageDailyVolume10Day struct {
				} `json:"averageDailyVolume10Day"`
				AverageDailyVolume3Month struct {
				} `json:"averageDailyVolume3Month"`
				RegularMarketPreviousClose struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketPreviousClose"`
				RegularMarketSource string `json:"regularMarketSource"`
				RegularMarketOpen   struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"regularMarketOpen"`
				StrikePrice struct {
				} `json:"strikePrice"`
				OpenInterest struct {
				} `json:"openInterest"`
				Exchange              string      `json:"exchange"`
				ExchangeName          string      `json:"exchangeName"`
				ExchangeDataDelayedBy int         `json:"exchangeDataDelayedBy"`
				MarketState           string      `json:"marketState"`
				QuoteType             string      `json:"quoteType"`
				Symbol                string      `json:"symbol"`
				UnderlyingSymbol      interface{} `json:"underlyingSymbol"`
				ShortName             string      `json:"shortName"`
				LongName              string      `json:"longName"`
				Currency              string      `json:"currency"`
				QuoteSourceName       string      `json:"quoteSourceName"`
				CurrencySymbol        string      `json:"currencySymbol"`
				FromCurrency          interface{} `json:"fromCurrency"`
				ToCurrency            interface{} `json:"toCurrency"`
				LastMarket            interface{} `json:"lastMarket"`
				Volume24Hr            struct {
				} `json:"volume24Hr"`
				VolumeAllCurrencies struct {
				} `json:"volumeAllCurrencies"`
				CirculatingSupply struct {
				} `json:"circulatingSupply"`
				MarketCap struct {
					Raw     int64  `json:"raw"`
					Fmt     string `json:"fmt"`
					LongFmt string `json:"longFmt"`
				} `json:"marketCap"`
			} `json:"price"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteSummary"`
}

func NewYahooAPI() *YahooAPI {
	return &YahooAPI{m: &sync.RWMutex{}, cache: map[string]cache{}}
}

func (a *YahooAPI) Price(symbol string) (float32, error) {
	a.m.RLock()
	if v, ok := a.cache[symbol]; ok && time.Since(v.updated) < CACHE_VALIDITY_PERIOD {
		defer a.m.RUnlock()
		return v.price, nil
	}
	a.m.RUnlock()

	a.m.Lock()
	defer a.m.Unlock()

	resp, err := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=price", symbol))
	if err != nil {
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}
	// log.Println(string(bytes))

	quote := new(yahooQuote)
	if err := json.Unmarshal(bytes, quote); err != nil {
		return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	}

	price := quote.QuoteSummary.Result[0].Price.RegularMarketPrice.Raw
	// if err != nil {
	// 	if v, ok := a.cache[symbol]; ok {
	// 		log.Printf("return expired cache value for %s", symbol)
	// 		return v.price, nil
	// 	}
	// 	return 0, fmt.Errorf("can't get price for %s: %w", symbol, err)
	// }

	a.cache[symbol] = cache{price: float32(price), updated: time.Now()}
	// log.Println(a.cache)

	return float32(price), nil
}
