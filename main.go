package main

import (
	"coin-show/bian"
	"coin-show/config"
	"coin-show/frames"
)

func main() {

	config.InitConfig()
	go func() {
		// 缓存币币种
		bian.LoadCoin()
		frames.LoadContainer()
		// 异步加载
		frames.AsyncPrice()
	}()
	frames.RunFrame()
}
