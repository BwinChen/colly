package btsow

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

var URL = "https://btsow.boats/tags/page/1"
var Cookie = ""

func ParseList(b *colly.HTMLElement) {
	b.ForEach(".tag", func(_ int, a *colly.HTMLElement) {
		if a.DOM.Get(0).LastChild.Data != "" {
			//fmt.Println("tag:", a.DOM.Get(0).LastChild.Data)
			if err := a.Request.Visit(a.Attr("href")); err != nil {
				log.Println(err)
			}
		}
	})
	b.ForEach(".pagination a", func(_ int, a *colly.HTMLElement) {
		//fmt.Println("page:", a.Attr("href"))
		if err := a.Request.Visit(a.Attr("href")); err != nil {
			log.Println(err)
		}
	})
}

func ParseInfo(b *colly.HTMLElement) {
	b.ForEach("#magnetLink", func(_ int, t *colly.HTMLElement) {
		fmt.Println("title:", t.DOM.Prev().Text(), "magnet:", t.Text)
	})
	b.ForEach(".field", func(_ int, div *colly.HTMLElement) {
		if strings.Contains(div.Text, "Torrent Description") {
			// 此分支放在ParseList将导致冲突
			for _, a := range div.DOM.Parent().Siblings().Find("a").Nodes {
				for _, attr := range a.Attr {
					if attr.Key == "href" {
						//fmt.Println("href", div.Request.AbsoluteURL(attr.Val))
						if err := div.Request.Visit(attr.Val); err != nil {
							log.Println(err)
						}
					}
				}
			}
		} else if strings.Contains(div.Text, "Content Size:") {
			fmt.Println("size:", div.DOM.Next().Text())
		} else if strings.Contains(div.Text, "Convert On:") {
			fmt.Println("day:", div.DOM.Next().Text())
		} else if strings.Contains(div.Text, "File Name") {
			for _, d := range div.DOM.Parent().Siblings().Find(".file").Nodes {
				fmt.Println("name:", d.LastChild.Data, "size:", d.NextSibling.FirstChild.Data)
			}
		}
	})
}
