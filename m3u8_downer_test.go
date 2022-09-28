package movie

import (
	"fmt"
	"testing"
)

func TestDowner(t *testing.T) {

	url := "https://cctvwbcdtxy.liveplay.myqcloud.com/cctvwbcd/cdrmjzcctv1_1_td.m3u8"

	err := Downer(url, "F:\\媒体库\\临时文件", "cctv")
	fmt.Println(err)
	//d := NewDownloader(
	//	WithSavePath("F:\\媒体库\\临时文件"),
	//	WithUrl("https://s3.fsvod1.com/20220318/cLs5CMOE/index.m3u8"),
	//	WithName("二十不惑2_第29集"),
	//)
	//d.Run()
	//
	////d.Run("https://iqiyi.sd-play.com/20220830/5aqT7RNG/1200kb/hls/index.m3u8", "二十不惑2_第23集")
}

func TestName(t *testing.T) {
}
