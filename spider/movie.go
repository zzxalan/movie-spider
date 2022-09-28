package spider

type Movie struct {
	Name    string    // 影片名称
	Chapter []Chapter // 章节
}

// Chapter  章节
type Chapter struct {
	Url      string
	M3u8Url  string // 媒体地址
	Order    int    // 剧集
	OrderStr string // 剧集名称
}
