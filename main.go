package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Md5TimeString struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
	Version string `json:"version"`
}

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

func getMd5TimeStrings(files []os.FileInfo, mts map[string]*Md5TimeString) error {
	for _, f := range files {
		if !f.IsDir() {
			fullName := f.Name()
			file, err := os.Open(fullName)
			if err != nil {
				return err
			}

			getMd5TimeString(fullName, file, f, mts)

			file.Close()
		}
	}
	return nil
}

func getMd5TimeString(fullName string, file *os.File, f os.FileInfo, mts map[string]*Md5TimeString) {
	names := strings.Split(fullName, ".")
	if len(names) == 2 && names[1] == "apk" {
		md5String := getMd5(file)
		mt := new(Md5TimeString)
		mt.Name = names[0]
		mt.Hash = md5String
		mt.Version = f.ModTime().Format("20060102150400")
		mts[names[0]] = mt
	}
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func createJsons(files []os.FileInfo, err error, path string, mts map[string]*Md5TimeString) {

	for _, f := range files {
		if f.IsDir() {
			file, _ := os.Open(f.Name())
			err = os.Chdir(path)
			checkErrf(err)

			names, err := file.Readdir(0)
			checkErrf(err)

			for _, s := range names {
				var mmt []*Md5TimeString
				if s.IsDir() {
					continue
				}

				jsonPath := f.Name() + "\\app.json"
				if FileExist(f.Name()+"\\app.apk") && FileExist(jsonPath) {
					file, err := os.Open(f.Name() + "\\app.apk")
					if err != nil {
						log.Println(err)
						continue
					}
					md5String := getMd5(file)
					timeString := f.ModTime().Format("20060102150400")
					file.Close()

					appJson, _ := ioutil.ReadFile(jsonPath)

					var jsondata interface{}
					_ = json.Unmarshal(appJson, &jsondata)
					jd, ok := jsondata.([]interface{})

					if ok {
						for k, v := range jd {
							switch v2 := v.(type) {
							case string:
								fmt.Println(k, "is string", v2)
							case int:
								fmt.Println(k, "is int", v2)
							case bool:
								fmt.Println(k, "is bool", v2)
							case []interface{}:
								fmt.Println(k, "is an array:")
								for i, iv := range v2 {
									fmt.Println(i, iv)
								}
							case map[string]interface{}:
								//fmt.Println(k, "is an map:")
								var mt Md5TimeString
								for i, iv := range v2 {
									if i == "name" {
										ivs, ok := iv.(string)
										if ok {
											mt.Name = ivs
										}
									} else if i == "hash" {
										ivs, ok := iv.(string)
										if ok {
											mt.Hash = ivs
										}
									} else if i == "version" {
										ivs, ok := iv.(string)
										if ok {
											mt.Version = ivs
										}
									}
								}

								// 赋值 app.apk
								if mt.Name == "app" {
									mt.Hash = md5String
									mt.Version = timeString
								}

								// 赋值其他 apk
								for name, item := range mts {
									if mt.Name == name {
										mt.Hash = item.Hash
										mt.Version = item.Version
									}
								}

								mmt = append(mmt, &mt)
							default:
								fmt.Println(k, "类型未知", v2)
							}
						}
					}

					data, _ := json.Marshal(mmt)
					ioutil.WriteFile(jsonPath, data, 0)
				}
			}
		}
	}

}

func main() {

	mts := map[string]*Md5TimeString{}

	path := GetCurPath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// 获取文件，并输出它们的名字
	err = getMd5TimeStrings(files, mts)
	if err != nil {
		panic(err)
	}

	//生成 .json 文件
	createJsons(files, err, path, mts)
}
