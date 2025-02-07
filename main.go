package main

import (
	"github.com/BwinChen/colly/src/sukebei"
	"github.com/BwinChen/colly/src/util"
	"github.com/gocolly/colly/v2"
	"log"
	"time"
)

func main() {
	c := colly.NewCollector()

	// 设置代理
	//proxyFunc, err := proxy.RoundRobinProxySwitcher("socks5://192.168.0.4:1070")
	//if err != nil {
	//	log.Fatalf("RoundRobinProxySwitcher Error: %v", err)
	//}
	//c.SetProxyFunc(proxyFunc)

	// 设置请求超时
	c.SetRequestTimeout(10000 * time.Millisecond)

	// 设置Request Header
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", util.UserAgent)
		r.Headers.Set("Cookie", sukebei.Cookie)
		log.Println("Visiting:", r.URL)
	})

	// 限制速率
	sukebei.Limit(c)

	// 处理错误
	c.OnError(sukebei.ErrorHandler)

	// 处理响应
	c.OnResponse(sukebei.Save)

	// 解析HTML
	c.OnHTML("body", sukebei.ParseHTML)

	// 开始爬取
	sukebei.Visit(c)
}
