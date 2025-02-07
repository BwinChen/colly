package bt4g

import (
	"fmt"
	util2 "github.com/BwinChen/colly/src/util"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"net/url"
	"strings"
	"testing"
	"time"
)

const key = "selenium:bt4g:urls"

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
	_, err = util2.Del(key)
	if err != nil {
		log.Fatalf("Del Error: %v\n", err)
	}
	for {
		lLen, err := util2.LLen(key)
		if err != nil {
			log.Fatalf("LLen Error: %v\n", err)
		}
		if lLen == 0 {
			_, err = util2.RPush(key, "https://bt4gprx.com/new")
			if err != nil {
				log.Fatalf("RPush Error: %v\n", err)
			}
		}
		next, err := util2.LPop(key)
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
				ok := util2.SetNX(fmt.Sprintf("selenium:bt4g:id:%s", id), href, 5*time.Minute)
				if ok {
					_, err = util2.RPush(key, href)
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
			_, err = util2.RPush(key, fmt.Sprintf("http://%s:8080/dht/torrent?infoHash=%s", util2.IP, infoHash))
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
