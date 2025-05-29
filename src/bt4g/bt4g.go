package bt4g

import (
	"fmt"
	"github.com/BwinChen/colly/src/util"
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
	"strings"
)

var Cookie = ""
var homeURL = "https://bt4gprx.com"
var newURL = fmt.Sprintf("%s/new", homeURL)
var hrefsKey = "colly:bt4g:hrefs"

func Limit(c *colly.Collector) {
	//err := c.Limit(&colly.LimitRule{
	//	DomainGlob: "*",
	//	Delay:      1 * time.Second,
	//})
	//if err != nil {
	//	log.Fatalf("Limit Error: %v", err)
	//}
}

func Visit(c *colly.Collector) {
	c.AllowURLRevisit = true
	for {
		href := util.SPop(hrefsKey)
		if href == "" || util.RandomInt(1, 10) == 5 {
			log.Printf("Visiting: %s", newURL)
			if err := c.Visit(buildSplashURL(newURL)); err != nil {
				log.Printf("Visit Error: %v\n", err)
			}
		} else {
			tmp := href
			if !strings.Contains(tmp, "infoHashes") {
				tmp = buildSplashURL(tmp)
			}
			log.Printf("Visiting: %s", href)
			if err := c.Visit(tmp); err != nil {
				log.Printf("Visit Error: %v\n", err)
				util.SAdd(hrefsKey, href)
			}
		}
	}
}

func ParseHTML(body *colly.HTMLElement) {
	u := body.Request.URL.String()
	if strings.Contains(u, "%2Fnew") {
		body.ForEach("div.list-group-item", func(_ int, div *colly.HTMLElement) {
			div.ForEach("a.text-decoration-none", func(_ int, a *colly.HTMLElement) {
				href := a.Attr("href")
				href = fmt.Sprintf("%s%s", homeURL, href)
				util.SAdd(hrefsKey, href)
			})
		})
	} else if strings.Contains(u, "%2Fmagnet%2F") {
		body.ForEach("a.btn-primary", func(_ int, a *colly.HTMLElement) {
			href := a.Attr("href")
			infoHash := href[strings.Index(href, "/hash/")+len("/hash/") : strings.Index(href, "?name=")]
			href = fmt.Sprintf("%s?infoHashes=%s", util.DhtTorrentURL, infoHash)
			util.SAdd(hrefsKey, href)
		})
	}
}

func Save(r *colly.Response) {}

func ErrorHandler(r *colly.Response, err error) {}

func buildSplashURL(u string) string {
	wait := "1.5"
	timeout := "30.0"
	luaScript := `
		function main(splash, args)
		  splash:set_user_agent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		  assert(splash:go(args.url))
		  assert(splash:wait(1.5))
		  return {
			html = splash:html(),
			png = splash:png(),
			har = splash:har(),
		  }
		end
	`
	return fmt.Sprintf("%s?wait=%s&images=1&timeout=%s&url=%s&lua_source=%s",
		util.SplashURL, wait, timeout, url.QueryEscape(u), url.QueryEscape(luaScript))
}
