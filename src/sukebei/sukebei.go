package sukebei

import (
	"fmt"
	"github.com/BwinChen/colly/src/util"
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

const idsKey = "colly:sukebei:ids"

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

func VisitPages(c *colly.Collector) {
	if err := c.Visit(fmt.Sprintf("https://sukebei.nyaa.si/?p=%d", page)); err != nil {
		log.Fatal(err)
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

func ParseHTML(body *colly.HTMLElement) {
	url := body.Request.URL.String()
	if strings.HasSuffix(url, "/") {
		body.ForEach("tr:nth-child(1) > td:nth-child(2) > a", func(_ int, a *colly.HTMLElement) {
			href := a.Attr("href")
			id, err := strconv.Atoi(href[strings.Index(href, "/view/")+len("/view/"):])
			if err != nil {
				log.Fatalf("Atoi Error: %v\n", err)
			}
			for i := id; i > 0; i-- {
				r, err := util.SIsMember(idsKey, strconv.Itoa(i))
				if err != nil {
					log.Printf("SIsMember Error: %v\n", err)
					continue
				}
				if r {
					// 去重
					log.Printf("duplicate ID %d\n", i)
					continue
				}
				err = body.Request.Visit(fmt.Sprintf("https://sukebei.nyaa.si/view/%d", i))
				if err != nil {
					log.Printf("Visit Error: %v\n", err)
				}
			}
		})
	} else if strings.Contains(url, "/view/") {
		var infoHash, id string
		body.ForEach(".row kbd", func(i int, kbd *colly.HTMLElement) {
			infoHash = kbd.Text
			kbd.Request.Ctx.Put("InfoHash", infoHash)
			url := kbd.Request.URL.String()
			id = url[strings.Index(url, "/view/")+len("/view/"):]
			kbd.Request.Ctx.Put("ID", id)
		})
		// es去重
		hit, err := util.SearchByInfoHash(infoHash)
		if err != nil {
			log.Printf("SearchByInfoHash Error: %v\n", err)
			return
		}
		if hit > 0 {
			r := util.SAdd(idsKey, id)
			if r == -1 {
				return
			}
			log.Printf("ID %s added to Redis\n", id)
			return
		}
		body.ForEach(".panel-footer > a", func(i int, a *colly.HTMLElement) {
			if i == 0 {
				href := a.Attr("href")
				if strings.HasPrefix(href, "magnet:?xt=urn:btih:") {
					err := a.Request.Visit(fmt.Sprintf("%s?infoHashes=%s", util.DhtTorrentURL, infoHash))
					if err != nil {
						log.Printf("Visit Error: %v\n", err)
						return
					}
					r := util.SAdd(idsKey, id)
					if r == -1 {
						return
					}
					log.Printf("ID %s added to Redis\n", id)
				} else {
					err := a.Request.Visit(a.Request.AbsoluteURL(href))
					if err != nil {
						log.Printf("Visit Error: %v\n", err)
					}
				}
			}
		})
	}
}

func Visit(c *colly.Collector) {
	err := c.Visit("https://sukebei.nyaa.si/")
	if err != nil {
		log.Fatalf("Visit Error: %v\n", err)
	}
}

func Save(r *colly.Response) {
	url := r.Request.URL.String()
	if strings.Contains(url, ".torrent") {
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
		// 删除文件
		defer func(f string) {
			err := os.Remove(f)
			if err != nil {
				log.Printf("Remove Error: %v\n", err)
				return
			}
			_ = os.Remove(filepath.Dir(f))
		}(f)
		// redis记录ID以去重
		defer func() {
			id := r.Ctx.Get("ID")
			r := util.SAdd(idsKey, id)
			if r == -1 {
				return
			}
			log.Printf("ID %s added to Redis\n", id)
		}()
		torrent, err := util.ParseTorrent(f)
		if err != nil {
			log.Printf("ParseTorrent Error: %v\n", err)
			return
		}
		// 索引torrent
		id, err := util.IndexTorrent(torrent)
		if err != nil {
			log.Printf("IndexTorrent Error: %v\n", err)
			return
		}
		log.Printf("ES Torrent id: %s\n", id)
	}
}

func ErrorHandler(r *colly.Response, err error) {
	url := r.Request.URL.String()
	if strings.Contains(url, "/view/") && r.StatusCode == 404 {
		id := url[strings.Index(url, "/view/")+len("/view/"):]
		r := util.SAdd(idsKey, id)
		if r == -1 {
			return
		}
		log.Printf("ID %s added to Redis\n", id)
	}
}

func Limit(c *colly.Collector) {}
