package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func NewPriceCollector() *PriceCollector {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	return &PriceCollector{
		Config: cfg,
	}
}

type PriceCollector struct {
	Config *Config
}

func (p *PriceCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, coin := range p.Config.Coins {
		coin.MetricDesc = prometheus.NewDesc(
			prometheus.BuildFQName("pricex", "", coin.Symbol+"_"+p.Config.Currencies),
			"coin price of "+coin.Name, []string{}, nil,
		)
		log.Infof("metric description for \"%s\" registerd", coin.Symbol)
	}
}

func (p *PriceCollector) Collect(ch chan<- prometheus.Metric) {

	if len(p.Config.Coins) == 0 {
		return
	}

	ids := p.Config.Coins[0].ID
	for i := 1; i < len(p.Config.Coins); i++ {
		ids += "," + p.Config.Coins[i].ID
	}

	resp, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", ids, p.Config.Currencies))
	if err != nil {
		log.Errorf("failed to fetch coin price: %v", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to parse coin gecko response: %v", err)
		return
	}

	var coinPrices map[string]map[string]float64
	err = json.Unmarshal(b, &coinPrices)
	if err != nil {
		log.Errorf("failed to unmarshal coin gecko response: %v", err)
		return
	}

	for _, coin := range p.Config.Coins {
		coinPrice, isExists := coinPrices[coin.ID]
		if isExists {
			for _, price := range coinPrice {
				ch <- prometheus.MustNewConstMetric(coin.MetricDesc, prometheus.GaugeValue, price)
			}
		}
	}
}
