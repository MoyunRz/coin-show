package cache

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

func DownloadFlatIconImages(keyword string) {
	mkdirSrc(".coin_show/icons")
	filePath := fmt.Sprintf(".coin_show/icons/%s.png", keyword)
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		return
	}

	// 构造client
	client := &http.Client{}
	// 构造Flaticon搜索URL
	searchURL := "https://www.flaticon.com/search?word=" + keyword
	// 获取搜索结果页内容
	response, err := client.Get(searchURL)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 获取li标签的正则表达式
	imgPattern := regexp.MustCompile(`https://cdn-icons-png.flaticon.com/.*?.png`)
	// 解析页面,获取li标签
	lis := imgPattern.FindAllStringSubmatch(string(body), -1)
	if len(lis) >= 3 {
		lis = lis[2:3]
	}
	// 遍历li标签
	for _, li := range lis {
		// 获取li标签内文本
		// 构造图片URL
		imgURL := li[0]
		// 下载图片
		response, err := http.Get(imgURL)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(filePath) //创建文件
		if err != nil {
			panic(err)
		}
		io.Copy(f, response.Body) //写入文件
	}
	time.Sleep(time.Second)
}

func WImg(imgURL, keyword string) {

	mkdirSrc(".coin_show/icons")
	filePath := fmt.Sprintf(".coin_show/icons/%s.png", keyword)
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		return
	}
	// 下载图片
	response, err := http.Get(imgURL)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(filePath) //创建文件
	if err != nil {
		panic(err)
	}
	io.Copy(f, response.Body) //写入文件
}
