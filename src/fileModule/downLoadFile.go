package fileModule

import (
	"io"
	"net/http"
	"os"
)

func downloadFile(url string, path string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 创建一个空的文件用于保存
	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 对接响应流和文件流
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}