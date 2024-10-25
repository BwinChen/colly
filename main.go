package main

import (
	"github.com/BwinChen/colly/sukebei"
	"github.com/BwinChen/colly/util"
	"github.com/gocolly/colly/v2"
	"log"
	"time"
)

func main() {
	c := colly.NewCollector()

	// 代理
	//proxyFunc, err := proxy.RoundRobinProxySwitcher("socks5://192.168.0.4:1070")
	//if err != nil {
	//	log.Fatalf("RoundRobinProxySwitcher Error: %v", err)
	//}
	//c.SetProxyFunc(proxyFunc)

	// 请求超时
	c.SetRequestTimeout(10000 * time.Millisecond)

	// 拦截
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", util.UserAgent)
		r.Headers.Set("Cookie", sukebei.Cookie)
		log.Println("Visiting:", r.URL)
	})

	c.OnError(sukebei.ErrorHandler)

	// 下载
	c.OnResponse(sukebei.Save)

	// 详情
	c.OnHTML("body", sukebei.ParseInfo)

	// 列表
	//c.OnHTML("body", sukebei.ParseList)

	// 入口
	sukebei.VisitViews(c)
}
