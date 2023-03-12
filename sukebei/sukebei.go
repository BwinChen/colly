package sukebei

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func ParseInfo(div *colly.HTMLElement) {
	title := div.ChildText("div > div:nth-child(5) > div.panel-heading > h3")
	if title != "" {
		fmt.Println("title:", title)
	}
	infoHash := div.ChildText("kbd")
	if infoHash != "" {
		fmt.Println("infoHash:", infoHash)
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
			fmt.Println("magnet:", a.Attr("href"))
		}
	})
	div.ForEach(".torrent-file-list i.fa-file", func(_ int, i *colly.HTMLElement) {
		fmt.Println("file:", i.DOM.Parent().Text())
	})
}

func Save(r *colly.Response) {
	if strings.Contains(r.Request.URL.String(), ".torrent") {
		n := time.Now()
		f := fmt.Sprintf("./torrents/%d/%d/%d", n.Year(), n.Month(), n.Day())
		if err := os.MkdirAll(f, 0777); err != nil {
			log.Println(err)
			return
		}
		f += fmt.Sprintf("/%s.torrent", r.Ctx.Get("InfoHash"))
		if err := r.Save(f); err != nil {
			log.Println(err)
		}
	}
}

func Deadline(duration string) int64 {
	n := time.Now()
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err)
	}
	n = n.Add(d)
	return n.Unix()
}

var deadline = Deadline("-12h")
