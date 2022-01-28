package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"os"
	"strings"
)

type setting struct {
	POST    string
	VERSION string
	UPOST   string
	DIR     string
	DIR0    string
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
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
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
		htmllen, err := ctx.HTML(outhtml)
		if err != nil {
			app.Logger().Error(htmllen)
			return
		}
	})

	err := app.Run(iris.Addr(AppStart.UPOST))
	if err != nil {
		return
	}
}
