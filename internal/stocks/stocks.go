package stocks

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"speedy-cli/internal/common"
)

type Stock struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	ChangePct float64   `json:"changePct"`
	Spark     string    `json:"spark"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type apiItem struct {
	Symbol string  `json:"symbol"`
	LTP    float64 `json:"ltp"`
	Point  float64 `json:"pointChange"`
	Pct    float64 `json:"percentageChange"`
}

func Get(symbol string, topN int) (common.Result, []Stock) {
	stocks, ok := fromAPI()
	result := common.Result{Status: common.StatusSuccess, Message: "live API"}
	if !ok {
		stocks = fallbackWithJitter()
		result.Status = common.StatusWarning
		result.Message = "fallback demo data"
		result.Suggestion = "Live NEPSE endpoint blocked/unreachable right now; showing realistic simulated movers"
	}

	if symbol != "" {
		symbol = strings.ToUpper(symbol)
		filtered := make([]Stock, 0, 1)
		for _, s := range stocks {
			if s.Symbol == symbol {
				filtered = append(filtered, s)
				break
			}
		}
		if len(filtered) == 0 {
			return common.Result{Status: common.StatusError, Message: "symbol not found", Suggestion: "Try a valid NEPSE symbol or remove --symbol"}, []Stock{}
		}
		return result, filtered
	}

	sort.Slice(stocks, func(i, j int) bool {
		return abs(stocks[i].ChangePct) > abs(stocks[j].ChangePct)
	})
	if topN <= 0 || topN > len(stocks) {
		topN = len(stocks)
	}
	return result, stocks[:topN]
}

func fromAPI() ([]Stock, bool) {
	client := &http.Client{Timeout: 7 * time.Second}
	// Endpoint can change; this is a common public mirror.
	req, _ := http.NewRequest(http.MethodGet, "https://nepseapi.surajr.com.np/api/live", nil)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return nil, false
	}
	defer resp.Body.Close()

	var raw []apiItem
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil || len(raw) == 0 {
		return nil, false
	}
	out := make([]Stock, 0, len(raw))
	for _, r := range raw {
		pct := r.Pct
		if pct == 0 && r.LTP != 0 {
			pct = (r.Point / r.LTP) * 100
		}
		out = append(out, Stock{Symbol: strings.ToUpper(r.Symbol), Price: r.LTP, ChangePct: pct, Spark: spark(pct), UpdatedAt: time.Now()})
	}
	return out, true
}

func fallbackWithJitter() []Stock {
	now := time.Now()
	base := []Stock{
		{Symbol: "NABIL", Price: 512.50, ChangePct: 1.2, Spark: "▁▂▄▆█", UpdatedAt: now},
		{Symbol: "NBL", Price: 350.00, ChangePct: -0.8, Spark: "█▆▄▃▁", UpdatedAt: now},
		{Symbol: "SBI", Price: 410.00, ChangePct: 2.1, Spark: "▁▃▅▇█", UpdatedAt: now},
		{Symbol: "NLIC", Price: 280.00, ChangePct: 0.6, Spark: "▂▃▄▅▆", UpdatedAt: now},
		{Symbol: "EBL", Price: 450.00, ChangePct: -1.0, Spark: "█▇▅▃▁", UpdatedAt: now},
		{Symbol: "SHIVM", Price: 468.30, ChangePct: 1.4, Spark: "▁▂▅▆█", UpdatedAt: now},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range base {
		priceJitter := (r.Float64() - 0.5) * 8.0
		pctJitter := (r.Float64() - 0.5) * 0.7
		base[i].Price = round2(base[i].Price + priceJitter)
		base[i].ChangePct = round2(base[i].ChangePct + pctJitter)
		base[i].Spark = spark(base[i].ChangePct)
		base[i].UpdatedAt = now
	}
	return base
}

func spark(pct float64) string {
	if pct >= 1.5 {
		return "▁▃▅▇█"
	}
	if pct > 0 {
		return "▂▃▄▅▆"
	}
	if pct <= -1.5 {
		return "█▇▅▃▁"
	}
	return "▆▅▄▃▂"
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
