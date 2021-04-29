package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Md5TimeString struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
	Version string `json:"version"`
}

func GetCurPath() (string, error) {
	str, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Replace(str, "\\", "/", -1), nil

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
			defer file.Close()

			mt := getMd5TimeString(fullName, file, f, mts)
			if mt == nil {
				continue
			}
			fmt.Printf("%s : %+v\n", mt.Name, mt)
			mts[mt.Name] = mt

		}
	}
	return nil
}

func getMd5TimeString(fullName string, file *os.File, f os.FileInfo, mts map[string]*Md5TimeString) *Md5TimeString {
	names := strings.Split(fullName, ".")
	if len(names) != 2 || names[1] != "apk" {
		return nil
	}
	if isExcludeFile(names[0]) {
		return nil
	}
	md5String := getMd5(file)
	mt := new(Md5TimeString)
	mt.Name = names[0]
	mt.Hash = md5String
	mt.Version = f.ModTime().Format("20060102150400")
	return mt
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func isExcludeDir(name string) bool {
	if *excludeDir == "" {
		return false
	}
	for _, dir := range strings.Split(*excludeDir, ",") {
		if dir == name {
			return true
		}
	}
	return false
}

func isExcludeFile(name string) bool {
	if *excludeFile == "" {
		return false
	}
	for _, dir := range strings.Split(*excludeFile, ",") {
		if dir == name {
			return true
		}
	}
	return false
}

func checkApp(f, s os.FileInfo) (string, string, error) {
	app, err := os.Open(filepath.Join(f.Name(), s.Name()))
	if err != nil {
		return "", "", err
	}
	defer app.Close()
	return getMd5(app), f.ModTime().Format("20060102150400"), nil

}

func updateJson(f os.FileInfo, path string, mts map[string]*Md5TimeString) error {

	if f.IsDir() && !isExcludeDir(f.Name()) {
		file, err := os.Open(f.Name())
		if err != nil {
			return fmt.Errorf("open dir %s error : %v", f.Name(), err)
		}
		err = os.Chdir(path)
		if err != nil {
			return fmt.Errorf("chdir %s error : %v", path, err)
		}

		names, err := file.Readdir(0)
		if err != nil {
			return fmt.Errorf("Readdir error : %v", err)
		}

		md5String := ""
		timeString := ""
		for _, s := range names {
			pathName := filepath.Join(f.Name(), s.Name())
			if s.IsDir() || s.Name() != "app.apk" || !FileExist(pathName) {
				continue
			}

			md5String, timeString, err = checkApp(f, s)
			if err != nil {
				fmt.Printf("ReadFile %s error : %v\n", pathName, err)
				continue
			}

			fmt.Printf("read %s \n", pathName)
		}

		for _, s := range names {
			var mmt []*Md5TimeString
			pathName := filepath.Join(f.Name(), s.Name())
			if s.IsDir() || !FileExist(pathName) || !strings.Contains(s.Name(), ".json") {
				continue
			}

			appJson, err := ioutil.ReadFile(pathName)
			if err != nil {
				fmt.Printf("ReadFile %s error : %v\n", pathName, err)
				continue
			}

			var jsondata interface{}
			err = json.Unmarshal(appJson, &jsondata)
			if err != nil {
				fmt.Printf("json unmarshal error : %v\n", err)
				continue
			}
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

			data, err := json.Marshal(mmt)
			if err != nil {
				fmt.Printf("json marshal error : %v\n", err)
				continue
			}
			err = ioutil.WriteFile(pathName, data, os.ModePerm)
			if err != nil {
				fmt.Printf("write %s error : %v\n", pathName, err)
				continue
			}
			fmt.Printf("write json %v \n", pathName)
		}
	}
	return nil
}

var (
	excludeDir  = flag.String("exclude_dir", "com.chindeo.launcher.app", "过滤目录名称，多个目录用 ，隔开")
	excludeFile = flag.String("exclude_file", "", "过滤apk文件名称，多个文件用 ，隔开")
)

func main() {
	flag.Parse()
	mts := map[string]*Md5TimeString{}
	path, err := GetCurPath()
	if err != nil {
		fmt.Printf("get current path error : %v", err)
		return
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("get current path error : %v", err)
		return
	}
	// 获取当前目录 apk 文件 ，并输获取对应文件md5值和文件更新时间
	err = getMd5TimeStrings(files, mts)
	if err != nil {
		fmt.Printf("get md5 time string error : %v", err)
		return
	}
	for _, file := range files {
		updateJson(file, path, mts)
	}
}
