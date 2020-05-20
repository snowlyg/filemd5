package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	callApkMd5String      string
	faceApkMd5String      string
	inputApkMd5String     string
	iptvApkMd5String      string
	ttsApkMd5String       string
	websocketApkMd5String string

	callApkTimeString      int64
	faceApkTimeString      int64
	inputApkTimeString     int64
	iptvApkTimeString      int64
	ttsApkTimeString       int64
	websocketApkTimeString int64
)

func checkErrf(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func GetCurPath() string {

	str, _ := os.Getwd()
	return strings.Replace(str, "\\", "/", -1)

}

func getMd5(f *os.File) string {
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	w := io.MultiWriter(md5, sha1, sha256)
	if _, err := io.Copy(w, f); err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("%X", md5.Sum(nil))

	//fmt.Printf("File: %s\n", f.Name())
	//fmt.Printf("MD5: %X\n", md5.Sum(nil))
	//fmt.Printf("SHA-1: %X\n", sha1.Sum(nil))
	//fmt.Printf("SHA-256: %X\n\n", sha256.Sum(nil))
}

func inDirArray(name string) bool {
	dirs := []string{
		"com.chindeo.bed.app",
		"com.chindeo.launcher.app",
		"com.chindeo.nursehost",
		"com.ktcp.launcher",
	}

	for _, dir := range dirs {
		if name == dir {
			return true
		}
	}

	return false
}

func main() {

	path := GetCurPath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// 获取文件，并输出它们的名字
	for _, f := range files {
		if !f.IsDir() {
			file, err := os.Open(f.Name())
			if err != nil {
				log.Println(err)
				return
			}

			md5String := getMd5(file)
			switch f.Name() {
			case "call.apk":
				callApkMd5String = md5String
				callApkTimeString = f.ModTime().Unix()
			case "face.apk":
				faceApkMd5String = md5String
				faceApkTimeString = f.ModTime().Unix()
			case "input.apk":
				inputApkMd5String = md5String
				inputApkTimeString = f.ModTime().Unix()
			case "iptv.apk":
				iptvApkMd5String = md5String
				iptvApkTimeString = f.ModTime().Unix()
			case "tts.apk":
				ttsApkMd5String = md5String
				ttsApkTimeString = f.ModTime().Unix()
			case "websocket.apk":
				websocketApkMd5String = md5String
				websocketApkTimeString = f.ModTime().Unix()
			case "main.go":
				continue
			default:
				log.Println(fmt.Sprintf("其他文件:%v", f.Name()))
			}
			file.Close()
		}
	}

	for _, f := range files {
		if f.IsDir() && inDirArray(f.Name()) {
			file, _ := os.Open(f.Name())
			err = os.Chdir(path)
			checkErrf(err)

			names, err := file.Readdir(0)
			checkErrf(err)

			for _, s := range names {
				if s.IsDir() {
					continue
				}

				file, err := os.Open(f.Name() + "\\app.apk")
				if err != nil {
					log.Println(err)
					continue
				}
				md5String := getMd5(file)
				file.Close()

				appJson, err := os.OpenFile(f.Name()+"\\app.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
				if err != nil {
					checkErrf(err)
				}

				var data string
				switch f.Name() {
				case "com.chindeo.bed.app":
					comChindeoBedAppJson := "[{\"name\":\"app\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"tts\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"face\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"call\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"websocket\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"iptv\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"input\",\"hash\":\"%v\",\"version\":\"%d\"}]"
					data = fmt.Sprintf(comChindeoBedAppJson, md5String, f.ModTime().Unix(), ttsApkMd5String, ttsApkTimeString, faceApkMd5String, faceApkTimeString, callApkMd5String, callApkTimeString, websocketApkMd5String, websocketApkTimeString, iptvApkMd5String, iptvApkTimeString, inputApkMd5String, inputApkTimeString)

				case "com.chindeo.launcher.app":
					comChindeoBedAppJson := "[{\"name\":\"app\",\"hash\":\"%v\",\"version\":\"%d\"}]"
					data = fmt.Sprintf(comChindeoBedAppJson, md5String, f.ModTime().Unix())

				case "com.chindeo.nursehost":
					comChindeoBedAppJson := "[{\"name\":\"app\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"tts\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"call\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"websocket\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"input\",\"hash\":\"%v\",\"version\":\"%d\"}]"
					data = fmt.Sprintf(comChindeoBedAppJson, md5String, f.ModTime().Unix(), ttsApkMd5String, ttsApkTimeString, callApkMd5String, callApkTimeString, websocketApkMd5String, websocketApkTimeString, inputApkMd5String, inputApkTimeString)

				case "com.ktcp.launcher":
					comChindeoBedAppJson := "[{\"name\":\"app\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"tts\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"face\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"input\",\"hash\":\"%v\",\"version\":\"%d\"},{\"name\":\"call\",\"hash\":\"%v\",\"version\":\"%d\"}]"
					data = fmt.Sprintf(comChindeoBedAppJson, md5String, f.ModTime().Unix(), ttsApkMd5String, ttsApkTimeString, faceApkMd5String, faceApkTimeString, inputApkMd5String, inputApkTimeString, callApkMd5String, callApkTimeString)
				default:
					panic(fmt.Sprintf("错误目录:%v", f.Name()))
				}

				appJson.WriteString(data)

				appJson.Close()
			}
		}
	}
}
