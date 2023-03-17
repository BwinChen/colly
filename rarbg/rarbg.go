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
var Cookie = "tzWHMELq=gkFrCnQx; tzWHMELq=gkFrCnQx; aby=2; skt=qz6kh97yrx; skt=qz6kh97yrx; expla=1; tcc"

// 差7个时区
var deadline = util.Deadline("-31h")

func ParseList(b *colly.HTMLElement) {
	b.ForEach("tr.lista2", func(_ int, tr *colly.HTMLElement) {
		var h string
		tr.ForEach("td.lista", func(i int, td *colly.HTMLElement) {
			if i == 1 {
				h = td.ChildAttr("a", "href")
			} else if i == 2 {
				t, err := time.ParseInLocation("2006-01-02 15:04:05", td.Text, time.Local)
				if err != nil {
					log.Println(err)
					return
				}
				if t.Unix() < deadline {
					log.Println("已到截止时间，爬取完成")
					os.Exit(0)
				}
				if util.Search(util.Checksum(td.Request.AbsoluteURL(h))) > 0 {
					// URL去重
					return
				}
				err = td.Request.Visit(h)
				if err != nil {
					log.Println(err)
				}
			}
		})
	})
	b.ForEach("div#pager_links > a", func(_ int, a *colly.HTMLElement) {
		if err := a.Request.Visit(a.Attr("href")); err != nil {
			log.Println(err)
		}
	})
}

func ParseInfo(b *colly.HTMLElement) {
	m := util.Magnet{}
	b.ForEach("table.lista-rounded td.header2[align='right']", func(_ int, td *colly.HTMLElement) {
		if strings.Contains(td.Text, "Torrent:") {
			var h string
			a := td.DOM.Next().Children().Get(1)
			m.Name = a.FirstChild.Data
			for _, attr := range a.Attr {
				if attr.Key == "href" {
					h = attr.Val
					break
				}
			}
			a = td.DOM.Next().Children().Get(2)
			for _, attr := range a.Attr {
				if attr.Key == "href" {
					m.Magnet = attr.Val
					m.Torrent = td.Request.AbsoluteURL(h)
					//每小时只能下载30个种子
					//td.Request.Ctx.Put("InfoHash", util.InfoHash(attr.Val))
					//if err := td.Request.Visit(h); err != nil {
					//	log.Println(err)
					//}
					break
				}
			}
		} else if strings.Contains(td.Text, "Size:") {
			m.Size, _ = util.ConvertSize(td.DOM.Next().Text())
		} else if strings.Contains(td.Text, "Show Files »") {
			for i, tr := range td.DOM.Next().Find("tr").Nodes {
				if i == 0 {
					continue
				}
				f := util.File{}
				f.Name = strings.TrimSpace(tr.FirstChild.LastChild.Data)
				f.Size, _ = util.ConvertSize(tr.LastChild.FirstChild.Data)
				m.Files = append(m.Files, f)
			}
		}
	})
	if m.Magnet != "" {
		m.URL = util.Checksum(b.Request.URL.String())
		//log.Println(m)
		util.IndexRequest(m)
	}
}
