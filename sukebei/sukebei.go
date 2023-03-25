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
var Cookie = ""

func ParseList(b *colly.HTMLElement) {
	b.ForEach("tr", func(i int, tr *colly.HTMLElement) {
		var h string
		tr.ForEach("td", func(j int, td *colly.HTMLElement) {
			if j == 1 {
				h = td.ChildAttr("a", "href")
				//err := view(i, h, td)
				//if err != nil {
				//	log.Println(err)
				//	return
				//}
			} else if j == 4 {
				t, err := strconv.ParseInt(td.Attr("data-timestamp"), 10, 64)
				if err != nil {
					log.Println(err)
					return
				}
				if deadline > t {
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
				//防止跳回首页
				return
			}
			err = a.Request.Visit(h)
			if err != nil {
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
	b.ForEach("div.col-md-1", func(_ int, div *colly.HTMLElement) {
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

func VisitPages(c *colly.Collector) {
	if err := c.Visit(fmt.Sprintf("https://sukebei.nyaa.si/?p=%d", page)); err != nil {
		log.Fatal(err)
	}
}

func VisitViews(c *colly.Collector) {
	//3770000 -> 3750000
	//3750000 -> 3700000
	for i := 3750000; i > 3700000; i-- {
		v := fmt.Sprintf("https://sukebei.nyaa.si/view/%d", i)
		if util.Search(util.Checksum(v)) > 0 {
			// URL去重
			continue
		}
		if err := c.Visit(v); err != nil {
			log.Println(err)
		}
	}
}

func view(i int, h string, td *colly.HTMLElement) error {
	p, err := strconv.Atoi(strings.Split(td.Request.URL.String(), "p=")[1])
	if err != nil {
		return err
	}
	if p == 100 && i == 75 {
		var v int
		v, err = strconv.Atoi(h[strings.Index(h, "view/")+5:])
		if err != nil {
			return err
		}
		//爬取100页之后
		for k := v; k > 1919904; k-- {
			vk := fmt.Sprintf("https://sukebei.nyaa.si/view/%d", k)
			if util.Search(util.Checksum(vk)) > 0 {
				// URL去重
				continue
			}
			err = td.Request.Visit(vk)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
