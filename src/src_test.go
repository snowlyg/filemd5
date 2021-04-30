package src

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_GetCurPath(t *testing.T) {
	t.Run("测试当前路径", func(t *testing.T) {
		currentPath, err := GetCurrPath()
		if err != nil {
			t.Fatal(err)
		}
		if currentPath != "D:/go/src/github.com/chindeo/filemd5/src" {
			t.Errorf("GetCurPath want '%s' and get %s", "D:/go/src/github.com/chindeo/filemd5/src", currentPath)
		}
	})
}

func Test_GetMd5TimeStrings(t *testing.T) {
	res := map[string]*Md5TimeString{
		"awake": &Md5TimeString{
			Name:    "awake",
			Hash:    "CE76E02F83ABAC2CC1A4898E2296F095",
			Version: "202011161004",
		},
		"call": &Md5TimeString{
			Name:    "call",
			Hash:    "C332F74197CFE06346B6D7F2A54234B8",
			Version: "202104200942",
		},
		"face": &Md5TimeString{
			Name:    "face",
			Hash:    "9D6C59B11AD2631D6E34F915E65E68EA",
			Version: "202101211433",
		},
		"input": &Md5TimeString{
			Name:    "input",
			Hash:    "7380448E9B40DEC4E839998CBA0BFEC7",
			Version: "202004231605",
		},
		"iptv": &Md5TimeString{
			Name:    "iptv",
			Hash:    "ABD153FE17CFAAF11CFA56C8AE68A257",
			Version: "202101210932",
		},
		"l": &Md5TimeString{
			Name:    "l",
			Hash:    "8B3CF88E0C943E546626244CD99A5371",
			Version: "202102261419",
		},
		"live": &Md5TimeString{
			Name:    "live",
			Hash:    "98A3B965B81C1223EC1C8BC88034F9B8",
			Version: "202012031447",
		},
		"monitor": &Md5TimeString{
			Name:    "monitor",
			Hash:    "E80B959A93CD2F47D08BCEE4EDF4967B",
			Version: "202103031927",
		},
		"monitorPartner": &Md5TimeString{
			Name:    "monitorPartner",
			Hash:    "71C0D14156CA3890EAA946CD6A701817",
			Version: "202103031956",
		},
		"pda": &Md5TimeString{
			Name:    "pda",
			Hash:    "039E9C8B2B7AAACA59FC1B6C3DE00F66",
			Version: "202101271626",
		},
		"tts": &Md5TimeString{
			Name:    "tts",
			Hash:    "B5A65F49E4156665E1044A09E57C876B",
			Version: "201911081749",
		},
		"watch": &Md5TimeString{
			Name:    "watch",
			Hash:    "5CBE56071DB4FD5C4333CE03187597F0",
			Version: "202101271648",
		},
		"websocket": &Md5TimeString{
			Name:    "websocket",
			Hash:    "1A1907C74666D1B186D77ADD3A4898C5",
			Version: "202103051018",
		},
		"websocket_old": &Md5TimeString{
			Name:    "websocket_old",
			Hash:    "C91278C5171C9BAEA94266D49B7A74D9",
			Version: "202101211552",
		},
		"x5resources": &Md5TimeString{
			Name:    "x5resources",
			Hash:    "FB703FCE977A5438606C0A4D3917A4F7",
			Version: "202104081917",
		},
		"x5resources_32": &Md5TimeString{
			Name:    "x5resources_32",
			Hash:    "745A10520C89D0E67129AE3FD18F4B42",
			Version: "202011241447",
		},
		"xwalk": &Md5TimeString{
			Name:    "xwalk",
			Hash:    "57943B84462308029C63A26DD9CCDEA5",
			Version: "202010200921",
		},
	}
	t.Run("测试 .apk hash", func(t *testing.T) {
		mts := map[string]*Md5TimeString{}
		path := "d:/go/src/github.com/chindeo/filemd5/data/android"
		files, _, err := GetDirs(path)
		if err != nil {
			t.Fatal(err)
		}
		err = GetMd5TimeStrings(path, files, mts)
		if err != nil {
			t.Fatal(err)
		}

		if len(mts) != len(res) {
			t.Errorf(".apk 文件个数 '%d' and get %d", len(res), len(mts))
		}

		for _, mt := range mts {
			for _, re := range res {
				if re.Name != mt.Name {
					continue
				}
				if re.Hash != mt.Hash {
					t.Errorf("%s hash want '%s' and get %s", re.Name, re.Hash, mt.Hash)
				}
			}
		}
	})
}

func Test_FileExist(t *testing.T) {
	args := []struct {
		Name string
		Path string
		B    bool
	}{
		{
			"awake.apk 测试文件是否存在",
			"D:/go/src/github.com/chindeo/filemd5/data/android/awake.apk",
			true,
		},
		{
			"awakedd.apk 测试文件是否存在",
			"D:/go/src/github.com/chindeo/filemd5/data/android/awakedd.apk",
			false,
		},
	}

	for _, arg := range args {
		t.Run(arg.Name, func(t *testing.T) {
			b := FileExist(arg.Path)
			if b != arg.B {
				t.Errorf("path '%s' not file exist", arg.Path)
			}
		})
	}

}

func Test_isExcludeDir(t *testing.T) {
	ExcludeDir = "com.chindeo.bed.app,com.chindeo.nursehost"
	args := []struct {
		Name string
		Path string
		B    bool
	}{
		{
			"测试过滤目录 com.chindeo.bed.app",
			"com.chindeo.bed.app",
			true,
		},
		{
			"测试过滤目录 com.chindeo.nursehost",
			"com.chindeo.nursehost",
			true,
		},
		{
			"测试过滤目录 com.chindeo.nurseho",
			"com.chindeo.nurseho",
			false,
		},
	}

	for _, arg := range args {
		t.Run(arg.Name, func(t *testing.T) {
			b := isExcludeDir(arg.Path)
			if b != arg.B {
				t.Errorf("dir '%s' not file exist", arg.Path)
			}
		})
	}
}

func Test_isExcludeFile(t *testing.T) {
	ExcludeFile = "awake.apk,call.apk"
	args := []struct {
		Name string
		Path string
		B    bool
	}{
		{
			"测试过滤 face.apk 文件",
			"face.apk",
			false,
		},
		{
			"测试过滤 tts.apk 文件",
			"tts.apk",
			false,
		},
		{
			"测试过滤 call.apk 文件",
			"call.apk",
			true,
		},
	}

	for _, arg := range args {
		t.Run(arg.Name, func(t *testing.T) {
			b := isExcludeFile(arg.Path)
			if b != arg.B {
				t.Errorf("dir '%s' not file exist", arg.Path)
			}
		})
	}
}

func Test_checkExt(t *testing.T) {
	args := []struct {
		Name     string
		FileName string
		FullName string
		Ext      string
		B        bool
	}{
		{
			"测试过滤 face.apk 文件",
			"face",
			"face.apk",
			"apk",
			true,
		},
		{
			"测试过滤 tts.apk 文件",
			"tts",
			"tts.apk",
			"apk",
			true,
		},
		{
			"测试过滤 call.apk 文件",
			"call",
			"call.json",
			"json",
			true,
		},
		{
			"测试过滤 face.apk 文件",
			"",
			"face.ts",
			"apk",
			false,
		},
		{
			"测试过滤 tts.apk 文件",
			"",
			"tts.ts",
			"apk",
			false,
		},
		{
			"测试过滤 call.apk 文件",
			"",
			"call.ts",
			"json",
			false,
		},
	}

	for _, arg := range args {
		t.Run(arg.Name, func(t *testing.T) {
			name, ok := checkExt(arg.FullName, arg.Ext)
			if ok != arg.B {
				t.Errorf("dir '%s' not file exist", arg.FullName)
			}
			if name != arg.FileName {
				t.Errorf("file '%s' not get %s", name, arg.FileName)
			}
		})
	}
}

func Test_GetDirs(t *testing.T) {
	path := "D:/go/src/github.com/chindeo/filemd5/data/android"
	dirNames := []string{"com.chindeo.bed.app", "com.chindeo.launcher.app", "com.chindeo.nursehost", "com.chindeo.webapp", "com.ktcp.launcher", "医院"}
	fileNames := map[string]string{
		"awake":          "awake.apk",
		"call":           "call.apk",
		"face":           "face.apk",
		"input":          "input.apk",
		"iptv":           "iptv.apk",
		"l":              "l.apk",
		"live":           "live.apk",
		"monitor":        "monitor.apk",
		"monitorPartner": "monitorPartner.apk",
		"pda":            "pda.apk",
		"tts":            "tts.apk",
		"watch":          "watch.apk",
		"websocket_old":  "websocket_old.apk",
		"websocket":      "websocket.apk",
		"x5resources_32": "x5resources_32.apk",
		"x5resources":    "x5resources.apk",
		"xwalk":          "xwalk.apk",
	}
	ExcludeDir = ""
	t.Run("测试 GetDirs", func(t *testing.T) {
		files, dirs, err := GetDirs(path)
		if err != nil {
			t.Fatalf("get dirs err %v", err)
		}
		if len(dirs) != 6 {
			t.Fatalf("dirs want %d and get %d", 6, len(dirs))
		}
		for i := 0; i < len(dirs); i++ {
			if dirs[i].Name() != dirNames[i] {
				t.Errorf("%s not eq %s", dirs[i].Name(), dirNames[i])
			}
		}
		if len(files) != 17 {
			t.Fatalf("files want %d and get %d", 17, len(files))
		}

		for name, _ := range files {
			if files[name].Name() != fileNames[name] {
				t.Errorf("%s not eq %s", files[name].Name(), fileNames[name])
			}
		}

	})
}

func Test_checkApp(t *testing.T) {
	res := []struct {
		Name    string
		Path    string
		Hash    string
		Version string
	}{
		{
			Name:    "检查 com.chindeo.bed.app",
			Path:    "com.chindeo.bed.app/app.apk",
			Hash:    "CE76E02F83ABAC2CC1A4898E2296F095",
			Version: "202011161004",
		},
		{
			Name:    "检查 com.chindeo.launcher.app",
			Path:    "com.chindeo.launcher.app/app.apk",
			Hash:    "CE76E02F83ABAC2CC1A4898E2296F095",
			Version: "202011161004",
		},
		{
			Name:    "检查 com.chindeo.nursehost",
			Path:    "com.chindeo.nursehost/app.apk",
			Hash:    "CE76E02F83ABAC2CC1A4898E2296F095",
			Version: "202011161004",
		},
		{
			Name:    "检查 com.chindeo.webapp",
			Path:    "com.chindeo.webapp/app.apk",
			Hash:    "CE76E02F83ABAC2CC1A4898E2296F095",
			Version: "202011161004",
		},
		{
			Name:    "检查 com.ktcp.launcher",
			Path:    "com.ktcp.launcher/app.apk",
			Hash:    "CE76E02F83ABAC2CC1A4898E2296F095",
			Version: "202011161004",
		},
	}
	path := "D:/go/src/github.com/chindeo/filemd5/data/android"
	_, dirs, err := GetDirs(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, re := range res {
		t.Run(re.Name, func(t *testing.T) {
			for _, f := range dirs {
				fullpath := filepath.Join(path, f.Name())
				file, err := os.Open(fullpath)
				if err != nil {
					t.Error(err)
				}

				err = os.Chdir(path)
				if err != nil {
					t.Error(err)
				}

				names, err := file.Readdir(0)
				if err != nil {
					t.Error(err)
				}
				md5String, timeString, err := checkApp(f, names)
				if err != nil {
					t.Error(err)
				}
				if md5String != re.Hash {
					t.Errorf("hash want '%s' and get %s", re.Hash, md5String)
				}
				if timeString != re.Version {
					t.Errorf("version want '%s' and get %s", re.Version, timeString)
				}

			}
		})
	}
}
