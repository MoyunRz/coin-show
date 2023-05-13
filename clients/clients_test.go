package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type PriceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func GetSymbol() {
	binance_url := "https://data.binance.com/api/v3/ticker/price"
	// 创建一个代理 URL
	proxyUrl, err := url.Parse("http://127.0.0.1:7890")
	if err != nil {
		panic(err)
	}

	// 创建一个自定义的 Transport 对象，并设置 Proxy 字段
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	resp, err := client.Get(binance_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var priceList []PriceResponse
	err = json.Unmarshal(body, &priceList)
	if err != nil {
		panic(err)
	}
	for _, price := range priceList {
		if price.Symbol == "BTCUSDT" {
			fmt.Printf("当前 BTC/USDT 市场价格为：%.2f USDT\n", price.Price)
		}
	}
}

func Test_InitBinanceClient(t *testing.T) {
	//GetSymbol()
	// https://data.binance.com/api/v3/ticker/price
	//ReqChainLink("BTC")
	f := 1.159000
	w := 18                        // int 类型的宽度
	s := fmt.Sprintf("%-*f", w, f) // 正确!
	fmt.Println(s)
}
