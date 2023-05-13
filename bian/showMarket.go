package bian

import (
	"coin-show/cache"
	"coin-show/clients"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
)

// 记录最新成交价格
var (
	StopC       chan struct{}
	Symbols     = map[string]string{}
	NewPriceMap = map[string]MySymbols{}
	showCoin    []string
)

type MySymbols struct {
	NewPrice  string          `json:"new_price"`
	BasePrice string          `json:"base_price"`
	Nums      decimal.Decimal `json:"nums"`
}

func LoadCoin() {
	res, _ := cache.CacheClient.Search("USDT")
	showCoin = []string{}
	for k, v := range res {
		var nm MySymbols

		err := json.Unmarshal([]byte(v), &nm)
		if err != nil {
			log.Println(err)
		}
		if nm.BasePrice == "" {
			nm.BasePrice = "0.0"
		}
		NewPriceMap[k] = nm
		showCoin = append(showCoin, k)
	}
}

func AsyncSymbolPrice() {

	plist := clients.ReqChainLinkByList(showCoin)
	for k, _ := range NewPriceMap {
		old := NewPriceMap[k]
		if old.BasePrice == "" {
			old.BasePrice = plist[k]
		}
		NewPriceMap[k] = MySymbols{
			NewPrice:  plist[k],
			BasePrice: old.BasePrice,
			Nums:      old.Nums,
		}
	}
}

func AddOrUpdateSymbolPrice(key, base, value string) {

	nums, err := decimal.NewFromString(value)
	if err != nil {
		log.Println(err)
		return
	}

	mySymbols := MySymbols{
		NewPrice:  clients.ReqChainLink(key),
		BasePrice: base,
		Nums:      nums,
	}

	addOrUpdateCoin(key, mySymbols)
}

func addOrUpdateCoin(key string, value MySymbols) {

	b, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}
	if _, k := NewPriceMap[key]; !k {
		err := cache.CacheClient.Set(key, string(b))
		if err != nil {
			log.Fatal(err)
		}
		showCoin = append(showCoin, key)
	} else {
		err := cache.CacheClient.Update(key, string(b))
		if err != nil {
			log.Fatal(err)
		}
	}
	NewPriceMap[key] = value
}

func DelCoin(key string) {
	if _, k := NewPriceMap[key]; k {
		err := cache.CacheClient.Delete(key)
		if err != nil {
			log.Fatal(err)
		}
	}
	delete(NewPriceMap, key)
	LoadCoin()
}

func GetLastPrice(symbol string) {

	// 获取最新成交价格
	prices, err := clients.GetBinanceClient().NewListPricesService().Symbol(symbol).Do(context.Background())

	if err != nil {
		if strings.Contains(err.Error(), "net/http: TLS handshake timeout") {
			return
		}
		log.Println(err)
		newPriceMap := NewPriceMap[symbol]
		newPriceMap.NewPrice = "0.0"
		addOrUpdateCoin(symbol, newPriceMap)
		return
	}
	newPriceMap := NewPriceMap[symbol]
	newPriceMap.NewPrice = prices[0].Price
	addOrUpdateCoin(symbol, newPriceMap)
}

func SubCoin(coin string, isDel bool) {
	if !isDel {
		Symbols[coin] = "5s"
	}
	stopSub()

	wsDepthHandler := func(event *binance.WsKlineEvent) {
		newPriceMap := NewPriceMap[event.Symbol]
		newPriceMap.NewPrice = event.Kline.ActiveBuyVolume

		addOrUpdateCoin(event.Symbol, newPriceMap)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, stopC, err := binance.WsCombinedKlineServe(Symbols, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	StopC = stopC
	// remove this if you do not want to be blocked here
	<-doneC
}

func stopSub() {
	// use stopC to exit
	go func() {
		time.Sleep(5 * time.Second)
		StopC <- struct{}{}
	}()
}
