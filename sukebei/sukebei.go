package sukebei

import (
	"fmt"
	"github.com/BwinChen/colly/util"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"path/filepath"
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
	if strings.Contains(b.Request.URL.String(), "/view/") {
		//m := util.Magnet{}
		//b.ForEach("h3", func(i int, h3 *colly.HTMLElement) {
		//	if i == 0 {
		//		m.Name = strings.Trim(h3.Text, "\n\t")
		//	}
		//})
		//var infoHash string
		//b.ForEach("div.col-md-1", func(_ int, div *colly.HTMLElement) {
		//	if strings.Contains(div.Text, "File size:") {
		//		m.Size, _ = util.ConvertSize(div.DOM.Next().Text())
		//	} else if strings.Contains(div.Text, "Info hash:") {
		//		infoHash = div.DOM.Next().Text()
		//	} else if strings.Contains(div.Text, "Date:") {
		//		ts, _ := div.DOM.Next().Attr("data-timestamp")
		//		sec, err := strconv.ParseInt(ts, 10, 64)
		//		if err != nil {
		//			log.Println(err)
		//			return
		//		}
		//		m.AddedTime = time.Unix(sec, 0).Format("2006-01-02 15:04:05")
		//	}
		//})
		//b.ForEach("div.panel-footer > a", func(i int, a *colly.HTMLElement) {
		//	h := a.Attr("href")
		//	if i == 0 {
		//		m.Torrent = a.Request.AbsoluteURL(h)
		//		a.Request.Ctx.Put("InfoHash", infoHash)
		//		//if err := a.Request.Visit(h); err != nil {
		//		//	log.Println(err)
		//		//}
		//	} else if i == 1 {
		//		m.Magnet = h
		//	}
		//})
		//b.ForEach(".torrent-file-list i.fa-file", func(_ int, i *colly.HTMLElement) {
		//	f := util.File{}
		//	f.Path = i.DOM.Get(0).NextSibling.Data
		//	f.Length, _ = util.ConvertSize(strings.Trim(i.DOM.Next().Text(), "()"))
		//	m.Files = append(m.Files, f)
		//})
		b.ForEach(".row kbd", func(i int, k *colly.HTMLElement) {
			k.Request.Ctx.Put("InfoHash", k.Text)
			u := strings.Split(k.Request.URL.String(), "/")
			k.Request.Ctx.Put("ID", u[len(u)-1])
		})
		b.ForEach(".panel-footer > a", func(i int, a *colly.HTMLElement) {
			if i == 0 {
				err := a.Request.Visit(a.Request.AbsoluteURL(a.Attr("href")))
				if err != nil {
					log.Println(err)
				}
			}
		})
	}
}

func VisitPages(c *colly.Collector) {
	if err := c.Visit(fmt.Sprintf("https://sukebei.nyaa.si/?p=%d", page)); err != nil {
		log.Fatal(err)
	}
}

func VisitViews(c *colly.Collector) {
	for i := 4192944; i > 4000000; i-- {
		r, err := util.SIsMember(strconv.Itoa(i))
		if err != nil {
			log.Printf("SIsMember Error: %v\n", err)
			continue
		}
		if r {
			// 去重
			continue
		}
		err = c.Visit(fmt.Sprintf("https://sukebei.nyaa.si/view/%d", i))
		if err != nil {
			log.Printf("Visit Error: %v\n", err)
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

func Save(r *colly.Response) {
	if strings.Contains(r.Request.URL.String(), ".torrent") {
		n := time.Now()
		f := fmt.Sprintf("./torrents/%d/%d/%d", n.Year(), n.Month(), n.Day())
		err := os.MkdirAll(f, 0777)
		if err != nil {
			log.Printf("MkdirAll Error: %v\n", err)
			return
		}
		f += fmt.Sprintf("/%s.torrent", r.Ctx.Get("InfoHash"))
		// 写入文件
		err = r.Save(f)
		if err != nil {
			log.Printf("Save Error: %v\n", err)
			return
		}
		t, err := util.ParseTorrent(f)
		if err != nil {
			log.Printf("ParseTorrent Error: %v\n", err)
			return
		}
		// es去重
		h, err := util.SearchByInfoHash(t.InfoHash)
		if err != nil {
			log.Printf("SearchByInfoHash Error: %v\n", err)
			return
		}
		if h > 0 {
			return
		}
		// 索引torrent
		i, err := util.IndexTorrent(*t)
		if err != nil {
			log.Printf("IndexTorrent Error: %v\n", err)
			return
		}
		log.Printf("ES Torrent id: %s\n", i)
		// redis记录ID以去重
		r, err := util.SAdd(r.Ctx.Get("ID"))
		if err != nil {
			log.Printf("SAdd Error: %v\n", err)
			return
		}
		log.Printf("%d ID added to Redis\n", r)
		// 删除文件
		err = os.Remove(f)
		if err != nil {
			log.Printf("Remove Error: %v\n", err)
			return
		}
		// 删除父目录
		_ = os.Remove(filepath.Dir(f))
	}
}
