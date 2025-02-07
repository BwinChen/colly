package u3c3

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strconv"
	"strings"
	"time"
)

var Cookie = ""

func Visit(c *colly.Collector) {
	i := 10229
	for {
		if err := c.Visit(fmt.Sprintf("https://u3c3.com/?p=%d", i)); err != nil {
			log.Printf("Visit Error: %v\n", err)
		}
		i++
	}
}

func ParseHTML(body *colly.HTMLElement) {
	if strings.Contains(body.Request.URL.String(), "?p=") {
		body.ForEach("tr.default", func(i int, tr *colly.HTMLElement) {
			if i == 2 {
				tr.ForEach("td.text-center", func(j int, td *colly.HTMLElement) {
					if j == 2 {
						url := td.Request.URL.String()
						p, err := strconv.Atoi(url[strings.Index(url, "?p=")+len("?p="):])
						if err != nil {
							log.Printf("Atoi Error: %v\n", err)
							return
						}
						if p > 1 && strings.Contains(td.Text, strconv.Itoa(time.Now().Year())) {
							// 跳回第一页，爬取完成
							log.Fatalf("Done\n")
						}
					}
				})
			}
		})
		body.ForEach("a[href$='.torrent']", func(i int, a *colly.HTMLElement) {
			if i < 2 {
				return
			}
			href := a.Attr("href")
			infoHash := href[strings.Index(href, "/torrent/")+len("/torrent/") : strings.Index(href,
				".torrent")]
			fmt.Println("infoHash:", infoHash)
		})
	}
}
