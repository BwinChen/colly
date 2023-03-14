package sukebei

import (
	"colly/util"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strconv"
)

var deadline = util.Deadline("-12h")
var URL = "https://sukebei.nyaa.si/?p=1"

func ParseList(b *colly.HTMLElement) {
	// 列表
	b.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
		t := tr.ChildAttr("td:nth-child(5)", "data-timestamp")
		if t != "" {
			t1, err := strconv.ParseInt(t, 10, 64)
			if err != nil {
				return
			}
			if deadline > t1 {
				// 早于截止时间，停止爬取
				os.Exit(0)
			}
		}
		tr.ForEach("td", func(i int, td *colly.HTMLElement) {
			if i == 1 {
				if err := td.Request.Visit(td.ChildAttr("a", "href")); err != nil {
					log.Println(err)
				}
			}
		})
	})
	// 分页
	b.ForEach("ul.pagination a", func(_ int, a *colly.HTMLElement) {
		h := a.Attr("href")
		if h != "#" {
			if err := a.Request.Visit(h); err != nil {
				log.Println(err)
			}
		}
	})
}

func ParseInfo(b *colly.HTMLElement) {
	title := b.ChildText("div > div:nth-child(5) > div.panel-heading > h3")
	if title != "" {
		fmt.Println("title:", title)
	}
	infoHash := b.ChildText("kbd")
	if infoHash != "" {
		fmt.Println("infoHash:", infoHash)
	}
	b.ForEach("div.panel-footer > a", func(i int, a *colly.HTMLElement) {
		if i == 0 {
			// 下载
			a.Request.Ctx.Put("InfoHash", infoHash)
			if err := a.Request.Visit(a.Attr("href")); err != nil {
				log.Println(err)
			}
		}
		if i == 1 {
			fmt.Println("magnet:", a.Attr("href"))
		}
	})
	b.ForEach(".torrent-file-list i.fa-file", func(_ int, i *colly.HTMLElement) {
		fmt.Println("file:", i.DOM.Parent().Text())
	})
}
