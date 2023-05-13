package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ReqChainLink(fsyms string) string {
	fsyms = strings.Split(fsyms, "USDT")[0]
	// 发送 HTTP GET 请求
	resp, err := http.Get(fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USDT", fsyms))
	if err != nil {
		log.Println("Error:", err)
		return "0.0"
	}
	defer resp.Body.Close()

	// 解析 JSON 响应
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println("Error:", err)
		return "0.0"
	}
	if data[fsyms] == nil {
		return "0.0"
	}
	// 输出美元价格
	usd := data[fsyms].(map[string]interface{})["USDT"]
	if usd == nil {
		return "0.0"
	}
	log.Printf("%s 美元价格：%f\n", fsyms, usd)
	return fmt.Sprintf("%f", usd)
}

func ReqChainLinkByList(fsyms []string) map[string]string {
	fsymsMap := map[string]string{}
	// 将数字转成字符串切片
	strSlice := make([]string, len(fsyms))

	for i, v := range fsyms {
		strSlice[i] = fmt.Sprintf("%s", strings.Split(v, "USDT")[0])
		fsymsMap[v] = "0.0"
	}
	// 将字符串切片拼接为一个字符串
	result := strings.Join(strSlice, ",")
	// 发送 HTTP GET 请求
	resp, err := http.Get(fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USDT", result))
	if err != nil {
		log.Println("Error:", err)
		return fsymsMap
	}
	defer resp.Body.Close()

	// 解析 JSON 响应
	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println("Error:", err)
		return fsymsMap
	}

	for j := 0; j < len(strSlice); j++ {
		// 输出美元价格
		if data == nil || data[strSlice[j]] == nil {
			continue
		}

		usd := data[strSlice[j]].(map[string]interface{})["USDT"]
		log.Printf("%s 美元价格：%f\n", strSlice[j], usd)
		fsymsMap[strSlice[j]+"USDT"] = fmt.Sprintf("%f", usd)
	}
	return fsymsMap
}
