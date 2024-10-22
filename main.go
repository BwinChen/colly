package main

import (
	"colly/sukebei"
	"colly/util"
	"github.com/gocolly/colly/v2"
	"log"
	"time"
)

func main() {
	c := colly.NewCollector()

	// 代理
	//pf, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:10808")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c.SetProxyFunc(pf)

	// 请求超时
	c.SetRequestTimeout(2000 * time.Millisecond)

	// 拦截
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", util.UserAgent)
		r.Headers.Set("Cookie", sukebei.Cookie)
		log.Println("Visiting:", r.URL)
	})

	// 下载
	c.OnResponse(util.Save)

	// 详情
	c.OnHTML("body", sukebei.ParseInfo)

	// 列表
	//c.OnHTML("body", sukebei.ParseList)

	// 入口
	sukebei.VisitViews(c)
}
