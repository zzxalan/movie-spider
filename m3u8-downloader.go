package movie

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/grafov/m3u8"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"movie-spider/decrypt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Downer(url, path, name string) error {
	now := time.Now()
	if path == "" {
		return fmt.Errorf("请输入保存路径")
	}
	if name == "" {
		return fmt.Errorf("请输入名称")
	}
	// 初始化下载ts的目录，后面所有的ts文件会保存在这里
	var download_dir = filepath.Join(path, name)
	outFile := download_dir + ".mp4"
	if isExist, _ := pathExists(outFile); isExist {
		fmt.Println("视频已经存在:", outFile)
		return nil
	}
	if isExist, _ := pathExists(download_dir); !isExist {
		_ = os.MkdirAll(download_dir, os.ModePerm)
	}

	// m3u8解析出来的数据
	var m *m3u8.MediaPlaylist
	var key *[]byte
	var bar *progressbar.ProgressBar

	spider := colly.NewCollector(
		colly.MaxBodySize(100 * 1024 * 1024),
	)

	spider.OnResponse(func(response *colly.Response) {
		if strings.HasSuffix(response.Request.URL.Path, ".m3u8") {
			buf := bytes.NewBuffer(response.Body)
			playlist, listType, err := m3u8.DecodeFrom(buf, false)
			if err != nil {
				fmt.Println(err)
				return
			}
			if listType == m3u8.MASTER {
				for _, variant := range playlist.(*m3u8.MasterPlaylist).Variants {
					response.Request.Visit(variant.URI)
					return
				}
			}
			m = playlist.(*m3u8.MediaPlaylist)
			bar = progressbar.Default(int64(m.Count()))
			if m.Key != nil {
				response.Request.Visit(m.Key.URI)
			} else {
				spider.Async = true // 开启异步
				for i, segment := range m.Segments {
					if segment == nil {
						break
					}
					response.Request.Ctx = colly.NewContext()
					response.Request.Ctx.Put("segmentIndex", i)
					response.Request.Visit(segment.URI)
				}
			}
		}
		if strings.HasSuffix(response.Request.URL.Path, ".key") {
			// key
			key = &response.Body

			spider.Async = true // 开启异步
			for i, segment := range m.Segments {
				if segment == nil {
					break
				}
				response.Request.Ctx = colly.NewContext()
				response.Request.Ctx.Put("segmentIndex", i)
				response.Request.Visit(segment.URI)
			}
		}
		if strings.HasSuffix(response.Request.URL.Path, ".ts") {
			origData := response.Body
			if key != nil {
				origData, _ = decrypt.Aes(origData, *key)
			}
			// 保存
			curr_path_file := filepath.Join(download_dir, fmt.Sprintf("%05d.ts", response.Ctx.GetAny("segmentIndex")))

			_ = ioutil.WriteFile(curr_path_file, origData, 0666)
			go bar.Add(1)
		}
	})

	spider.Visit(url)

	spider.Wait()

	// 合并ts文件
	outMv, _ := os.Create(outFile)
	defer outMv.Close()
	writer := bufio.NewWriter(outMv)
	_ = filepath.Walk(download_dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() || filepath.Ext(path) != ".ts" {
			return nil
		}
		bytes, _ := ioutil.ReadFile(path)
		_, err = writer.Write(bytes)
		return err
	})
	_ = writer.Flush()
	os.RemoveAll(download_dir)

	//5、输出下载视频信息
	fmt.Printf("\n[Success] 下载保存路径：%s | 共耗时: %6.2fs\n", outFile, time.Now().Sub(now).Seconds())
	return nil
}

// 判断文件是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
