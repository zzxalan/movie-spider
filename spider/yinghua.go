package spider

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

type Yinghua struct {
}

func (r Yinghua) Run(url string, start int) (Movie, error) {
	movie := Movie{}
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("#playlist1", func(e *colly.HTMLElement) {
		if e.Request.Ctx.GetAny("chapter") != nil {
			return
		}
		e.ForEach("a", func(i int, a *colly.HTMLElement) {
			if i >= start {
				chapter := Chapter{Order: i, OrderStr: a.Text}
				a.Request.Ctx.Put("chapter", chapter)
				a.Request.Visit(a.Attr("href"))
			}
		})
	})

	c.OnHTML(".myui-content__detail .title", func(e *colly.HTMLElement) {
		movie.Name = e.Text
	})
	c.OnHTML(".embed-responsive", func(e *colly.HTMLElement) {
		jsonStr := strings.TrimLeft(e.Text, "\n\t\t\tvar player_aaaa=")
		m := make(map[string]any)
		if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
			fmt.Println(err)
			return
		}
		chapter := e.Request.Ctx.GetAny("chapter").(Chapter)
		chapter.M3u8Url = m["url"].(string)
		movie.Chapter = append(movie.Chapter, chapter)
	})

	if err := c.Visit(url); err != nil {
		return movie, err
	}

	return movie, nil
}
