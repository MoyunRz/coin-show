package bian

import (
	"coin-show/config"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type PriceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func GetSymbol() []PriceResponse {
	var priceList []PriceResponse
	binance_url := "https://data.binance.com/api/v3/ticker/price"
	// 创建一个代理 URL
	proxyUrl, err := url.Parse(config.Cfg.Vpn.Url)
	if err != nil {
		log.Println(err)
		return priceList
	}

	// 创建一个自定义的 Transport 对象，并设置 Proxy 字段
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	resp, err := client.Get(binance_url)
	if err != nil {
		log.Println(err)
		return priceList
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return priceList
	}
	err = json.Unmarshal(body, &priceList)
	if err != nil {
		log.Println(err)
		return priceList
	}

	// for _, price := range priceList {
	// 	if price.Symbol == "BTCUSDT" {
	// 		fmt.Printf("当前 BTC/USDT 市场价格为：%.6f USDT\n", price.Price)
	// 	}
	// }

	return priceList
}
