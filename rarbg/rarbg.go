package rarbg

import (
	"colly/util"
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strings"
	"time"
)

var URL = "https://rarbgprx.org/torrents.php?page=1"
var Cookie = "tzWHMELq=gkFrCnQx; tzWHMELq=gkFrCnQx; aby=2; tcc; skt=iupm6xlpqa; skt=iupm6xlpqa; expla=2"

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
				// 详情页
				err = td.Request.Visit(h)
				if err != nil {
					return
				}
			}
		})
	})
	b.ForEach("div#pager_links > a", func(_ int, a *colly.HTMLElement) {
		// 翻页
		if err := a.Request.Visit(a.Attr("href")); err != nil {
			return
		}
	})
}

func ParseInfo(b *colly.HTMLElement) {
	b.ForEach("table.lista-rounded td.header2[align='right']", func(_ int, td *colly.HTMLElement) {
		if strings.Contains(td.Text, "Torrent:") {
			var h string
			for i, a := range td.DOM.Siblings().Find("a").Nodes {
				if i == 0 {
					fmt.Println("title:", a.FirstChild.Data)
					for _, attr := range a.Attr {
						if attr.Key == "href" {
							h = attr.Val
						}
					}
				}
				if i == 1 {
					for _, attr := range a.Attr {
						if attr.Key == "href" {
							fmt.Println("magnet:", attr.Val)
							s := strings.Index(attr.Val, "btih:")
							e := strings.Index(attr.Val, "&dn=")
							td.Request.Ctx.Put("InfoHash", attr.Val[s+5:e])
							// 下载种子
							if err := td.Request.Visit(h); err != nil {
								return
							}
						}
					}
				}
			}
		}
		if strings.Contains(td.Text, "Size:") {
			fmt.Println("size:", td.DOM.Siblings().Text())
		}
		if strings.Contains(td.Text, "Show Files »") {
			for i, tr := range td.DOM.Siblings().Find("tr").Nodes {
				if i == 0 {
					continue
				}
				fmt.Println("file:", tr.FirstChild.LastChild.Data)
				fmt.Println("size:", tr.LastChild.FirstChild.Data)
			}
		}
	})
}
