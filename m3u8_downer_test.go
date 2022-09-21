package movie

import (
	"fmt"
	"testing"
)

func TestDowner(t *testing.T) {
	err := Downer("https://m1.taopianplay1.com:43333/taopian/ecd7f271-487e-48d6-9873-9edc06e79ce8/6898c0aa-5d5a-44c4-b24b-d89aee30fc68/36491/e4ac6b92-c5f8-4ddb-8940-f71054f99212/SD/playlist.m3u8", "F:\\媒体库\\临时文件", "test")
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
