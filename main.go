package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"
)

type setting struct {
	POST          string
	VERSION       string
	UPOST         string
	DIR           string
	DIR0          string
	POWERWORD     string
	CookieExpires int
}

func main() {
	var AppStart setting
	var js, _ = ioutil.ReadFile("./App.json")
	var jsonerr = json.Unmarshal(js, &AppStart)
	if jsonerr != nil {
		fmt.Println(jsonerr)
		var goin string
		scanln, err := fmt.Scanln(&goin)
		if err != nil {
			fmt.Println(scanln)
			return
		}
		return
	}

	var CookieExpires time.Duration = time.Duration(AppStart.CookieExpires)
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		if ctx.URLParam("powerword") == AppStart.POWERWORD {
			ctx.UpsertCookie(&http.Cookie{Name: "Download_Licence", Value: randStr(32), Expires: time.Now().Add(CookieExpires * time.Second)})
		}
		UserIP := ctx.Request().RemoteAddr
		app.Logger().Info(ctx.Path(), UserIP)
		FileList, err := ioutil.ReadDir(AppStart.DIR)
		if err != nil {
			return
		}

		var outhtml = "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>DownloadPage By Go</title>\n</head>\n<body>\n<h1>DownloadPage By Go</h1>\n<hr>\n<ul>\n    "
		for _, onefile := range FileList {
			if onefile.IsDir() {
				outhtml = outhtml + "<li><a href=\"" + "\\DIR?dir=" + onefile.Name() + "\" style=\"background-color: yellow\">" + onefile.Name() + "</a></li>\n"
			} else {
				file, err := os.Stat(AppStart.DIR0 + onefile.Name())
				if err == nil {
					outhtml = outhtml + "<li><a href=\"" + "\\DW?file=" + onefile.Name() + "\">" + onefile.Name() + "------" + fmt.Sprintf("%d", file.Size()) + "B</a></li>\n"
				} else {
					return
				}
			}
		}

		outhtml = outhtml + "</ul>\n</body>\n</html>"
		htmllen, err := ctx.HTML(outhtml)
		if err != nil {
			app.Logger().Error(htmllen)
			return
		}
	})

	app.Get("/DW", func(ctx iris.Context) {
		file := ctx.URLParam("file")
		if ctx.URLParam("powerword") == AppStart.POWERWORD {
			ctx.UpsertCookie(&http.Cookie{Name: "Download_Licence", Value: randStr(32), Expires: time.Now().Add(CookieExpires * time.Second)})
		}
		CookieValue := ctx.GetCookie("Download_Licence")
		if CookieValue == "" {
			var outhtml = "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>DownloadPage By Go</title>\n</head>\n<body>\n<h1>DownloadPage By Go</h1>\n<hr>\n<h1>Incorrect password</h1>\n<h1>密码不正确</h1>\n</body>\n</html>"
			_, err := ctx.HTML(outhtml)
			if err != nil {
				return
			}
			return
		}
		app.Logger().Info(ctx.Path(), file, ctx.Request().RemoteAddr)

		fs := strings.Split(file, "\\")
		err := ctx.SendFile(AppStart.DIR0+file, fs[len(fs)-1])
		if err != nil {
			app.Logger().Error(err.Error())
		}
	})

	app.Get("/DIR", func(ctx iris.Context) {
		dir := ctx.URLParam("dir")
		app.Logger().Info(ctx.Path(), dir, ctx.Request().RemoteAddr)

		FileList, err := ioutil.ReadDir(AppStart.DIR0 + dir)
		if err != nil {
			return
		}

		var outhtml = "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>DownloadPage By Go</title>\n</head>\n<body>\n<h1>DownloadPage By Go</h1>\n<hr>\n<ul>\n    "
		for _, onefile := range FileList {
			if onefile.IsDir() {
				outhtml = outhtml + "<li><a href=\"" + "\\DIR?dir=" + dir + "\\" + onefile.Name() + "\" style=\"background-color: yellow\">" + onefile.Name() + "</a></li>\n"
			} else {
				file, err := os.Stat(AppStart.DIR0 + dir + "\\" + onefile.Name())
				if err == nil {
					outhtml = outhtml + "<li><a href=\"" + "\\DW?file=" + dir + "\\" + onefile.Name() + "\">" + onefile.Name() + "------" + fmt.Sprintf("%d", file.Size()) + "B</a></li>\n"
				} else {
					return
				}
			}
		}

		outhtml = outhtml + "</ul>\n</body>\n</html>"
		_, err = ctx.HTML(outhtml)
		if err != nil {
			return
		}
	})

	app.Get("/del-cookie", func(ctx iris.Context) {
		ctx.RemoveCookie("Download_Licence")
		_, err := ctx.WriteString("Download_Licence removed")
		if err != nil {
			return
		}
	})

	err := app.Run(iris.Addr(AppStart.UPOST))
	if err != nil {
		return
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
