package rarbg

import (
	"colly/util"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strings"
	"time"
)

var URL = "https://rarbgprx.org/torrents.php?page=1"

// Cookie 绕过验证码
var Cookie = "tzWHMELq=gkFrCnQx; tzWHMELq=gkFrCnQx; aby=2; tcc; skt=iupm6xlpqa; skt=iupm6xlpqa"

// 差7个时区
var deadline = util.Deadline("-19h")

func ParseList(b *colly.HTMLElement) {
	b.ForEach("tr.lista2", func(_ int, tr *colly.HTMLElement) {
		var h string
		tr.ForEach("td.lista", func(i int, td *colly.HTMLElement) {
			if i == 1 {
				h = td.ChildAttr("a", "href")
			}
			if i == 2 {
				t, err := time.ParseInLocation("2006-01-02 15:04:05", td.Text, time.Local)
				if err != nil {
					return
				}
				if t.Unix() < deadline {
					// 早于截止时间，停止爬取
					os.Exit(0)
				}
				err = td.Request.Visit(h)
				if err != nil {
					return
				}
			}
		})
	})
	b.ForEach("div#pager_links > a", func(_ int, a *colly.HTMLElement) {
		if err := a.Request.Visit(a.Attr("href")); err != nil {
			return
		}
	})
}

func ParseInfo(b *colly.HTMLElement) {
	b.ForEach("table.lista-rounded td.header2[align='right']", func(_ int, td *colly.HTMLElement) {
		if strings.Contains(td.Text, "Torrent:") {
			var h string
			a := td.DOM.Next().Children().Get(1)
			log.Println("title:", a.FirstChild.Data)
			for _, attr := range a.Attr {
				if attr.Key == "href" {
					h = attr.Val
				}
			}
			a = td.DOM.Next().Children().Get(2)
			for _, attr := range a.Attr {
				if attr.Key == "href" {
					log.Println("magnet:", attr.Val)
					s := strings.Index(attr.Val, "btih:")
					e := strings.Index(attr.Val, "&dn=")
					td.Request.Ctx.Put("InfoHash", attr.Val[s+5:e])
					// 每小时只能下载30个种子
					//if err := td.Request.Visit(h); err != nil {
					//	return
					//}
					log.Println("torrent:", td.Request.AbsoluteURL(h))
				}
			}
		}
		if strings.Contains(td.Text, "Size:") {
			log.Println("size:", td.DOM.Next().Text())
		}
		if strings.Contains(td.Text, "Show Files »") {
			for i, tr := range td.DOM.Next().Find("tr").Nodes {
				if i == 0 {
					continue
				}
				log.Println("file:", strings.TrimSpace(tr.FirstChild.LastChild.Data))
				log.Println("size:", tr.LastChild.FirstChild.Data)
			}
		}
	})
}
