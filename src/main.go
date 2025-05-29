package main

import (
	"github.com/BwinChen/colly/src/bt4g"
	"github.com/BwinChen/colly/src/util"
	"github.com/gocolly/colly/v2"
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
		r.Headers.Set("Cookie", bt4g.Cookie)
		//log.Println("Visiting:", r.URL)
	})

	// 限制速率
	bt4g.Limit(c)

	// 处理错误
	c.OnError(bt4g.ErrorHandler)

	// 处理响应
	c.OnResponse(bt4g.Save)

	// 解析HTML
	c.OnHTML("body", bt4g.ParseHTML)

	// 开始爬取
	bt4g.Visit(c)
}
