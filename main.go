package main

import (
	"colly/sukebei"
	"colly/util"
	"fmt"
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
		r.Headers.Set("Cookie", sukebei.Cookie)
		fmt.Println("Visiting:", r.URL)
	})

	// 下载
	c.OnResponse(util.Save)

	// 详情
	c.OnHTML("body", sukebei.ParseInfo)

	// 列表
	c.OnHTML("body", sukebei.ParseList)

	// 入口
	err = c.Visit(sukebei.URL)
	if err != nil {
		log.Fatal(err)
	}
}
