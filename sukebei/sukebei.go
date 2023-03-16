package sukebei

import (
	"colly/util"
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
			} else if i == 4 {
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
	m := util.Magnet{}
	b.ForEach("h3", func(i int, h3 *colly.HTMLElement) {
		if i == 0 {
			m.Name = strings.Trim(h3.Text, "\n\t")
		}
	})
	var infoHash string
	b.ForEach("div.col-md-1", func(i int, div *colly.HTMLElement) {
		if strings.Contains(div.Text, "File size:") {
			m.Size, _ = util.ConvertSize(div.DOM.Next().Text())
		} else if strings.Contains(div.Text, "Info hash:") {
			infoHash = div.DOM.Next().Text()
			m.InfoHash = infoHash
		}
	})
	b.ForEach("div.panel-footer > a", func(i int, a *colly.HTMLElement) {
		h := a.Attr("href")
		if i == 0 {
			m.Torrent = a.Request.AbsoluteURL(h)
			a.Request.Ctx.Put("InfoHash", infoHash)
			if err := a.Request.Visit(h); err != nil {
				log.Println(err)
			}
		} else if i == 1 {
			m.Magnet = h
		}
	})
	b.ForEach(".torrent-file-list i.fa-file", func(_ int, i *colly.HTMLElement) {
		f := util.File{}
		f.Name = i.DOM.Get(0).NextSibling.Data
		f.Size, _ = util.ConvertSize(strings.Trim(i.DOM.Next().Text(), "()"))
		m.Files = append(m.Files, f)
	})
	if m.InfoHash != "" {
		log.Println(m)
		util.IndexRequest(m)
	}
}
