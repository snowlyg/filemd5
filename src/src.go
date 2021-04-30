package src

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	ExcludeDir  string
	ExcludeFile string
)

type Md5TimeString struct {
	Name    string `json:"name"`
	Hash    string `json:"hash"`
	Version string `json:"version"`
}

func GetCurrPath() (string, error) {
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

func GetMd5TimeStrings(path string, files map[string]os.FileInfo, mts map[string]*Md5TimeString) error {
	for name, f := range files {
		fullName := filepath.Join(path, f.Name())
		file, err := os.Open(fullName)
		if err != nil {
			return err
		}
		defer file.Close()
		version := f.ModTime().Format("20060102150400")
		mt := getMd5TimeString(name, version, file, mts)
		if mt == nil {
			continue
		}
		fmt.Printf("%s : %+v\n", mt.Name, mt)
		mts[mt.Name] = mt

	}
	return nil
}

func getMd5TimeString(name, version string, file *os.File, mts map[string]*Md5TimeString) *Md5TimeString {

	md5String := getMd5(file)
	mt := new(Md5TimeString)
	mt.Name = name
	mt.Hash = md5String
	mt.Version = version
	return mt
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func isExcludeDir(name string) bool {
	if ExcludeDir == "" {
		return false
	}
	for _, dir := range strings.Split(ExcludeDir, ",") {
		if dir == name {
			return true
		}
	}
	return false
}

func isExcludeFile(name string) bool {
	if ExcludeFile == "" {
		return false
	}
	for _, dir := range strings.Split(ExcludeFile, ",") {
		if dir == name {
			return true
		}
	}
	return false
}

func GetDirs(path string) (map[string]fs.FileInfo, []fs.FileInfo, error) {
	var dirs []fs.FileInfo
	files := make(map[string]fs.FileInfo)
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return files, dirs, err
	}

	for _, f := range fs {
		if !f.IsDir() || isExcludeDir(f.Name()) {
			continue
		}
		dirs = append(dirs, f)
	}
	for _, f := range fs {
		name, ok := checkExt(f.Name(), "apk")
		if f.IsDir() || !ok {
			continue
		}
		files[name] = f
	}
	return files, dirs, nil
}

func checkExt(name, ext string) (string, bool) {
	names := strings.Split(name, ".")
	if len(names) != 2 || names[1] != ext {
		return "", false
	}
	if isExcludeFile(names[0]) {
		return "", false
	}
	return names[0], true
}

func checkApp(f os.FileInfo, names []os.FileInfo) (string, string, error) {
	md5String := ""
	timeString := ""
	for _, s := range names {
		app, err := os.Open(filepath.Join(f.Name(), s.Name()))
		if err != nil {
			return "", "", err
		}
		defer app.Close()
		md5String = getMd5(app)
		timeString = f.ModTime().Format("20060102150400")
	}

	return md5String, timeString, nil

}

func UpdateJson(f os.FileInfo, path string, mts map[string]*Md5TimeString) error {
	fullpath := filepath.Join(path, f.Name())
	file, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("open dir %s error : %v", f.Name(), err)
	}
	err = os.Chdir(path)
	if err != nil {
		return fmt.Errorf("chdir %s error : %v", path, err)
	}

	names, err := file.Readdir(0)
	if err != nil {
		return fmt.Errorf("readdir error : %w", err)
	}

	md5String, timeString, err := checkApp(f, names)
	if err != nil {
		return fmt.Errorf("checkApp error : %w", err)
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

	return nil
}
