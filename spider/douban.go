package spider

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"goApp/pkg"
	netUrl "net/url"
	"regexp"
)

type DouBan struct {
}

func (r DouBan) Run(url, tag string) (Movie, error) {
	movie := Movie{}
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	reg := regexp.MustCompile("https://img3.doubanio.com/misc/mixed_static/.*\\.js")
	c.OnHTML("script", func(e *colly.HTMLElement) {
		if e.Request.Ctx.GetAny("chapter") == nil {
			return
		}
		src := e.Attr("src")
		if src != "" {
			if reg.FindString(src) != "" {
				e.Request.Ctx.Put("js", true)
				e.Request.Visit(src)
			}
		}
	})

	c.OnHTML(".episode_list", func(e *colly.HTMLElement) {
		if e.Request.Ctx.GetAny("chapter") != nil {
			return
		}
		e.ForEach(".item", func(i int, a *colly.HTMLElement) {
			chapter := Chapter{Order: i, OrderStr: a.Text}
			a.Request.Ctx.Put("chapter", chapter)
			a.Request.Visit(a.Attr("href"))
		})
	})

	jsonReg := regexp.MustCompile("videos = (.*),\n")
	urlReg := regexp.MustCompile("url=(.*)")
	c.OnResponse(func(response *colly.Response) {
		if response.Ctx.GetAny("js") == nil {
			return
		}
		jsonStr := pkg.Bytes2String(response.Body)
		s := jsonReg.FindStringSubmatch(jsonStr)
		if s != nil && len(s) > 0 {
			jsonStr = s[len(s)-1]
			m := make(map[string]any)
			if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
				fmt.Println(err)
				return
			}

			list := m["data"].([]any)
			for _, item := range list {
				from := item.(map[string]any)["source"].(map[string]any)["name"].(string)
				if tag == "" || tag == from {
					u := item.(map[string]any)["play_link"].([]any)[0].(string)
					ux := urlReg.FindStringSubmatch(u)

					chapter := response.Ctx.GetAny("chapter").(Chapter)
					chapter.Url, _ = netUrl.QueryUnescape(ux[len(ux)-1])
					movie.Chapter = append(movie.Chapter, chapter)
				}
			}
		}
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Host", "movie.douban.com")
		//request.Headers.Set("Cookie", "bid=vbrQpYJJrEM; douban-fav-remind=1; ll=\"118201\"; __utmc=30149280; __utmc=223695111; _vwo_uuid_v2=D36AB948AB503E4ED24B03B1F48F2E812|5f2f76fbd6e7dfc3e52089f1a99859e8; __utmz=30149280.1662628288.3.3.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmz=223695111.1662628296.2.2.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/search; __utma=30149280.1378295524.1662618689.1662628288.1662639743.4; __utmb=30149280.0.10.1662639743; __utma=223695111.1274697979.1662623163.1662628296.1662639743.3; __utmb=223695111.0.10.1662639743; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1662639744%2C%22https%3A%2F%2Fwww.douban.com%2Fsearch%3Fq%3D%25E4%25BA%258C%25E5%258D%2581%25E4%25B8%258D%25E6%2583%2591%22%5D; _pk_ses.100001.4cf6=*; _pk_id.100001.4cf6=d31863d138a7e4d3.1662623163.3.1662640183.1662629750.; ct=y")
	})

	if err := c.Visit(url); err != nil {
		return movie, err
	}

	return movie, nil
}
