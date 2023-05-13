# 币种价格变化监控

这是一个用golang写的桌面虚拟货币价格监控应用 


## 使用的框架
Fyne v2: <https://github.com/fyne-io/fyne>

go版本：1.20

## 运行效果
币种名称｜价格｜收益

![img.png](doc%2Fimgs%2Fimg.png)

## 查询接口

币安API (已经不用): <https://data.binance.com/api/v3/ticker/price>
chainLink(现在版本使用): <https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USDT>


## 打包
>  fyne package -os windows -icon img.png

os 后的windows 是系统：
> android, android/arm, android/arm64, android/amd64, android/386, ios, iossimulator, wasm, gopherjs, web
