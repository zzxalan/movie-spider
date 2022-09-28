package spider

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// 全民影视解析 https://www.daga.cc/
type DaGa struct {
}

func (r DaGa) Run(url string) (Movie, error) {
	movie := Movie{}
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	c.OnResponse(func(response *colly.Response) {
		data := daGaData{}
		if err := json.Unmarshal(response.Body, &data); err != nil {
			fmt.Println(err)
			return
		}
		movieData := data.Data[0]
		movie.Name = movieData.Name
		for i, ep := range movieData.Source.Eps {
			movie.Chapter = append(movie.Chapter, Chapter{
				Order:    i,
				OrderStr: ep.Name,
				M3u8Url:  ep.Url,
			})
		}
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Host", "a1.m1907.cn:404")
		request.Headers.Set("Origin", "https://z2.m1907.cn:404")
		request.Headers.Set("Referer", "https://z2.m1907.cn:404/")
	})

	if err := c.Visit(fmt.Sprintf(urlTemp, url)); err != nil {
		return movie, err
	}

	return movie, nil
}

var urlTemp = "https://a1.m1907.cn:404/api/v/?z=1030abade6431f748eb120d39e107603&jx=%s&s1ig=11402"

type daGaData struct {
	Type string
	Data []struct {
		Name   string
		Year   string
		Source struct {
			Eps []struct {
				Name string
				Url  string
			}
		}
	}
}
