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

func (a *YahooAPI) PriceHistorical(symbol string, t time.Time) (float32, error) {

	type chart struct {
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
					DataGranularity string   `json:"dataGranularity"`
					Range           string   `json:"range"`
					ValidRanges     []string `json:"validRanges"`
				} `json:"meta"`
				Timestamp  []int `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						High   []float64 `json:"high"`
						Volume []float64 `json:"volume"`
						Open   []float64 `json:"open"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
					} `json:"quote"`
					Adjclose []struct {
						Adjclose []float64 `json:"adjclose"`
					} `json:"adjclose"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"chart"`
	}

	resp, err := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?symbol=%s&period1=%d&period2=%d&interval=1d", symbol, symbol, t.Unix(), t.Add(24*time.Hour).Unix()))
	if err != nil {
		return 0, fmt.Errorf("can't get price at %s for %s: %w", t.String(), symbol, err)
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("can't get price at %s for %s: %w", t.String(), symbol, err)
	}

	c := new(chart)
	if err := json.Unmarshal(bytes, c); err != nil {
		return 0, fmt.Errorf("can't get price at %s for %s: %w", t.String(), symbol, err)
	}

	price := c.Chart.Result[0].Indicators.Quote[0].Open[0]

	return float32(price), nil
}
