package main

import (
	"colly/rarbg"
	"colly/util"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"log"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// 代理
	pf, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:10808")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(pf)

	// 拦截
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69")
		r.Headers.Set("Cookie", rarbg.Cookie)
		fmt.Println("Visiting:", r.URL)
	})

	// 保存文件
	c.OnResponse(util.Save)

	// 详情
	c.OnHTML("body", rarbg.ParseInfo)

	// 列表
	c.OnHTML("body", rarbg.ParseList)

	// 入口
	if err := c.Visit(rarbg.URL); err != nil {
		log.Fatal(err)
	}
}
