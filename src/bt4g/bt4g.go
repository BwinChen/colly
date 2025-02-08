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

var Cookie = "ge_js_validator_28=1730522647@28@b4e9ac25624c8970a90e01a008c04d86"
var startURL = "https://bt4gprx.com/new"
var splashURL = "http://192.168.1.60:8050/render.html"

func Limit(c *colly.Collector) {
	//err := c.Limit(&colly.LimitRule{
	//	DomainGlob: "*",
	//	Delay:      2 * time.Second,
	//})
	//if err != nil {
	//	log.Fatalf("Limit Error: %v", err)
	//}
}

func Visit(c *colly.Collector) {
	c.AllowURLRevisit = true
	if err := c.Visit(buildSplashURL(startURL)); err != nil {
		log.Printf("Visit Error: %v\n", err)
	}
}

func ParseHTML(body *colly.HTMLElement) {
	u := body.Request.URL.String()
	if strings.Contains(u, "%2Fnew") {
		var count int
		body.ForEach("div.list-group-item", func(_ int, _ *colly.HTMLElement) {
			count++
		})
		body.ForEach("div.list-group-item", func(i int, div *colly.HTMLElement) {
			var creationTime string
			div.ForEach("p.mb-1 > span:nth-child(2)", func(_ int, span *colly.HTMLElement) {
				creationTime = strings.TrimSpace(span.Text)
			})
			div.ForEach("a.text-decoration-none", func(_ int, a *colly.HTMLElement) {
				href := a.Attr("href")
				id := href[strings.Index(href, "/magnet/")+len("/magnet/"):]
				ok := util.SetNX(fmt.Sprintf("colly:bt4g:%s", id), a.Attr("title"), 5*time.Minute)
				if !ok {
					return
				}
				log.Printf("%s\n", creationTime)
				if err := a.Request.Visit(buildSplashURL(buildAbsoluteURL(href))); err != nil {
					log.Printf("Visit Error: %v\n", err)
				}
			})
			if i == count-1 {
				if err := div.Request.Visit(buildSplashURL(startURL)); err != nil {
					log.Printf("Visit Error: %v\n", err)
				}
			}
		})
	} else if strings.Contains(u, "%2Fmagnet%2F") {
		body.ForEach("a.btn-primary", func(i int, a *colly.HTMLElement) {
			href := a.Attr("href")
			infoHash := href[strings.Index(href, "/hash/")+len("/hash/") : strings.Index(href, "?name=")]
			if err := a.Request.Visit(fmt.Sprintf("http://%s:8080/dht/torrent?infoHash=%s", util.IP, infoHash)); err != nil {
				log.Printf("Visit Error: %v\n", err)
			}
		})
	}
}

func Save(r *colly.Response) {}

func ErrorHandler(r *colly.Response, err error) {}

func buildSplashURL(u string) string {
	wait := "5.0"
	timeout := "90.0"
	luaScript := `
function main(splash, args)
  splash:set_user_agent("Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko")
  assert(splash:go(args.url))
  assert(splash:wait(5))
  return {
    html = splash:html(),
    png = splash:png(),
    har = splash:har(),
  }
end
`
	return fmt.Sprintf("%s?wait=%s&images=1&timeout=%s&url=%s&lua_source=%s",
		splashURL, wait, timeout, url.QueryEscape(u), url.QueryEscape(luaScript))
}

func buildAbsoluteURL(u string) string {
	// 基础 URL
	baseURL, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("url.Parse(%s) error: %v", u, err)
		return ""
	}
	// 相对 URL
	relativeURL, err := url.Parse(u)
	if err != nil {
		fmt.Printf("url.Parse(%s) error: %v", u, err)
		return ""
	}
	// 将相对 URL 转换为绝对 URL
	return baseURL.ResolveReference(relativeURL).String()
}
