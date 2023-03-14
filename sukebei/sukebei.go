package sukebei

import (
	"colly/util"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

var deadline = util.Deadline("-12h")
var URL = "https://sukebei.nyaa.si/?p=1"
var Cookie = ""

func ParseList(b *colly.HTMLElement) {
	b.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
		var h string
		tr.ForEach("td", func(i int, td *colly.HTMLElement) {
			if i == 1 {
				h = td.ChildAttr("a", "href")
			}
			if i == 4 {
				t := td.Attr("data-timestamp")
				t1, err := strconv.ParseInt(t, 10, 64)
				if err != nil {
					return
				}
				if deadline > t1 {
					// 早于截止时间，停止爬取
					os.Exit(0)
				}
				err = td.Request.Visit(h)
				if err != nil {
					log.Println(err)
				}
			}
		})
	})
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
	b.ForEach("h3", func(i int, h3 *colly.HTMLElement) {
		if i == 0 {
			fmt.Println("title:", strings.Trim(h3.Text, "\n\t"))
		}
	})
	var infoHash string
	b.ForEach("div.col-md-1", func(i int, div *colly.HTMLElement) {
		if strings.Contains(div.Text, "File size:") {
			fmt.Println("size:", div.DOM.Siblings().Get(0).FirstChild.Data)
		}
		if strings.Contains(div.Text, "Info hash:") {
			infoHash = div.DOM.Siblings().Get(0).FirstChild.FirstChild.Data
			fmt.Println("infoHash:", infoHash)
		}
	})
	b.ForEach("div.panel-footer > a", func(i int, a *colly.HTMLElement) {
		if i == 0 {
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
