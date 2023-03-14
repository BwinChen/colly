package util

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strings"
	"time"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69"

func Deadline(duration string) int64 {
	n := time.Now()
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err)
	}
	n = n.Add(d)
	return n.Unix()
}

func Save(r *colly.Response) {
	if strings.Contains(r.Request.URL.String(), ".torrent") {
		n := time.Now()
		f := fmt.Sprintf("./torrents/%d/%d/%d", n.Year(), n.Month(), n.Day())
		err := os.MkdirAll(f, 0777)
		if err != nil {
			log.Println(err)
			return
		}
		f += fmt.Sprintf("/%s.torrent", r.Ctx.Get("InfoHash"))
		err = r.Save(f)
		if err != nil {
			log.Println(err)
		}
	}
}
