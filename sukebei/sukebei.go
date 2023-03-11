package sukebei

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

func Save(r *colly.Response) {
	if strings.Contains(r.Request.URL.String(), ".torrent") {
		if err := r.Save(fmt.Sprintf("./%s.torrent", r.Ctx.Get("InfoHash"))); err != nil {
			log.Println(err)
		}
	}
}

func ParseList(b *colly.HTMLElement) {
	// 列表
	b.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
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

func ParseInfo(div *colly.HTMLElement) {
	title := div.ChildText("div > div:nth-child(5) > div.panel-heading > h3")
	if title != "" {
		fmt.Println("title", title)
	}
	infoHash := div.ChildText("kbd")
	if infoHash != "" {
		fmt.Println("infoHash", infoHash)
	}
	div.ForEach("div.panel-footer > a", func(i int, a *colly.HTMLElement) {
		if i == 0 {
			// 下载
			a.Request.Ctx.Put("InfoHash", infoHash)
			if err := a.Request.Visit(a.Attr("href")); err != nil {
				log.Println(err)
			}
		}
		if i == 1 {
			fmt.Println("magnet: ", a.Attr("href"))
		}
	})
	div.ForEach(".torrent-file-list i.fa-file", func(_ int, i *colly.HTMLElement) {
		fmt.Println("file:", i.DOM.Parent().Text())
	})
}
