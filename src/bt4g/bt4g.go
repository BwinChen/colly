package bt4g

import (
	"fmt"
	util2 "github.com/BwinChen/colly/src/util"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"time"
)

var Cookie = "ge_js_validator_28=1730522647@28@b4e9ac25624c8970a90e01a008c04d86"

func Limit(c *colly.Collector) {
	err := c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      2 * time.Second,
	})
	if err != nil {
		log.Fatalf("Limit Error: %v", err)
	}
}

func Visit(c *colly.Collector) {
	c.AllowURLRevisit = true
	if err := c.Visit("https://bt4gprx.com/new"); err != nil {
		log.Printf("Visit Error: %v\n", err)
	}
}

func ParseHTML(body *colly.HTMLElement) {
	url := body.Request.URL.String()
	if strings.Contains(url, "/new") {
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
				ok := util2.SetNX(fmt.Sprintf("colly:bt4g:%s", id), a.Attr("title"), 5*time.Minute)
				if !ok {
					return
				}
				log.Printf("%s\n", creationTime)
				if err := a.Request.Visit(a.Request.AbsoluteURL(href)); err != nil {
					log.Printf("Visit Error: %v\n", err)
				}
			})
			if i == count-1 {
				if err := div.Request.Visit("https://bt4gprx.com/new"); err != nil {
					log.Printf("Visit Error: %v\n", err)
				}
			}
		})
	} else if strings.Contains(url, "/magnet/") {
		body.ForEach("a.btn-primary", func(i int, a *colly.HTMLElement) {
			href := a.Attr("href")
			infoHash := href[strings.Index(href, "/hash/")+len("/hash/") : strings.Index(href, "?name=")]
			if err := a.Request.Visit(fmt.Sprintf("http://%s:8080/dht/torrent?infoHash=%s", util2.IP, infoHash)); err != nil {
				log.Printf("Visit Error: %v\n", err)
			}
		})
	}
}

func Save(r *colly.Response) {}

func ErrorHandler(r *colly.Response, err error) {}
