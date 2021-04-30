package main

import (
	"flag"
	"fmt"

	"github.com/chindeo/filemd5/src"
)

var (
	excludeDir  = flag.String("exclude_dir", "com.chindeo.launcher.app", "过滤目录名称，多个目录用 ，隔开")
	excludeFile = flag.String("exclude_file", "", "过滤apk文件名称，多个文件用 ，隔开")
)

func main() {
	flag.Parse()

	src.ExcludeDir = *excludeDir
	src.ExcludeDir = *excludeFile

	mts := map[string]*src.Md5TimeString{}
	path, err := src.GetCurrPath()
	if err != nil {
		fmt.Printf("get current path error : %v", err)
		return
	}

	files, dirs, err := src.GetDirs(path)
	if err != nil {
		fmt.Printf("get current path error : %v", err)
		return
	}

	// 获取当前目录 apk 文件 ，并输获取对应文件md5值和文件更新时间
	err = src.GetMd5TimeStrings(path, files, mts)
	if err != nil {
		fmt.Printf("get md5 time string error : %v", err)
		return
	}

	for _, dir := range dirs {
		src.UpdateJson(dir, path, mts)
	}
}
