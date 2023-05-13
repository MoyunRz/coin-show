package clients

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
)

var (
	binanceClient *ClientInfo

	// 测试环境
	apiKey    = "2B2lRdERmYmo99V40SbI169FkXAO3XyQhvipnLmW4YgiBZESB7kPgIg6ORNzdeKW"
	secretKey = "GgaNI3ym2dqg5bwsbrmMRngn0epHvvruzKKmOOZ1Tzk0rQw61mUd0sqI2pyEscGd"
)

type ApiKey struct {
	apiKey    string
	secretKey string
	isTest    bool
}

type ClientInfo struct {
	Client         *binance.Client
	FuturesClient  *futures.Client
	DeliveryClient *delivery.Client
}

func getAPIByEnv(env string) ApiKey {
	if env == "main" {
		return ApiKey{
			apiKey:    "",
			secretKey: "",
			isTest:    false,
		}
	}

	if env == "test" {
		return ApiKey{
			apiKey:    "2B2lRdERmYmo99V40SbI169FkXAO3XyQhvipnLmW4YgiBZESB7kPgIg6ORNzdeKW",
			secretKey: "GgaNI3ym2dqg5bwsbrmMRngn0epHvvruzKKmOOZ1Tzk0rQw61mUd0sqI2pyEscGd",
			isTest:    true,
		}
	}
	return ApiKey{}
}

func InitBinanceClient(env string) {

	if binanceClient == nil || binanceClient.Client == nil {
		apis := getAPIByEnv(env)
		binance.UseTestnet = apis.isTest
		binanceClient = &ClientInfo{
			Client:         binance.NewClient(apis.apiKey, apis.secretKey),
			FuturesClient:  binance.NewFuturesClient(apis.apiKey, apis.secretKey),
			DeliveryClient: binance.NewDeliveryClient(apis.apiKey, apis.secretKey),
		}
	}
	err := binanceClient.Client.NewPingService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func GetBinanceClient() *binance.Client {

	return binanceClient.Client
}

func GetFuturesClient() *futures.Client {
	return binanceClient.FuturesClient
}

func GetDeliveryClient() *delivery.Client {
	return binanceClient.DeliveryClient
}
