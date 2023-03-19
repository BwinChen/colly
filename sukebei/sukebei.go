package sukebei

import (
	"colly/util"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var deadline = util.Deadline(fmt.Sprintf("-%dh", 24*365))
var page = 1
var URL = fmt.Sprintf("https://sukebei.nyaa.si/?p=%d", page)
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
					log.Println(err)
					return
				}
				if deadline > t1 {
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
	b.ForEach("ul.pagination a", func(_ int, a *colly.HTMLElement) {
		h := a.Attr("href")
		if h != "#" {
			p, err := strconv.Atoi(strings.Split(h, "p=")[1])
			if err != nil || p < page {
				//防止重复爬取
				return
			}
			err = a.Request.Visit(h)
			if err != nil {
				log.Println(err)
			}
		}
	})
}

func ParseView(b *colly.HTMLElement) {
	b.ForEachWithBreak("ul.pagination a", func(_ int, a *colly.HTMLElement) bool {
		//1919904
		for v := 3839925; v > 1919904; v-- {
			h := fmt.Sprintf("https://sukebei.nyaa.si/view/%d", v)
			if util.Search(util.Checksum(h)) > 0 {
				// URL去重
				continue
			}
			if err := b.Request.Visit(h); err != nil {
				log.Println(err)
			}
		}
		return false
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
		} else if strings.Contains(div.Text, "Date:") {
			ts, _ := div.DOM.Next().Attr("data-timestamp")
			sec, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			m.AddedTime = time.Unix(sec, 0).Format("2006-01-02 15:04:05")
		}
	})
	b.ForEach("div.panel-footer > a", func(i int, a *colly.HTMLElement) {
		h := a.Attr("href")
		if i == 0 {
			m.Torrent = a.Request.AbsoluteURL(h)
			a.Request.Ctx.Put("InfoHash", infoHash)
			//if err := a.Request.Visit(h); err != nil {
			//	log.Println(err)
			//}
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
	if m.Magnet != "" {
		m.URL = util.Checksum(b.Request.URL.String())
		//log.Println(m)
		util.IndexRequest(m)
	}
}
