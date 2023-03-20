package main

import (
	"colly/btsow"
	"colly/util"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"log"
)

func main() {
	c := colly.NewCollector()

	// 代理
	pf, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:10808")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(pf)

	// 拦截
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", util.UserAgent)
		r.Headers.Set("Cookie", btsow.Cookie)
		log.Println("Visiting:", r.URL)
	})

	// 下载
	c.OnResponse(util.Save)

	// 详情
	c.OnHTML("body", btsow.ParseInfo)

	// 列表
	c.OnHTML("body", btsow.ParseList)

	// 入口
	err = c.Visit(btsow.URL)
	if err != nil {
		log.Fatal(err)
	}
}
