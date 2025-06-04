package bt4g

import (
	"fmt"
	"github.com/BwinChen/colly/src/util"
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
	"strings"
	"time"
)

var Cookie = ""
var homeURL = "https://bt4gprx.com"
var newURL = fmt.Sprintf("%s/new", homeURL)
var key = "colly:bt4g"
var hrefsKey = fmt.Sprintf("%s:hrefs", key)
var infoHashesKey = fmt.Sprintf("%s:infoHashes", key)

func Limit(c *colly.Collector) {
	err := c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      10 * time.Second, // 据估算每10秒一个请求可规避 502
	})
	if err != nil {
		log.Fatalf("Limit Error: %v", err)
	}
}

func Visit(c *colly.Collector) {
	c.AllowURLRevisit = true
	i := 0
	for {
		i++
		href := util.SPop(hrefsKey)
		if href == "" || i == 6 {
			i = 0
			log.Printf("Visiting: %s", newURL)
			if err := c.Visit(buildSplashURL(newURL)); err != nil {
				log.Printf("Visit Error: %v\n", err)
			}
		} else {
			log.Printf("Visiting: %s", href)
			if err := c.Visit(buildSplashURL(href)); err != nil {
				log.Printf("Visit Error: %v\n", err)
				util.SAdd(hrefsKey, href)
			}
		}
	}
}

func ParseHTML(body *colly.HTMLElement) {
	u := body.Request.URL
	if strings.Contains(u.String(), "%2Fnew") {
		body.ForEach("div.list-group-item", func(_ int, div *colly.HTMLElement) {
			div.ForEach("a.text-decoration-none", func(_ int, a *colly.HTMLElement) {
				util.SAdd(hrefsKey, fmt.Sprintf("%s%s", homeURL, a.Attr("href")))
			})
		})
	} else if strings.Contains(u.String(), "%2Fmagnet%2F") {
		count := 0
		body.ForEach("a.btn-primary", func(_ int, _ *colly.HTMLElement) {
			count++
		})
		if count == 0 {
			// 网络问题
			log.Printf("count: %d", count)
			util.SAdd(hrefsKey, u.Query().Get("url"))
			return
		}
		body.ForEach("a.btn-primary", func(_ int, a *colly.HTMLElement) {
			href := a.Attr("href")
			infoHash := href[strings.Index(href, "/hash/")+len("/hash/") : strings.Index(href, "?name=")]
			log.Printf("infoHash: %s", infoHash)
			util.SAdd(infoHashesKey, infoHash)
		})
	}
}

func Save(r *colly.Response) {}

func ErrorHandler(r *colly.Response, err error) {}

func buildSplashURL(u string) string {
	wait := "1.5"
	timeout := "30.0"
	// 必须设置user-agent
	luaScript := fmt.Sprintf(`
    function main(splash, args)
        splash:set_user_agent(%s)
        assert(splash:go(args.url))
        assert(splash:wait(%s))
        return {
            html = splash:html(),
            png = splash:png(),
            har = splash:har(),
        }
    end
`, util.UserAgent, wait)
	return fmt.Sprintf("%s?wait=%s&images=1&timeout=%s&url=%s&lua_source=%s",
		util.SplashURL, wait, timeout, url.QueryEscape(u), url.QueryEscape(luaScript))
}
