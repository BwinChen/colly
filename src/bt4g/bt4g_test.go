package bt4g

import (
	"fmt"
	"github.com/BwinChen/colly/src/util"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"net/url"
	"strings"
	"testing"
	"time"
)

const urlsKey = "selenium:bt4g:urls"

func TestSelenium(t *testing.T) {
	wd, err := setUp()
	if err != nil {
		log.Fatalf("setUp Error: %v\n", err)
	}
	defer func(wd selenium.WebDriver) {
		// 关闭浏览器
		if err := wd.Quit(); err != nil {
			log.Fatalf("Quit Error: %v", err)
		}
	}(wd)
	// 设置隐式等待时间
	err = wd.SetImplicitWaitTimeout(5 * time.Second)
	if err != nil {
		log.Fatalf("SetImplicitWaitTimeout Error: %v\n", err)
	}
	_, err = util.Del(urlsKey)
	if err != nil {
		log.Fatalf("Del Error: %v\n", err)
	}
	for {
		lLen, err := util.LLen(urlsKey)
		if err != nil {
			log.Fatalf("LLen Error: %v\n", err)
		}
		if lLen == 0 {
			_, err = util.RPush(urlsKey, "https://bt4gprx.com/new")
			if err != nil {
				log.Fatalf("RPush Error: %v\n", err)
			}
		}
		next, err := util.LPop(urlsKey)
		if err != nil {
			log.Fatalf("LPop Error: %v\n", err)
		}
		err = wd.Get(next)
		if err != nil {
			fmt.Printf("Get Error: %v\n", err)
			continue
		}
		currentURL, err := wd.CurrentURL()
		if err != nil {
			log.Printf("CurrentURL Error: %v\n", err)
			continue
		}
		if strings.Contains(currentURL, "/new") {
			divs, err := wd.FindElements(selenium.ByCSSSelector, "div.list-group-item")
			if err != nil {
				log.Printf("FindElements Error: %v\n", err)
				continue
			}
			for _, div := range divs {
				a, err := div.FindElement(selenium.ByCSSSelector, "a.text-decoration-none")
				if err != nil {
					log.Printf("FindElement Error: %v\n", err)
					continue
				}
				href, err := a.GetAttribute("href")
				if err != nil {
					log.Printf("GetAttribute Error: %v\n", err)
					continue
				}
				href, err = resolveRelativeURL(href)
				if err != nil {
					log.Printf("ResolveRelativeURL Error: %v\n", err)
					continue
				}
				id := href[strings.Index(href, "/magnet/")+len("/magnet/"):]
				ok := util.SetNX(fmt.Sprintf("selenium:bt4g:id:%s", id), href, 5*time.Minute)
				if ok {
					_, err = util.RPush(urlsKey, href)
					if err != nil {
						log.Printf("RPush Error: %v\n", err)
					}
				}
			}
		} else if strings.Contains(currentURL, "/magnet/") {
			a, err := wd.FindElement(selenium.ByCSSSelector, "a.btn-primary")
			if err != nil {
				log.Printf("FindElement Error: %v\n", err)
				continue
			}
			href, err := a.GetAttribute("href")
			if err != nil {
				log.Printf("GetAttribute Error: %v\n", err)
				continue
			}
			infoHash := href[strings.Index(href, "/hash/")+len("/hash/") : strings.Index(href, "?name=")]
			log.Printf("infoHash: %s\n", infoHash)
			_, err = util.RPush(urlsKey, fmt.Sprintf("%s?infoHashes=%s", util.DhtTorrentURL, infoHash))
			if err != nil {
				log.Printf("RPush Error: %v\n", err)
			}
		}
	}
}

func setUp() (selenium.WebDriver, error) {
	// 设置 Edge 选项
	capabilities := selenium.Capabilities{"browserName": "MicrosoftEdge"}
	f := chrome.Capabilities{
		Args: []string{
			//"--proxy-server=socks5://192.168.0.7:10808",
			"--start-maximized",
			"--disable-extensions",
		},
	}
	capabilities.AddChrome(f)
	// 启动 Selenium WebDriver 服务
	wd, err := selenium.NewRemote(capabilities, "http://127.0.0.1:4444/wd/hub")
	if err != nil {
		return nil, err
	}
	return wd, nil
}

// 将相对URL转换为绝对URL
func resolveRelativeURL(rawURL string) (string, error) {
	baseURL, err := url.Parse("https://bt4gprx.com/")
	if err != nil {
		return "", err
	}
	relativeURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(relativeURL).String(), nil
}
