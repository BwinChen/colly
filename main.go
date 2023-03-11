package main

import (
	"colly/sukebei"
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
		fmt.Println("Visiting:", r.URL)
	})

	// 保存文件
	c.OnResponse(sukebei.Save)

	// 详情
	c.OnHTML("body", sukebei.ParseInfo)

	// 列表
	c.OnHTML("body", sukebei.ParseList)

	// 入口
	if err := c.Visit("https://sukebei.nyaa.si/?p=1"); err != nil {
		log.Fatal(err)
	}
}
