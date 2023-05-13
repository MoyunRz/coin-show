package bian

import (
	"coin-show/clients"
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

// GetAccount
// 获取用户
func GetAccount() []binance.Balance {

	res, err := clients.GetBinanceClient().NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Println(err)
		return nil
	}
	// 过滤
	balances := filterAndSort(res.Balances, 0.0)
	for i := 0; i < len(balances); i++ {
		fmt.Println(balances[i])
	}
	return balances
}

// StartUserStreamService
// 启动用户流
func StartUserStreamService() string {

	res, err := clients.GetBinanceClient().NewStartUserStreamService().Do(context.Background())
	if err != nil {
		log.Println(err)
		return res
	}
	log.Println(res)
	return res
}

func filterAndSort(binances []binance.Balance, free float64) []binance.Balance {
	// 过滤条件:Age >= minAge
	filtered := make([]binance.Balance, 0)
	for i := 0; i < len(binances); i++ {
		bs := binances[i]
		num, _ := strconv.ParseFloat(bs.Free, 64)
		if num > free {
			filtered = append(filtered, bs)
		}
	}
	// 按Name排序
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Free > filtered[j].Free
	})

	return filtered
}
